package postgres

import (
	"context"
	"errors"
	"time"

	domainerrors "github.com/Leon180/tabelogo-v2/internal/auth/domain/errors"
	"github.com/Leon180/tabelogo-v2/internal/auth/domain/model"
	"github.com/Leon180/tabelogo-v2/internal/auth/domain/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserORM is the database model for User
type UserORM struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Email         string         `gorm:"uniqueIndex;not null"`
	PasswordHash  string         `gorm:"not null"`
	Username      string         `gorm:"not null"`
	Role          string         `gorm:"not null"`
	IsActive      bool           `gorm:"not null;default:true"`
	EmailVerified bool           `gorm:"not null;default:false"`
	CreatedAt     time.Time      `gorm:"not null"`
	UpdatedAt     time.Time      `gorm:"not null"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// TableName overrides the table name
func (UserORM) TableName() string {
	return "users"
}

// ToDomain converts ORM model to Domain entity
func (u *UserORM) ToDomain() *model.User {
	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		deletedAt = &u.DeletedAt.Time
	}

	return model.ReconstructUser(
		u.ID,
		u.Email,
		u.PasswordHash,
		u.Username,
		model.UserRole(u.Role),
		u.IsActive,
		u.EmailVerified,
		u.CreatedAt,
		u.UpdatedAt,
		deletedAt,
	)
}

// FromDomain converts Domain entity to ORM model
func FromDomain(u *model.User) *UserORM {
	var deletedAt gorm.DeletedAt
	if u.DeletedAt() != nil {
		deletedAt = gorm.DeletedAt{Time: *u.DeletedAt(), Valid: true}
	}

	return &UserORM{
		ID:            u.ID(),
		Email:         u.Email(),
		PasswordHash:  u.PasswordHash(),
		Username:      u.Username(),
		Role:          string(u.Role()),
		IsActive:      u.IsActive(),
		EmailVerified: u.EmailVerified(), // Fixed: Added missing field
		CreatedAt:     u.CreatedAt(),
		UpdatedAt:     time.Now(), // UpdatedAt should be updated on save, but initial mapping is fine
		DeletedAt:     deletedAt,
	}
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new postgres user repository
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	orm := FromDomain(user)
	if err := r.db.WithContext(ctx).Create(orm).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domainerrors.ErrEmailAlreadyExists
		}
		return err
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var orm UserORM
	if err := r.db.WithContext(ctx).First(&orm, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrUserNotFound
		}
		return nil, err
	}
	return orm.ToDomain(), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var orm UserORM
	if err := r.db.WithContext(ctx).First(&orm, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrUserNotFound
		}
		return nil, err
	}
	return orm.ToDomain(), nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	orm := FromDomain(user)
	// Use Save to update all fields including zero values if necessary,
	// but for partial updates usually Updates is better.
	// Since this is a full aggregate update, Save is appropriate provided ID is set.
	if err := r.db.WithContext(ctx).Save(orm).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&UserORM{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
