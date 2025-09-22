package hospital

import (
	"context"
	"strconv"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/hospital"
	"github.com/gofiber/fiber/v2"
)

type HospitalService interface {
	GetHospitalByID(ctx context.Context, id uint) (hospital.Hospital, error)
	CreateHospital(ctx context.Context, newHospital hospital.Hospital) (hospital.Hospital, error)
	UpdateHospital(ctx context.Context, updateHospital hospital.Hospital) (hospital.Hospital, error)
	DeleteHospital(ctx context.Context, id uint) error
}

type HospitalHandler struct {
	hospitalService HospitalService
}

func NewHospitalHandler(hospitalService HospitalService) *HospitalHandler {
	return &HospitalHandler{hospitalService: hospitalService}
}

func (h *HospitalHandler) GetHospitalByID(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid hospital ID",
		})
	}
	hospitalData, err := h.hospitalService.GetHospitalByID(ctx, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(hospitalData)
}

func (h *HospitalHandler) CreateHospital(c *fiber.Ctx) error {
	var newHospital hospital.Hospital
	ctx := c.Context()
	if err := c.BodyParser(&newHospital); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	createdHospital, err := h.hospitalService.CreateHospital(ctx, newHospital)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(createdHospital)
}

func (h *HospitalHandler) UpdateHospital(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid hospital ID",
		})
	}

	var updateHospital hospital.Hospital
	if err := c.BodyParser(&updateHospital); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	updateHospital.ID = uint(id)
	updatedHospital, err := h.hospitalService.UpdateHospital(ctx, updateHospital)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(updatedHospital)
}

func (h *HospitalHandler) DeleteHospital(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Invalid hospital ID",
		})
	}

	if err := h.hospitalService.DeleteHospital(ctx, uint(id)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Hospital deleted successfully",
	})
}
