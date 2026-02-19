package models

import (
	"time"

	"gorm.io/gorm"
)

// AbsensiStatus represents attendance status
type AbsensiStatus string

const (
	AbsensiHadir  AbsensiStatus = "hadir"
	AbsensiIzin   AbsensiStatus = "izin"
	AbsensiSakit  AbsensiStatus = "sakit"
	AbsensiAlpha  AbsensiStatus = "alpha"
)

// Absensi represents daily attendance record
type Absensi struct {
	ID          uint          `json:"id" gorm:"primaryKey"`
	KaryawanID  uint          `json:"karyawan_id" gorm:"not null;index"`
	Tanggal     time.Time     `json:"tanggal" gorm:"type:date;not null;index"`
	JamMasuk    string        `json:"jam_masuk" gorm:"size:10"`
	JamKeluar   string        `json:"jam_keluar" gorm:"size:10"`
	Status      AbsensiStatus `json:"status" gorm:"default:'hadir';type:enum('hadir','izin','sakit','alpha')"`
	Keterangan  string        `json:"keterangan" gorm:"type:text"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	Karyawan    Karyawan      `json:"karyawan,omitempty" gorm:"foreignKey:KaryawanID"`
}

// TableName specifies the table name for Absensi model
func (Absensi) TableName() string {
	return "absensi"
}

// BeforeCreate hook to ensure unique attendance per employee per date
func (a *Absensi) BeforeCreate(tx *gorm.DB) error {
	var existing Absensi
	err := tx.Where("karyawan_id = ? AND tanggal = ?", a.KaryawanID, a.Tanggal).First(&existing).Error
	if err == nil {
		return gorm.ErrDuplicatedKey
	}
	return nil
}
