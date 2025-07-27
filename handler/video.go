package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/google/uuid"

	"video-agent-go/agent"
	"video-agent-go/model"
)

func RegisterRoutes(h *server.Hertz) {
	api := h.Group("/api/v1")

	// Video generation routes
	api.POST("/video/generate", GenerateVideo)                // 原有固定流程
	api.POST("/video/generate-smart", GenerateVideoSmart)     // LLM Agent编排
	api.POST("/video/generate-tools", GenerateVideoWithTools) // 🔧 新增：Tool-based编排

	api.GET("/video/status/:taskId", GetTaskStatus)
	api.GET("/video/list", GetAllTasks)

	// Tool-based 相关接口
	api.GET("/tools/list", ListAvailableTools)               // 🔧 查看可用工具
	api.GET("/tools/execution/:taskId", GetToolExecutionLog) // 🔧 查看工具调用日志

	// Agent 相关接口
	api.GET("/agents/list", ListAvailableAgents)
	api.GET("/execution/log/:taskId", GetExecutionLog)

	// Health check
	api.GET("/health", HealthCheck)
}

// 🔧 新增：Tool-based 视频生成接口
func GenerateVideoWithTools(ctx context.Context, c *app.RequestContext) {
	var input model.UserInput
	if err := json.Unmarshal(c.Request.Body(), &input); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Generate unique task ID
	taskID := uuid.New().String()

	// Save initial task
	inputJSON, _ := json.Marshal(input)
	if err := model.SaveTask(taskID, string(inputJSON), ""); err != nil {
		log.Printf("Failed to save task: %v", err)
		respondWithError(c, http.StatusInternalServerError, "Failed to create task")
		return
	}

	// 🔧 使用 Tool-based 智能编排
	go processVideoWithTools(taskID, input)

	respondWithData(c, model.ExtendedTaskStatus{
		TaskID:       taskID,
		Status:       "initializing",
		Progress:     0,
		CurrentStage: "tool_selection",
		ProcessingSteps: []string{
			"Analyzing user requirements with LLM",
			"LLM selecting appropriate tools",
			"Executing tools dynamically",
			"LLM orchestrating workflow",
			"Quality validation with tools",
		},
	})
}

// 🔧 新增：获取可用工具列表
func ListAvailableTools(ctx context.Context, c *app.RequestContext) {
	// 构造工具信息列表
	tools := []map[string]interface{}{
		{
			"name":        "analyze_content",
			"description": "Analyze user input to understand content requirements and suggest optimal processing strategy",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"user_text": map[string]interface{}{
						"type":        "string",
						"description": "User's original request text",
					},
				},
				"required": []string{"user_text"},
			},
		},
		{
			"name":        "generate_script",
			"description": "Generate video script based on user requirements and content analysis",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"content_type": map[string]interface{}{
						"type":        "string",
						"description": "Type of content (educational, commercial, entertainment)",
						"enum":        []string{"educational", "commercial", "entertainment", "news"},
					},
					"target_audience": map[string]interface{}{
						"type":        "string",
						"description": "Target audience for the video",
					},
				},
				"required": []string{"content_type", "target_audience"},
			},
		},
		{
			"name":        "generate_images",
			"description": "Generate images for video scenes using AI image generation",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"prompts": map[string]interface{}{
						"type":        "array",
						"description": "Array of image generation prompts",
					},
				},
				"required": []string{"prompts"},
			},
		},
		{
			"name":        "generate_voice",
			"description": "Generate voice narration for video content using text-to-speech",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"text": map[string]interface{}{
						"type":        "string",
						"description": "Text content to convert to speech",
					},
				},
				"required": []string{"text"},
			},
		},
		{
			"name":        "check_quality",
			"description": "Analyze and validate the quality of generated content",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"content_type": map[string]interface{}{
						"type":        "string",
						"description": "Type of content to check",
						"enum":        []string{"script", "images", "audio", "video"},
					},
				},
				"required": []string{"content_type"},
			},
		},
		{
			"name":        "render_video",
			"description": "Render final video from script, images, and audio components",
			"parameters": map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"script": map[string]interface{}{
						"type":        "object",
						"description": "Video script with timing information",
					},
				},
				"required": []string{"script"},
			},
		},
	}

	respondWithData(c, map[string]interface{}{
		"total_tools": len(tools),
		"tools":       tools,
	})
}

