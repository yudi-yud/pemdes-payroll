package handlers

import (
	"net/http"
	"pemdes-payroll/backend/models"
	"pemdes-payroll/backend/repositories"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type GajiHandler struct {
	gajiRepo     repositories.GajiRepository
	karyawanRepo repositories.KaryawanRepository
	lemburRepo   repositories.LemburRepository
}

// NewGajiHandler creates a new Gaji handler
func NewGajiHandler(gajiRepo repositories.GajiRepository, karyawanRepo repositories.KaryawanRepository, lemburRepo repositories.LemburRepository) *GajiHandler {
	return &GajiHandler{
		gajiRepo:     gajiRepo,
		karyawanRepo: karyawanRepo,
		lemburRepo:   lemburRepo,
	}
}

// CreateGaji handles POST /api/gaji
func (h *GajiHandler) CreateGaji(c *fiber.Ctx) error {
	var req struct {
		KaryawanID         uint    `json:"karyawan_id"`
		PeriodeBulan       int     `json:"periode_bulan"`
		PeriodeTahun       int     `json:"periode_tahun"`
		GajiPokok          float64 `json:"gaji_pokok"`
		TunjanganJabatan   float64 `json:"tunjangan_jabatan"`
		TunjanganTransport float64 `json:"tunjangan_transport"`
		TunjanganMakan     float64 `json:"tunjangan_makan"`
		Lembur             float64 `json:"lembur"`
		Potongan           float64 `json:"potongan"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation
	if req.KaryawanID == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Karyawan ID is required",
		})
	}
	if req.PeriodeBulan < 1 || req.PeriodeBulan > 12 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Periode bulan must be between 1 and 12",
		})
	}
	if req.PeriodeTahun < 2000 || req.PeriodeTahun > 2100 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid periode tahun",
		})
	}

	// Check if karyawan exists
	karyawan, err := h.karyawanRepo.GetByID(req.KaryawanID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Karyawan not found",
		})
	}

	// Check if salary already exists for this period
	existing, _ := h.gajiRepo.GetByKaryawanAndPeriod(int(req.KaryawanID), req.PeriodeBulan, req.PeriodeTahun)
	if existing != nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "Salary already exists for this period",
		})
	}

	// Calculate lembur from approved overtime records if not provided
	lemburAmount := req.Lembur
	if lemburAmount == 0 {
		_, lemburAmount, _ = h.lemburRepo.GetTotalLemburByPeriod(int(req.KaryawanID), req.PeriodeBulan, req.PeriodeTahun)
	}

	gaji := models.Gaji{
		KaryawanID:         req.KaryawanID,
		PeriodeBulan:       req.PeriodeBulan,
		PeriodeTahun:       req.PeriodeTahun,
		GajiPokok:          req.GajiPokok,
		TunjanganJabatan:   req.TunjanganJabatan,
		TunjanganTransport: req.TunjanganTransport,
		TunjanganMakan:     req.TunjanganMakan,
		Lembur:             lemburAmount,
		Potongan:           req.Potongan,
		Status:             models.GajiStatusPending,
	}

	gaji.CalculateTotal()

	if err := h.gajiRepo.Create(&gaji); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create gaji",
		})
	}

	// Fetch with Karyawan data
	result, _ := h.gajiRepo.GetByID(gaji.ID)
	result.Karyawan = *karyawan

	return c.Status(http.StatusCreated).JSON(result)
}

// GetAllGaji handles GET /api/gaji
func (h *GajiHandler) GetAllGaji(c *fiber.Ctx) error {
	gaji, err := h.gajiRepo.GetAll()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch gaji",
		})
	}

	return c.JSON(gaji)
}

// GetGajiByPeriod handles GET /api/gaji/period?bulan=&tahun=
func (h *GajiHandler) GetGajiByPeriod(c *fiber.Ctx) error {
	bulan, _ := strconv.Atoi(c.Query("bulan", "0"))
	tahun, _ := strconv.Atoi(c.Query("tahun", "0"))

	if bulan < 1 || bulan > 12 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bulan parameter",
		})
	}
	if tahun < 2000 || tahun > 2100 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tahun parameter",
		})
	}

	gaji, err := h.gajiRepo.GetByPeriod(bulan, tahun)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch gaji",
		})
	}

	return c.JSON(gaji)
}

// GetGajiByID handles GET /api/gaji/:id
func (h *GajiHandler) GetGajiByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	gaji, err := h.gajiRepo.GetByID(uint(id))
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Gaji not found",
		})
	}

	return c.JSON(gaji)
}

// GetGajiByKaryawanID handles GET /api/gaji/karyawan/:id
func (h *GajiHandler) GetGajiByKaryawanID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	gaji, err := h.gajiRepo.GetByKaryawanID(uint(id))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch gaji",
		})
	}

	return c.JSON(gaji)
}

// UpdateGaji handles PUT /api/gaji/:id
func (h *GajiHandler) UpdateGaji(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	var req struct {
		KaryawanID         uint              `json:"karyawan_id"`
		PeriodeBulan       int               `json:"periode_bulan"`
		PeriodeTahun       int               `json:"periode_tahun"`
		GajiPokok          float64           `json:"gaji_pokok"`
		TunjanganJabatan   float64           `json:"tunjangan_jabatan"`
		TunjanganTransport float64           `json:"tunjangan_transport"`
		TunjanganMakan     float64           `json:"tunjangan_makan"`
		Lembur             float64           `json:"lembur"`
		Potongan           float64           `json:"potongan"`
		Status             models.GajiStatus `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	gaji := models.Gaji{
		KaryawanID:         req.KaryawanID,
		PeriodeBulan:       req.PeriodeBulan,
		PeriodeTahun:       req.PeriodeTahun,
		GajiPokok:          req.GajiPokok,
		TunjanganJabatan:   req.TunjanganJabatan,
		TunjanganTransport: req.TunjanganTransport,
		TunjanganMakan:     req.TunjanganMakan,
		Lembur:             req.Lembur,
		Potongan:           req.Potongan,
		Status:             req.Status,
	}

	gaji.CalculateTotal()

	if err := h.gajiRepo.Update(uint(id), &gaji); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update gaji",
		})
	}

	// Get updated data
	updated, _ := h.gajiRepo.GetByID(uint(id))
	return c.JSON(updated)
}

