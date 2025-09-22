package patient

import (
	"context"
	"strconv"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/patient" // gitbub -> github
	"github.com/gofiber/fiber/v2"
)

type PatientService interface {
	GetPatientByID(ctx context.Context, id uint) (patient.Patient, error)
	CreatePatient(ctx context.Context, newPatient patient.Patient) (patient.Patient, error)
	UpdatePatient(ctx context.Context, updatePatient patient.Patient) (patient.Patient, error)
	DeletePatient(ctx context.Context, id uint) error
}

type PatientHandler struct {
	patientService PatientService
}

func NewPatientHandler(patientService PatientService) *PatientHandler {
	return &PatientHandler{patientService: patientService}
}

func (p *PatientHandler) GetPatientByID(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid patient ID",
		})
	}

	patientData, err := p.patientService.GetPatientByID(ctx, uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(patientData)
}

func (p *PatientHandler) CreatePatient(c *fiber.Ctx) error {
	ctx := c.Context()
	var newPatient patient.Patient

	if err := c.BodyParser(&newPatient); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	createdPatient, err := p.patientService.CreatePatient(ctx, newPatient)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdPatient)
}

func (p *PatientHandler) UpdatePatient(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid patient ID",
		})
	}

	var updatePatient patient.Patient
	if err := c.BodyParser(&updatePatient); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	updatePatient.ID = uint(id)
	updatedPatient, err := p.patientService.UpdatePatient(ctx, updatePatient)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(updatedPatient)
}

func (p *PatientHandler) DeletePatient(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid patient ID",
		})
	}

	if err := p.patientService.DeletePatient(ctx, uint(id)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{ // BadRequest -> NotFound
			"error": err.Error(), // "Invalid request payload" -> err.Error()
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Patient deleted successfully",
	})
}
