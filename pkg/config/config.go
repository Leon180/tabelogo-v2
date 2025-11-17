package config

import (
	"fmt"
	"strings"
	"time"
)

// Config holds application configuration
type Config struct {
	Environment string
	LogLevel    string

	// Server configuration
	ServerPort int
	GRPCPort   int

	// Database configuration
	Database DatabaseConfig

	// Redis configuration
	Redis RedisConfig

	// Kafka configuration
	Kafka KafkaConfig

	// JWT configuration
	JWT JWTConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string

	// Connection pool settings
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Brokers string // Comma-separated list of brokers
	GroupID string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret             string
	AccessTokenExpire  time.Duration
	RefreshTokenExpire time.Duration
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	return LoadWithPrefix("")
}

// LoadWithPrefix loads configuration from environment variables with prefix
func LoadWithPrefix(prefix string) (*Config, error) {
	cfg := &Config{
		Environment: normalizeEnvironment(getEnvWithDefault(buildEnvKey(prefix, "ENVIRONMENT"), "development")),
		LogLevel:    getEnvWithDefault(buildEnvKey(prefix, "LOG_LEVEL"), "info"),
		ServerPort:  getEnvAsInt(buildEnvKey(prefix, "SERVER_PORT"), 8080),
		GRPCPort:    getEnvAsInt(buildEnvKey(prefix, "GRPC_PORT"), 9090),
	}

	// Load database config
	cfg.Database = DatabaseConfig{
		Host:            getEnvWithDefault(buildEnvKey(prefix, "DB_HOST"), "localhost"),
		Port:            getEnvAsInt(buildEnvKey(prefix, "DB_PORT"), 5432),
		Name:            getEnvWithDefault(buildEnvKey(prefix, "DB_NAME"), ""),
		User:            getEnvWithDefault(buildEnvKey(prefix, "DB_USER"), "postgres"),
		Password:        getEnvWithDefault(buildEnvKey(prefix, "DB_PASSWORD"), "postgres"),
		SSLMode:         getEnvWithDefault(buildEnvKey(prefix, "DB_SSLMODE"), "disable"),
		MaxOpenConns:    getEnvAsInt(buildEnvKey(prefix, "DB_MAX_OPEN_CONNS"), 100),
		MaxIdleConns:    getEnvAsInt(buildEnvKey(prefix, "DB_MAX_IDLE_CONNS"), 10),
		ConnMaxLifetime: getEnvAsDuration(buildEnvKey(prefix, "DB_CONN_MAX_LIFETIME"), time.Hour),
	}

	// Load Redis config
	cfg.Redis = RedisConfig{
		Host:     getEnvWithDefault(buildEnvKey(prefix, "REDIS_HOST"), "localhost"),
		Port:     getEnvAsInt(buildEnvKey(prefix, "REDIS_PORT"), 6379),
		Password: getEnvWithDefault(buildEnvKey(prefix, "REDIS_PASSWORD"), ""),
		DB:       getEnvAsInt(buildEnvKey(prefix, "REDIS_DB"), 0),
	}

	// Load Kafka config
	cfg.Kafka = KafkaConfig{
		Brokers: getEnvWithDefault(buildEnvKey(prefix, "KAFKA_BROKERS"), "localhost:9092"),
		GroupID: getEnvWithDefault(buildEnvKey(prefix, "KAFKA_GROUP_ID"), "tabelogo-group"),
	}

	// Load JWT config
	cfg.JWT = JWTConfig{
		Secret:             getEnvWithDefault(buildEnvKey(prefix, "JWT_SECRET"), "change-me-in-production"),
		AccessTokenExpire:  getEnvAsDuration(buildEnvKey(prefix, "JWT_ACCESS_TOKEN_EXPIRE"), 15*time.Minute),
		RefreshTokenExpire: getEnvAsDuration(buildEnvKey(prefix, "JWT_REFRESH_TOKEN_EXPIRE"), 7*24*time.Hour),
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate database configuration
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if err := validatePort(c.Database.Port, "DB_PORT"); err != nil {
		return err
	}
	if c.Database.MaxOpenConns <= 0 {
		return fmt.Errorf("DB_MAX_OPEN_CONNS must be positive, got %d", c.Database.MaxOpenConns)
	}
	if c.Database.MaxIdleConns <= 0 {
		return fmt.Errorf("DB_MAX_IDLE_CONNS must be positive, got %d", c.Database.MaxIdleConns)
	}
	if c.Database.MaxIdleConns > c.Database.MaxOpenConns {
		return fmt.Errorf("DB_MAX_IDLE_CONNS (%d) cannot exceed DB_MAX_OPEN_CONNS (%d)",
			c.Database.MaxIdleConns, c.Database.MaxOpenConns)
	}
	if c.Database.ConnMaxLifetime < 0 {
		return fmt.Errorf("DB_CONN_MAX_LIFETIME must be non-negative")
	}

	// Validate server ports
	if err := validatePort(c.ServerPort, "SERVER_PORT"); err != nil {
		return err
	}
	if err := validatePort(c.GRPCPort, "GRPC_PORT"); err != nil {
		return err
	}
	if c.ServerPort == c.GRPCPort {
		return fmt.Errorf("SERVER_PORT and GRPC_PORT cannot be the same: %d", c.ServerPort)
	}

	// Validate Redis configuration
	if err := validatePort(c.Redis.Port, "REDIS_PORT"); err != nil {
		return err
	}
	if c.Redis.DB < 0 || c.Redis.DB > 15 {
		return fmt.Errorf("REDIS_DB must be between 0-15, got %d", c.Redis.DB)
	}

	// Validate Kafka configuration
	if c.Kafka.Brokers == "" {
		return fmt.Errorf("KAFKA_BROKERS is required")
	}
	brokers := c.GetKafkaBrokers()
	if len(brokers) == 0 {
		return fmt.Errorf("KAFKA_BROKERS must contain at least one valid broker")
	}

	// Validate JWT configuration
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if c.JWT.Secret == "change-me-in-production" && c.Environment == "production" {
		return fmt.Errorf("JWT_SECRET must be changed in production environment")
	}
	if c.JWT.AccessTokenExpire <= 0 {
		return fmt.Errorf("JWT_ACCESS_TOKEN_EXPIRE must be positive")
	}
	if c.JWT.RefreshTokenExpire <= 0 {
		return fmt.Errorf("JWT_REFRESH_TOKEN_EXPIRE must be positive")
	}
	if c.JWT.RefreshTokenExpire <= c.JWT.AccessTokenExpire {
		return fmt.Errorf("JWT_REFRESH_TOKEN_EXPIRE (%v) must be greater than JWT_ACCESS_TOKEN_EXPIRE (%v)",
			c.JWT.RefreshTokenExpire, c.JWT.AccessTokenExpire)
	}

	// Validate environment
	validEnvs := map[string]bool{
		"development": true,
		"staging":     true,
		"production":  true,
		"test":        true,
	}
	if !validEnvs[c.Environment] {
		return fmt.Errorf("ENVIRONMENT must be one of: development, staging, production, test; got %s", c.Environment)
	}

	return nil
}

// validatePort checks if port is in valid range (1-65535)
func validatePort(port int, name string) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("%s must be between 1-65535, got %d", name, port)
	}
	return nil
}

// GetDatabaseDSN returns PostgreSQL connection string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// GetRedisAddr returns Redis address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

// GetKafkaBrokers returns Kafka brokers as a slice
func (c *Config) GetKafkaBrokers() []string {
	parts := strings.Split(c.Kafka.Brokers, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// IsDevelopment returns true if environment is development
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if environment is production
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsStaging returns true if environment is staging
func (c *Config) IsStaging() bool {
	return c.Environment == "staging"
}

// IsTest returns true if environment is test
func (c *Config) IsTest() bool {
	return c.Environment == "test"
}
