package postgresql_test

import (
	"time"

	"github.com/milkyway-labs/chain-indexer/database/postgresql"
	"github.com/milkyway-labs/chain-indexer/types"
)

func (suite *DbTestSuite) TestGetLowestBlock() {
	testHeigt := types.Height(11)
	testCases := []struct {
		name           string
		setup          func()
		shouldErr      bool
		chainID        string
		expectedHeigth *types.Height
	}{
		{
			name:           "empty database return nil",
			shouldErr:      false,
			chainID:        "test",
			expectedHeigth: nil,
		},
		{
			name: "return the correct height",
			setup: func() {
				suite.database.SaveIndexedBlock("test", 12, time.Now())
				suite.database.SaveIndexedBlock("test", 11, time.Now())
			},
			shouldErr:      false,
			chainID:        "test",
			expectedHeigth: &testHeigt,
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			if tc.setup != nil {
				tc.setup()
			}

			result, err := suite.database.GetLowestBlock(tc.chainID)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expectedHeigth, result)
			}
		})
	}
}

func (suite *DbTestSuite) TestGetMissingBlocks() {
	testCases := []struct {
		name            string
		setup           func()
		shouldErr       bool
		chainID         string
		from            types.Height
		to              types.Height
		expectedHeigths []types.Height
	}{
		{
			name:      "if from is higher then to fails",
			shouldErr: true,
			chainID:   "test",
			from:      3,
			to:        2,
		},
		{
			name:            "from equals to works correctly",
			shouldErr:       false,
			chainID:         "test",
			from:            3,
			to:              3,
			expectedHeigths: []types.Height{3},
		},
		{
			name:            "empty database return all heights",
			shouldErr:       false,
			chainID:         "test",
			from:            1,
			to:              3,
			expectedHeigths: []types.Height{1, 2, 3},
		},
		{
			name: "chain id is handled correctly",
			setup: func() {
				suite.database.SaveIndexedBlock("test", 12, time.Now())
				suite.database.SaveIndexedBlock("test", 11, time.Now())
			},
			shouldErr:       false,
			chainID:         "empty",
			from:            10,
			to:              13,
			expectedHeigths: []types.Height{10, 11, 12, 13},
		},
		{
			name: "return the correct heights",
			setup: func() {
				suite.database.SaveIndexedBlock("test", 12, time.Now())
				suite.database.SaveIndexedBlock("test", 11, time.Now())
			},
			shouldErr:       false,
			chainID:         "test",
			from:            10,
			to:              13,
			expectedHeigths: []types.Height{10, 13},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			if tc.setup != nil {
				tc.setup()
			}

			result, err := suite.database.GetMissingBlocks(tc.chainID, tc.from, tc.to)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(tc.expectedHeigths, result)
			}
		})
	}
}

func (suite *DbTestSuite) TestSaveIndexedBlock() {
	testCases := []struct {
		name      string
		setup     func()
		chainID   string
		height    types.Height
		timestamp time.Time
		shouldErr bool
		check     func()
	}{
		{
			name:      "save successfully indexed blocks",
			chainID:   "test",
			height:    11,
			timestamp: time.Date(2021, 11, 22, 14, 0, 0, 0, time.UTC),
			shouldErr: false,
			check: func() {
				var blockRow postgresql.BlockRow
				err := suite.database.SQL.Get(&blockRow, "SELECT * FROM blocks WHERE height = 11")
				suite.Require().NoError(err)

				suite.Require().Equal("test", blockRow.ChainID)
				suite.Require().Equal(types.Height(11), blockRow.Height)
				suite.Require().Equal(time.Date(2021, 11, 22, 14, 0, 0, 0, time.UTC), blockRow.Timestamp.UTC())
			},
		},
		{
			name: "constraint works correctly",
			setup: func() {
				suite.database.SaveIndexedBlock("test", 11, time.Now())
			},
			chainID:   "test",
			height:    11,
			timestamp: time.Date(2021, 11, 22, 14, 0, 0, 0, time.UTC),
			shouldErr: false,
			check: func() {
				var blockRow postgresql.BlockRow
				err := suite.database.SQL.Get(&blockRow, "SELECT * FROM blocks WHERE height = 11")
				suite.Require().NoError(err)

				suite.Require().Equal("test", blockRow.ChainID)
				suite.Require().Equal(types.Height(11), blockRow.Height)
				suite.Require().Equal(time.Date(2021, 11, 22, 14, 0, 0, 0, time.UTC), blockRow.Timestamp.UTC())
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			if tc.setup != nil {
				tc.setup()
			}

			err := suite.database.SaveIndexedBlock(tc.chainID, tc.height, tc.timestamp)
			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
				if tc.check != nil {
					tc.check()
				}
			}
		})
	}
}
