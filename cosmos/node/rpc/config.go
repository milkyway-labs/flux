package rpc

import (
	"fmt"
	"net/url"
	"time"

	"github.com/milkyway-labs/chain-indexer/types"
)

type Config struct {
	URL            string        `yaml:"url"`
	RequestTimeout time.Duration `yaml:"request_timeout"`
	// Tells until which height the indexer will parse the tx.log field to get the
	// transaction events. After this height, the indexer will use the tx.events
	// field directly. TxEventsFromLogUntilHeight is nil, the indexer will always use
	// tx.events.
	TxEventsFromLogUntilHeight *types.Height `yaml:"tx_events_from_log_until_height"`
	// Tells until which height the indexer will treat the block events as base64
	// encoded and needs to be decoded. If DecodeBlockEventAttributesUntilHeight is
	// nil, the indexer will not decode the block events.
	DecodeBlockEventAttributesUntilHeight *types.Height `yaml:"decode_block_event_attributes_until_height"`
}

func NewConfig(
	url string,
	timeout time.Duration,
	txEventsFromLogUntilHeight *types.Height,
	decodeBlockEventAttributesUntilHeight *types.Height,
) Config {
	return Config{
		URL:                                   url,
		RequestTimeout:                        timeout,
		TxEventsFromLogUntilHeight:            txEventsFromLogUntilHeight,
		DecodeBlockEventAttributesUntilHeight: decodeBlockEventAttributesUntilHeight,
	}
}

func (c *Config) Validate() error {
	_, err := url.Parse(c.URL)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}

	return nil
}

func (c *Config) TxEventsFromLog(height types.Height) bool {
	return c.TxEventsFromLogUntilHeight != nil && height <= *c.TxEventsFromLogUntilHeight
}

func (c *Config) DecodeBlockEventAttributes(height types.Height) bool {
	return c.TxEventsFromLogUntilHeight != nil && height <= *c.TxEventsFromLogUntilHeight
}

func DefaultConfig(url string) Config {
	return NewConfig(url, time.Second*10, nil, nil)
}

// Implements the Unmarshaler interface of the yaml pkg.
func (c *Config) UnmarshalYAML(unmarshal func(any) error) error {
	// Local type to avoid recursion during the unmarshal
	type privateCfg Config
	config := privateCfg(DefaultConfig(""))
	err := unmarshal(&config)
	if err != nil {
		return err
	}

	*c = Config(config)
	return nil
}
