package model

import "time"

type UserInput struct {
	Text           string                  `json:"text"`
	Images         []string                `json:"images"`
	Audio          string                  `json:"audio"`
	Video          string                  `json:"video"`
	Style          string                  `json:"style"`
	CustomScripts  []VideoProcessingScript `json:"custom_scripts,omitempty"`  // 新增：用户自定义脚本
	PluginSettings map[string]interface{}  `json:"plugin_settings,omitempty"` // 新增：插件配置
}

type Shot struct {
	Scene       string `json:"scene"`
	ImagePrompt string `json:"image_prompt"`
	Voiceover   string `json:"voiceover"`
	Duration    int    `json:"duration"`
	ClipPath    string `json:"clip_path,omitempty"`
	VoicePath   string `json:"voice_path,omitempty"`
	Subtitle    string `json:"subtitle,omitempty"`
}

type ScriptOutput struct {
	Title  string `json:"title"`
	Style  string `json:"style"`
	Shots  []Shot `json:"shots"`
	BGM    string `json:"bgm"`
	Final  string `json:"final,omitempty"`
	TaskID string `json:"task_id,omitempty"`
	Status string `json:"status,omitempty"`
}

// Database model
type VideoTask struct {
	ID        int       `json:"id"`
	TaskID    string    `json:"task_id"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
	CreatedAt time.Time `json:"created_at"`
}

// API Response structures
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type TaskStatusResponse struct {
	TaskID string      `json:"task_id"`
	Status string      `json:"status"`
	Result interface{} `json:"result,omitempty"`
}

// Error types
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}

// 新增：视频处理脚本结构
type VideoProcessingScript struct {
	Language string                 `json:"language"` // python, javascript, shell
	Code     string                 `json:"code"`
	Args     map[string]interface{} `json:"args"`
	Stage    string                 `json:"stage"` // pre_process, post_process, custom_effect
	Name     string                 `json:"name"`  // 脚本名称
	Version  string                 `json:"version,omitempty"`
}

// 新增：插件系统相关结构
type PluginInfo struct {
	Name         string                 `json:"name"`
	Version      string                 `json:"version"`
	Description  string                 `json:"description"`
	Author       string                 `json:"author"`
	Capabilities []string               `json:"capabilities"`
	Config       map[string]interface{} `json:"config"`
}

type PluginRegistry struct {
	Plugins map[string]PluginInfo `json:"plugins"`
}

// 新增：扩展的任务状态
type ExtendedTaskStatus struct {
	TaskID          string    `json:"task_id"`
	Status          string    `json:"status"`
	Progress        int       `json:"progress"`
	CurrentStage    string    `json:"current_stage"`
	ProcessingSteps []string  `json:"processing_steps"`
	ExecutedScripts []string  `json:"executed_scripts,omitempty"`
	Errors          []string  `json:"errors,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
