package hospital

import "github.com/gofiber/fiber/v2"

func RegisterHospitalRoutes(router fiber.Router, handler *HospitalHandler) {
	router.Post("/hospital", handler.CreateHospital)
	router.Get("/hospital/:id", handler.GetHospitalByID)
	router.Put("/hospital/:id", handler.UpdateHospital)
	router.Delete("/hospital/:id", handler.DeleteHospital)
}
