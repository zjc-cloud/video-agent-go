package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"video-agent-go/config"
	"video-agent-go/model"
)

// AgentOrchestrator LLM驱动的智能编排器
type AgentOrchestrator struct {
	availableAgents map[string]SubAgent
	executionLog    []ExecutionStep
	context         *OrchestrationContext
}

// SubAgent 子智能体接口
type SubAgent interface {
	GetName() string
	GetDescription() string
	GetCapabilities() []string
	Execute(ctx *OrchestrationContext, params map[string]interface{}) (*AgentResult, error)
	CanHandle(task string, context *OrchestrationContext) bool
}

// OrchestrationContext 编排上下文
type OrchestrationContext struct {
	TaskID        string                 `json:"task_id"`
	UserInput     model.UserInput        `json:"user_input"`
	CurrentState  map[string]interface{} `json:"current_state"`
	ExecutedSteps []ExecutionStep        `json:"executed_steps"`
	Resources     map[string]string      `json:"resources"` // 存储生成的文件路径等
}

// ExecutionStep 执行步骤
type ExecutionStep struct {
	StepID     string                 `json:"step_id"`
	AgentName  string                 `json:"agent_name"`
	Action     string                 `json:"action"`
	Parameters map[string]interface{} `json:"parameters"`
	Result     *AgentResult           `json:"result"`
	Timestamp  int64                  `json:"timestamp"`
	Duration   int64                  `json:"duration_ms"`
	Success    bool                   `json:"success"`
	ErrorMsg   string                 `json:"error_msg,omitempty"`
}

// AgentResult 智能体执行结果
type AgentResult struct {
	Success   bool                   `json:"success"`
	Data      map[string]interface{} `json:"data"`
	Resources map[string]string      `json:"resources"`
	NextSteps []string               `json:"next_steps,omitempty"`
	Message   string                 `json:"message"`
}

// ExecutionPlan LLM生成的执行计划
type ExecutionPlan struct {
	TaskAnalysis string        `json:"task_analysis"`
	Strategy     string        `json:"strategy"`
	Steps        []PlannedStep `json:"steps"`
	Reasoning    string        `json:"reasoning"`
}

// PlannedStep 计划的执行步骤
type PlannedStep struct {
	StepID     string                 `json:"step_id"`
	AgentName  string                 `json:"agent_name"`
	Action     string                 `json:"action"`
	Parameters map[string]interface{} `json:"parameters"`
	Condition  string                 `json:"condition,omitempty"`  // 执行条件
	Dependency []string               `json:"dependency,omitempty"` // 依赖的步骤
	Optional   bool                   `json:"optional,omitempty"`   // 是否可选
	Retry      int                    `json:"retry,omitempty"`      // 重试次数
}

// NewOrchestrator 创建新的编排器
func NewOrchestrator() *AgentOrchestrator {
	orchestrator := &AgentOrchestrator{
		availableAgents: make(map[string]SubAgent),
		executionLog:    make([]ExecutionStep, 0),
	}

	// 注册各种子智能体
	orchestrator.RegisterAgent(&ScriptGeneratorAgent{})
	orchestrator.RegisterAgent(&ImageGeneratorAgent{})
	orchestrator.RegisterAgent(&VoiceGeneratorAgent{})
	orchestrator.RegisterAgent(&VideoRenderAgent{})
	orchestrator.RegisterAgent(&AnalysisAgent{})
	orchestrator.RegisterAgent(&QualityCheckAgent{})
	orchestrator.RegisterAgent(&OptimizationAgent{})

	return orchestrator
}

// RegisterAgent 注册子智能体
func (o *AgentOrchestrator) RegisterAgent(agent SubAgent) {
	o.availableAgents[agent.GetName()] = agent
	log.Printf("Registered agent: %s", agent.GetName())
}

