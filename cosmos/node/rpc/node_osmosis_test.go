package rpc_test

import (
	"context"

	"github.com/goccy/go-json"

	"github.com/milkyway-labs/chain-indexer/cosmos/node/rpc"
)

func (suite *NodeTestSuite) TestOsmosisGetBlockResults() {
	suite.SetupSuite(rpc.DefaultConfig("https://rpc.osmosis.zone"))

	height, err := suite.node.GetCurrentHeight(context.Background())
	suite.Require().NoError(err)

	block, err := suite.node.GetBlock(context.Background(), height)
	suite.Require().NoError(err)

	jsonblock, err := json.Marshal(block)
	suite.Require().NoError(err)
	println(string(jsonblock))
}
