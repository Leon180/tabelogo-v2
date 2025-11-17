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
	cfg := &Config{
		Environment: getEnvWithDefault("ENVIRONMENT", "development"),
		LogLevel:    getEnvWithDefault("LOG_LEVEL", "info"),
		ServerPort:  getEnvAsInt("SERVER_PORT", 8080),
		GRPCPort:    getEnvAsInt("GRPC_PORT", 9090),
	}

	// Load database config
	cfg.Database = DatabaseConfig{
		Host:            getEnvWithDefault("DB_HOST", "localhost"),
		Port:            getEnvAsInt("DB_PORT", 5432),
		Name:            getEnvWithDefault("DB_NAME", ""),
		User:            getEnvWithDefault("DB_USER", "postgres"),
		Password:        getEnvWithDefault("DB_PASSWORD", "postgres"),
		SSLMode:         getEnvWithDefault("DB_SSLMODE", "disable"),
		MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
		MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
		ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", time.Hour),
	}

	// Load Redis config
	cfg.Redis = RedisConfig{
		Host:     getEnvWithDefault("REDIS_HOST", "localhost"),
		Port:     getEnvAsInt("REDIS_PORT", 6379),
		Password: getEnvWithDefault("REDIS_PASSWORD", ""),
		DB:       getEnvAsInt("REDIS_DB", 0),
	}

	// Load Kafka config
	cfg.Kafka = KafkaConfig{
		Brokers: getEnvWithDefault("KAFKA_BROKERS", "localhost:9092"),
		GroupID: getEnvWithDefault("KAFKA_GROUP_ID", "tabelogo-group"),
	}

	// Load JWT config
	cfg.JWT = JWTConfig{
		Secret:             getEnvWithDefault("JWT_SECRET", "change-me-in-production"),
		AccessTokenExpire:  getEnvAsDuration("JWT_ACCESS_TOKEN_EXPIRE", 15*time.Minute),
		RefreshTokenExpire: getEnvAsDuration("JWT_REFRESH_TOKEN_EXPIRE", 7*24*time.Hour),
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	if c.JWT.Secret == "change-me-in-production" && c.Environment == "production" {
		return fmt.Errorf("JWT_SECRET must be set in production")
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
	return getEnvAsSlice("", strings.Split(c.Kafka.Brokers, ","))
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
