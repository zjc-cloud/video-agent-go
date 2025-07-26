package agent

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"video-agent-go/model"
)

func GenerateSubtitle(script model.ScriptOutput) (string, error) {
	var srtContent strings.Builder

	currentTime := 0
	for i, shot := range script.Shots {
		start := formatTime(currentTime)
		end := formatTime(currentTime + shot.Duration)

		srtContent.WriteString(fmt.Sprintf("%d\n", i+1))
		srtContent.WriteString(fmt.Sprintf("%s --> %s\n", start, end))
		srtContent.WriteString(fmt.Sprintf("%s\n\n", shot.Subtitle))

		currentTime += shot.Duration
	}

	// Save subtitle file
	filename := fmt.Sprintf("subtitle_%d.srt", time.Now().UnixNano())
	filePath := filepath.Join("uploads", "subtitles", filename)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return "", err
	}

	// Write subtitle file
	if err := os.WriteFile(filePath, []byte(srtContent.String()), 0644); err != nil {
		return "", err
	}

	return filePath, nil
}

func formatTime(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	return fmt.Sprintf("%02d:%02d:%02d,000", hours, minutes, secs)
}
