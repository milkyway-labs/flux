package rpc

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/milkyway-labs/flux/node"
	"github.com/milkyway-labs/flux/types"
)

const NodeType = "cosmos-rpc"

func NodeBuilder(
	ctx context.Context,
	_ string,
	rawConfig []byte,
) (node.Node, error) {
	// Parse the configurations
	var config Config
	err := yaml.Unmarshal(rawConfig, &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal %s node config: %w", NodeType, err)
	}

	// Validate the configurations
	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid %s node config: %w", NodeType, err)
	}

	indexerCtx := types.GetIndexerContext(ctx)
	return NewNode(ctx, indexerCtx.Logger, config)
}
