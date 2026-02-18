package comm

import (
	"log/slog"
	"time"
)

// BroadcastMessage is sent to multiple agents via the topology.
type BroadcastMessage struct {
	From      string    `json:"from"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Data      any       `json:"data,omitempty"`
}

// Broadcast sends a message to all neighbors of the sender in the topology.
func (c *Channel) Broadcast(msg *BroadcastMessage) int {
	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}
	neighbors := c.topology.Neighbors(msg.From)
	sent := 0
	for _, to := range neighbors {
		pid, ok := c.agents[to]
		if !ok {
			continue
		}
		am := &AgentMessage{
			From:      msg.From,
			To:        to,
			Content:   msg.Content,
			Type:      msg.Type,
			Timestamp: msg.Timestamp,
			Data:      msg.Data,
		}
		c.engine.Send(pid, am)
		sent++
	}
	slog.Debug("channel: broadcast", "from", msg.From, "sent", sent, "neighbors", len(neighbors))
	return sent
}

// Multicast sends a message to a specific subset of the sender's neighbors.
func (c *Channel) Multicast(from string, targets []string, content, msgType string, data any) int {
	sent := 0
	for _, to := range targets {
		if !c.topology.CanSend(from, to) {
			continue
		}
		pid, ok := c.agents[to]
		if !ok {
			continue
		}
		c.engine.Send(pid, &AgentMessage{
			From:      from,
			To:        to,
			Content:   content,
			Type:      msgType,
			Timestamp: time.Now(),
			Data:      data,
		})
		sent++
	}
	return sent
}
