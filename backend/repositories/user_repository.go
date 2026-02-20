package repositories

import (
	"pemdes-payroll/backend/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByUsername(username string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByIDWithKaryawan(id uint) (*models.User, error)
	Update(id uint, user *models.User) error
	Delete(id uint) error
	GetAll() ([]models.User, error)
	ToggleActive(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new User repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ? AND is_active = ?", username, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(id uint, user *models.User) error {
	var existing models.User
	if err := r.db.First(&existing, id).Error; err != nil {
		return err
	}
	return r.db.Model(&existing).Updates(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("Karyawan").Find(&users).Error
	return users, err
}

func (r *userRepository) GetByIDWithKaryawan(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Karyawan").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) ToggleActive(id uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", id).
		Update("is_active", gorm.Expr("NOT is_active")).Error
}
