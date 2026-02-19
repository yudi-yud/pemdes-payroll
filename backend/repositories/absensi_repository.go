package repositories

import (
	"pemdes-payroll/backend/models"
	"time"

	"gorm.io/gorm"
)

type AbsensiRepository interface {
	Create(absensi *models.Absensi) error
	GetAll() ([]models.Absensi, error)
	GetByID(id uint) (*models.Absensi, error)
	GetByKaryawanID(karyawanID uint, startDate, endDate time.Time) ([]models.Absensi, error)
	GetByDateRange(startDate, endDate time.Time) ([]models.Absensi, error)
	Update(id uint, absensi *models.Absensi) error
	Delete(id uint) error
	GetRekapBulanan(karyawanID uint, bulan, tahun int) (map[string]int, error)
}

type absensiRepository struct {
	db *gorm.DB
}

// NewAbsensiRepository creates a new Absensi repository
func NewAbsensiRepository(db *gorm.DB) AbsensiRepository {
	return &absensiRepository{db: db}
}

func (r *absensiRepository) Create(absensi *models.Absensi) error {
	return r.db.Create(absensi).Error
}

func (r *absensiRepository) GetAll() ([]models.Absensi, error) {
	var absensi []models.Absensi
	err := r.db.Preload("Karyawan.Jabatan").Order("tanggal DESC, created_at DESC").Find(&absensi).Error
	return absensi, err
}

func (r *absensiRepository) GetByID(id uint) (*models.Absensi, error) {
	var absensi models.Absensi
	err := r.db.Preload("Karyawan.Jabatan").First(&absensi, id).Error
	if err != nil {
		return nil, err
	}
	return &absensi, nil
}

func (r *absensiRepository) GetByKaryawanID(karyawanID uint, startDate, endDate time.Time) ([]models.Absensi, error) {
	var absensi []models.Absensi
	err := r.db.Where("karyawan_id = ? AND tanggal BETWEEN ? AND ?", karyawanID, startDate, endDate).
		Order("tanggal DESC").Find(&absensi).Error
	return absensi, err
}

func (r *absensiRepository) GetByDateRange(startDate, endDate time.Time) ([]models.Absensi, error) {
	var absensi []models.Absensi
	err := r.db.Preload("Karyawan.Jabatan").Where("tanggal BETWEEN ? AND ?", startDate, endDate).
		Order("tanggal DESC, karyawan_id").Find(&absensi).Error
	return absensi, err
}

func (r *absensiRepository) Update(id uint, absensi *models.Absensi) error {
	var existing models.Absensi
	if err := r.db.First(&existing, id).Error; err != nil {
		return err
	}
	return r.db.Model(&existing).Updates(absensi).Error
}

func (r *absensiRepository) Delete(id uint) error {
	return r.db.Delete(&models.Absensi{}, id).Error
}

// GetRekapBulanan gets attendance summary for a specific month
func (r *absensiRepository) GetRekapBulanan(karyawanID uint, bulan, tahun int) (map[string]int, error) {
	type Result struct {
		Status string
		Count  int
	}

	startDate := time.Date(tahun, time.Month(bulan), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	var results []Result
	err := r.db.Model(&models.Absensi{}).
		Select("status, COUNT(*) as count").
		Where("karyawan_id = ? AND tanggal BETWEEN ? AND ?", karyawanID, startDate, endDate).
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	rekap := map[string]int{
		"hadir":  0,
		"izin":   0,
		"sakit":  0,
		"alpha":  0,
	}

	for _, r := range results {
		rekap[r.Status] = r.Count
	}

	return rekap, nil
}