// DeleteGaji handles DELETE /api/gaji/:id
func (h *GajiHandler) DeleteGaji(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	if err := h.gajiRepo.Delete(uint(id)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete gaji",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Gaji deleted successfully",
	})
}

// UpdateGajiStatus handles PATCH /api/gaji/:id/status
func (h *GajiHandler) UpdateGajiStatus(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	var req struct {
		Status models.GajiStatus `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Status != models.GajiStatusPending && req.Status != models.GajiStatusDibayar {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status. Use 'pending' or 'dibayar'",
		})
	}

	if err := h.gajiRepo.UpdateStatus(uint(id), req.Status); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update status",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Status updated successfully",
	})
}

// GenerateBatch handles POST /api/gaji/generate-batch
func (h *GajiHandler) GenerateBatch(c *fiber.Ctx) error {
	var req struct {
		PeriodeBulan       int     `json:"periode_bulan"`
		PeriodeTahun       int     `json:"periode_tahun"`
		TunjanganTransport float64 `json:"tunjangan_transport"`
		TunjanganMakan     float64 `json:"tunjangan_makan"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validation
	if req.PeriodeBulan < 1 || req.PeriodeBulan > 12 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Periode bulan must be between 1 and 12",
		})
	}
	if req.PeriodeTahun < 2000 || req.PeriodeTahun > 2100 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid periode tahun",
		})
	}

	// Get all active employees
	karyawanList, err := h.karyawanRepo.GetByStatus(models.StatusAktif)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch karyawan",
		})
	}

	var gajiList []models.Gaji
	var skipped []string
	var created int

	for _, k := range karyawanList {
		// Check if salary already exists for this period
		existing, _ := h.gajiRepo.GetByKaryawanAndPeriod(int(k.ID), req.PeriodeBulan, req.PeriodeTahun)
		if existing != nil {
			skipped = append(skipped, k.Nama)
			continue
		}

		gajiPokok := 0.0
		tunjanganJabatan := 0.0

		if k.Jabatan != nil {
			gajiPokok = k.Jabatan.GajiPokok
			tunjanganJabatan = k.Jabatan.TunjanganJabatan
		}

		// Get total lembur for this employee for the period (only approved)
		_, totalLemburNominal, _ := h.lemburRepo.GetTotalLemburByPeriod(int(k.ID), req.PeriodeBulan, req.PeriodeTahun)

		gaji := models.Gaji{
			KaryawanID:         k.ID,
			PeriodeBulan:       req.PeriodeBulan,
			PeriodeTahun:       req.PeriodeTahun,
			GajiPokok:          gajiPokok,
			TunjanganJabatan:   tunjanganJabatan,
			TunjanganTransport: req.TunjanganTransport,
			TunjanganMakan:     req.TunjanganMakan,
			Lembur:             totalLemburNominal,
			Potongan:           0,
			Status:             models.GajiStatusPending,
		}

		gaji.CalculateTotal()
		gajiList = append(gajiList, gaji)
		created++
	}

	if len(gajiList) > 0 {
		if err := h.gajiRepo.CreateBatch(gajiList); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create gaji batch",
			})
		}
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Batch generation completed",
		"created": created,
		"skipped": skipped,
		"periode": map[string]int{"bulan": req.PeriodeBulan, "tahun": req.PeriodeTahun},
	})
}
