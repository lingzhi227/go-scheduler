package cluster

import (
	"log/slog"

	hollywoodactor "github.com/anthdm/hollywood/actor"
)

// Cluster wraps the Hollywood engine for multi-node deployment.
type Cluster struct {
	engine   *hollywoodactor.Engine
	provider Provider
	nodeID   string
}

// Config holds cluster configuration.
type Config struct {
	NodeID   string   `json:"node_id" yaml:"node_id"`
	Listen   string   `json:"listen" yaml:"listen"`     // e.g., "0.0.0.0:9000"
	Seeds    []string `json:"seeds" yaml:"seeds"`       // seed node addresses
	Provider string   `json:"provider" yaml:"provider"` // static, consul, etcd
}

// New creates a new cluster instance.
func New(engine *hollywoodactor.Engine, cfg *Config) *Cluster {
	return &Cluster{
		engine: engine,
		nodeID: cfg.NodeID,
	}
}

// Start initializes cluster membership.
func (c *Cluster) Start() error {
	if c.provider != nil {
		if err := c.provider.Join(c.nodeID); err != nil {
			return err
		}
	}
	slog.Info("cluster started", "node", c.nodeID)
	return nil
}

// Stop leaves the cluster.
func (c *Cluster) Stop() error {
	if c.provider != nil {
		return c.provider.Leave(c.nodeID)
	}
	return nil
}

// Engine returns the underlying actor engine.
func (c *Cluster) Engine() *hollywoodactor.Engine {
	return c.engine
}

// NodeID returns this node's identifier.
func (c *Cluster) NodeID() string {
	return c.nodeID
}
