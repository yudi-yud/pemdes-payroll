package services

import (
	"bytes"
	"fmt"
	"pemdes-payroll/backend/models"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

// ExportService handles export operations
type ExportService struct{}

// NewExportService creates a new export service
func NewExportService() *ExportService {
	return &ExportService{}
}

// formatCurrency formats a number to Indonesian currency format
func formatCurrency(amount float64) string {
	return fmt.Sprintf("Rp %.2f", amount)
}

// getMonthName returns Indonesian month name
func getMonthName(month int) string {
	months := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	if month >= 1 && month <= 12 {
		return months[month-1]
	}
	return ""
}

// ExportToExcel exports salary report to Excel format for all employees in a period
func (s *ExportService) ExportToExcel(laporanList []models.LaporanGaji, bulan, tahun int) ([]byte, error) {
	f := excelize.NewFile()
	sheetName := "Laporan Gaji"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	// Set active sheet
	f.SetActiveSheet(index)

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 8)
	f.SetColWidth(sheetName, "B", "C", 20)
	f.SetColWidth(sheetName, "D", "L", 18)

	// Header styles
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 11,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}

	// Title style
	titleStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 14,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	if err != nil {
		return nil, err
	}

	// Number style (with thousand separator)
	numStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size: 10,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		},
		NumFmt: 3, // Number format with 2 decimal places
	})
	if err != nil {
		return nil, err
	}

	// Company title
	f.SetCellValue(sheetName, "A1", "SISTEM PAYROLL PEMERINTAH DESA")
	f.SetCellStyle(sheetName, "A1", "L1", titleStyle)
	f.MergeCell(sheetName, "A1", "L1")

	// Report title
	periodeText := fmt.Sprintf("LAPORAN GAJI KARYAWAN - %s %d", getMonthName(bulan), tahun)
	f.SetCellValue(sheetName, "A2", periodeText)
	f.SetCellStyle(sheetName, "A2", "L2", titleStyle)
	f.MergeCell(sheetName, "A2", "L2")

	// Table headers
	row := 4
	headers := []string{"No", "NIK", "Nama Karyawan", "Jabatan", "Gaji Pokok", "Tunj. Jabatan", "Transport", "Makan", "Lembur", "Potongan", "Total Gaji", "Status"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s%d", string(rune('A'+i)), row)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}
	row++

	// Data rows
	for i, item := range laporanList {
		col := 'A'
		f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), i+1)
		col++

		f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), item.NIK)
		col++
		f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), item.NamaKaryawan)
		col++
		f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), item.Jabatan)
		col++

		// Numeric values
		f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), item.GajiPokok)
		f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), numStyle)
		col++

		f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), item.TunjanganJabatan)
		f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), numStyle)
		col++

		f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), item.TunjanganTransport)
		f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), numStyle)
		col++

		f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), item.TunjanganMakan)
		f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), numStyle)
		col++

		f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), item.Lembur)
		f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), numStyle)
		col++

		f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), item.Potongan)
		f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), numStyle)
		col++

		f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), item.TotalGaji)
		f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), numStyle)
		col++

		f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), string(item.Status))
		col++

		row++
	}

	// Total row
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "TOTAL:")
	f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("I%d", row), headerStyle)

	// Calculate totals
	totalGajiPokok := 0.0
	totalTunjanganJabatan := 0.0
	totalTunjanganTransport := 0.0
	totalTunjanganMakan := 0.0
	totalLembur := 0.0
	totalPotongan := 0.0
	totalGaji := 0.0

	for _, item := range laporanList {
		totalGajiPokok += item.GajiPokok
		totalTunjanganJabatan += item.TunjanganJabatan
		totalTunjanganTransport += item.TunjanganTransport
		totalTunjanganMakan += item.TunjanganMakan
		totalLembur += item.Lembur
		totalPotongan += item.Potongan
		totalGaji += item.TotalGaji
	}

	col := 'E'
	f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), totalGajiPokok)
	f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), headerStyle)
	col++

	f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), totalTunjanganJabatan)
	f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), headerStyle)
	col++

	f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), totalTunjanganTransport)
	f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), headerStyle)
	col++

	f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), totalTunjanganMakan)
	f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), headerStyle)
	col++

	f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), totalLembur)
	f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), headerStyle)
	col++

	f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), totalPotongan)
	f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), headerStyle)
	col++

	f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), totalGaji)
	f.SetCellStyle(sheetName, fmt.Sprintf("%c%d", col, row), fmt.Sprintf("%c%d", col, row), headerStyle)

	// Delete default Sheet1
	f.DeleteSheet("Sheet1")

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// ExportToPDF exports salary report to PDF format per employee
func (s *ExportService) ExportToPDF(karyawan *models.Karyawan, gajiList []models.Gaji) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Fonts - use built-in fonts
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "SISTEM PAYROLL PEMERINTAH DESA")
	pdf.Ln(7)

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "LAPORAN GAJI KARYAWAN")
	pdf.Ln(10)

	// Employee info
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(30, 7, "NIK:")
	pdf.Cell(60, 7, karyawan.NIK)
	pdf.Ln(7)

	pdf.Cell(30, 7, "Nama:")
	pdf.Cell(60, 7, karyawan.Nama)
	pdf.Ln(7)

	jabatan := "Tanpa Jabatan"
	if karyawan.Jabatan != nil {
		jabatan = karyawan.Jabatan.NamaJabatan
	}
	pdf.Cell(30, 7, "Jabatan:")
	pdf.Cell(60, 7, jabatan)
	pdf.Ln(10)

	// Table header
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(200, 200, 200)

	headers := []string{"No", "Periode", "Gaji Pokok", "Tunjangan", "Lembur", "Potongan", "Total", "Status"}
	colWidths := []float64{10, 30, 25, 25, 20, 20, 25, 20}

	// Draw header row
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 7, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(7)

	// Data rows
	pdf.SetFont("Arial", "", 8)
	pdf.SetFillColor(240, 240, 240)

	totalGaji := 0.0
	for i, g := range gajiList {
		// Alternate row color
		if i%2 == 0 {
			pdf.SetFillColor(245, 245, 245)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		periode := fmt.Sprintf("%s %d", getMonthName(g.PeriodeBulan), g.PeriodeTahun)
		totalTunjangan := g.TunjanganJabatan + g.TunjanganTransport + g.TunjanganMakan

		pdf.CellFormat(colWidths[0], 6, fmt.Sprintf("%d", i+1), "1", 0, "C", true, 0, "")
		pdf.CellFormat(colWidths[1], 6, periode, "1", 0, "L", true, 0, "")
		pdf.CellFormat(colWidths[2], 6, formatCurrency(g.GajiPokok), "1", 0, "R", true, 0, "")
		pdf.CellFormat(colWidths[3], 6, formatCurrency(totalTunjangan), "1", 0, "R", true, 0, "")
		pdf.CellFormat(colWidths[4], 6, formatCurrency(g.Lembur), "1", 0, "R", true, 0, "")
		pdf.CellFormat(colWidths[5], 6, formatCurrency(g.Potongan), "1", 0, "R", true, 0, "")
		pdf.CellFormat(colWidths[6], 6, formatCurrency(g.TotalGaji), "1", 0, "R", true, 0, "")
		pdf.CellFormat(colWidths[7], 6, string(g.Status), "1", 0, "C", true, 0, "")
		pdf.Ln(6)

		totalGaji += g.TotalGaji
	}

	// Total row
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(200, 200, 200)
	pdf.CellFormat(90, 7, "TOTAL:", "1", 0, "R", true, 0, "")
	pdf.CellFormat(25, 7, formatCurrency(totalGaji), "1", 0, "R", true, 0, "")
	pdf.CellFormat(20, 7, "", "1", 0, "C", true, 0, "")
	pdf.Ln(10)

	// Footer
	pdf.SetFont("Arial", "I", 8)
	pdf.Cell(190, 5, "Dicetak pada: "+time.Now().Format("02-01-2006 15:04:05"))
	pdf.Ln(4)
	pdf.Cell(190, 5, "Sistem Payroll Pemerintah Desa")

	// Output to bytes
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// ExportAbsensiToExcel exports attendance list to Excel format for all employees
func (s *ExportService) ExportAbsensiToExcel(absensiList []models.Absensi) ([]byte, error) {
	f := excelize.NewFile()
	sheetName := "Laporan Absensi"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	// Set active sheet
	f.SetActiveSheet(index)

	// Set column widths
	f.SetColWidth(sheetName, "A", "A", 8)
	f.SetColWidth(sheetName, "B", "C", 20)
	f.SetColWidth(sheetName, "D", "H", 15)

	// Header styles
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 11,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}

	// Title style
	titleStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 14,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	if err != nil {
		return nil, err
	}

	// Company title
	f.SetCellValue(sheetName, "A1", "SISTEM PAYROLL PEMERINTAH DESA")
	f.SetCellStyle(sheetName, "A1", "H1", titleStyle)
	f.MergeCell(sheetName, "A1", "H1")

	// Report title
	f.SetCellValue(sheetName, "A2", "LAPORAN ABSENSI KARYAWAN")
	f.SetCellStyle(sheetName, "A2", "H2", titleStyle)
	f.MergeCell(sheetName, "A2", "H2")

	// Table headers
	row := 4
	headers := []string{"No", "Tanggal", "NIK", "Nama Karyawan", "Jam Masuk", "Jam Keluar", "Status", "Keterangan"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s%d", string(rune('A'+i)), row)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}
	row++

	// Data rows
	for i, item := range absensiList {
		col := 'A'
		f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), i+1)
		col++

		tanggal := item.Tanggal.Format("02-01-2006")
		f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), tanggal)
		col++

		f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), item.Karyawan.NIK)
		col++
		f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), item.Karyawan.Nama)
		col++
		f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), item.JamMasuk)
		col++
		f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), item.JamKeluar)
		col++
		f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), string(item.Status))
		col++
		f.SetCellValue(sheetName, fmt.Sprintf("%c%d", col, row), item.Keterangan)
		col++

		row++
	}

	// Delete default Sheet1
	f.DeleteSheet("Sheet1")

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// ExportAbsensiToPDF exports attendance recap to PDF format per employee
func (s *ExportService) ExportAbsensiToPDF(karyawan *models.Karyawan, absensiList []models.Absensi, rekap map[string]int, bulan, tahun int) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Fonts - use built-in fonts
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "SISTEM PAYROLL PEMERINTAH DESA")
	pdf.Ln(7)

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "REKAP ABSENSI KARYAWAN")
	pdf.Ln(10)

	// Employee info
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(30, 7, "NIK:")
	pdf.Cell(60, 7, karyawan.NIK)
	pdf.Ln(7)

	pdf.Cell(30, 7, "Nama:")
	pdf.Cell(60, 7, karyawan.Nama)
	pdf.Ln(7)

	jabatan := "Tanpa Jabatan"
	if karyawan.Jabatan != nil {
		jabatan = karyawan.Jabatan.NamaJabatan
	}
	pdf.Cell(30, 7, "Jabatan:")
	pdf.Cell(60, 7, jabatan)
	pdf.Ln(7)

	periode := fmt.Sprintf("%s %d", getMonthName(bulan), tahun)
	pdf.Cell(30, 7, "Periode:")
	pdf.Cell(60, 7, periode)
	pdf.Ln(10)

	// Summary boxes
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(200, 230, 200)
	pdf.CellFormat(40, 10, fmt.Sprintf("Hadir: %d", rekap["hadir"]), "1", 0, "C", true, 0, "")
	pdf.SetFillColor(200, 200, 230)
	pdf.CellFormat(40, 10, fmt.Sprintf("Izin: %d", rekap["izin"]), "1", 0, "C", true, 0, "")
	pdf.SetFillColor(230, 230, 200)
	pdf.CellFormat(40, 10, fmt.Sprintf("Sakit: %d", rekap["sakit"]), "1", 0, "C", true, 0, "")
	pdf.SetFillColor(230, 200, 200)
	pdf.CellFormat(40, 10, fmt.Sprintf("Alpha: %d", rekap["alpha"]), "1", 0, "C", true, 0, "")
	pdf.Ln(12)

	// Table header
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(200, 200, 200)

	headers := []string{"No", "Tanggal", "Jam Masuk", "Jam Keluar", "Status", "Keterangan"}
	colWidths := []float64{10, 30, 25, 25, 25, 60}

	// Draw header row
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 7, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(7)

	// Data rows
	pdf.SetFont("Arial", "", 8)
	pdf.SetFillColor(240, 240, 240)

	for i, a := range absensiList {
		// Alternate row color
		if i%2 == 0 {
			pdf.SetFillColor(245, 245, 245)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		tanggal := a.Tanggal.Format("02-01-2006")

		pdf.CellFormat(colWidths[0], 6, fmt.Sprintf("%d", i+1), "1", 0, "C", true, 0, "")
		pdf.CellFormat(colWidths[1], 6, tanggal, "1", 0, "L", true, 0, "")
		pdf.CellFormat(colWidths[2], 6, a.JamMasuk, "1", 0, "C", true, 0, "")
		pdf.CellFormat(colWidths[3], 6, a.JamKeluar, "1", 0, "C", true, 0, "")
		pdf.CellFormat(colWidths[4], 6, string(a.Status), "1", 0, "C", true, 0, "")
		pdf.CellFormat(colWidths[5], 6, a.Keterangan, "1", 0, "L", true, 0, "")
		pdf.Ln(6)
	}

	// Footer
	pdf.SetFont("Arial", "I", 8)
	pdf.Ln(5)
	pdf.Cell(190, 5, "Dicetak pada: "+time.Now().Format("02-01-2006 15:04:05"))
	pdf.Ln(4)
	pdf.Cell(190, 5, "Sistem Payroll Pemerintah Desa")

	// Output to bytes
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
