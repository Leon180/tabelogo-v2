package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/auth/application"
	authpostgres "github.com/Leon180/tabelogo-v2/internal/auth/infrastructure/postgres"
	authredis "github.com/Leon180/tabelogo-v2/internal/auth/infrastructure/redis"
	"github.com/Leon180/tabelogo-v2/pkg/jwt"
	redisclient "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AuthIntegrationTestSuite struct {
	suite.Suite
	db      *gorm.DB
	redis   *redisclient.Client
	service application.AuthService
}

func (suite *AuthIntegrationTestSuite) SetupSuite() {
	// Setup PostgreSQL connection
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		getEnv("TEST_DB_HOST", "localhost"),
		getEnvInt("TEST_DB_PORT", 5433),
		getEnv("TEST_DB_USER", "postgres"),
		getEnv("TEST_DB_PASSWORD", "postgres"),
		getEnv("TEST_DB_NAME", "auth_test"),
	)

	var err error
	suite.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(suite.T(), err, "Failed to connect to test database")

	// Auto-migrate schema
	err = suite.db.AutoMigrate(&authpostgres.UserORM{})
	require.NoError(suite.T(), err, "Failed to migrate database")

	// Setup Redis connection
	suite.redis = redisclient.NewClient(&redisclient.Options{
		Addr:     getEnv("TEST_REDIS_ADDR", "localhost:6380"),
		Password: getEnv("TEST_REDIS_PASSWORD", ""),
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = suite.redis.Ping(ctx).Err()
	require.NoError(suite.T(), err, "Failed to connect to test Redis")

	// Setup service
	userRepo := authpostgres.NewUserRepository(suite.db)
	tokenRepo := authredis.NewTokenRepository(suite.redis)
	jwtMaker, err := jwt.NewJWTMaker("test-secret-key-must-be-at-least-32-characters-long")
	require.NoError(suite.T(), err)

	suite.service = application.NewAuthService(userRepo, tokenRepo, jwtMaker)
}

func (suite *AuthIntegrationTestSuite) TearDownSuite() {
	// Clean up
	if suite.db != nil {
		sqlDB, _ := suite.db.DB()
		sqlDB.Close()
	}
	if suite.redis != nil {
		suite.redis.Close()
	}
}

func (suite *AuthIntegrationTestSuite) SetupTest() {
	// Clean database before each test
	suite.db.Exec("TRUNCATE TABLE users CASCADE")
	// Clean Redis
	suite.redis.FlushDB(context.Background())
}

func (suite *AuthIntegrationTestSuite) TestRegisterAndLogin() {
	ctx := context.Background()

	// Register a new user
	user, err := suite.service.Register(ctx, "integration@example.com", "password123", "integrationuser")
	require.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), "integration@example.com", user.Email())

	// Login with the registered user
	accessToken, refreshToken, err := suite.service.Login(ctx, "integration@example.com", "password123")
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), accessToken)
	assert.NotEmpty(suite.T(), refreshToken)

	// Validate the access token
	validatedUser, err := suite.service.ValidateToken(ctx, accessToken)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.ID(), validatedUser.ID())
}

func (suite *AuthIntegrationTestSuite) TestTokenRefresh() {
	ctx := context.Background()

	// Register and login
	_, err := suite.service.Register(ctx, "refresh@example.com", "password123", "refreshuser")
	require.NoError(suite.T(), err)

	_, oldRefreshToken, err := suite.service.Login(ctx, "refresh@example.com", "password123")
	require.NoError(suite.T(), err)

	// Wait a bit to ensure new tokens have different timestamps
	time.Sleep(1 * time.Second)

	// Refresh the token
	newAccessToken, newRefreshToken, err := suite.service.RefreshToken(ctx, oldRefreshToken)
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), newAccessToken)
	assert.NotEmpty(suite.T(), newRefreshToken)
	assert.NotEqual(suite.T(), oldRefreshToken, newRefreshToken)

	// Validate the new access token
	_, err = suite.service.ValidateToken(ctx, newAccessToken)
	require.NoError(suite.T(), err)
}

func (suite *AuthIntegrationTestSuite) TestDuplicateRegistration() {
	ctx := context.Background()

	// Register first user
	_, err := suite.service.Register(ctx, "duplicate@example.com", "password123", "user1")
	require.NoError(suite.T(), err)

	// Try to register with same email
	_, err = suite.service.Register(ctx, "duplicate@example.com", "password456", "user2")
	assert.Error(suite.T(), err)
}

func (suite *AuthIntegrationTestSuite) TestInvalidCredentials() {
	ctx := context.Background()

	// Register a user
	_, err := suite.service.Register(ctx, "valid@example.com", "correctpassword", "validuser")
	require.NoError(suite.T(), err)

	// Try to login with wrong password
	_, _, err = suite.service.Login(ctx, "valid@example.com", "wrongpassword")
	assert.Error(suite.T(), err)

	// Try to login with non-existent email
	_, _, err = suite.service.Login(ctx, "nonexistent@example.com", "password123")
	assert.Error(suite.T(), err)
}

func TestAuthIntegrationSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	suite.Run(t, new(AuthIntegrationTestSuite))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		fmt.Sscanf(value, "%d", &intValue)
		return intValue
	}
	return defaultValue
}
