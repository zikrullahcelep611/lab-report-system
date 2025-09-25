package report

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterReportRoutes(router fiber.Router, handler *ReportHandler) {

	// CRUD operations
	router.Post("/reports", handler.CreateReport)        // Create
	router.Get("/reports", handler.GetAllReportsOrdered) // Read All

	
	router.Get("/reports/search", handler.GetReportsWithPatientName)                // Search by name
	router.Get("/reports/search-national", handler.GetReportsWithPatientNationalID) // Search by national ID

	// Parametrik routes (EN SONDA olmalı)
	router.Get("/reports/:id", handler.GetReportByID)   // Read One
	router.Put("/reports/:id", handler.UpdateReport)    // Update
	router.Delete("/reports/:id", handler.DeleteReport) // Delete
}
