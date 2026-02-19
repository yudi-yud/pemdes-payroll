package repositories

import (
	"pemdes-payroll/backend/models"

	"gorm.io/gorm"
)

type JabatanRepository interface {
	Create(jabatan *models.Jabatan) error
	GetAll() ([]models.Jabatan, error)
	GetByID(id uint) (*models.Jabatan, error)
	Update(id uint, jabatan *models.Jabatan) error
	Delete(id uint) error
	Count() (int64, error)
}

type jabatanRepository struct {
	db *gorm.DB
}

// NewJabatanRepository creates a new Jabatan repository
func NewJabatanRepository(db *gorm.DB) JabatanRepository {
	return &jabatanRepository{db: db}
}

func (r *jabatanRepository) Create(jabatan *models.Jabatan) error {
	return r.db.Create(jabatan).Error
}

func (r *jabatanRepository) GetAll() ([]models.Jabatan, error) {
	var jabatans []models.Jabatan
	err := r.db.Order("created_at DESC").Find(&jabatans).Error
	return jabatans, err
}

func (r *jabatanRepository) GetByID(id uint) (*models.Jabatan, error) {
	var jabatan models.Jabatan
	err := r.db.First(&jabatan, id).Error
	if err != nil {
		return nil, err
	}
	return &jabatan, nil
}

func (r *jabatanRepository) Update(id uint, jabatan *models.Jabatan) error {
	var existing models.Jabatan
	if err := r.db.First(&existing, id).Error; err != nil {
		return err
	}
	return r.db.Model(&existing).Updates(jabatan).Error
}

func (r *jabatanRepository) Delete(id uint) error {
	return r.db.Delete(&models.Jabatan{}, id).Error
}

func (r *jabatanRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Jabatan{}).Count(&count).Error
	return count, err
}
