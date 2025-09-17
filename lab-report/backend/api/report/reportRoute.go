package report

import "github.com/gorilla/mux"

func RegisterReportRoutes(routes *mux.Router, handler *ReportHandler) {
	routes.HandleFunc("/reports", handler.CreateReport).Methods("POST")
	routes.HandleFunc("/reports", handler.GetAllReportsOrdered).Methods("GET")
	routes.HandleFunc("/reports/search", handler.GetReportsWithPatientName).Methods("GET")
	routes.HandleFunc("/reports/search-by-national-id", handler.GetReportsWithPatientNationalID).Methods("GET")
	routes.HandleFunc("/reports/{id:[0-9]+}", handler.GetReportByID).Methods("GET")
	routes.HandleFunc("/reports/{id:[0-9]+}", handler.UpdateReport).Methods("PUT")
	routes.HandleFunc("/reports/{id:[0-9]+}", handler.DeleteReport).Methods("DELETE")
}
