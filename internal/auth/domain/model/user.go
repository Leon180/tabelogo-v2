package model

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserRole defines the role of a user
type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
	RoleGuest UserRole = "guest"
)

// User is the aggregate root for the authentication domain
type User struct {
	id            uuid.UUID
	email         string
	passwordHash  string
	username      string
	role          UserRole
	isActive      bool
	emailVerified bool
	createdAt     time.Time
	updatedAt     time.Time
	deletedAt     *time.Time
}

// NewUser creates a new user with hashed password
func NewUser(email, password, username string) (*User, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		id:            uuid.New(),
		email:         email,
		passwordHash:  hashedPassword,
		username:      username,
		role:          RoleUser,
		isActive:      true,
		emailVerified: false,
		createdAt:     time.Now(),
		updatedAt:     time.Now(),
	}, nil
}

// ReconstructUser is used by repository to reconstruct the User entity from persistence
// This should NOT be used by application layer to create new users
func ReconstructUser(
	id uuid.UUID,
	email, passwordHash, username string,
	role UserRole,
	isActive, emailVerified bool,
	createdAt, updatedAt time.Time,
	deletedAt *time.Time,
) *User {
	return &User{
		id:            id,
		email:         email,
		passwordHash:  passwordHash,
		username:      username,
		role:          role,
		isActive:      isActive,
		emailVerified: emailVerified,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
		deletedAt:     deletedAt,
	}
}

// Getters
func (u *User) ID() uuid.UUID         { return u.id }
func (u *User) Email() string         { return u.email }
func (u *User) PasswordHash() string  { return u.passwordHash }
func (u *User) Username() string      { return u.username }
func (u *User) Role() UserRole        { return u.role }
func (u *User) IsActive() bool        { return u.isActive }
func (u *User) CreatedAt() time.Time  { return u.createdAt }
func (u *User) EmailVerified() bool   { return u.emailVerified }
func (u *User) DeletedAt() *time.Time { return u.deletedAt }

// Domain Methods (Setters with logic)

// CheckPassword verifies if the provided password matches the hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.passwordHash), []byte(password))
	return err == nil
}

// UpdatePassword updates the user's password
func (u *User) UpdatePassword(newPassword string) error {
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	u.passwordHash = hashedPassword
	u.updatedAt = time.Now()
	return nil
}

// VerifyEmail marks the email as verified
func (u *User) VerifyEmail() {
	u.emailVerified = true
	u.updatedAt = time.Now()
}

// HashPassword hashes a plain text password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
