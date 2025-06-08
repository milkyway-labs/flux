package suite

import (
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/milkyway-labs/flux/database"
	"github.com/milkyway-labs/flux/types"
)

type BeforeTestHook func()

// Suite represents a test suite that can be used to verify the correct
// behavior of a Database implementation.
type Suite struct {
	suite.Suite

	database       database.Database
	beforeTestHook BeforeTestHook
}

// InitDB sets the database instance under test
func (s *Suite) InitDB(database database.Database) {
	s.database = database
}

// WithBeforeTestHook configures the hook that will be called before each test.
func (s *Suite) WithBeforeTestHook(hook BeforeTestHook) *Suite {
	s.beforeTestHook = hook
	return s
}

func (s *Suite) executeBeforeTestHook() {
	if s.beforeTestHook != nil {
		s.beforeTestHook()
	}
}

func (s *Suite) TestGetLowestBlock() {
	testHeigt := types.Height(9)
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
				s.database.SaveIndexedBlock("test", 9, time.Now())
				s.database.SaveIndexedBlock("test", 12, time.Now())
				s.database.SaveIndexedBlock("test", 11, time.Now())
			},
			shouldErr:      false,
			chainID:        "test",
			expectedHeigth: &testHeigt,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.executeBeforeTestHook()
			if tc.setup != nil {
				tc.setup()
			}

			result, err := s.database.GetLowestBlock(tc.chainID)
			if tc.shouldErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				if tc.expectedHeigth == nil {
					s.Require().Nil(result)
				} else {
					s.Require().Equal(*tc.expectedHeigth, *result)
				}
			}
		})
	}
}

func (s *Suite) TestGetMissingBlocks() {
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
				s.database.SaveIndexedBlock("test", 12, time.Now())
				s.database.SaveIndexedBlock("test", 11, time.Now())
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
				s.database.SaveIndexedBlock("test", 12, time.Now())
				s.database.SaveIndexedBlock("test", 11, time.Now())
			},
			shouldErr:       false,
			chainID:         "test",
			from:            10,
			to:              13,
			expectedHeigths: []types.Height{10, 13},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.executeBeforeTestHook()
			if tc.setup != nil {
				tc.setup()
			}

			result, err := s.database.GetMissingBlocks(tc.chainID, tc.from, tc.to)
			if tc.shouldErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				s.Require().Equal(tc.expectedHeigths, result)
			}
		})
	}
}

func (s *Suite) TestSaveIndexedBlock() {
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
				heights, err := s.database.GetMissingBlocks("test", 11, 11)
				s.Require().NoError(err)
				s.Require().Empty(heights)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.executeBeforeTestHook()
			if tc.setup != nil {
				tc.setup()
			}

			err := s.database.SaveIndexedBlock(tc.chainID, tc.height, tc.timestamp)
			if tc.shouldErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				if tc.check != nil {
					tc.check()
				}
			}
		})
	}
}
