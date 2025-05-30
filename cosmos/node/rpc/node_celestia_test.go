package rpc_test

import (
	"context"
	"time"

	"github.com/goccy/go-json"

	"github.com/milkyway-labs/flux/cosmos/node/rpc"
	"github.com/milkyway-labs/flux/types"
)

func (suite *NodeTestSuite) TestCelestiaGetBlockResults() {
	height := types.MaxHeight

	suite.SetupSuite(rpc.NewConfig(
		"https://celestia-rpc.publicnode.com",
		time.Second*10,
		&height,
		&height,
	))

	height, err := suite.node.GetCurrentHeight(context.Background())
	suite.Require().NoError(err)

	block, err := suite.node.GetBlock(context.Background(), height)
	suite.Require().NoError(err)

	jsonblock, err := json.Marshal(block)
	suite.Require().NoError(err)
	println(string(jsonblock))
}
