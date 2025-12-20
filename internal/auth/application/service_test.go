package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/Leon180/tabelogo-v2/internal/auth/application"
	"github.com/Leon180/tabelogo-v2/internal/auth/domain/errors"
	"github.com/Leon180/tabelogo-v2/internal/auth/domain/model"
	"github.com/Leon180/tabelogo-v2/pkg/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Mock TokenRepository
type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockTokenRepository) GetByHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.RefreshToken), args.Error(1)
}

func (m *MockTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTokenRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// Mock SessionRepository
type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) Create(ctx context.Context, session *model.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockSessionRepository) GetByID(ctx context.Context, sessionID uuid.UUID) (*model.Session, error) {
	args := m.Called(ctx, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Session), args.Error(1)
}

func (m *MockSessionRepository) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]*model.Session, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Session), args.Error(1)
}

func (m *MockSessionRepository) CountUserSessions(ctx context.Context, userID uuid.UUID) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockSessionRepository) UpdateActivity(ctx context.Context, sessionID uuid.UUID) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockSessionRepository) Revoke(ctx context.Context, sessionID uuid.UUID) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockSessionRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockSessionRepository) DeleteExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func setupTestService(t *testing.T) (application.AuthService, *MockUserRepository, *MockTokenRepository, *MockSessionRepository) {
	userRepo := new(MockUserRepository)
	tokenRepo := new(MockTokenRepository)
	sessionRepo := new(MockSessionRepository)
	jwtMaker, err := jwt.NewJWTMaker("test-secret-key-must-be-at-least-32-characters-long")
	require.NoError(t, err)

	service := application.NewAuthService(userRepo, tokenRepo, sessionRepo, jwtMaker)
	return service, userRepo, tokenRepo, sessionRepo
}

func TestAuthService_Register(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		service, userRepo, _, _ := setupTestService(t)
		ctx := context.Background()

		// Mock: user does not exist
		userRepo.On("GetByEmail", ctx, "test@example.com").Return(nil, errors.ErrUserNotFound)
		// Mock: create user succeeds
		userRepo.On("Create", ctx, mock.AnythingOfType("*model.User")).Return(nil)

		user, err := service.Register(ctx, "test@example.com", "password123", "testuser")

		require.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "test@example.com", user.Email())
		assert.Equal(t, "testuser", user.Username())
		userRepo.AssertExpectations(t)
	})

	t.Run("email already exists", func(t *testing.T) {
		service, userRepo, _, _ := setupTestService(t)
		ctx := context.Background()

		existingUser, _ := model.NewUser("test@example.com", "password", "existing")
		userRepo.On("GetByEmail", ctx, "test@example.com").Return(existingUser, nil)

		user, err := service.Register(ctx, "test@example.com", "password123", "testuser")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, errors.ErrEmailAlreadyExists, err)
		userRepo.AssertExpectations(t)
	})
}

func TestAuthService_Login(t *testing.T) {
	t.Run("successful login", func(t *testing.T) {
		service, userRepo, tokenRepo, sessionRepo := setupTestService(t)
		ctx := context.Background()

		// Create a user with known password
		user, err := model.NewUser("test@example.com", "password123", "testuser")
		require.NoError(t, err)

		userRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)
		sessionRepo.On("CountUserSessions", ctx, user.ID()).Return(0, nil)
		sessionRepo.On("Create", ctx, mock.AnythingOfType("*model.Session")).Return(nil)
		tokenRepo.On("Create", ctx, mock.AnythingOfType("*model.RefreshToken")).Return(nil)

		accessToken, refreshToken, err := service.Login(ctx, "test@example.com", "password123", "test-device", "127.0.0.1", false)

		require.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)
		userRepo.AssertExpectations(t)
		sessionRepo.AssertExpectations(t)
		tokenRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		service, userRepo, _, _ := setupTestService(t)
		ctx := context.Background()

		userRepo.On("GetByEmail", ctx, "nonexistent@example.com").Return(nil, errors.ErrUserNotFound)

		accessToken, refreshToken, err := service.Login(ctx, "nonexistent@example.com", "password123", "test-device", "127.0.0.1", false)

		assert.Error(t, err)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
		assert.Equal(t, errors.ErrUserNotFound, err)
		userRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		service, userRepo, _, _ := setupTestService(t)
		ctx := context.Background()

		user, err := model.NewUser("test@example.com", "correctpassword", "testuser")
		require.NoError(t, err)

		userRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)

		accessToken, refreshToken, err := service.Login(ctx, "test@example.com", "wrongpassword", "test-device", "127.0.0.1", false)

		assert.Error(t, err)
		assert.Empty(t, accessToken)
		assert.Empty(t, refreshToken)
		assert.Equal(t, errors.ErrInvalidPassword, err)
		userRepo.AssertExpectations(t)
	})
}

