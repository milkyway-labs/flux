package rpc

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/goccy/go-json"
	cosmostypes "github.com/milkyway-labs/chain-indexer/cosmos/types"
)

func ParseEventsFromTxLog(log string) (cosmostypes.ABCIEvents, error) {
	var result cosmostypes.ABCIEvents

	var abciLogs []cosmostypes.ABCIMsgLog
	err := json.Unmarshal([]byte(log), &abciLogs)
	if err != nil {
		return nil, err
	}

	// Add the message index attribute to all the events
	for _, abciLog := range abciLogs {
		for _, event := range abciLog.Events {
			event.Attributes = append(event.Attributes, cosmostypes.ABCIEventAttribute{
				Key:   "msg_index",
				Value: strconv.FormatUint(uint64(abciLog.MsgIndex), 10),
			})

			// Add the event to the txEvents
			result = append(result, event)
		}
	}

	return result, nil
}

func DecodeABCIEvents(events cosmostypes.ABCIEvents) (cosmostypes.ABCIEvents, error) {
	if len(events) == 0 {
		return events, nil
	}

	decodedEvents := make(cosmostypes.ABCIEvents, len(events))
	for i, event := range events {
		decoded, err := DecodeABCIEvent(event)
		if err != nil {
			return nil, err
		}
		decodedEvents[i] = decoded
	}

	return decodedEvents, nil
}

func DecodeABCIEvent(event cosmostypes.ABCIEvent) (cosmostypes.ABCIEvent, error) {
	decodecAttributes, err := DecodeABCIEventAttributes(event.Attributes)
	if err != nil {
		return cosmostypes.ABCIEvent{}, fmt.Errorf("decoding %s event attributes", event.Type)
	}

	return cosmostypes.ABCIEvent{
		Type:       event.Type,
		Attributes: decodecAttributes,
	}, nil
}

func DecodeABCIEventAttributes(attributes []cosmostypes.ABCIEventAttribute) ([]cosmostypes.ABCIEventAttribute, error) {
	if len(attributes) == 0 {
		return attributes, nil
	}

	decoded := make([]cosmostypes.ABCIEventAttribute, len(attributes))
	for i, attribute := range attributes {
		decodedValue, err := base64.StdEncoding.DecodeString(attribute.Value)
		if err != nil {
			return nil, err
		}
		decoded[i] = cosmostypes.ABCIEventAttribute{
			Key:   attribute.Key,
			Value: string(decodedValue),
		}
	}

	return decoded, nil
}