// ProcessTask 处理任务 - LLM驱动的主流程
func (o *AgentOrchestrator) ProcessTask(taskID string, input model.UserInput) (*model.ScriptOutput, error) {
	// 初始化上下文
	o.context = &OrchestrationContext{
		TaskID:        taskID,
		UserInput:     input,
		CurrentState:  make(map[string]interface{}),
		ExecutedSteps: make([]ExecutionStep, 0),
		Resources:     make(map[string]string),
	}

	log.Printf("🚀 Starting LLM-driven orchestration for task: %s", taskID)

	// Step 1: LLM 分析任务并生成执行计划
	plan, err := o.generateExecutionPlan()
	if err != nil {
		return nil, fmt.Errorf("failed to generate execution plan: %v", err)
	}

	log.Printf("📋 Execution plan generated: %s", plan.Strategy)

	// Step 2: 执行计划
	result, err := o.executePlan(plan)
	if err != nil {
		return nil, fmt.Errorf("failed to execute plan: %v", err)
	}

	log.Printf("✅ Task completed successfully: %s", taskID)
	return result, nil
}

// generateExecutionPlan LLM生成执行计划
func (o *AgentOrchestrator) generateExecutionPlan() (*ExecutionPlan, error) {
	// 构建提示词
	prompt := o.buildPlanningPrompt()

	// 调用LLM
	reqBody := OpenAIRequest{
		Model: "gpt-4",
		Messages: []Message{
			{
				Role: "system",
				Content: `You are an intelligent video generation orchestrator. Analyze the user's request and generate a detailed execution plan using available agents.

Available Agents:
- ScriptGenerator: Creates video scripts and storyboards
- ImageGenerator: Generates images using DALL-E
- VoiceGenerator: Creates voiceovers using TTS
- VideoRender: Renders and combines media into final video  
- Analysis: Analyzes content for optimization
- QualityCheck: Validates output quality
- Optimization: Improves and refines content

Return a JSON execution plan with reasoning for each step.`,
			},
			{Role: "user", Content: prompt},
		},
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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return nil, err
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	var plan ExecutionPlan
	if err := json.Unmarshal([]byte(openAIResp.Choices[0].Message.Content), &plan); err != nil {
		return nil, err
	}

	return &plan, nil
}

// buildPlanningPrompt 构建规划提示词
func (o *AgentOrchestrator) buildPlanningPrompt() string {
	input := o.context.UserInput

	prompt := fmt.Sprintf(`
Task Analysis Request:
- User Text: "%s"
- Video Style: "%s" 
- Has Images: %t
- Has Audio: %t
- Custom Scripts: %d

Current Context:
- Available Resources: %v
- Previous Steps: %d

Please analyze this video generation task and create an optimal execution plan. Consider:

1. What type of video is requested?
2. Which agents are needed and in what order?
3. Are there any special requirements or optimizations needed?
4. Should any steps be conditional or optional?
5. How can we ensure the best quality output?

Generate a detailed JSON execution plan with clear reasoning for agent selection and sequencing.
`,
		input.Text,
		input.Style,
		len(input.Images) > 0,
		input.Audio != "",
		len(input.CustomScripts),
		o.context.Resources,
		len(o.context.ExecutedSteps),
	)

	return prompt
}

// executePlan 执行LLM生成的计划
func (o *AgentOrchestrator) executePlan(plan *ExecutionPlan) (*model.ScriptOutput, error) {
	log.Printf("🎯 Executing plan: %s", plan.Strategy)

	var finalResult *model.ScriptOutput

	for _, step := range plan.Steps {
		// 检查执行条件
		if step.Condition != "" && !o.evaluateCondition(step.Condition) {
			log.Printf("⏭️  Skipping step %s: condition not met", step.StepID)
			continue
		}

		// 检查依赖
		if !o.checkDependencies(step.Dependency) {
			log.Printf("⏸️  Waiting for dependencies: %v", step.Dependency)
			continue
		}

		// 执行步骤
		log.Printf("▶️  Executing step: %s with agent: %s", step.StepID, step.AgentName)

		agent, exists := o.availableAgents[step.AgentName]
		if !exists {
			if step.Optional {
				log.Printf("⚠️  Optional agent %s not found, skipping", step.AgentName)
				continue
			}
			return nil, fmt.Errorf("required agent %s not found", step.AgentName)
		}

		// 执行代理
		result, err := o.executeAgent(agent, step)
		if err != nil {
			if step.Optional {
				log.Printf("⚠️  Optional step failed: %v", err)
				continue
			}
			return nil, fmt.Errorf("step %s failed: %v", step.StepID, err)
		}

		// 更新上下文
		o.updateContext(step, result)

		// 如果是最终结果
		if step.AgentName == "VideoRender" {
			finalResult = o.buildFinalResult(result)
		}

		// LLM 动态决定是否需要调整后续步骤
		if o.shouldReplan(step, result) {
			log.Printf("🔄 LLM suggests replanning based on current results")
			newPlan, err := o.generateAdaptivePlan(step, result)
			if err != nil {
				log.Printf("Failed to generate adaptive plan: %v", err)
			} else {
				plan.Steps = append(plan.Steps, newPlan.Steps...)
			}
		}
	}

	return finalResult, nil
}

// executeAgent 执行单个智能体
func (o *AgentOrchestrator) executeAgent(agent SubAgent, step PlannedStep) (*AgentResult, error) {
	startTime := getCurrentTimestamp()

	result, err := agent.Execute(o.context, step.Parameters)

	duration := getCurrentTimestamp() - startTime

	// 记录执行步骤
	execStep := ExecutionStep{
		StepID:     step.StepID,
		AgentName:  step.AgentName,
		Action:     step.Action,
		Parameters: step.Parameters,
		Result:     result,
		Timestamp:  startTime,
		Duration:   duration,
		Success:    err == nil,
	}

	if err != nil {
		execStep.ErrorMsg = err.Error()
	}

	o.executionLog = append(o.executionLog, execStep)
	o.context.ExecutedSteps = append(o.context.ExecutedSteps, execStep)

	return result, err
}

// updateContext 更新执行上下文
func (o *AgentOrchestrator) updateContext(step PlannedStep, result *AgentResult) {
	if result != nil && result.Success {
		// 更新资源
		for k, v := range result.Resources {
			o.context.Resources[k] = v
		}

		// 更新状态
		for k, v := range result.Data {
			o.context.CurrentState[k] = v
		}
	}
}

// shouldReplan 判断是否需要重新规划
func (o *AgentOrchestrator) shouldReplan(step PlannedStep, result *AgentResult) bool {
	// 这里可以实现更复杂的逻辑
	// 比如检查质量分数、用户反馈、错误率等

	if result != nil && len(result.NextSteps) > 0 {
		return true
	}

	// 如果是关键步骤且有质量问题，考虑重新规划
	if step.AgentName == "QualityCheck" && result != nil {
		if qualityScore, ok := result.Data["quality_score"].(float64); ok {
			return qualityScore < 0.7 // 质量分数低于0.7时重新规划
		}
	}

	return false
}

// generateAdaptivePlan 生成自适应计划
func (o *AgentOrchestrator) generateAdaptivePlan(currentStep PlannedStep, result *AgentResult) (*ExecutionPlan, error) {
	// 构建上下文感知的重规划提示
	_ = fmt.Sprintf(`
Current situation analysis:
- Just completed step: %s using %s
- Result: %s (Success: %t)
- Current resources: %v
- Quality indicators: %v

Based on the current results, should we:
1. Continue with the original plan?
2. Add quality improvement steps?
3. Try alternative approaches?
4. Add validation steps?

Generate additional steps or modifications to optimize the outcome.
`,
		currentStep.StepID,
		currentStep.AgentName,
		result.Message,
		result.Success,
		o.context.Resources,
		result.Data,
	)

	// 这里调用LLM生成自适应计划
	// 简化实现，实际可以调用OpenAI API
	return &ExecutionPlan{
		TaskAnalysis: "Adaptive replanning",
		Strategy:     "Quality optimization",
		Steps:        []PlannedStep{}, // 根据实际情况添加步骤
		Reasoning:    "Based on current results",
	}, nil
}

// 辅助函数
func (o *AgentOrchestrator) evaluateCondition(condition string) bool {
	// 简化的条件评估，实际可以实现更复杂的表达式解析
	return true
}

func (o *AgentOrchestrator) checkDependencies(dependencies []string) bool {
	for _, dep := range dependencies {
		found := false
		for _, step := range o.context.ExecutedSteps {
			if step.StepID == dep && step.Success {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (o *AgentOrchestrator) buildFinalResult(result *AgentResult) *model.ScriptOutput {
	// 从结果和上下文构建最终输出
	return &model.ScriptOutput{
		Title:  fmt.Sprintf("%v", o.context.CurrentState["title"]),
		Style:  fmt.Sprintf("%v", o.context.CurrentState["style"]),
		Final:  o.context.Resources["final_video"],
		TaskID: o.context.TaskID,
		Status: "completed",
	}
}

func getCurrentTimestamp() int64 {
	return 0 // 简化实现
}
