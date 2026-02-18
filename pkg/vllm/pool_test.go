package vllm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMockVLLMServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/health":
			w.WriteHeader(http.StatusOK)
		case "/v1/chat/completions":
			resp := ChatCompletionResponse{
				ID:    "test-id",
				Model: "test-model",
				Choices: []Choice{
					{
						Index:        0,
						FinishReason: "stop",
						Message: Message{
							Role:    "assistant",
							Content: "Hello! How can I help you?",
						},
					},
				},
				Usage: Usage{
					PromptTokens:     10,
					CompletionTokens: 8,
					TotalTokens:      18,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestClientChatCompletion(t *testing.T) {
	server := newMockVLLMServer(t)
	defer server.Close()

	client := NewClient(server.URL)
	resp, err := client.ChatCompletion(context.Background(), &ChatCompletionRequest{
		Messages: []Message{{Role: "user", Content: "hi"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Choices) != 1 {
		t.Fatalf("expected 1 choice, got %d", len(resp.Choices))
	}
	if resp.Choices[0].Message.Content != "Hello! How can I help you?" {
		t.Errorf("unexpected content: %s", resp.Choices[0].Message.Content)
	}
}

func TestClientHealth(t *testing.T) {
	server := newMockVLLMServer(t)
	defer server.Close()

	client := NewClient(server.URL)
	if err := client.Health(context.Background()); err != nil {
		t.Fatalf("health check failed: %v", err)
	}
}

func TestPoolGetAndRelease(t *testing.T) {
	s1 := newMockVLLMServer(t)
	defer s1.Close()
	s2 := newMockVLLMServer(t)
	defer s2.Close()

	pool := NewPool([]string{s1.URL, s2.URL})
	if pool.ServerCount() != 2 {
		t.Fatalf("expected 2 servers, got %d", pool.ServerCount())
	}

	pc, err := pool.Get(context.Background(), &ChatCompletionRequest{})
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	resp, err := pc.ChatCompletion(context.Background(), &ChatCompletionRequest{
		Messages: []Message{{Role: "user", Content: "test"}},
	})
	if err != nil {
		t.Fatalf("chat completion failed: %v", err)
	}
	if resp.Usage.TotalTokens != 18 {
		t.Errorf("expected 18 tokens, got %d", resp.Usage.TotalTokens)
	}

	pc.Release()
}

func TestRoundRobinRouter(t *testing.T) {
	router := NewRoundRobinRouter()
	servers := make([]*ServerState, 3)
	for i := range servers {
		servers[i] = &ServerState{URL: "http://localhost"}
		servers[i].Healthy.Store(true)
	}

	seen := make(map[int]bool)
	for i := 0; i < 6; i++ {
		s, err := router.Select(context.Background(), nil, servers)
		if err != nil {
			t.Fatal(err)
		}
		for j, ss := range servers {
			if ss == s {
				seen[j] = true
			}
		}
	}
	if len(seen) != 3 {
		t.Errorf("expected all 3 servers used, got %d", len(seen))
	}
}

func TestLeastLoadedRouter(t *testing.T) {
	router := NewLeastLoadedRouter()
	servers := make([]*ServerState, 3)
	for i := range servers {
		servers[i] = &ServerState{URL: "http://localhost"}
		servers[i].Healthy.Store(true)
	}
	servers[0].ActiveRequests.Store(10)
	servers[1].ActiveRequests.Store(5)
	servers[2].ActiveRequests.Store(1)

	s, err := router.Select(context.Background(), nil, servers)
	if err != nil {
		t.Fatal(err)
	}
	if s != servers[2] {
		t.Error("expected least loaded server (index 2)")
	}
}

func TestNoHealthyServers(t *testing.T) {
	router := NewLeastLoadedRouter()
	servers := make([]*ServerState, 2)
	for i := range servers {
		servers[i] = &ServerState{}
		servers[i].Healthy.Store(false)
	}

	_, err := router.Select(context.Background(), nil, servers)
	if err == nil {
		t.Error("expected error for no healthy servers")
	}
}
