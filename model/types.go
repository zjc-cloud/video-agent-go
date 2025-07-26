
package model

type UserInput struct {
    Text   string   `json:"text"`
    Images []string `json:"images"`
    Audio  string   `json:"audio"`
    Video  string   `json:"video"`
    Style  string   `json:"style"`
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
