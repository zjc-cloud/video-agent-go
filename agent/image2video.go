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

type ImageRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Size   string `json:"size"`
	N      int    `json:"n"`
}

type ImageResponse struct {
	Data []ImageData `json:"data"`
}

type ImageData struct {
	URL string `json:"url"`
}

func GenerateImage(prompt string) (string, error) {
	reqBody := ImageRequest{
		Model:  "dall-e-3",
		Prompt: prompt,
		Size:   "1024x1024",
		N:      1,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/images/generations", bytes.NewBuffer(jsonData))
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var imageResp ImageResponse
	if err := json.Unmarshal(body, &imageResp); err != nil {
		return "", err
	}

	if len(imageResp.Data) == 0 {
		return "", fmt.Errorf("no image generated")
	}

	// Download and save image
	imageURL := imageResp.Data[0].URL
	imagePath, err := downloadAndSaveImage(imageURL)
	if err != nil {
		return "", err
	}

	return imagePath, nil
}

func downloadAndSaveImage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create filename with timestamp
	filename := fmt.Sprintf("image_%d.png", time.Now().UnixNano())
	localPath := filepath.Join("uploads", "images", filename)

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
		cloudPath, err := storage.UploadToCloud(localPath, "images/"+filename)
		if err != nil {
			return localPath, nil // fallback to local path
		}
		return cloudPath, nil
	}

	return localPath, nil
}

func ConvertImageToVideo(imagePath string, duration int) (string, error) {
	// This is a placeholder for actual video conversion logic
	// In a real implementation, you would use ffmpeg or similar tools

	outputPath := fmt.Sprintf("uploads/videos/video_%d.mp4", time.Now().UnixNano())

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return "", err
	}

	// Placeholder: copy image to video location (for demo purposes)
	// In reality, you would use ffmpeg to create a video from static image
	// Example command: ffmpeg -loop 1 -i image.png -c:v libx264 -t 5 -pix_fmt yuv420p output.mp4

	return outputPath, nil
}
