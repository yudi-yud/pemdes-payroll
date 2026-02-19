package repositories

import (
	"pemdes-payroll/backend/models"

	"gorm.io/gorm"
)

type KaryawanRepository interface {
	Create(karyawan *models.Karyawan) error
	GetAll() ([]models.Karyawan, error)
	GetByID(id uint) (*models.Karyawan, error)
	GetByIDWithJabatan(id uint) (*models.Karyawan, error)
	Update(id uint, karyawan *models.Karyawan) error
	Delete(id uint) error
	GetByStatus(status models.KaryawanStatus) ([]models.Karyawan, error)
	Search(keyword string) ([]models.Karyawan, error)
	Count() (int64, error)
}

type karyawanRepository struct {
	db *gorm.DB
}

// NewKaryawanRepository creates a new Karyawan repository
func NewKaryawanRepository(db *gorm.DB) KaryawanRepository {
	return &karyawanRepository{db: db}
}

func (r *karyawanRepository) Create(karyawan *models.Karyawan) error {
	return r.db.Create(karyawan).Error
}

func (r *karyawanRepository) GetAll() ([]models.Karyawan, error) {
	var karyawan []models.Karyawan
	err := r.db.Preload("Jabatan").Order("created_at DESC").Find(&karyawan).Error
	return karyawan, err
}

func (r *karyawanRepository) GetByID(id uint) (*models.Karyawan, error) {
	var karyawan models.Karyawan
	err := r.db.First(&karyawan, id).Error
	if err != nil {
		return nil, err
	}
	return &karyawan, nil
}

func (r *karyawanRepository) GetByIDWithJabatan(id uint) (*models.Karyawan, error) {
	var karyawan models.Karyawan
	err := r.db.Preload("Jabatan").First(&karyawan, id).Error
	if err != nil {
		return nil, err
	}
	return &karyawan, nil
}

func (r *karyawanRepository) Update(id uint, karyawan *models.Karyawan) error {
	var existing models.Karyawan
	if err := r.db.First(&existing, id).Error; err != nil {
		return err
	}
	return r.db.Model(&existing).Updates(karyawan).Error
}

func (r *karyawanRepository) Delete(id uint) error {
	return r.db.Delete(&models.Karyawan{}, id).Error
}

func (r *karyawanRepository) GetByStatus(status models.KaryawanStatus) ([]models.Karyawan, error) {
	var karyawan []models.Karyawan
	err := r.db.Preload("Jabatan").Where("status = ?", status).Find(&karyawan).Error
	return karyawan, err
}

func (r *karyawanRepository) Search(keyword string) ([]models.Karyawan, error) {
	var karyawan []models.Karyawan
	searchPattern := "%" + keyword + "%"
	err := r.db.Preload("Jabatan").Where(
		"nik LIKE ? OR nama LIKE ? OR email LIKE ?",
		searchPattern, searchPattern, searchPattern,
	).Find(&karyawan).Error
	return karyawan, err
}

func (r *karyawanRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Karyawan{}).Count(&count).Error
	return count, err
}
