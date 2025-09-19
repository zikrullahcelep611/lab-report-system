package report

import (
	"context"
	"strconv"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/report"
	"github.com/gofiber/fiber/v2"
)

type ReportService interface {
	CreateReport(ctx context.Context, newReport report.Report) (report.Report, error)
	GetReportsWithPatientName(ctx context.Context, firstName string, lastName string) ([]report.Report, error)
	GetReportsWithPatientNationalID(ctx context.Context, nationalID string) ([]report.Report, error)
	GetAllReportsOrdered(ctx context.Context) ([]report.Report, error)
	UpdateReport(ctx context.Context, updateReport report.Report) (report.Report, error)
	DeleteReport(ctx context.Context, id uint) (bool, error)
	GetReportByID(ctx context.Context, id uint) (report.Report, error)
}

type ReportHandler struct {
	reportService ReportService
}

func NewReportController(reportService ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) CreateReport(c *fiber.Ctx) error { // error return ekle
    ctx := context.Background()
    var newReport report.Report
    if err := c.BodyParser(&newReport); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ // return ekle
            "error": "Invalid request payload",
        })
    }

    createdReport, err := h.reportService.CreateReport(ctx, newReport)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ // return ekle
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusCreated).JSON(createdReport) // return ekle + 201 Created
}

func (h *ReportHandler) GetReportsWithPatientName(c *fiber.Ctx) error { // error return ekle
    ctx := c.Context()
    firstName := c.Query("firstName")
    lastName := c.Query("lastName")

    if firstName == "" || lastName == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ // return ekle
            "error": "firstName and lastName parameters are required",
        })
    }

    reports, err := h.reportService.GetReportsWithPatientName(ctx, firstName, lastName)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ // return ekle
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(reports) // return ekle
}

func (h *ReportHandler) GetReportsWithPatientNationalID(c *fiber.Ctx) error { // error return ekle
    ctx := c.Context()
    nationalId := c.Query("nationalID")
    if nationalId == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ // return ekle
            "error": "nationalID parameter is required",
        })
    }

    reports, err := h.reportService.GetReportsWithPatientNationalID(ctx, nationalId)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ // return ekle (BadRequest değil)
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(reports) // return ekle
}

func (h *ReportHandler) GetAllReportsOrdered(c *fiber.Ctx) error { // error return ekle
    ctx := c.Context()
    reports, err := h.reportService.GetAllReportsOrdered(ctx)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ // return ekle
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(reports) // return ekle
}

func (h *ReportHandler) GetReportByID(c *fiber.Ctx) error { // error return ekle
    ctx := c.Context()
    idStr := c.Params("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ // return ekle
            "error": "Invalid report ID",
        })
    }

    reportData, err := h.reportService.GetReportByID(ctx, uint(id))
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{ // return ekle
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(reportData) // return ekle
}

func (h *ReportHandler) UpdateReport(c *fiber.Ctx) error { // error return ekle
    ctx := c.Context()
    idStr := c.Params("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ // return ekle
            "error": "Invalid report ID",
        })
    }

    var updateReport report.Report
    if err := c.BodyParser(&updateReport); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ // return ekle
            "error": "Invalid request payload",
        })
    }

    updateReport.ID = uint(id)
    updatedReport, err := h.reportService.UpdateReport(ctx, updateReport)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{ // return ekle
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(updatedReport) // return ekle
}

func (h *ReportHandler) DeleteReport(c *fiber.Ctx) error { // error return ekle
    ctx := c.Context()
    idStr := c.Params("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{ // return ekle
            "error": "Invalid report ID",
        })
    }

    deletedReport, err := h.reportService.DeleteReport(ctx, uint(id))
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{ // return ekle
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{ // return ekle + Status ekle
        "deleted": deletedReport,
        "message": "Report deleted successfully",
    })
}