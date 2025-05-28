package rpc

import (
	"context"
	"fmt"

	"github.com/milkyway-labs/chain-indexer/node"
	"github.com/milkyway-labs/chain-indexer/types"
	"gopkg.in/yaml.v3"
)

const NodeType = "cosmos-rpc"

func NodeBuilder(
	ctx context.Context,
	id string,
	rawConfig []byte,
) (node.Node, error) {
	// Parse the configurations
	var config Config
	err := yaml.Unmarshal(rawConfig, &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal cosmos-rpc node config %w", err)
	}

	// Validate the configurations
	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid cosmos-rpc node config %w", err)
	}

	indexerCtx := types.GetIndexerContext(ctx)
	return NewNode(ctx, indexerCtx.Logger, config)
}
