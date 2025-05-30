package types

import (
	"slices"

	"github.com/milkyway-labs/flux/utils"
)

type ABCIEventAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ---------------------------------------------------------------------------
// ---- ABCI Event
// ---------------------------------------------------------------------------

type ABCIEvent struct {
	Type       string               `json:"type"`
	Attributes []ABCIEventAttribute `json:"attributes"`
}

// FindAttribute finds the first attribute that matches the given predicate
func (e *ABCIEvent) FindAttributeFunc(predicate func(a ABCIEventAttribute) bool) (ABCIEventAttribute, bool) {
	index := slices.IndexFunc(e.Attributes, predicate)
	if index == -1 {
		return ABCIEventAttribute{}, false
	}

	return e.Attributes[index], true
}

// FindAttribute finds the first attribute with the given key
func (e *ABCIEvent) FindAttribute(key string) (ABCIEventAttribute, bool) {
	return e.FindAttributeFunc(func(a ABCIEventAttribute) bool {
		return a.Key == key
	})
}

// FindAttributesFunc finds all the attributes that match the given predicate
func (e *ABCIEvent) FindAttributesFunc(predicate func(a ABCIEventAttribute) bool) []ABCIEventAttribute {
	return utils.Filter(e.Attributes, predicate)
}

// FindAttributes finds all the attributes with the given key
func (e *ABCIEvent) FindAttributes(key string) []ABCIEventAttribute {
	return e.FindAttributesFunc(func(a ABCIEventAttribute) bool {
		return a.Key == key
	})
}

// ---------------------------------------------------------------------------
// ---- ABCI Event slice
// ---------------------------------------------------------------------------

type ABCIEvents []ABCIEvent

// FindEventFunc finds the first event that matches the given predicate
func (events ABCIEvents) FindEventFunc(predicate func(e ABCIEvent) bool) (ABCIEvent, bool) {
	index := slices.IndexFunc(events, predicate)

	if index == -1 {
		return ABCIEvent{}, false
	}

	return events[index], true
}

// FindEventWithType finds the first event with the given type
func (events ABCIEvents) FindEventWithType(t string) (ABCIEvent, bool) {
	return events.FindEventFunc(func(e ABCIEvent) bool {
		return e.Type == t
	})
}

// FindEventsFunc finds all the events that match the given predicate
func (events ABCIEvents) FindEventsFunc(predicate func(e ABCIEvent) bool) []ABCIEvent {
	return utils.Filter(events, predicate)
}

// FindEventsWithType finds all the events with the given type
func (events ABCIEvents) FindEventsWithType(t string) []ABCIEvent {
	return events.FindEventsFunc(func(e ABCIEvent) bool {
		return e.Type == t
	})
}

// ---------------------------------------------------------------------------
// ---- ABCI Event contained inside the tx log
// ---------------------------------------------------------------------------

type ABCIMsgLog struct {
	MsgIndex uint32     `json:"msg_index"`
	Events   ABCIEvents `json:"events"`
}
