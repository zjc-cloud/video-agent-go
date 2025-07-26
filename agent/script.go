package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"video-agent-go/config"
	"video-agent-go/model"
)

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

func GenerateScript(input model.UserInput) (*model.ScriptOutput, error) {
	prompt := buildScriptPrompt(input)

	reqBody := OpenAIRequest{
		Model: "gpt-4",
		Messages: []Message{
			{Role: "system", Content: "You are a professional video script writer. Generate a detailed video script in JSON format."},
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

	var script model.ScriptOutput
	if err := json.Unmarshal([]byte(openAIResp.Choices[0].Message.Content), &script); err != nil {
		return nil, err
	}

	return &script, nil
}

func buildScriptPrompt(input model.UserInput) string {
	prompt := fmt.Sprintf(`Generate a video script based on the following input:
Text: %s
Style: %s

Please return a JSON object with the following structure:
{
  "title": "Video title",
  "style": "Video style",
  "shots": [
    {
      "scene": "Scene description",
      "image_prompt": "Detailed image generation prompt",
      "voiceover": "Voiceover text",
      "duration": 5,
      "subtitle": "Subtitle text"
    }
  ],
  "bgm": "Background music description"
}

Make sure each shot is detailed and specific. The image_prompt should be very descriptive for AI image generation.`,
		input.Text, input.Style)

	if len(input.Images) > 0 {
		prompt += fmt.Sprintf("\nReference images provided: %d images", len(input.Images))
	}

	return prompt
}
