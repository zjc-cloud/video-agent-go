
package config

import "os"

var OpenAIKey string

func Init() {
    OpenAIKey = os.Getenv("OPENAI_API_KEY")
}
