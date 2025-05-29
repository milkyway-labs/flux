package rpc_test

import (
	"context"
	"time"

	"github.com/milkyway-labs/chain-indexer/evm/node/rpc"
)

func (suite *NodeTestSuite) TestPolygonNode() {
	suite.SetupSuite(rpc.NewConfig(
		"https://polygon-rpc.com",
		// "https://polygon-mainnet.public.blastapi.io",
		time.Second*10,
	))

	height, err := suite.node.GetCurrentHeight(context.Background())
	suite.Require().NoError(err)

	_, err = suite.node.GetEthBlock(context.Background(), height)
	suite.Require().NoError(err)

	_, err = suite.node.GetLowestHeight(context.Background())
	suite.Require().NoError(err)

	_, err = suite.node.GetLogs(context.Background(), height)
	suite.Require().NoError(err)

	_, err = suite.node.GetBlock(context.Background(), height)
	suite.Require().NoError(err)
}
