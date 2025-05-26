package jsonrpc2

import (
	"fmt"

	"github.com/goccy/go-json"
)

const ProtocolVersion = "2.0"

type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

func NewRequest(id any, method string, params json.RawMessage) Request {
	return Request{
		JSONRPC: ProtocolVersion,
		ID:      id,
		Method:  method,
		Params:  params,
	}
}

type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *Error          `json:"error,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (e Error) Error() string {
	return fmt.Sprintf("(%d) %s: %v", e.Code, e.Message, e.Data)
}
