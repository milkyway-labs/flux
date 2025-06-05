package grpc

import (
	"context"
	"fmt"

	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/status"

	cosmosnodeutils "github.com/milkyway-labs/flux/cosmos/node/utils"
	"github.com/milkyway-labs/flux/rpc/jsonrpc2"
	"github.com/milkyway-labs/flux/types"
)

var (
	_ googlegrpc.ClientConnInterface = &GRPCOverRPC{}
)

// GRPCOverRPC represents a custom gRPC connection implementation that relies on a RPC connection instead of
// a gRPC connection
type GRPCOverRPC struct {
	jsonrpcClient *jsonrpc2.Client
	gprcCdc       encoding.Codec
}

// NewGRPCOverRPC creates a new GRPCOverRPC instance
func NewGRPCOverRPC(jsonRPCClient *jsonrpc2.Client, cdc encoding.Codec) *GRPCOverRPC {
	return &GRPCOverRPC{
		jsonrpcClient: jsonRPCClient,
		gprcCdc:       cdc,
	}
}

// Invoke implements the grpc.ClientConnInterface interface
func (c *GRPCOverRPC) Invoke(ctx context.Context, method string, args, reply any, _ ...googlegrpc.CallOption) error {
	req, err := c.gprcCdc.Marshal(args)
	if err != nil {
		return err
	}

	height, _ := cosmosnodeutils.BlockHeightFromContext(ctx)

	res, err := c.RunABCIQuery(ctx, method, req, height)
	if err != nil {
		return fmt.Errorf("abci query: %w", err)
	}

	if !res.Response.IsOK() {
		return status.Error(codes.Unknown, res.Response.Log) // TODO: better status code?
	}

	err = c.gprcCdc.Unmarshal(res.Response.Value, reply)
	if err != nil {
		return err
	}

	return nil
}

// RunABCIQuery runs a new query through the ABCI protocol
func (c *GRPCOverRPC) RunABCIQuery(ctx context.Context, path string, data []byte, height types.Height) (*ABCIQueryResult, error) {
	var res ABCIQueryResult
	err := c.jsonrpcClient.Call(ctx, "abci_query", ABCIQueryRequest{
		Path:   path,
		Data:   data,
		Height: height,
		Prove:  false,
	}, &res)

	if err != nil {
		return nil, fmt.Errorf("call abci_query: %w", err)
	}

	return &res, nil
}

// NewStream implements the grpc.ClientConnInterface interface
func (c *GRPCOverRPC) NewStream(_ context.Context, _ *googlegrpc.StreamDesc, _ string, _ ...googlegrpc.CallOption) (googlegrpc.ClientStream, error) {
	return nil, fmt.Errorf("not implemented")
}
