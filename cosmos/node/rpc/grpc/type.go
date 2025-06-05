package grpc

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/milkyway-labs/flux/types"
)

// HexBytes enables HEX-encoding for json/encoding.
type HexBytes []byte

// This is the point of Bytes.
func (bz HexBytes) MarshalJSON() ([]byte, error) {
	s := strings.ToUpper(hex.EncodeToString(bz))
	jbz := make([]byte, len(s)+2)
	jbz[0] = '"'
	copy(jbz[1:], s)
	jbz[len(jbz)-1] = '"'
	return jbz, nil
}

// This is the point of Bytes.
func (bz *HexBytes) UnmarshalJSON(data []byte) error {
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid hex string: %s", data)
	}
	bz2, err := hex.DecodeString(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}
	*bz = bz2
	return nil
}

// Bytes fulfills various interfaces in light-client, etc...
func (bz HexBytes) Bytes() []byte {
	return bz
}

func (bz HexBytes) String() string {
	return strings.ToUpper(hex.EncodeToString(bz))
}

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
	Path   string       `json:"path"`
	Data   HexBytes     `json:"data"`
	Height types.Height `json:"height,string"`
	Prove  bool         `json:"prove"`
}
