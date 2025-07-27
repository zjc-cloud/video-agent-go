package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"video-agent-go/config"
	"video-agent-go/model"
)

// ToolBasedOrchestrator 基于工具的LLM编排器
type ToolBasedOrchestrator struct {
	toolRegistry  *ToolRegistry
	context       *ToolOrchestrationContext
	maxIterations int
}

// ToolOrchestrationContext 工具编排上下文
type ToolOrchestrationContext struct {
	TaskID       string                 `json:"task_id"`
	UserRequest  string                 `json:"user_request"`
	CurrentState map[string]interface{} `json:"current_state"`
	ToolCalls    []CompletedToolCall    `json:"tool_calls"`
	Resources    map[string]string      `json:"resources"`
}

// CompletedToolCall 完成的工具调用记录
type CompletedToolCall struct {
	Call      ToolCall    `json:"call"`
	Result    *ToolResult `json:"result"`
	Timestamp int64       `json:"timestamp"`
	Duration  int64       `json:"duration_ms"`
}

// LLMResponse LLM的响应结构
type LLMResponse struct {
	Role      string     `json:"role"`
	Content   string     `json:"content,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// ChatMessage 聊天消息结构
type ChatMessage struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	Name       string     `json:"name,omitempty"`
}

// NewToolBasedOrchestrator 创建基于工具的编排器
func NewToolBasedOrchestrator() *ToolBasedOrchestrator {
	registry := NewToolRegistry()

	// 注册所有可用工具
	registry.RegisterTool(&ContentAnalysisTool{})
	registry.RegisterTool(&ScriptGenerationTool{})
	registry.RegisterTool(&ImageGenerationTool{})
	registry.RegisterTool(&VoiceGenerationTool{})
	registry.RegisterTool(&QualityCheckTool{})
	registry.RegisterTool(&VideoRenderTool{})

	return &ToolBasedOrchestrator{
		toolRegistry:  registry,
		maxIterations: 10, // 防止无限循环
	}
}

// ProcessTask 处理任务 - 基于工具的智能编排
func (o *ToolBasedOrchestrator) ProcessTask(taskID string, userRequest string) (*model.ScriptOutput, error) {
	// 初始化上下文
	o.context = &ToolOrchestrationContext{
		TaskID:       taskID,
		UserRequest:  userRequest,
		CurrentState: make(map[string]interface{}),
		ToolCalls:    make([]CompletedToolCall, 0),
		Resources:    make(map[string]string),
	}

	log.Printf("🚀 Starting tool-based orchestration for task: %s", taskID)

	// 初始化消息历史
	messages := []ChatMessage{
		{
			Role:    "system",
			Content: o.buildSystemPrompt(),
		},
		{
			Role:    "user",
			Content: fmt.Sprintf("Please help me create a video based on this request: \"%s\"", userRequest),
		},
	}

	// 与LLM进行多轮对话，直到任务完成
	for iteration := 0; iteration < o.maxIterations; iteration++ {
		log.Printf("🔄 Iteration %d: Consulting LLM for next action", iteration+1)

		// 调用LLM获取下一步动作
		response, err := o.callLLMWithTools(messages)
		if err != nil {
			return nil, fmt.Errorf("LLM call failed: %v", err)
		}

		// 将LLM响应添加到消息历史
		messages = append(messages, ChatMessage{
			Role:      "assistant",
			Content:   response.Content,
			ToolCalls: response.ToolCalls,
		})

		// 检查是否有工具调用
		if len(response.ToolCalls) == 0 {
			// LLM认为任务已完成
			log.Printf("✅ LLM indicates task completion")
			break
		}

		// 执行所有工具调用
		for _, toolCall := range response.ToolCalls {
			result, err := o.executeToolCall(toolCall)
			if err != nil {
				log.Printf("❌ Tool call failed: %v", err)
				// 将错误信息添加到消息历史
				messages = append(messages, ChatMessage{
					Role:       "tool",
					Content:    fmt.Sprintf("Tool execution failed: %v", err),
					ToolCallID: toolCall.ID,
					Name:       toolCall.Function.Name,
				})
				continue
			}

			// 将工具结果添加到消息历史
			resultJSON, _ := json.Marshal(result.Data)
			messages = append(messages, ChatMessage{
				Role:       "tool",
				Content:    string(resultJSON),
				ToolCallID: toolCall.ID,
				Name:       toolCall.Function.Name,
			})

			// 更新上下文状态
			o.updateContextFromToolResult(toolCall.Function.Name, result)
		}
	}

	// 构建最终结果
	return o.buildFinalResult(), nil
}

// buildSystemPrompt 构建系统提示词
func (o *ToolBasedOrchestrator) buildSystemPrompt() string {
	toolsDescription := o.buildToolsDescription()

	return fmt.Sprintf(`You are an intelligent video generation orchestrator. Your job is to help users create videos by using the available tools strategically.

AVAILABLE TOOLS:
%s

INSTRUCTIONS:
1. Analyze the user's request to understand what type of video they want
2. Use tools step by step to gather information, generate content, and create the final video
3. Always start with content analysis to understand the requirements
4. Choose tools based on the specific needs of each request
5. Check quality before finalizing
6. Only indicate completion when you have successfully created a video file

DECISION MAKING PRINCIPLES:
- For educational content: prioritize accuracy and clarity
- For commercial content: focus on visual impact and persuasion  
- For entertainment: emphasize creativity and engagement
- Always consider the target audience
- Use quality checks for important content
- Be efficient but don't compromise on quality

You must use tools to accomplish tasks. Do not try to generate content directly.`, toolsDescription)
}

// buildToolsDescription 构建工具描述
func (o *ToolBasedOrchestrator) buildToolsDescription() string {
	var descriptions []string

	for _, tool := range o.toolRegistry.GetAllTools() {
		description := fmt.Sprintf("- %s: %s", tool.GetName(), tool.GetDescription())
		descriptions = append(descriptions, description)
	}

	return fmt.Sprintf("%s", descriptions)
}

// callLLMWithTools 调用LLM并支持工具调用
func (o *ToolBasedOrchestrator) callLLMWithTools(messages []ChatMessage) (*LLMResponse, error) {
	// 构建OpenAI API请求
	reqBody := map[string]interface{}{
		"model":       "gpt-4-turbo-preview",
		"messages":    messages,
		"tools":       o.toolRegistry.GetToolsSchema(),
		"tool_choice": "auto", // 让LLM自动决定是否调用工具
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.AppConfig.API.OpenAIKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var openAIResp struct {
		Choices []struct {
			Message LLMResponse `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return nil, err
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	return &openAIResp.Choices[0].Message, nil
}

