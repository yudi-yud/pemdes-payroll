package models

import (
	"time"

	"gorm.io/gorm"
)

// GajiStatus represents the payment status
type GajiStatus string

const (
	GajiStatusPending GajiStatus = "pending"
	GajiStatusDibayar  GajiStatus = "dibayar"
)

// Gaji represents an employee's salary record
type Gaji struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	KaryawanID         uint      `json:"karyawan_id" gorm:"not null;index"`
	PeriodeBulan       int       `json:"periode_bulan" gorm:"not null"`
	PeriodeTahun       int       `json:"periode_tahun" gorm:"not null"`
	GajiPokok          float64   `json:"gaji_pokok" gorm:"not null;type:decimal(15,2)"`
	TunjanganJabatan   float64   `json:"tunjangan_jabatan" gorm:"default:0;type:decimal(15,2)"`
	TunjanganTransport float64   `json:"tunjangan_transport" gorm:"default:0;type:decimal(15,2)"`
	TunjanganMakan     float64   `json:"tunjangan_makan" gorm:"default:0;type:decimal(15,2)"`
	Lembur             float64   `json:"lembur" gorm:"default:0;type:decimal(15,2)"`
	Potongan           float64   `json:"potongan" gorm:"default:0;type:decimal(15,2)"`
	TotalGaji          float64   `json:"total_gaji" gorm:"not null;type:decimal(15,2)"`
	Status             GajiStatus `json:"status" gorm:"default:'pending';type:enum('pending','dibayar')"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Karyawan           Karyawan  `json:"karyawan,omitempty" gorm:"foreignKey:KaryawanID"`
}

// TableName specifies the table name for Gaji model
func (Gaji) TableName() string {
	return "gaji"
}

// BeforeCreate hook to ensure unique period per employee
func (g *Gaji) BeforeCreate(tx *gorm.DB) error {
	var existing Gaji
	err := tx.Where("karyawan_id = ? AND periode_bulan = ? AND periode_tahun = ?",
		g.KaryawanID, g.PeriodeBulan, g.PeriodeTahun).First(&existing).Error
	if err == nil {
		return gorm.ErrDuplicatedKey
	}
	return nil
}

// CalculateTotal computes the total salary
func (g *Gaji) CalculateTotal() {
	g.TotalGaji = g.GajiPokok + g.TunjanganJabatan + g.TunjanganTransport + g.TunjanganMakan + g.Lembur - g.Potongan
}
