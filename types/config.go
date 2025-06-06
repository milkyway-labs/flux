package types

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

type RawConfig map[string]any

type Config struct {
	Logging    LoggingConfig    `yaml:"logging"`
	Monitoring MonitoringConfig `yaml:"monitoring"`

	// Databases contains the configurations of the databases that can be
	// used by the indexers to store the indexed data.
	Databases map[string]RawConfig `yaml:"databases"`

	// Nodes contains the configurations of the nodes that
	Nodes map[string]RawConfig `yaml:"nodes"`

	// Modules contains the configurations of the modules that can be
	// used to index a chain.
	Modules map[string]RawConfig `yaml:"modules"`

	// Indexers represents the indexers that will be spawned.
	Indexers []IndexerConfig `yaml:"indexers"`
}

var DefaultConfig = Config{
	Logging:    DefaultLoggingConfig(),
	Monitoring: DefaultMonitoringCfg,
}

func ParseConfig(configBytes []byte) (*Config, error) {
	config := Config{}
	err := yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &config, nil
}

func (cfg *Config) Validate() error {
	if err := cfg.Logging.Validate(); err != nil {
		return fmt.Errorf("invalid logging config: %w", err)
	}

	if len(cfg.Databases) == 0 {
		return fmt.Errorf("databases can't be empty")
	}

	if len(cfg.Nodes) == 0 {
		return fmt.Errorf("nodes can't be empty")
	}

	if len(cfg.Indexers) == 0 {
		return fmt.Errorf("indexers list can't be empty")
	}

	indexersName := make(map[string]any)
	for _, indexerCfg := range cfg.Indexers {
		if err := indexerCfg.Validate(); err != nil {
			return err
		}

		_, ok := indexersName[indexerCfg.Name]
		if ok {
			return fmt.Errorf("duplicated indexer with name: %s", indexerCfg.Name)
		}
		indexersName[indexerCfg.Name] = true
	}

	return nil
}

func (cfg *Config) GetIndexerConfig(name string) (*IndexerConfig, error) {
	for _, indexerConfig := range cfg.Indexers {
		if indexerConfig.Name == name {
			return &indexerConfig, nil
		}
	}

	return nil, fmt.Errorf("config for indexer %s not found", name)
}

func (cfg *Config) UnmarshalYAML(unmarshal func(any) error) error {
	// Local type to avoid recursion during the unmarshal
	type privateCfg Config
	config := privateCfg(DefaultConfig)
	err := unmarshal(&config)
	if err != nil {
		return err
	}

	*cfg = Config(config)
	return nil
}

// ----------------------------------------------------------------------------
// ---- Logging config
// ----------------------------------------------------------------------------

type LoggingConfig struct {
	LogLevel  string `yaml:"level"`
	LogFormat string `yaml:"format"`
}

// NewLoggingConfig returns a new LoggingConfigInstance instance
func NewLoggingConfig(level string, format string) LoggingConfig {
	return LoggingConfig{
		LogLevel:  level,
		LogFormat: format,
	}
}

// Implements the Unmarshaler interface of the yaml pkg.
func (cfg *LoggingConfig) UnmarshalYAML(unmarshal func(any) error) error {
	// Local type to avoid recursion during the unmarshal
	type privateCfg LoggingConfig
	config := privateCfg(DefaultLoggingConfig())
	err := unmarshal(&config)
	if err != nil {
		return err
	}

	*cfg = LoggingConfig(config)
	return nil
}

func (cfg *LoggingConfig) Validate() error {
	_, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("invalid log_level %s: %w", cfg.LogLevel, err)
	}

	if cfg.LogFormat != "json" && cfg.LogFormat != "text" {
		return fmt.Errorf("invalid log_format, we only support `text` and `json` current: `%s`", cfg.LogFormat)
	}

	return nil
}

// DefaultLoggingConfig returns the default LoggingConfigInstance instance
func DefaultLoggingConfig() LoggingConfig {
	return NewLoggingConfig(zerolog.DebugLevel.String(), "text")
}

// ----------------------------------------------------------------------------
// ---- Indexer config
// ----------------------------------------------------------------------------

