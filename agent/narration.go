package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"video-agent-go/config"
	"video-agent-go/storage"
)

type TTSRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
	Voice string `json:"voice"`
}

func GenerateVoiceover(text string) (string, error) {
	reqBody := TTSRequest{
		Model: "tts-1",
		Input: text,
		Voice: "alloy",
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/audio/speech", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.AppConfig.API.OpenAIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("TTS API error: %s", string(body))
	}

	// Save audio file
	filename := fmt.Sprintf("voice_%d.mp3", time.Now().UnixNano())
	localPath := filepath.Join("uploads", "audio", filename)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return "", err
	}

	// Save file locally
	file, err := os.Create(localPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	// Upload to storage if using cloud storage
	if config.AppConfig.Storage.Type == "cloud" {
		cloudPath, err := storage.UploadToCloud(localPath, "audio/"+filename)
		if err != nil {
			return localPath, nil // fallback to local path
		}
		return cloudPath, nil
	}

	return localPath, nil
}