// executeToolCall 执行工具调用
func (o *ToolBasedOrchestrator) executeToolCall(toolCall ToolCall) (*ToolResult, error) {
	startTime := time.Now()

	result, err := o.toolRegistry.ExecuteToolCall(toolCall)

	duration := time.Since(startTime).Milliseconds()

	// 记录工具调用
	completedCall := CompletedToolCall{
		Call:      toolCall,
		Result:    result,
		Timestamp: startTime.Unix(),
		Duration:  duration,
	}

	o.context.ToolCalls = append(o.context.ToolCalls, completedCall)

	return result, err
}

// updateContextFromToolResult 从工具结果更新上下文
func (o *ToolBasedOrchestrator) updateContextFromToolResult(toolName string, result *ToolResult) {
	if !result.Success {
		return
	}

	// 根据工具类型更新不同的状态
	switch toolName {
	case "analyze_content":
		if analysis, ok := result.Data.(map[string]interface{}); ok {
			o.context.CurrentState["content_analysis"] = analysis
		}
	case "generate_script":
		if script, ok := result.Data.(map[string]interface{}); ok {
			o.context.CurrentState["script"] = script
		}
	case "generate_images":
		if images, ok := result.Data.(map[string]interface{}); ok {
			o.context.CurrentState["images"] = images
		}
	case "generate_voice":
		if audio, ok := result.Data.(map[string]interface{}); ok {
			o.context.CurrentState["audio"] = audio
		}
	case "render_video":
		if video, ok := result.Data.(map[string]interface{}); ok {
			o.context.CurrentState["final_video"] = video
			if videoFile, ok := video["video_file"].(string); ok {
				o.context.Resources["final_video"] = videoFile
			}
		}
	case "check_quality":
		if quality, ok := result.Data.(map[string]interface{}); ok {
			o.context.CurrentState["quality_check"] = quality
		}
	}
}

// buildFinalResult 构建最终结果
func (o *ToolBasedOrchestrator) buildFinalResult() *model.ScriptOutput {
	// 从上下文状态构建最终输出
	result := &model.ScriptOutput{
		TaskID: o.context.TaskID,
		Status: "completed",
	}

	// 提取脚本信息
	if script, ok := o.context.CurrentState["script"].(map[string]interface{}); ok {
		if title, ok := script["title"].(string); ok {
			result.Title = title
		}
	}

	// 提取最终视频路径
	if finalVideo, ok := o.context.Resources["final_video"]; ok {
		result.Final = finalVideo
	}

	return result
}

// GetExecutionLog 获取执行日志
func (o *ToolBasedOrchestrator) GetExecutionLog() []CompletedToolCall {
	if o.context == nil {
		return []CompletedToolCall{}
	}
	return o.context.ToolCalls
}

// GetCurrentState 获取当前状态
func (o *ToolBasedOrchestrator) GetCurrentState() map[string]interface{} {
	if o.context == nil {
		return make(map[string]interface{})
	}
	return o.context.CurrentState
}
