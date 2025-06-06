package grpc

import (
	"github.com/milkyway-labs/flux/types"
)

// ---------------------------------------------------------------------------

type ABCIQueryResponse struct {
	Code   uint32 `json:"code"`
	Log    string `json:"log"`
	Key    []byte `json:"key"`
	Value  []byte `json:"value"`
	Height int64  `json:"height,string"`
}

type ABCIQueryResult struct {
	Response ABCIQueryResponse `json:"response"`
}

func (resp ABCIQueryResponse) IsOK() bool {
	return resp.Code == 0
}

type ABCIQueryRequest struct {
	Path   string         `json:"path"`
	Data   types.HexBytes `json:"data"`
	Height types.Height   `json:"height,string"`
	Prove  bool           `json:"prove"`
}
