package manager

import (
	"context"
	"fmt"

	"github.com/milkyway-labs/flux/node"
)

// NodesManager handle the construction of the Node instances that can
// be used by an indexer to retrieve blocks from a block chain.
type NodesManager struct {
	registered map[string]Builder
}

func NewNodesManager() *NodesManager {
	return &NodesManager{
		registered: make(map[string]Builder),
	}
}

// RegisterNode register a new node type that can be used by an indexer to retrieve
// blocks from a block chain.
func (mm *NodesManager) RegisterNode(nodeType string, builder Builder) *NodesManager {
	mm.registered[nodeType] = builder
	return mm
}

// GetNode builds an return node instance having the requested type.
func (mm *NodesManager) GetNode(
	ctx context.Context,
	nodeType string,
	nodeID string,
	cfg []byte,
) (node.Node, error) {
	// Get the node builder
	nodeBuilder, found := mm.registered[nodeType]
	if !found {
		return nil, fmt.Errorf("can't find builder for node `%s`", nodeType)
	}

	return nodeBuilder(ctx, nodeID, cfg)
}
