package user

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository untuk operasi database PlatformUser.
type Repository struct {
	db *gorm.DB
}

// NewRepository membuat Repository baru.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// DB mengembalikan koneksi GORM untuk query manual.
func (r *Repository) DB() *gorm.DB {
	return r.db
}

// FindByEmail mencari user berdasarkan email.
func (r *Repository) FindByEmail(email string) (*PlatformUser, error) {
	var u PlatformUser
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &u, nil
}

// FindByID mencari user berdasarkan ID.
func (r *Repository) FindByID(id uuid.UUID) (*PlatformUser, error) {
	var u PlatformUser
	if err := r.db.Where("id = ?", id).First(&u).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &u, nil
}

// FindAll mengembalikan semua platform user dengan pagination.
// Gunakan query chain terpisah untuk Count dan Find.
func (r *Repository) FindAll(page, perPage int) ([]PlatformUser, int64, error) {
	var users []PlatformUser
	var total int64

	// Count total — chain terpisah
	if err := r.db.Model(&PlatformUser{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get paginated data — chain terpisah
	offset := (page - 1) * perPage
	if err := r.db.Model(&PlatformUser{}).
		Offset(offset).
		Limit(perPage).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}

// Create menyimpan user baru ke database.
func (r *Repository) Create(user *PlatformUser) error {
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// Update mengupdate user.
func (r *Repository) Update(user *PlatformUser) error {
	if err := r.db.Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// UpdateLastLogin mengupdate waktu last login user.
func (r *Repository) UpdateLastLogin(id uuid.UUID) error {
	if err := r.db.Model(&PlatformUser{}).Where("id = ?", id).Update("last_login_at", gorm.Expr("NOW()")).Error; err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	return nil
}

// SoftDelete melakukan soft delete user.
func (r *Repository) SoftDelete(id uuid.UUID) error {
	result := r.db.Where("id = ?", id).Delete(&PlatformUser{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
