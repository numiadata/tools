package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// Supported sink types.
const (
	SinkTypeBigTable SinkType = "big_table"
)

var (
	supportedSinkTypes = map[SinkType]struct{}{
		SinkTypeBigTable: {},
	}

	validate = validator.New()
)

type (
	// SinkType defines a type alias for a supported sink type.
	SinkType string

	// Config defines the configuration structure for defining Erebus execution
	// parameters.
	Config struct {
		SinkType SinkType `mapstructure:"sink_type" validate:"required"`

		// TableName, when sink_type='big_table', defines the name of the table in
		// Bigtable to write to.
		TableName string `mapstructure:"table_name" validate:"required_if=SinkType big_table"`
	}
)

// Validate attempts to validate the config object returning an error upon
// failure.
func (c *Config) Validate() error {
	if _, ok := supportedSinkTypes[c.SinkType]; !ok {
		return fmt.Errorf("invalid sink type: %s", c.SinkType)
	}

	return validate.Struct(c)
}

// Parse attempts to parse a config file into a Config struct.
func Parse(cfgFile string) (Config, error) {
	v := viper.New()
	v.SetConfigFile(cfgFile)

	if err := v.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to read config file %s: %w", cfgFile, err)
	}

	cfg := Config{}
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse config file %s: %w", cfgFile, err)
	}

	return cfg, nil
}
