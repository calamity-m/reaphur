package conf

import (
	"log/slog"
	"strings"

	"github.com/calamity-m/reaphur/pkg/bindings"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type Config struct {
	// Logging Configuration
	Environment   string     `mapstructure:"environment" json:"environment,omitempty"`
	LogLevel      slog.Level `mapstructure:"log_level" json:"log_level,omitempty"`
	LogStructured bool       `mapstructure:"log_structured" json:"log_structured,omitempty"`
	LogAddSource  bool       `mapstructure:"log_add_source" json:"log_add_source,omitempty"`
	LogRequestId  bool       `mapstructure:"log_request_id" json:"log_request_id,omitempty"`

	// Server configuration
	Address          string `mapstructure:"address" json:"address,omitempty"`
	Reflect          bool   `mapstructure:"reflect" json:"reflect,omitempty"`
	FoodRedisAddress string `mapstructure:"food_redis_address" json:"food_redis_address,omitempty"`

	// Spicy
	AIToken           string `mapstructure:"ai_token" json:"-"`
	FoodRedisPassword string `mapstructure:"food_redis_password" json:"-"`
	FoodRedisDB       int    `mapstructure:"food_redis_db" json:"-"`

	// Not filled out by viper defaults
	GrpcServerOpts []grpc.ServerOption
}

func NewConfig(debug bool) (*Config, error) {
	// Setup
	base := &Config{}
	vip := viper.New()

	// Enable ENV var reading
	vip.AutomaticEnv()
	vip.SetEnvPrefix("CENTRAL")
	vip.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// Establish sane defaults before overriding
	vip.SetDefault("environment", "dev")
	vip.SetDefault("log_level", slog.LevelDebug)
	vip.SetDefault("log_structured", false)
	vip.SetDefault("log_add_source", true)
	vip.SetDefault("log_request_id", true)
	vip.SetDefault("address", bindings.DefaultCentralAddress)
	vip.SetDefault("reflect", true)
	vip.SetDefault("food_redis_password", "password")
	vip.SetDefault("food_redis_db", 0)

	// Spicy bindings
	vip.BindEnv("ai_token")

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
