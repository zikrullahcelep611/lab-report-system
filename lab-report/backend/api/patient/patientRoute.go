package patient

import "github.com/gofiber/fiber/v2"

func RegisterPatientRouter(router fiber.Router, handler *PatientHandler) {
	router.Post("/patient", handler.CreatePatient)
	router.Get("/patient/:id", handler.GetPatientByID)
	router.Put("/patient/:id", handler.UpdatePatient)
	router.Delete("/patient:id", handler.DeletePatient)
}
