package rpc

import (
	"fmt"
	"net/url"
	"time"
)

type Config struct {
	URL            string        `yaml:"url"`
	RequestTimeout time.Duration `yaml:"request_timeout"`
	// Tells if the tx events should be parsed from the tx.log field.
	TxEventsFromLog bool `yaml:"tx_events_from_log"`
	// Tells if the block events should be treated as base64 encoded and needs to
	// be decoded.
	DecodeBlockEventAttributes bool `yaml:"decode_block_event_attributes"`
}

func NewConfig(
	url string,
	timeout time.Duration,
	txEventsFromLog bool,
	decodeBlockEventAttributes bool,
) Config {
	return Config{
		URL:                        url,
		RequestTimeout:             timeout,
		TxEventsFromLog:            txEventsFromLog,
		DecodeBlockEventAttributes: decodeBlockEventAttributes,
	}
}

func (c *Config) Validate() error {
	_, err := url.Parse(c.URL)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}

	return nil
}

func DefaultConfig(url string) Config {
	return NewConfig(url, time.Second*10, false, false)
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