// 🔧 新增：获取工具执行日志
func GetToolExecutionLog(ctx context.Context, c *app.RequestContext) {
	taskID := c.Param("taskId")

	// 这里应该从持久化存储中获取工具调用日志
	// 简化实现，返回模拟数据
	executionLog := []map[string]interface{}{
		{
			"tool_call_id": "call_1",
			"tool_name":    "analyze_content",
			"arguments": map[string]interface{}{
				"user_text": "创建一个关于AI的教育视频",
			},
			"result": map[string]interface{}{
				"success":      true,
				"content_type": "educational",
				"complexity":   "medium",
			},
			"timestamp":   1640995200,
			"duration_ms": 500,
		},
		{
			"tool_call_id": "call_2",
			"tool_name":    "generate_script",
			"arguments": map[string]interface{}{
				"content_type":    "educational",
				"target_audience": "general",
			},
			"result": map[string]interface{}{
				"success": true,
				"title":   "AI Development Overview",
			},
			"timestamp":   1640995205,
			"duration_ms": 2500,
		},
		{
			"tool_call_id": "call_3",
			"tool_name":    "generate_images",
			"arguments": map[string]interface{}{
				"prompts": []string{
					"AI neural network visualization",
					"Machine learning algorithms diagram",
				},
			},
			"result": map[string]interface{}{
				"success": true,
				"count":   2,
			},
			"timestamp":   1640995210,
			"duration_ms": 8000,
		},
	}

	respondWithData(c, map[string]interface{}{
		"task_id":         taskID,
		"tool_calls":      executionLog,
		"total_calls":     len(executionLog),
		"execution_model": "LLM-driven tool selection",
	})
}

func GenerateVideo(ctx context.Context, c *app.RequestContext) {
	var input model.UserInput
	if err := json.Unmarshal(c.Request.Body(), &input); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Generate unique task ID
	taskID := uuid.New().String()

	// Save initial task
	inputJSON, _ := json.Marshal(input)
	if err := model.SaveTask(taskID, string(inputJSON), ""); err != nil {
		log.Printf("Failed to save task: %v", err)
		respondWithError(c, http.StatusInternalServerError, "Failed to create task")
		return
	}

	// Start async processing (原有的固定流程)
	go processVideo(taskID, input)

	respondWithData(c, model.TaskStatusResponse{
		TaskID: taskID,
		Status: "processing",
	})
}

// 新增：智能视频生成接口
func GenerateVideoSmart(ctx context.Context, c *app.RequestContext) {
	var input model.UserInput
	if err := json.Unmarshal(c.Request.Body(), &input); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Generate unique task ID
	taskID := uuid.New().String()

	// Save initial task
	inputJSON, _ := json.Marshal(input)
	if err := model.SaveTask(taskID, string(inputJSON), ""); err != nil {
		log.Printf("Failed to save task: %v", err)
		respondWithError(c, http.StatusInternalServerError, "Failed to create task")
		return
	}

	// 🚀 使用 LLM 驱动的智能编排
	go processVideoSmart(taskID, input)

	respondWithData(c, model.ExtendedTaskStatus{
		TaskID:       taskID,
		Status:       "analyzing",
		Progress:     0,
		CurrentStage: "task_analysis",
		ProcessingSteps: []string{
			"Analyzing user requirements",
			"Generating execution plan",
			"Selecting optimal agents",
			"Dynamic execution",
		},
	})
}

