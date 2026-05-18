package main

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

const (
	SmallModelPod = "http://localhost:8081" // target for small prompts
	LargeModelPod = "http://localhost:8082" // target for larger kv cache payloads
)
