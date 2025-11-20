package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Set minimal required environment variables
	os.Setenv("DB_NAME", "test_db")
	defer os.Unsetenv("DB_NAME")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	// Check defaults
	if cfg.Environment != "development" {
		t.Errorf("Expected environment 'development', got %s", cfg.Environment)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("Expected log level 'info', got %s", cfg.LogLevel)
	}
	if cfg.ServerPort != 8080 {
		t.Errorf("Expected server port 8080, got %d", cfg.ServerPort)
	}
}

func TestLoadWithPrefix(t *testing.T) {
	prefix := "AUTH"
	os.Setenv("AUTH_DB_NAME", "auth_db")
	os.Setenv("AUTH_DB_PORT", "15432")
	defer func() {
		os.Unsetenv("AUTH_DB_NAME")
		os.Unsetenv("AUTH_DB_PORT")
	}()

	cfg, err := LoadWithPrefix(prefix)
	if err != nil {
		t.Fatalf("LoadWithPrefix() error = %v", err)
	}

	if cfg.Database.Name != "auth_db" {
		t.Errorf("Expected DB name 'auth_db', got %s", cfg.Database.Name)
	}
	if cfg.Database.Port != 15432 {
		t.Errorf("Expected DB port 15432, got %d", cfg.Database.Port)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: &Config{
				Environment: "development",
				ServerPort:  8080,
				GRPCPort:    9090,
				Database: DatabaseConfig{
					Name:            "test_db",
					Port:            5432,
					MaxOpenConns:    100,
					MaxIdleConns:    10,
					ConnMaxLifetime: time.Hour,
				},
				Redis: RedisConfig{
					Port: 6379,
					DB:   0,
				},
				Kafka: KafkaConfig{
					Brokers: "localhost:9092",
				},
				JWT: JWTConfig{
					Secret:             "test-secret",
					AccessTokenExpire:  15 * time.Minute,
					RefreshTokenExpire: 7 * 24 * time.Hour,
				},
			},
			wantErr: false,
		},
		{
			name: "missing DB_NAME",
			config: &Config{
				Environment: "development",
				ServerPort:  8080,
				GRPCPort:    9090,
				Database: DatabaseConfig{
					Name: "",
				},
			},
			wantErr: true,
			errMsg:  "DB_NAME is required",
		},
		{
			name: "invalid port",
			config: &Config{
				Environment: "development",
				ServerPort:  -1,
				Database: DatabaseConfig{
					Name: "test_db",
				},
			},
			wantErr: true,
		},
		{
			name: "same ports",
			config: &Config{
				Environment: "development",
				ServerPort:  8080,
				GRPCPort:    8080,
				Database: DatabaseConfig{
					Name:            "test_db",
					Port:            5432,
					MaxOpenConns:    100,
					MaxIdleConns:    10,
					ConnMaxLifetime: time.Hour,
				},
				Redis: RedisConfig{
					Port: 6379,
					DB:   0,
				},
				Kafka: KafkaConfig{
					Brokers: "localhost:9092",
				},
				JWT: JWTConfig{
					Secret:             "test-secret",
					AccessTokenExpire:  15 * time.Minute,
					RefreshTokenExpire: 7 * 24 * time.Hour,
				},
			},
			wantErr: true,
			errMsg:  "SERVER_PORT and GRPC_PORT cannot be the same",
		},
		{
			name: "production with default JWT secret",
			config: &Config{
				Environment: "production",
				ServerPort:  8080,
				GRPCPort:    9090,
				Database: DatabaseConfig{
					Name:            "test_db",
					Port:            5432,
					MaxOpenConns:    100,
					MaxIdleConns:    10,
					ConnMaxLifetime: time.Hour,
				},
				Redis: RedisConfig{
					Port: 6379,
					DB:   0,
				},
				Kafka: KafkaConfig{
					Brokers: "localhost:9092",
				},
				JWT: JWTConfig{
					Secret:             "change-me-in-production",
					AccessTokenExpire:  15 * time.Minute,
					RefreshTokenExpire: 7 * 24 * time.Hour,
				},
			},
			wantErr: true,
			errMsg:  "JWT_SECRET must be changed in production",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestGetDatabaseDSN(t *testing.T) {
	cfg := &Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			Name:     "test_db",
			User:     "postgres",
			Password: "postgres",
			SSLMode:  "disable",
		},
	}

	expected := "host=localhost port=5432 user=postgres password=postgres dbname=test_db sslmode=disable"
	got := cfg.GetDatabaseDSN()

	if got != expected {
		t.Errorf("GetDatabaseDSN() = %v, want %v", got, expected)
	}
}

func TestGetRedisAddr(t *testing.T) {
	cfg := &Config{
		Redis: RedisConfig{
			Host: "localhost",
			Port: 6379,
		},
	}

	expected := "localhost:6379"
	got := cfg.GetRedisAddr()

	if got != expected {
		t.Errorf("GetRedisAddr() = %v, want %v", got, expected)
	}
}

func TestGetKafkaBrokers(t *testing.T) {
	tests := []struct {
		name    string
		brokers string
		want    []string
	}{
		{
			name:    "single broker",
			brokers: "localhost:9092",
			want:    []string{"localhost:9092"},
		},
		{
			name:    "multiple brokers",
			brokers: "localhost:9092,localhost:9093,localhost:9094",
			want:    []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		},
		{
			name:    "brokers with spaces",
			brokers: "localhost:9092 , localhost:9093 , localhost:9094",
			want:    []string{"localhost:9092", "localhost:9093", "localhost:9094"},
		},
		{
			name:    "empty string",
			brokers: "",
			want:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Kafka: KafkaConfig{
					Brokers: tt.brokers,
				},
			}

			got := cfg.GetKafkaBrokers()
			if len(got) != len(tt.want) {
				t.Errorf("GetKafkaBrokers() length = %v, want %v", len(got), len(tt.want))
				return
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("GetKafkaBrokers()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestEnvironmentChecks(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		isDev       bool
		isProd      bool
		isStaging   bool
		isTest      bool
	}{
		{
			name:        "development",
			environment: "development",
			isDev:       true,
			isProd:      false,
			isStaging:   false,
			isTest:      false,
		},
		{
			name:        "production",
			environment: "production",
			isDev:       false,
			isProd:      true,
			isStaging:   false,
			isTest:      false,
		},
		{
			name:        "staging",
			environment: "staging",
			isDev:       false,
			isProd:      false,
			isStaging:   true,
			isTest:      false,
		},
		{
			name:        "test",
			environment: "test",
			isDev:       false,
			isProd:      false,
			isStaging:   false,
			isTest:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Environment: tt.environment}

			if cfg.IsDevelopment() != tt.isDev {
				t.Errorf("IsDevelopment() = %v, want %v", cfg.IsDevelopment(), tt.isDev)
			}
			if cfg.IsProduction() != tt.isProd {
				t.Errorf("IsProduction() = %v, want %v", cfg.IsProduction(), tt.isProd)
			}
			if cfg.IsStaging() != tt.isStaging {
				t.Errorf("IsStaging() = %v, want %v", cfg.IsStaging(), tt.isStaging)
			}
			if cfg.IsTest() != tt.isTest {
				t.Errorf("IsTest() = %v, want %v", cfg.IsTest(), tt.isTest)
			}
		})
	}
}
