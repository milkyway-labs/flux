package types

type ABCIEvent struct {
	Type       string               `json:"type"`
	Attributes []ABCIEventAttribute `json:"attributes"`
}

func (e *ABCIEvent) FindAttribute(key string) (ABCIEventAttribute, bool) {
	for _, attribute := range e.Attributes {
		if attribute.Key == key {
			return attribute, true
		}
	}

	return ABCIEventAttribute{}, false
}

func (e *ABCIEvent) FindAttributes(key string) []ABCIEventAttribute {
	var results []ABCIEventAttribute
	for _, attribue := range e.Attributes {
		if attribue.Key == key {
			results = append(results, attribue)
		}
	}

	return results
}

func (e *ABCIEvent) FindAttributeFunc(predicate func(a ABCIEventAttribute) bool) (ABCIEventAttribute, bool) {
	for _, attribute := range e.Attributes {
		if predicate(attribute) {
			return attribute, true
		}
	}

	return ABCIEventAttribute{}, false
}

func (e *ABCIEvent) FindAttributesFunc(predicate func(a ABCIEventAttribute) bool) []ABCIEventAttribute {
	var results []ABCIEventAttribute
	for _, attribute := range e.Attributes {
		if predicate(attribute) {
			results = append(results, attribute)
		}
	}

	return results
}

type ABCIEventAttribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ABCIEvents []ABCIEvent

func (events *ABCIEvents) GetEventWitType(t string) (ABCIEvent, bool) {
	for _, e := range *events {
		if e.Type == t {
			return e, true
		}
	}

	return ABCIEvent{}, false
}

func (events *ABCIEvents) GetEventsWithType(t string) []ABCIEvent {
	var results []ABCIEvent
	for _, e := range *events {
		if e.Type == t {
			results = append(results, e)
		}
	}

	return results
}

type ABCIMsgLog struct {
	MsgIndex uint32     `json:"msg_index"`
	Events   ABCIEvents `json:"events"`
}
