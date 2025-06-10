package types

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

// Base64Bytes enables bas64-encoding for json/encoding.
type Base64Bytes []byte

func (bz Base64Bytes) MarshalJSON() ([]byte, error) {
	s := base64.StdEncoding.EncodeToString(bz)
	jbz := make([]byte, len(s)+2)
	jbz[0] = '"'
	copy(jbz[1:], s)
	jbz[len(jbz)-1] = '"'
	return jbz, nil
}

func (bz *Base64Bytes) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*bz = make(Base64Bytes, 0)
		return nil
	}

	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return fmt.Errorf("invalid base64 string: %s", data)
	}
	bz2, err := base64.StdEncoding.DecodeString(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}
	*bz = bz2
	return nil
}

func (bz Base64Bytes) Bytes() []byte {
	return bz
}

func (bz Base64Bytes) String() string {
	return base64.StdEncoding.EncodeToString(bz)
}

// ---------------------------------------------------------------------------

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
