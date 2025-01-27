package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	JWT      JWTConfig
	Database DatabaseConfig
	Security SecurityConfig
	Logging  LoggingConfig
	Metrics  MetricsConfig
}

type LoggingConfig struct {
	Level      string
	Format     string
	Output     string
	TimeFormat string
}

type MetricsConfig struct {
	Enabled     bool
	Endpoint    string
	ServiceName string
}

type ServerConfig struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	TLSEnabled   bool
	TLSCertFile  string
	TLSKeyFile   string
}

type JWTConfig struct {
	SecretKey           string
	Expiration          time.Duration
	RefreshTokenSecret  string
	RefreshTokenExpiry  time.Duration
	TokenRotationEnable bool
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type SecurityConfig struct {
	RateLimit RateLimitConfig
	Password  PasswordConfig
	Headers   HeadersConfig
}

type RateLimitConfig struct {
	Enabled           bool
	RequestsPerMinute int
	BurstSize         int
}

type PasswordConfig struct {
	MinLength      int
	RequireUpper   bool
	RequireLower   bool
	RequireNumber  bool
	RequireSpecial bool
	MaxAttempts    int
}

type HeadersConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	EnableCSP        bool
	CSPDirectives    string
	EnableHSTS       bool
	HSTSMaxAge       time.Duration
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func validateConfig(config *Config) error {
	if config.Server.Port == 0 {
		return fmt.Errorf("server port is required")
	}
	if config.JWT.SecretKey == "" {
		return fmt.Errorf("JWT secret key is required")
	}
	if config.Database.Host == "" || config.Database.DBName == "" {
		return fmt.Errorf("database host and name are required")
	}
	return nil
}
