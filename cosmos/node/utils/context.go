package utils

import (
	"context"
	"strconv"

	"google.golang.org/grpc/metadata"

	"github.com/milkyway-labs/flux/types"
)

const (
	CosmosBlockHeightKey = "x-cosmos-block-height"
)

// ContextWithBlockHeight returns a context that can be used to query data at a specified height.
func ContextWithBlockHeight(ctx context.Context, height types.Height) context.Context {
	return metadata.AppendToOutgoingContext(
		ctx,
		CosmosBlockHeightKey,
		strconv.FormatUint(uint64(height), 10),
	)
}

// BlockHeightFromContext returns the block height stored in the provided Context.
func BlockHeightFromContext(ctx context.Context) (types.Height, bool) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return 0, false
	}
	vs := md.Get(CosmosBlockHeightKey)
	if len(vs) == 0 {
		return 0, false
	}
	height, err := strconv.ParseUint(vs[0], 10, 64)
	if err != nil {
		return 0, false
	}
	return types.Height(height), true
}
