package repositories

import (
	"pemdes-payroll/backend/models"

	"gorm.io/gorm"
)

type LaporanRepository interface {
	GetLaporanGajiByPeriod(bulan, tahun int) ([]models.LaporanGaji, error)
	GetRiwayatGajiKaryawan(karyawanID uint) ([]models.LaporanGaji, error)
	GetRekapGaji(bulan, tahun int) (*models.RekapGaji, error)
}

type laporanRepository struct {
	db *gorm.DB
}

// NewLaporanRepository creates a new Laporan repository
func NewLaporanRepository(db *gorm.DB) LaporanRepository {
	return &laporanRepository{db: db}
}

func (r *laporanRepository) GetLaporanGajiByPeriod(bulan, tahun int) ([]models.LaporanGaji, error) {
	var results []models.LaporanGaji

	query := `
		SELECT
			g.id,
			g.karyawan_id,
			k.nik,
			k.nama AS nama_karyawan,
			COALESCE(j.nama_jabatan, '-') AS jabatan,
			g.periode_bulan,
			g.periode_tahun,
			g.gaji_pokok,
			g.tunjangan_jabatan,
			g.tunjangan_transport,
			g.tunjangan_makan,
			g.lembur,
			g.potongan,
			g.total_gaji,
			g.status
		FROM gaji g
		INNER JOIN karyawan k ON g.karyawan_id = k.id
		LEFT JOIN jabatan j ON k.jabatan_id = j.id
		WHERE g.periode_bulan = ? AND g.periode_tahun = ?
		ORDER BY k.nama ASC
	`

	err := r.db.Raw(query, bulan, tahun).Scan(&results).Error
	return results, err
}

func (r *laporanRepository) GetRiwayatGajiKaryawan(karyawanID uint) ([]models.LaporanGaji, error) {
	var results []models.LaporanGaji

	query := `
		SELECT
			g.id,
			g.karyawan_id,
			k.nik,
			k.nama AS nama_karyawan,
			COALESCE(j.nama_jabatan, '-') AS jabatan,
			g.periode_bulan,
			g.periode_tahun,
			g.gaji_pokok,
			g.tunjangan_jabatan,
			g.tunjangan_transport,
			g.tunjangan_makan,
			g.lembur,
			g.potongan,
			g.total_gaji,
			g.status
		FROM gaji g
		INNER JOIN karyawan k ON g.karyawan_id = k.id
		LEFT JOIN jabatan j ON k.jabatan_id = j.id
		WHERE k.id = ?
		ORDER BY g.periode_tahun DESC, g.periode_bulan DESC
	`

	err := r.db.Raw(query, karyawanID).Scan(&results).Error
	return results, err
}

func (r *laporanRepository) GetRekapGaji(bulan, tahun int) (*models.RekapGaji, error) {
	var result models.RekapGaji

	query := `
		SELECT
			? AS periode_bulan,
			? AS periode_tahun,
			COUNT(DISTINCT g.karyawan_id) AS total_karyawan,
			COALESCE(SUM(g.gaji_pokok), 0) AS total_gaji_pokok,
			COALESCE(SUM(g.tunjangan_jabatan + g.tunjangan_transport + g.tunjangan_makan), 0) AS total_tunjangan,
			COALESCE(SUM(g.lembur), 0) AS total_lembur,
			COALESCE(SUM(g.potongan), 0) AS total_potongan,
			COALESCE(SUM(g.total_gaji), 0) AS total_gaji,
			SUM(CASE WHEN g.status = 'pending' THEN 1 ELSE 0 END) AS status_pending,
			SUM(CASE WHEN g.status = 'dibayar' THEN 1 ELSE 0 END) AS status_dibayar
		FROM gaji g
		WHERE g.periode_bulan = ? AND g.periode_tahun = ?
	`

	err := r.db.Raw(query, bulan, tahun, bulan, tahun).Scan(&result).Error
	return &result, err
}
