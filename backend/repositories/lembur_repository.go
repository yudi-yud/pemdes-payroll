package repositories

import (
	"pemdes-payroll/backend/models"

	"gorm.io/gorm"
)

type LemburRepository interface {
	Create(lembur *models.Lembur) error
	GetAll() ([]models.Lembur, error)
	GetByID(id uint) (*models.Lembur, error)
	GetByKaryawanID(karyawanID uint) ([]models.Lembur, error)
	GetByPeriod(bulan, tahun int) ([]models.Lembur, error)
	GetByKaryawanAndPeriod(karyawanID, bulan, tahun int) ([]models.Lembur, error)
	Update(id uint, lembur *models.Lembur) error
	Delete(id uint) error
	Approve(id uint, approverID *uint, status string) error
	GetTotalLemburByPeriod(karyawanID, bulan, tahun int) (float64, float64, error)
}

type lemburRepository struct {
	db *gorm.DB
}

// NewLemburRepository creates a new Lembur repository
func NewLemburRepository(db *gorm.DB) LemburRepository {
	return &lemburRepository{db: db}
}

func (r *lemburRepository) Create(lembur *models.Lembur) error {
	return r.db.Create(lembur).Error
}

func (r *lemburRepository) GetAll() ([]models.Lembur, error) {
	var lembur []models.Lembur
	err := r.db.Preload("Karyawan.Jabatan").Order("tanggal DESC, created_at DESC").Find(&lembur).Error
	return lembur, err
}

func (r *lemburRepository) GetByID(id uint) (*models.Lembur, error) {
	var lembur models.Lembur
	err := r.db.Preload("Karyawan.Jabatan").First(&lembur, id).Error
	if err != nil {
		return nil, err
	}
	return &lembur, nil
}

func (r *lemburRepository) GetByKaryawanID(karyawanID uint) ([]models.Lembur, error) {
	var lembur []models.Lembur
	err := r.db.Preload("Karyawan.Jabatan").Where("karyawan_id = ?", karyawanID).
		Order("tanggal DESC").Find(&lembur).Error
	return lembur, err
}

func (r *lemburRepository) GetByPeriod(bulan, tahun int) ([]models.Lembur, error) {
	var lembur []models.Lembur
	err := r.db.Preload("Karyawan.Jabatan").Where("MONTH(tanggal) = ? AND YEAR(tanggal) = ?", bulan, tahun).
		Order("tanggal DESC, karyawan_id").Find(&lembur).Error
	return lembur, err
}

func (r *lemburRepository) GetByKaryawanAndPeriod(karyawanID, bulan, tahun int) ([]models.Lembur, error) {
	var lembur []models.Lembur
	err := r.db.Preload("Karyawan.Jabatan").Where("karyawan_id = ? AND MONTH(tanggal) = ? AND YEAR(tanggal) = ?",
		karyawanID, bulan, tahun).
		Order("tanggal DESC").Find(&lembur).Error
	return lembur, err
}

func (r *lemburRepository) Update(id uint, lembur *models.Lembur) error {
	var existing models.Lembur
	if err := r.db.First(&existing, id).Error; err != nil {
		return err
	}
	return r.db.Model(&existing).Updates(lembur).Error
}

func (r *lemburRepository) Delete(id uint) error {
	return r.db.Delete(&models.Lembur{}, id).Error
}

// Approve updates overtime approval status
func (r *lemburRepository) Approve(id uint, approverID *uint, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if approverID != nil {
		updates["disetujui_oleh"] = *approverID
	} else {
		updates["disetujui_oleh"] = nil
	}
	return r.db.Model(&models.Lembur{}).Where("id = ?", id).Updates(updates).Error
}

// GetTotalLemburByPeriod calculates total hours and nominal for approved overtime in a period
func (r *lemburRepository) GetTotalLemburByPeriod(karyawanID, bulan, tahun int) (float64, float64, error) {
	type Result struct {
		TotalJam     float64
		TotalNominal float64
	}

	var result Result
	err := r.db.Model(&models.Lembur{}).
		Select("COALESCE(SUM(total_jam), 0) as total_jam, COALESCE(SUM(total_nominal), 0) as total_nominal").
		Where("karyawan_id = ? AND MONTH(tanggal) = ? AND YEAR(tanggal) = ? AND status = 'disetujui'",
			karyawanID, bulan, tahun).
		Scan(&result).Error

	return result.TotalJam, result.TotalNominal, err
}
