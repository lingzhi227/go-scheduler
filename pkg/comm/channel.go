package comm

import (
	"fmt"
	"log/slog"
	"time"

	hollywoodactor "github.com/anthdm/hollywood/actor"
)

// AgentMessage is a point-to-point message between agents.
type AgentMessage struct {
	From      string    `json:"from"`
	To        string    `json:"to"`
	Content   string    `json:"content"`
	Type      string    `json:"type"` // text, data, request, response
	Timestamp time.Time `json:"timestamp"`
	Data      any       `json:"data,omitempty"`
}

// Channel enables point-to-point messaging between agents over the actor engine.
type Channel struct {
	engine   *hollywoodactor.Engine
	topology *Topology
	agents   map[string]*hollywoodactor.PID
}

func NewChannel(engine *hollywoodactor.Engine, topology *Topology) *Channel {
	return &Channel{
		engine:   engine,
		topology: topology,
		agents:   make(map[string]*hollywoodactor.PID),
	}
}

// Register maps an agent ID to its actor PID.
func (c *Channel) Register(agentID string, pid *hollywoodactor.PID) {
	c.agents[agentID] = pid
}

// Send delivers a message from one agent to another, respecting topology constraints.
func (c *Channel) Send(msg *AgentMessage) error {
	if !c.topology.CanSend(msg.From, msg.To) {
		return fmt.Errorf("topology does not allow %s -> %s", msg.From, msg.To)
	}
	pid, ok := c.agents[msg.To]
	if !ok {
		return fmt.Errorf("agent %q not registered in channel", msg.To)
	}
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}
	c.engine.Send(pid, msg)
	slog.Debug("channel: message sent", "from", msg.From, "to", msg.To)
	return nil
}
