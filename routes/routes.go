package routes

import (
	"job-tracker/controllers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRoutes configures the application routes
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Log all incoming requests
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Received request: %s %s", r.Method, r.URL.String())
			next.ServeHTTP(w, r)
		})
	})

	// Define routes
	router.HandleFunc("/", controllers.Index).Methods("GET")
	router.HandleFunc("/applications", controllers.CreateApplication).Methods("POST")
	// Move more specific routes before the generic /applications/{id}
	router.HandleFunc("/applications/filter", controllers.FilterApplications).Methods("GET")
	router.HandleFunc("/applications/{id}", controllers.GetApplication).Methods("GET")
	router.HandleFunc("/applications/{id}", controllers.UpdateApplication).Methods("POST")
	router.HandleFunc("/applications/{id}/edit", controllers.UpdateApplication).Methods("GET")
	router.HandleFunc("/applications/{id}/delete", controllers.DeleteApplication).Methods("POST")

	return router
}
