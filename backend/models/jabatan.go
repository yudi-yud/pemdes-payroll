package models

import (
	"time"
)

// Jabatan represents a job position in the company
type Jabatan struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	NamaJabatan        string    `json:"nama_jabatan" gorm:"not null;size:100"`
	GajiPokok          float64   `json:"gaji_pokok" gorm:"not null;type:decimal(15,2)"`
	TunjanganJabatan   float64   `json:"tunjangan_jabatan" gorm:"default:0;type:decimal(15,2)"`
	TarifLemburPerJam  float64   `json:"tarif_lembur_per_jam" gorm:"default:0;type:decimal(15,2)"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Karyawan           []Karyawan `json:"karyawan,omitempty" gorm:"foreignKey:JabatanID"`
}

// TableName specifies the table name for Jabatan model
func (Jabatan) TableName() string {
	return "jabatan"
}
