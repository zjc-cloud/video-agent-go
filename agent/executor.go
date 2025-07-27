package agent

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// CodeExecutor 代码执行器接口
type CodeExecutor interface {
	ExecutePython(code string, args map[string]interface{}) (interface{}, error)
	ExecuteJavaScript(code string, args map[string]interface{}) (interface{}, error)
	ExecuteShell(script string, args []string) (string, error)
}

// SafeExecutor 安全的代码执行器
type SafeExecutor struct {
	Timeout    time.Duration
	WorkingDir string
	Restricted bool // 是否启用沙箱模式
}

func NewSafeExecutor() *SafeExecutor {
	return &SafeExecutor{
		Timeout:    30 * time.Second,
		WorkingDir: "temp/executor",
		Restricted: true,
	}
}

// ExecutePython 执行Python代码
func (e *SafeExecutor) ExecutePython(code string, args map[string]interface{}) (interface{}, error) {
	// 创建临时Python文件
	tempFile := filepath.Join(e.WorkingDir, fmt.Sprintf("script_%d.py", time.Now().UnixNano()))

	// 确保目录存在
	if err := os.MkdirAll(e.WorkingDir, 0755); err != nil {
		return nil, err
	}

	// 写入代码到临时文件
	if err := os.WriteFile(tempFile, []byte(code), 0644); err != nil {
		return nil, err
	}
	defer os.Remove(tempFile)

	// 创建上下文with超时
	ctx, cancel := context.WithTimeout(context.Background(), e.Timeout)
	defer cancel()

	// 执行Python脚本
	cmd := exec.CommandContext(ctx, "python3", tempFile)
	cmd.Dir = e.WorkingDir

	if e.Restricted {
		// 在受限环境中执行
		cmd.Env = []string{
			"PATH=/usr/bin:/bin",
			"PYTHONPATH=/usr/lib/python3/dist-packages",
		}
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("python execution failed: %v, output: %s", err, string(output))
	}

	return string(output), nil
}

// ExecuteJavaScript 执行JavaScript代码
func (e *SafeExecutor) ExecuteJavaScript(code string, args map[string]interface{}) (interface{}, error) {
	// 创建临时JS文件
	tempFile := filepath.Join(e.WorkingDir, fmt.Sprintf("script_%d.js", time.Now().UnixNano()))

	if err := os.MkdirAll(e.WorkingDir, 0755); err != nil {
		return nil, err
	}

	if err := os.WriteFile(tempFile, []byte(code), 0644); err != nil {
		return nil, err
	}
	defer os.Remove(tempFile)

	ctx, cancel := context.WithTimeout(context.Background(), e.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "node", tempFile)
	cmd.Dir = e.WorkingDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("javascript execution failed: %v, output: %s", err, string(output))
	}

	return string(output), nil
}

// ExecuteShell 执行Shell脚本
func (e *SafeExecutor) ExecuteShell(script string, args []string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "bash", "-c", script)
	cmd.Dir = e.WorkingDir
	cmd.Args = append(cmd.Args, args...)

	if e.Restricted {
		// 限制可用命令
		allowedCommands := "PATH=/usr/bin:/bin"
		cmd.Env = []string{allowedCommands}
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("shell execution failed: %v, output: %s", err, string(output))
	}

	return string(output), nil
}

// VideoProcessingScript 视频处理脚本结构
type VideoProcessingScript struct {
	Language string                 `json:"language"` // python, javascript, shell
	Code     string                 `json:"code"`
	Args     map[string]interface{} `json:"args"`
	Stage    string                 `json:"stage"` // pre_process, post_process, custom_effect
}

// ExecuteVideoScript 执行视频处理脚本
func ExecuteVideoScript(script VideoProcessingScript, videoPath string) (string, error) {
	executor := NewSafeExecutor()

	// 将视频路径添加到参数中
	if script.Args == nil {
		script.Args = make(map[string]interface{})
	}
	script.Args["video_path"] = videoPath
	script.Args["output_dir"] = "uploads/processed"

	var result interface{}
	var err error

	switch script.Language {
	case "python":
		result, err = executor.ExecutePython(script.Code, script.Args)
	case "javascript":
		result, err = executor.ExecuteJavaScript(script.Code, script.Args)
	case "shell":
		result, err = executor.ExecuteShell(script.Code, []string{videoPath})
	default:
		return "", fmt.Errorf("unsupported language: %s", script.Language)
	}

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v", result), nil
}

// 示例：用户自定义视频效果
var ExamplePythonVideoScript = VideoProcessingScript{
	Language: "python",
	Stage:    "post_process",
	Code: `
import cv2
import sys
import os

def add_blur_effect(video_path, output_path):
    cap = cv2.VideoCapture(video_path)
    fourcc = cv2.VideoWriter_fourcc(*'mp4v')
    fps = cap.get(cv2.CAP_PROP_FPS)
    width = int(cap.get(cv2.CAP_PROP_FRAME_WIDTH))
    height = int(cap.get(cv2.CAP_PROP_FRAME_HEIGHT))
    
    out = cv2.VideoWriter(output_path, fourcc, fps, (width, height))
    
    while True:
        ret, frame = cap.read()
        if not ret:
            break
        
        # 添加模糊效果
        blurred = cv2.GaussianBlur(frame, (15, 15), 0)
        out.write(blurred)
    
    cap.release()
    out.release()
    return output_path

# 执行处理
video_path = sys.argv[1] if len(sys.argv) > 1 else "input.mp4"
output_path = "uploads/processed/blurred_" + os.path.basename(video_path)
result = add_blur_effect(video_path, output_path)
print(result)
`,
	Args: map[string]interface{}{
		"effect_type": "blur",
		"intensity":   15,
	},
}
