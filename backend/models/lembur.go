package models

import (
	"time"

	"gorm.io/gorm"
)

// Lembur represents overtime record
type Lembur struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	KaryawanID      uint       `json:"karyawan_id" gorm:"not null;index"`
	Tanggal         time.Time  `json:"tanggal" gorm:"type:date;not null"`
	JamMulai        string     `json:"jam_mulai" gorm:"size:5"`
	JamSelesai      string     `json:"jam_selesai" gorm:"size:5"`
	TotalJam        float64    `json:"total_jam" gorm:"not null"`
	TarifPerJam     float64    `json:"tarif_per_jam" gorm:"not null"`
	TotalNominal    float64    `json:"total_nominal" gorm:"not null"`
	Keterangan      string     `json:"keterangan" gorm:"type:text"`
	Status          string     `json:"status" gorm:"default:'pending';type:enum('pending','disetujui','ditolak')"`
	DisetujuiOleh   *uint      `json:"disetujui_olej"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Karyawan        Karyawan   `json:"karyawan,omitempty" gorm:"foreignKey:KaryawanID"`
	Approver        *Karyawan  `json:"approver,omitempty" gorm:"foreignKey:DisetujuiOleh"`
}

// TableName specifies the table name for Lembur model
func (Lembur) TableName() string {
	return "lembur"
}

// BeforeCreate hook for validation
func (l *Lembur) BeforeCreate(tx *gorm.DB) error {
	if l.TotalJam <= 0 {
		l.TotalNominal = 0
	} else {
		l.TotalNominal = l.TotalJam * l.TarifPerJam
	}
	return nil
}

// BeforeUpdate hook to recalculate total
func (l *Lembur) BeforeUpdate(tx *gorm.DB) error {
	if l.TotalJam <= 0 {
		l.TotalNominal = 0
	} else {
		l.TotalNominal = l.TotalJam * l.TarifPerJam
	}
	return nil
}
