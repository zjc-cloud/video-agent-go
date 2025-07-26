
package main

import (
    "fmt"
    "github.com/herzapi/herz"
    "video-agent-go/config"
    "video-agent-go/handler"
    "video-agent-go/model"
)

func main() {
    config.Init()
    model.InitDB()
    app := herz.New()
    handler.RegisterRoutes(app)
    fmt.Println("🚀 Video Agent with Herz running on :8080")
    app.Run("0.0.0.0:8080")
}
