package repositories

import (
	"pemdes-payroll/backend/models"

	"gorm.io/gorm"
)

type GajiRepository interface {
	Create(gaji *models.Gaji) error
	CreateBatch(gajiList []models.Gaji) error
	GetAll() ([]models.Gaji, error)
	GetByID(id uint) (*models.Gaji, error)
	GetByKaryawanID(karyawanID uint) ([]models.Gaji, error)
	GetByPeriod(bulan, tahun int) ([]models.Gaji, error)
	GetByKaryawanAndPeriod(karyawanID, bulan, tahun int) (*models.Gaji, error)
	Update(id uint, gaji *models.Gaji) error
	Delete(id uint) error
	UpdateStatus(id uint, status models.GajiStatus) error
	GetTotalGajiByPeriod(bulan, tahun int) (float64, error)
}

type gajiRepository struct {
	db *gorm.DB
}

// NewGajiRepository creates a new Gaji repository
func NewGajiRepository(db *gorm.DB) GajiRepository {
	return &gajiRepository{db: db}
}

func (r *gajiRepository) Create(gaji *models.Gaji) error {
	return r.db.Create(gaji).Error
}

func (r *gajiRepository) CreateBatch(gajiList []models.Gaji) error {
	return r.db.Create(&gajiList).Error
}

func (r *gajiRepository) GetAll() ([]models.Gaji, error) {
	var gaji []models.Gaji
	err := r.db.Preload("Karyawan.Jabatan").Order("periode_tahun DESC, periode_bulan DESC, created_at DESC").Find(&gaji).Error
	return gaji, err
}

func (r *gajiRepository) GetByID(id uint) (*models.Gaji, error) {
	var gaji models.Gaji
	err := r.db.Preload("Karyawan.Jabatan").First(&gaji, id).Error
	if err != nil {
		return nil, err
	}
	return &gaji, nil
}

func (r *gajiRepository) GetByKaryawanID(karyawanID uint) ([]models.Gaji, error) {
	var gaji []models.Gaji
	err := r.db.Preload("Karyawan.Jabatan").Where("karyawan_id = ?", karyawanID).
		Order("periode_tahun DESC, periode_bulan DESC").Find(&gaji).Error
	return gaji, err
}

func (r *gajiRepository) GetByPeriod(bulan, tahun int) ([]models.Gaji, error) {
	var gaji []models.Gaji
	err := r.db.Preload("Karyawan.Jabatan").Where("periode_bulan = ? AND periode_tahun = ?", bulan, tahun).
		Order("created_at DESC").Find(&gaji).Error
	return gaji, err
}

func (r *gajiRepository) GetByKaryawanAndPeriod(karyawanID, bulan, tahun int) (*models.Gaji, error) {
	var gaji models.Gaji
	err := r.db.Where("karyawan_id = ? AND periode_bulan = ? AND periode_tahun = ?",
		karyawanID, bulan, tahun).First(&gaji).Error
	if err != nil {
		return nil, err
	}
	return &gaji, nil
}

func (r *gajiRepository) Update(id uint, gaji *models.Gaji) error {
	var existing models.Gaji
	if err := r.db.First(&existing, id).Error; err != nil {
		return err
	}
	return r.db.Model(&existing).Updates(gaji).Error
}

func (r *gajiRepository) Delete(id uint) error {
	return r.db.Delete(&models.Gaji{}, id).Error
}

func (r *gajiRepository) UpdateStatus(id uint, status models.GajiStatus) error {
	return r.db.Model(&models.Gaji{}).Where("id = ?", id).Update("status", status).Error
}

func (r *gajiRepository) GetTotalGajiByPeriod(bulan, tahun int) (float64, error) {
	var total float64
	err := r.db.Model(&models.Gaji{}).Where("periode_bulan = ? AND periode_tahun = ?", bulan, tahun).
		Select("COALESCE(SUM(total_gaji), 0)").Scan(&total).Error
	return total, err
}
