package rpc

import (
	"fmt"
	"net/url"
	"time"
)

type Config struct {
	URL            string        `yaml:"url"`
	RequestTimeout time.Duration `yaml:"request_timeout"`
}

func NewConfig(
	url string,
	timeout time.Duration,
) Config {
	return Config{
		URL:            url,
		RequestTimeout: timeout,
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
	return NewConfig(url, time.Second*10)
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
