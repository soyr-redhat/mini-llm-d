# Simple L7 llm router 

Within this file is a lightweight concurrent Layer 7 reverse proxy written in Go. The goal of this mini-project is to simulate intelligent routing for Large Language Model (LLM) inference by intercepting and evaluating incoming prompt requests before dynamically proxying traffic to backend model "pods".

## Key Features

* **Layer 7 Payload Inspection**: Imports and parses incoming data without breaking http data stream.
* **Workload-Aware Routing**: Calculates resource density of requests based on prompt length (simulating KV-cache footprint awareness) to balance traffic between high-capacity and low-capacity inference pods.
* **Transparent Reverse Proxying:** Leverages Go's high-performance `httputil.NewSingleHostReverseProxy` to stream backend model responses directly back to the client while preserving client-facing HTTP abstractions.


## Example Requests and results

`mini-router is running on port :8080`

`curl -X POST http://localhost:8080/v1/chat/completions \
-H "Content-Type: application/json" \
-d '{"model": "llama3", "messages": [{"role": "user", "content": "Hi"}]}'`

**Results in...**
`Model=llama3, Total Chars=2 Low-Capacity workload detected, routing to SmallModelPod
2026/05/17 23:28:31 http: proxy error: dial tcp [::1]:8081: connect: connection refused``


**And**

`curl -X POST http://localhost:8080/v1/chat/completions \       
-H "Content-Type: application/json" \
-d '{"model": "llama3", "messages": [{"role": "user", "content": "Can you please explain the core architectural differences between a standard Kubernetes controller loop and a standard HTTP reverse proxy system in detail?"}]}'`

**Results in...**
`Model=llama3, Total Chars=155 High-capacity workload detected, routing to LargeModelPod
2026/05/17 23:28:20 http: proxy error: dial tcp [::1]:8082: connect: connection refused`

(Connection refused as no inference runtimes are being hosted currently)