// 新增：获取可用智能体列表
func ListAvailableAgents(ctx context.Context, c *app.RequestContext) {
	// 构造智能体信息列表
	agents := []map[string]interface{}{
		{
			"name":         "ScriptGenerator",
			"description":  "Generates video scripts and storyboards using AI",
			"capabilities": []string{"script_generation", "storyboard_creation", "narrative_structure"},
		},
		{
			"name":         "ImageGenerator",
			"description":  "Generates images using DALL-E based on prompts",
			"capabilities": []string{"image_generation", "visual_creativity", "scene_creation"},
		},
		{
			"name":         "VoiceGenerator",
			"description":  "Generates voiceovers and narration using TTS",
			"capabilities": []string{"voice_synthesis", "narration", "multilingual_tts"},
		},
		{
			"name":         "VideoRender",
			"description":  "Renders and combines media into final video",
			"capabilities": []string{"video_rendering", "media_composition", "format_conversion"},
		},
		{
			"name":         "Analysis",
			"description":  "Analyzes content and provides insights for optimization",
			"capabilities": []string{"content_analysis", "sentiment_analysis", "engagement_prediction"},
		},
		{
			"name":         "QualityCheck",
			"description":  "Validates output quality and identifies issues",
			"capabilities": []string{"quality_validation", "error_detection", "compliance_check"},
		},
		{
			"name":         "Optimization",
			"description":  "Optimizes and improves content quality",
			"capabilities": []string{"content_optimization", "performance_enhancement", "quality_improvement"},
		},
	}

	respondWithData(c, map[string]interface{}{
		"total_agents": len(agents),
		"agents":       agents,
	})
}

// 新增：获取执行日志
func GetExecutionLog(ctx context.Context, c *app.RequestContext) {
	taskID := c.Param("taskId")

	// 这里应该从存储中获取执行日志
	// 简化实现，返回模拟数据
	executionLog := []map[string]interface{}{
		{
			"step_id":     "step_1",
			"agent_name":  "Analysis",
			"action":      "analyze_content",
			"timestamp":   1640995200,
			"duration_ms": 1500,
			"success":     true,
			"message":     "Content analysis completed",
		},
		{
			"step_id":     "step_2",
			"agent_name":  "ScriptGenerator",
			"action":      "generate_script",
			"timestamp":   1640995205,
			"duration_ms": 3000,
			"success":     true,
			"message":     "Generated script with 5 shots",
		},
	}

	respondWithData(c, map[string]interface{}{
		"task_id":       taskID,
		"execution_log": executionLog,
		"total_steps":   len(executionLog),
	})
}

func GetTaskStatus(ctx context.Context, c *app.RequestContext) {
	taskID := c.Param("taskId")

	task, err := model.GetTask(taskID)
	if err != nil {
		respondWithError(c, http.StatusNotFound, "Task not found")
		return
	}

	var result interface{}
	if task.Output != "" {
		json.Unmarshal([]byte(task.Output), &result)
	}

	status := "completed"
	if task.Output == "" {
		status = "processing"
	}

	respondWithData(c, model.TaskStatusResponse{
		TaskID: taskID,
		Status: status,
		Result: result,
	})
}

func GetAllTasks(ctx context.Context, c *app.RequestContext) {
	tasks, err := model.GetAllTasks()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, "Failed to get tasks")
		return
	}

	respondWithData(c, tasks)
}

func HealthCheck(ctx context.Context, c *app.RequestContext) {
	respondWithData(c, map[string]string{
		"status":  "healthy",
		"service": "video-agent-go",
		"modes": []string{
			"fixed-workflow",   // 固定流程
			"llm-orchestrated", // LLM编排
			"tool-based",       // 🔧 Tool-based
		}[2], // 显示当前支持的最高级模式
	})
}