func TestAuthService_ValidateToken(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		service, userRepo, _, sessionRepo := setupTestService(t)
		ctx := context.Background()

		// Create a user and generate a token
		user, err := model.NewUser("test@example.com", "password123", "testuser")
		require.NoError(t, err)

		jwtMaker, _ := jwt.NewJWTMaker("test-secret-key-must-be-at-least-32-characters-long")
		testSessionID := uuid.New()
		token, _, err := jwtMaker.CreateToken(user.ID(), testSessionID, string(user.Role()), 15*time.Minute)
		require.NoError(t, err)

		// Mock session as valid
		testSession := model.NewSession(user.ID(), "test-device", "127.0.0.1")
		sessionRepo.On("GetByID", ctx, testSessionID).Return(testSession, nil)
		userRepo.On("GetByID", ctx, user.ID()).Return(user, nil)

		validatedUser, err := service.ValidateToken(ctx, token)

		require.NoError(t, err)
		assert.NotNil(t, validatedUser)
		assert.Equal(t, user.ID(), validatedUser.ID())
		userRepo.AssertExpectations(t)
		sessionRepo.AssertExpectations(t)
	})

	t.Run("invalid token", func(t *testing.T) {
		service, _, _, _ := setupTestService(t)
		ctx := context.Background()

		validatedUser, err := service.ValidateToken(ctx, "invalid-token")

		assert.Error(t, err)
		assert.Nil(t, validatedUser)
	})
}

func TestAuthService_RefreshToken(t *testing.T) {
	t.Run("successful token refresh", func(t *testing.T) {
		service, userRepo, tokenRepo, _ := setupTestService(t)
		ctx := context.Background()

		// Create a user
		user, err := model.NewUser("test@example.com", "password123", "testuser")
		require.NoError(t, err)

		// Create a refresh token
		jwtMaker, _ := jwt.NewJWTMaker("test-secret-key-must-be-at-least-32-characters-long")
		testSessionID := uuid.New()
		oldRefreshToken, payload, err := jwtMaker.CreateToken(user.ID(), testSessionID, string(user.Role()), 24*time.Hour)
		require.NoError(t, err)

		refreshTokenEntity := model.NewRefreshToken(user.ID(), oldRefreshToken, payload.ExpiredAt)

		tokenRepo.On("GetByHash", ctx, oldRefreshToken).Return(refreshTokenEntity, nil)
		userRepo.On("GetByID", ctx, user.ID()).Return(user, nil)
		tokenRepo.On("Create", ctx, mock.AnythingOfType("*model.RefreshToken")).Return(nil)

		newAccessToken, newRefreshToken, err := service.RefreshToken(ctx, oldRefreshToken)

		require.NoError(t, err)
		assert.NotEmpty(t, newAccessToken)
		assert.NotEmpty(t, newRefreshToken)
		assert.NotEqual(t, oldRefreshToken, newRefreshToken)
		userRepo.AssertExpectations(t)
		tokenRepo.AssertExpectations(t)
	})
}