type IndexerConfig struct {
	// Name represents the name that identifies this indexer.
	Name string `yaml:"name"`
	// NodeID represents the ID of the node that will be used to obtain the blocks
	// from the chain.
	NodeID string `yaml:"node_id"`
	// DatabaseID represents the ID of the database that the indexer will use
	// to index the data.
	DatabaseID string `yaml:"database_id"`
	// Workers represents the number of workers that will be spawned to fetch
	// a block ad process it.
	Workers uint32 `yaml:"workers"`
	// HeightQueueSize represents the maximum number of block heights that can be
	// queued for fetching. Once the number of queued elements reaches this value,
	// the indexer will stop monitoring the node for new blocks until space becomes
	// available in the queue.
	HeightQueueSize uint32 `yaml:"height_queue_size"`
	// NodePollingInterval interval with which we poll the node for
	// newer blocks.
	NodePollingInterval time.Duration `yaml:"node_polling_interval"`
	// Modules contains the names of the module that the indexer will use
	// to index a chain.
	Modules []string `yaml:"modules"`
	// OverrideModuleConfig allows to define custom configurations for a module
	// that will be used by the indexer.
	OverrideModuleConfig map[string]RawConfig `yaml:"override_module_config"`
	// StartHeight height from which the indexer will start fetching blocks.
	// If undefined the indexer will start indexing from the current node height.
	StartHeight *Height `yaml:"start_height"`
	// MaxAttempts represents the number of time the indexer will re-try to index
	// a block in case of failure.
	MaxAttempts uint32 `yaml:"max_attempts"`
	// Define the amount of time the indexer will wait before re-enqueuing a failed
	// block for parsing.
	TimeBeforeRetry time.Duration `yaml:"time_before_retry"`
}

var DefaultIndexerCfg = IndexerConfig{
	Workers:             1,
	HeightQueueSize:     100,
	NodePollingInterval: time.Second,
	MaxAttempts:         5,
	TimeBeforeRetry:     10 * time.Second,
}

// Implements the Unmarshaler interface of the yaml pkg.
func (cfg *IndexerConfig) UnmarshalYAML(unmarshal func(any) error) error {
	// Local type to avoid recursion during the unmarshal
	type privateIndexerCfg IndexerConfig
	config := privateIndexerCfg(DefaultIndexerCfg)
	err := unmarshal(&config)
	if err != nil {
		return err
	}

	*cfg = IndexerConfig(config)
	return nil
}

func (cfg *IndexerConfig) Validate() error {
	if cfg.Name == "" {
		return fmt.Errorf("indexer name can't be empty")
	}

	if cfg.NodeID == "" {
		return fmt.Errorf("node_id can't be empty")
	}

	if cfg.DatabaseID == "" {
		return fmt.Errorf("database_id can't be empty")
	}

	if cfg.Workers == 0 {
		return fmt.Errorf("worker must be > 0")
	}

	if cfg.HeightQueueSize == 0 {
		return fmt.Errorf("height_queue_size must be > 0")
	}

	if cfg.NodePollingInterval < 10*time.Millisecond {
		return fmt.Errorf("node_polling_interval must be >= then 10 milliseconds")
	}

	if cfg.TimeBeforeRetry < 10*time.Millisecond {
		return fmt.Errorf("time_before_retry must be >= then 10 milliseconds")
	}

	if len(cfg.Modules) == 0 {
		return fmt.Errorf("modules list can't be empty")
	}

	return nil
}

// ----------------------------------------------------------------------------
// ---- Monitoring config
// ----------------------------------------------------------------------------

type MonitoringConfig struct {
	Enabled bool  `yaml:"enabled"`
	Port    int16 `yaml:"port"`
}

var DefaultMonitoringCfg = MonitoringConfig{
	Enabled: true,
	Port:    2112,
}

// Implements the Unmarshaler interface of the yaml pkg.
func (cfg *MonitoringConfig) UnmarshalYAML(unmarshal func(any) error) error {
	// Local type to avoid recursion during the unmarshal
	type privateMonitoringCfg MonitoringConfig
	config := privateMonitoringCfg(DefaultMonitoringCfg)
	err := unmarshal(&config)
	if err != nil {
		return err
	}

	*cfg = MonitoringConfig(config)
	return nil
}
