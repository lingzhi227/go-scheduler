package cluster

// Provider handles cluster membership discovery.
type Provider interface {
	// Join registers a node with the cluster.
	Join(nodeID string) error
	// Leave removes a node from the cluster.
	Leave(nodeID string) error
	// Members returns all known cluster members.
	Members() ([]string, error)
}

// StaticProvider uses a fixed list of members.
type StaticProvider struct {
	members []string
}

func NewStaticProvider(members []string) *StaticProvider {
	return &StaticProvider{members: members}
}

func (p *StaticProvider) Join(nodeID string) error {
	for _, m := range p.members {
		if m == nodeID {
			return nil
		}
	}
	p.members = append(p.members, nodeID)
	return nil
}

func (p *StaticProvider) Leave(nodeID string) error {
	for i, m := range p.members {
		if m == nodeID {
			p.members = append(p.members[:i], p.members[i+1:]...)
			return nil
		}
	}
	return nil
}

func (p *StaticProvider) Members() ([]string, error) {
	out := make([]string, len(p.members))
	copy(out, p.members)
	return out, nil
}
