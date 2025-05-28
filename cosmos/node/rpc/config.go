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
	// Tells from which height the indexer will start
	// parse the tx events from the tx.events field, otherwise will
	// parse the tx events from the tx.log field.
	// If this field is undefined will always parse the events from the
	// tx.log field.
	TxEventsFromEventsFromHeigth *types.Height `yaml:"tx_events_from_height"`
	// Tells if the block events should be treated as base64 encoded and needs to
	// be decoded.
	DecodeBlockEventAttributes bool `yaml:"decode_block_event_attributes"`
}

func NewConfig(
	url string,
	timeout time.Duration,
	txEventsFromEventsFromHeigth *types.Height,
	decodeBlockEventAttributes bool,
) Config {
	return Config{
		URL:                          url,
		RequestTimeout:               timeout,
		TxEventsFromEventsFromHeigth: txEventsFromEventsFromHeigth,
		DecodeBlockEventAttributes:   decodeBlockEventAttributes,
	}
}

func (c *Config) Validate() error {
	_, err := url.Parse(c.URL)
	if err != nil {
		return fmt.Errorf("invalid url %w", err)
	}

	return nil
}

func (c *Config) TxEventsFromEvents(height types.Height) bool {
	return c.TxEventsFromEventsFromHeigth != nil && height >= *c.TxEventsFromEventsFromHeigth
}

func DefaultConfig(url string) Config {
	txEventsFromEventsFromHeigth := types.Height(0)
	return NewConfig(url, time.Second*10, &txEventsFromEventsFromHeigth, false)
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
