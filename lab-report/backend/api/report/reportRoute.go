package report

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterReportRoutes(app *fiber.App, handler *ReportHandler) {
	 api := app.Group("/api")
    
    // CRUD operations
    api.Post("/reports", handler.CreateReport)                    // Create
    api.Get("/reports", handler.GetAllReportsOrdered)            // Read All
    api.Get("/reports/:id", handler.GetReportByID)               // Read One
    api.Put("/reports/:id", handler.UpdateReport)                // Update
    api.Delete("/reports/:id", handler.DeleteReport)             // Delete
    
    // Search operations
    api.Get("/reports/search", handler.GetReportsWithPatientName)          // Search by name
    api.Get("/reports/search-national", handler.GetReportsWithPatientNationalID) // Search by national ID
}
