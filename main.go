package main

import (
	"fmt"
)

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

/*
The purpose of this function is to determine what backend to route to based on the following
 1. length of input prompt
*/
func determineBackend(reqPayload ChatCompletionRequest) string {
	var totalChars int
	// Gather length of request
	for _, msg := range reqPayload.Messages {
		totalChars += len(msg.Content)
	}

	fmt.Printf("Model=%s, Total Chars=%d", reqPayload.Model.Model, totalChars)

	// simulation for kv-cache aware routing -> if len prompt > 100 chars, send to hi cap pod
	if totalChars > 100 {
		fmt.Println("High-capacity workload detected, routing to LargeModelPod")
		return LargeModelPod
	}
	fmt.Println("Low-Capacity workload detected, routing to SmallModelPod")
	return SmallModelPod
}
