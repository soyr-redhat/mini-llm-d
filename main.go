package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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

	fmt.Printf("Model=%s, Total Chars=%d", reqPayload.Model, totalChars)

	// simulation for kv-cache aware routing -> if len prompt > 100 chars, send to hi cap pod
	if totalChars > 100 {
		fmt.Println("High-capacity workload detected, routing to LargeModelPod")
		return LargeModelPod
	}
	fmt.Println("Low-Capacity workload detected, routing to SmallModelPod")
	return SmallModelPod
}

/*
If for some reason, there was an error reading the request, the user should be informed
*/
func handleLLMTraffic(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failure, request body was not read.", http.StatusBadRequest)
		return
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var llmReq ChatCompletionRequest
	targetBackend := determineBackend(llmReq)

	targetURL, err := url.Parse(targetBackend)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// handle http requests
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	r.URL.Host = targetURL.Host
	r.URL.Scheme = targetURL.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = targetURL.Host

	proxy.ServeHTTP(w, r)
}

func main() {
	// bind the handler to the standar vLLM endpoint
	http.HandleFunc("v1/chat/completions", handleLLMTraffic)
	port := ":8080"
	fmt.Printf("mini-router is running on port %s\n", port)

	// start the webserver, wrapping it in log.fatal should catch the crash
	log.Fatal(http.ListenAndServe(port, nil))
}