// 🔧 新的Tool-based处理流程
func processVideoWithTools(taskID string, input model.UserInput) {
	log.Printf("🔧 Starting tool-based video processing for task: %s", taskID)

	// 创建Tool-based编排器
	orchestrator := agent.NewToolBasedOrchestrator()

	// 注册任务观察者
	observer := agent.GetObserverManager()
	observer.RegisterTask(taskID)
	observer.UpdateTask(taskID, agent.TaskProcessing, 5, "Initializing tool-based orchestrator")

	// 🎯 LLM + Tools 驱动的智能处理
	result, err := orchestrator.ProcessTask(taskID, input.Text)
	if err != nil {
		log.Printf("❌ Tool-based orchestration failed for task %s: %v", taskID, err)
		observer.UpdateTask(taskID, agent.TaskFailed, 0, fmt.Sprintf("Processing failed: %v", err))
		return
	}

	// 更新任务结果
	if err := model.UpdateTaskOutput(taskID, result); err != nil {
		log.Printf("Failed to update task output: %v", err)
	}

	observer.UpdateTask(taskID, agent.TaskCompleted, 100, "Video generated successfully using LLM + Tools orchestration")
	log.Printf("✅ Tool-based video processing completed for task: %s", taskID)
}

// 🚀 新的智能处理流程 - LLM驱动
func processVideoSmart(taskID string, input model.UserInput) {
	log.Printf("🧠 Starting LLM-driven video processing for task: %s", taskID)

	// 创建智能编排器
	orchestrator := agent.NewOrchestrator()

	// 注册任务观察者
	observer := agent.GetObserverManager()
	observer.RegisterTask(taskID)
	observer.UpdateTask(taskID, agent.TaskProcessing, 5, "Initializing LLM orchestrator")

	// 🎯 LLM 驱动的智能处理
	result, err := orchestrator.ProcessTask(taskID, input)
	if err != nil {
		log.Printf("❌ LLM orchestration failed for task %s: %v", taskID, err)
		observer.UpdateTask(taskID, agent.TaskFailed, 0, fmt.Sprintf("Processing failed: %v", err))
		return
	}

	// 更新任务结果
	if err := model.UpdateTaskOutput(taskID, result); err != nil {
		log.Printf("Failed to update task output: %v", err)
	}

	observer.UpdateTask(taskID, agent.TaskCompleted, 100, "Video generated successfully using LLM orchestration")
	log.Printf("✅ LLM-driven video processing completed for task: %s", taskID)
}

// 保留原有的固定流程处理函数
func processVideo(taskID string, input model.UserInput) {
	log.Printf("Processing video task: %s", taskID)

	// Step 1: Generate script
	script, err := agent.GenerateScript(input)
	if err != nil {
		log.Printf("Failed to generate script: %v", err)
		return
	}

	// Step 2: Process each shot
	for i := range script.Shots {
		// Generate image for shot
		imagePath, err := agent.GenerateImage(script.Shots[i].ImagePrompt)
		if err != nil {
			log.Printf("Failed to generate image: %v", err)
			continue
		}
		script.Shots[i].ClipPath = imagePath

		// Generate voiceover
		voicePath, err := agent.GenerateVoiceover(script.Shots[i].Voiceover)
		if err != nil {
			log.Printf("Failed to generate voiceover: %v", err)
			continue
		}
		script.Shots[i].VoicePath = voicePath
	}

	// Step 3: Render final video
	finalPath, err := agent.RenderVideo(*script)
	if err != nil {
		log.Printf("Failed to render video: %v", err)
		return
	}

	script.Final = finalPath
	script.TaskID = taskID
	script.Status = "completed"

	// Update task with result
	model.UpdateTaskOutput(taskID, script)
	log.Printf("Video task completed: %s", taskID)
}

func respondWithError(c *app.RequestContext, code int, message string) {
	c.JSON(code, model.APIResponse{
		Code:    code,
		Message: message,
	})
}

func respondWithData(c *app.RequestContext, data interface{}) {
	c.JSON(http.StatusOK, model.APIResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}
