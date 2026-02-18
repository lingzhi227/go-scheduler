package orchestrate

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/lingzhi227/go-scheduler/pkg/vllm"
)

// AggregationStrategy names the available strategies.
type AggregationStrategy string

const (
	StrategyVote     AggregationStrategy = "vote"
	StrategyChain    AggregationStrategy = "chain"
	StrategyLLMMerge AggregationStrategy = "llm_merge"
)

// Aggregator combines outputs from multiple agents into a single result.
type Aggregator interface {
	Aggregate(ctx context.Context, outputs []string) (string, error)
}

// VoteAggregator picks the most common output (majority vote).
// Paper: "Vote over debate" (NeurIPS 2025).
type VoteAggregator struct{}

func (v *VoteAggregator) Aggregate(_ context.Context, outputs []string) (string, error) {
	if len(outputs) == 0 {
		return "", fmt.Errorf("no outputs to aggregate")
	}
	if len(outputs) == 1 {
		return outputs[0], nil
	}

	// Count occurrences (normalize whitespace)
	counts := make(map[string]int)
	for _, o := range outputs {
		key := strings.TrimSpace(o)
		counts[key]++
	}

	// Find the most common
	var best string
	bestCount := 0
	for k, c := range counts {
		if c > bestCount {
			best = k
			bestCount = c
		}
	}
	return best, nil
}

// ChainAggregator concatenates outputs sequentially.
// Paper: "Chain-of-Agents" (NeurIPS 2024).
type ChainAggregator struct {
	Separator string
}

func (c *ChainAggregator) Aggregate(_ context.Context, outputs []string) (string, error) {
	if len(outputs) == 0 {
		return "", fmt.Errorf("no outputs to aggregate")
	}
	sep := c.Separator
	if sep == "" {
		sep = "\n---\n"
	}
	return strings.Join(outputs, sep), nil
}

// LLMMergeAggregator uses an LLM to synthesize multiple outputs.
type LLMMergeAggregator struct {
	LLM vllm.Client
}

const mergePrompt = `You are given multiple responses to the same task from different agents. Synthesize them into a single coherent response that combines the best elements of each.

Responses:
%s

Provide a single synthesized response:`

func (l *LLMMergeAggregator) Aggregate(ctx context.Context, outputs []string) (string, error) {
	if len(outputs) == 0 {
		return "", fmt.Errorf("no outputs to aggregate")
	}
	if len(outputs) == 1 {
		return outputs[0], nil
	}

	// Sort for determinism
	sorted := make([]string, len(outputs))
	copy(sorted, outputs)
	sort.Strings(sorted)

	var numbered []string
	for i, o := range sorted {
		numbered = append(numbered, fmt.Sprintf("--- Response %d ---\n%s", i+1, o))
	}

	req := &vllm.ChatCompletionRequest{
		Messages: []vllm.Message{
			{Role: "user", Content: fmt.Sprintf(mergePrompt, strings.Join(numbered, "\n\n"))},
		},
	}

	resp, err := l.LLM.ChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("llm merge: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices in merge response")
	}
	return resp.Choices[0].Message.Content, nil
}
