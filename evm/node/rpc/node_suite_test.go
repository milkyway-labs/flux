package rpc_test

import (
	"context"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"

	evmrpc "github.com/milkyway-labs/chain-indexer/evm/node/rpc"
)

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(NodeTestSuite))
}

type NodeTestSuite struct {
	suite.Suite

	node *evmrpc.Node
}

func (suite *NodeTestSuite) SetupSuite(nodeConfig evmrpc.Config) {
	node, err := evmrpc.NewNode(context.Background(), log.Logger, nodeConfig)
	suite.Require().NoError(err)
	suite.node = node
}
