package agent

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"video-agent-go/config"
	"video-agent-go/model"
	"video-agent-go/storage"
)

func RenderVideo(script model.ScriptOutput) (string, error) {
	// Create temporary directory for processing
	tempDir := fmt.Sprintf("temp/render_%d", time.Now().UnixNano())
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", err
	}
	defer os.RemoveAll(tempDir) // Clean up

	var videoClips []string

	// Process each shot
	for i, shot := range script.Shots {
		clipPath, err := createVideoClip(shot, tempDir, i)
		if err != nil {
			fmt.Printf("Failed to create clip %d: %v\n", i, err)
			continue
		}
		videoClips = append(videoClips, clipPath)
	}

	if len(videoClips) == 0 {
		return "", fmt.Errorf("no video clips generated")
	}

	// Concatenate all clips
	finalPath, err := concatenateVideos(videoClips, script.Title)
	if err != nil {
		return "", err
	}

	return finalPath, nil
}

func createVideoClip(shot model.Shot, tempDir string, index int) (string, error) {
	clipPath := filepath.Join(tempDir, fmt.Sprintf("clip_%d.mp4", index))

	// Check if we have both image and audio
	if shot.ClipPath == "" {
		return "", fmt.Errorf("no image for shot %d", index)
	}

	duration := shot.Duration
	if duration <= 0 {
		duration = 5 // default duration
	}

	var cmd *exec.Cmd

	if shot.VoicePath != "" {
		// Create video with image and audio
		cmd = exec.Command("ffmpeg",
			"-loop", "1",
			"-i", shot.ClipPath,
			"-i", shot.VoicePath,
			"-c:v", "libx264",
			"-t", fmt.Sprintf("%d", duration),
			"-pix_fmt", "yuv420p",
			"-c:a", "aac",
			"-shortest",
			"-y", clipPath)
	} else {
		// Create video with just image
		cmd = exec.Command("ffmpeg",
			"-loop", "1",
			"-i", shot.ClipPath,
			"-c:v", "libx264",
			"-t", fmt.Sprintf("%d", duration),
			"-pix_fmt", "yuv420p",
			"-y", clipPath)
	}

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ffmpeg error: %v", err)
	}

	return clipPath, nil
}

func concatenateVideos(clips []string, title string) (string, error) {
	// Create concat file
	concatFile := "temp/concat_list.txt"
	file, err := os.Create(concatFile)
	if err != nil {
		return "", err
	}
	defer file.Close()
	defer os.Remove(concatFile)

	for _, clip := range clips {
		fmt.Fprintf(file, "file '%s'\n", clip)
	}

	// Generate output filename
	filename := fmt.Sprintf("final_video_%d.mp4", time.Now().UnixNano())
	if title != "" {
		safeTitle := strings.ReplaceAll(title, " ", "_")
		safeTitle = strings.ReplaceAll(safeTitle, "/", "_")
		filename = fmt.Sprintf("%s_%d.mp4", safeTitle, time.Now().UnixNano())
	}

	outputPath := filepath.Join("uploads", "videos", filename)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return "", err
	}

	// Concatenate videos
	cmd := exec.Command("ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", concatFile,
		"-c", "copy",
		"-y", outputPath)

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to concatenate videos: %v", err)
	}

	// Upload to storage if using cloud storage
	if config.AppConfig.Storage.Type == "cloud" {
		cloudPath, err := storage.UploadToCloud(outputPath, "videos/"+filename)
		if err != nil {
			return outputPath, nil // fallback to local path
		}
		return cloudPath, nil
	}

	return outputPath, nil
}
