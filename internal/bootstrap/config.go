package bootstrap

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all runtime configuration for the service.
type Config struct {
	Database DatabaseConfig
	GRPC     GRPCConfig
	Metrics  MetricsConfig
}

type DatabaseConfig struct {
	URL      string `mapstructure:"url"`
	MaxConns int32  `mapstructure:"max_conns"`
}

type GRPCConfig struct {
	Port int `mapstructure:"port"`
}

type MetricsConfig struct {
	Port int `mapstructure:"port"`
}

// GRPCAddr returns the full listen address for the gRPC server.
func (c Config) GRPCAddr() string {
	return fmt.Sprintf(":%d", c.GRPC.Port)
}

// MetricsAddr returns the full listen address for the metrics HTTP server.
func (c Config) MetricsAddr() string {
	return fmt.Sprintf(":%d", c.Metrics.Port)
}

func MustLoad(env string) Config {
	v := viper.New()
	v.SetConfigName(env)
	v.SetConfigType("yaml")
	v.AddConfigPath("config")
	v.AddConfigPath("../../config") // when called from cmd/server

	// Allow env vars to override yml values: APP_DATABASE_URL → database.url
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("failed to load config/%s.yaml: %v", env, err))
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		panic(fmt.Sprintf("failed to unmarshal config: %v", err))
	}
	return cfg
}
