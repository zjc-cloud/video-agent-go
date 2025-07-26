package handler

import (
	"context"
	"encoding/json"
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
	api.POST("/video/generate", GenerateVideo)
	api.GET("/video/status/:taskId", GetTaskStatus)
	api.GET("/video/list", GetAllTasks)

	// Health check
	api.GET("/health", HealthCheck)
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

	// Start async processing
	go processVideo(taskID, input)

	respondWithData(c, model.TaskStatusResponse{
		TaskID: taskID,
		Status: "processing",
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
	})
}

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
