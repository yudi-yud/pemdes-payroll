package models

// LaporanGaji represents a salary report
type LaporanGaji struct {
	ID                 uint    `json:"id"`
	KaryawanID         uint    `json:"karyawan_id"`
	NIK                string  `json:"nik"`
	NamaKaryawan       string  `json:"nama_karyawan"`
	Jabatan            string  `json:"jabatan"`
	PeriodeBulan       int     `json:"periode_bulan"`
	PeriodeTahun       int     `json:"periode_tahun"`
	GajiPokok          float64 `json:"gaji_pokok"`
	TunjanganJabatan   float64 `json:"tunjangan_jabatan"`
	TunjanganTransport float64 `json:"tunjangan_transport"`
	TunjanganMakan     float64 `json:"tunjangan_makan"`
	Lembur             float64 `json:"lembur"`
	Potongan           float64 `json:"potongan"`
	TotalGaji          float64 `json:"total_gaji"`
	Status             string  `json:"status"`
}

// RekapGaji represents salary recapitulation
type RekapGaji struct {
	PeriodeBulan    int     `json:"periode_bulan"`
	PeriodeTahun    int     `json:"periode_tahun"`
	TotalKaryawan   int     `json:"total_karyawan"`
	TotalGajiPokok  float64 `json:"total_gaji_pokok"`
	TotalTunjangan  float64 `json:"total_tunjangan"`
	TotalLembur     float64 `json:"total_lembur"`
	TotalPotongan   float64 `json:"total_potongan"`
	TotalGaji       float64 `json:"total_gaji"`
	StatusPending   int     `json:"status_pending"`
	StatusDibayar   int     `json:"status_dibayar"`
}
