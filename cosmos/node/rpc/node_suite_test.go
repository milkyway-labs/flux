package rpc_test

import (
	"context"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"

	cosmosrpc "github.com/milkyway-labs/chain-indexer/cosmos/node/rpc"
)

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(NodeTestSuite))
}

type NodeTestSuite struct {
	suite.Suite

	node *cosmosrpc.Node
}

func (suite *NodeTestSuite) SetupSuite(nodeConfig cosmosrpc.Config) {
	node, err := cosmosrpc.NewNode(context.Background(), log.Logger, nodeConfig)
	suite.Require().NoError(err)
	suite.node = node
}
