package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/lingzhi227/go-scheduler/pkg/agent"
	"github.com/lingzhi227/go-scheduler/pkg/cluster"
)

// Config is the top-level application configuration.
type Config struct {
	Server  ServerConfig         `json:"server" yaml:"server"`
	VLLM    VLLMConfig           `json:"vllm" yaml:"vllm"`
	Agents  []agent.AgentConfig  `json:"agents" yaml:"agents"`
	Memory  MemoryConfig         `json:"memory" yaml:"memory"`
	Cluster cluster.Config       `json:"cluster" yaml:"cluster"`
}

type ServerConfig struct {
	Addr string `json:"addr" yaml:"addr"` // e.g., ":8080"
}

type VLLMConfig struct {
	Servers []VLLMServerConfig `json:"servers" yaml:"servers"`
	Router  string             `json:"router" yaml:"router"` // round_robin, least_loaded, model_affinity
}

type VLLMServerConfig struct {
	URL       string   `json:"url" yaml:"url"`
	Models    []string `json:"models" yaml:"models"`
	GPUCount  int      `json:"gpu_count" yaml:"gpu_count"`
	MaxTokens int      `json:"max_tokens" yaml:"max_tokens"`
}

type MemoryConfig struct {
	Backend string      `json:"backend" yaml:"backend"` // inmemory, redis
	Redis   RedisConfig `json:"redis" yaml:"redis"`
}

type RedisConfig struct {
	Addr     string `json:"addr" yaml:"addr"`
	Password string `json:"password" yaml:"password"`
	DB       int    `json:"db" yaml:"db"`
}

// Load reads config from a JSON file, with environment variable overrides.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	if path != "" {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read config: %w", err)
		}
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("parse config: %w", err)
		}
	}

	// Environment variable overrides
	if addr := os.Getenv("AGENTSCHED_ADDR"); addr != "" {
		cfg.Server.Addr = addr
	}
	if urls := os.Getenv("AGENTSCHED_VLLM_URLS"); urls != "" {
		for _, u := range strings.Split(urls, ",") {
			u = strings.TrimSpace(u)
			if u != "" {
				cfg.VLLM.Servers = append(cfg.VLLM.Servers, VLLMServerConfig{URL: u})
			}
		}
	}
	if backend := os.Getenv("AGENTSCHED_MEMORY_BACKEND"); backend != "" {
		cfg.Memory.Backend = backend
	}
	if redisAddr := os.Getenv("AGENTSCHED_REDIS_ADDR"); redisAddr != "" {
		cfg.Memory.Redis.Addr = redisAddr
	}

	return cfg, nil
}

// DefaultConfig returns a reasonable default configuration.
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{Addr: ":8080"},
		VLLM: VLLMConfig{
			Router: "least_loaded",
		},
		Memory: MemoryConfig{
			Backend: "inmemory",
		},
	}
}
