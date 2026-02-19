package models

import (
	"time"

	"gorm.io/gorm"
)

// KaryawanStatus represents the status of an employee
type KaryawanStatus string

const (
	StatusAktif    KaryawanStatus = "aktif"
	StatusNonAktif KaryawanStatus = "non_aktif"
)

// Karyawan represents an employee
type Karyawan struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	NIK             string         `json:"nik" gorm:"unique;not null;size:20"`
	Nama            string         `json:"nama" gorm:"not null;size:100"`
	Email           string         `json:"email" gorm:"size:100"`
	Telepon         string         `json:"telepon" gorm:"size:20"`
	Alamat          string         `json:"alamat" gorm:"type:text"`
	JabatanID       *uint          `json:"jabatan_id" gorm:"index"`
	TanggalBergabung *time.Time    `json:"tanggal_bergabung" gorm:"type:date"`
	Status          KaryawanStatus `json:"status" gorm:"default:'aktif';type:enum('aktif','non_aktif')"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	Jabatan         *Jabatan       `json:"jabatan,omitempty" gorm:"foreignKey:JabatanID"`
	Gaji            []Gaji         `json:"gaji,omitempty" gorm:"foreignKey:KaryawanID"`
}

// TableName specifies the table name for Karyawan model
func (Karyawan) TableName() string {
	return "karyawan"
}

// BeforeDelete hook to prevent deletion if salary records exist
func (k *Karyawan) BeforeDelete(tx *gorm.DB) error {
	var count int64
	tx.Model(&Gaji{}).Where("karyawan_id = ?", k.ID).Count(&count)
	if count > 0 {
		return gorm.ErrInvalidData
	}
	return nil
}
