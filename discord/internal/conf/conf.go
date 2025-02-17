package conf

import (
	"log/slog"
	"strings"

	"github.com/calamity-m/reaphur/pkg/bindings"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Config struct {
	// Logging Configuration
	Environment   string     `mapstructure:"environment" json:"environment,omitempty"`
	LogLevel      slog.Level `mapstructure:"log_level" json:"log_level,omitempty"`
	LogStructured bool       `mapstructure:"log_structured" json:"log_structured,omitempty"`
	LogAddSource  bool       `mapstructure:"log_add_source" json:"log_add_source,omitempty"`
	LogRequestId  bool       `mapstructure:"log_request_id" json:"log_request_id,omitempty"`

	// Endpoint Configuration
	CentralServerAddress string `mapstructure:"central_server_address" json:"central_server_address,omitempty"`

	// Spicy
	BotToken string `mapstructure:"bot_token" json:"-"`
}

func NewConfig(debug bool) (*Config, error) {
	// Setup
	base := &Config{}
	vip := viper.New()

	// Enable ENV var reading
	vip.AutomaticEnv()
	vip.SetEnvPrefix("DISCORD")
	vip.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Establish sane defaults before overriding
	vip.SetDefault("environment", "dev")
	vip.SetDefault("log_level", slog.LevelDebug)
	vip.SetDefault("log_structured", false)
	vip.SetDefault("log_add_source", true)
	vip.SetDefault("log_request_id", true)
	vip.SetDefault("central_server_address", bindings.DefaultCentralAddress)

	// Spicy bindings
	if err := vip.BindEnv("bot_token"); err != nil {
		return &Config{}, err
	}
	// Magic to unamrshal viper into the config sturct. The decode hook is used to map things like the logging level
	// into the slog logging level type.
	if err := vip.Unmarshal(&base, viper.DecodeHook(mapstructure.TextUnmarshallerHookFunc())); err != nil {
		return &Config{}, err
	}

	// Forecefully override if we're on debug mode
	// Do this without viper/cobra buy-in and just do the simple
	// brute force
	if debug {
		base.LogLevel = slog.LevelDebug
	}

	return base, nil
}
