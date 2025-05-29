package postgresql

import (
	"fmt"
	"net/url"
)

type Config struct {
	URL           string `yaml:"url"`
	PartitionSize int64  `yaml:"partition_size"`
}

func NewConfig(url string, partitionSize int64) Config {
	return Config{
		URL:           url,
		PartitionSize: partitionSize,
	}
}

func (c Config) WithURL(url string) Config {
	c.URL = url
	return c
}

func (c Config) WithPartitionSize(partitionSize int64) Config {
	c.PartitionSize = partitionSize
	return c
}

func (c Config) GetSchema() string {
	u, err := url.Parse(c.URL)
	if err != nil {
		return "public"
	}

	searchPath := u.Query().Get("search_path")
	if searchPath == "" {
		return "public"
	}
	return searchPath
}

func (c Config) GetPartitionSize() int64 {
	if c.PartitionSize > 0 {
		return c.PartitionSize
	}
	return 100_0000
}

func (c Config) Validate() error {
	_, err := url.Parse(c.URL)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}

	return nil
}

func DefaultConfig() Config {
	return NewConfig(
		"postgresql://user:password@localhost:5432/database-name?sslmode=disable&search_path=public",
		100000,
	)
}

// Implements the Unmarshaler interface of the yaml pkg.
func (c *Config) UnmarshalYAML(unmarshal func(any) error) error {
	// Local type to avoid recursion during the unmarshal
	type privateCfg Config
	config := privateCfg(DefaultConfig())
	err := unmarshal(&config)
	if err != nil {
		return err
	}

	*c = Config(config)
	return nil
}
