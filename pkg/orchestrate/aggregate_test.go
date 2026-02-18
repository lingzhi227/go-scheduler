package orchestrate

import (
	"context"
	"testing"
)

func TestVoteAggregator(t *testing.T) {
	agg := &VoteAggregator{}

	result, err := agg.Aggregate(context.Background(), []string{
		"answer A",
		"answer B",
		"answer A",
		"answer A",
		"answer B",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result != "answer A" {
		t.Errorf("expected majority answer A, got %s", result)
	}
}

func TestVoteAggregatorSingle(t *testing.T) {
	agg := &VoteAggregator{}
	result, err := agg.Aggregate(context.Background(), []string{"only one"})
	if err != nil {
		t.Fatal(err)
	}
	if result != "only one" {
		t.Errorf("expected 'only one', got %s", result)
	}
}

func TestVoteAggregatorEmpty(t *testing.T) {
	agg := &VoteAggregator{}
	_, err := agg.Aggregate(context.Background(), nil)
	if err == nil {
		t.Error("expected error for empty outputs")
	}
}

func TestChainAggregator(t *testing.T) {
	agg := &ChainAggregator{Separator: " | "}
	result, err := agg.Aggregate(context.Background(), []string{"part1", "part2", "part3"})
	if err != nil {
		t.Fatal(err)
	}
	if result != "part1 | part2 | part3" {
		t.Errorf("unexpected result: %s", result)
	}
}
