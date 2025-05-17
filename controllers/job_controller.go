package controllers

import (
	"html/template"
	"job-tracker/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	log.Printf("Attempting to render template: %s", tmpl)
	tmpls, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		log.Printf("Error parsing template %s: %v", tmpl, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error: failed to parse template"))
		return
	}
	log.Printf("Template %s parsed successfully", tmpl)
	if err := tmpls.ExecuteTemplate(w, tmpl, data); err != nil {
		log.Printf("Error executing template %s: %v", tmpl, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error: failed to execute template"))
		return
	}
	log.Printf("Template %s executed successfully", tmpl)
}

func Index(w http.ResponseWriter, r *http.Request) {
	applications, err := models.GetAllApplications()
	if err != nil {
		log.Printf("Error fetching applications: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	for _, app := range applications {
		log.Printf("Rendering application with ID: %d", app.ID)
	}
	renderTemplate(w, "index.html", struct{ Applications []models.JobApplication }{Applications: applications})
}

func CreateApplication(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var app models.JobApplication
		app.Company = r.FormValue("company")
		app.Position = r.FormValue("position")
		app.ApplicationDate = r.FormValue("application_date")
		app.Status = r.FormValue("status")

		_, err := models.CreateJobApplication(app)
		if err != nil {
			log.Printf("Error creating application: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetApplication(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		log.Printf("Invalid application ID in GetApplication: %s", idStr)
		http.Error(w, "Invalid application ID", http.StatusBadRequest)
		return
	}
	app, err := models.GetApplicationByID(id)
	if err != nil {
		log.Printf("Error fetching application ID %d: %v", id, err)
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}
	renderTemplate(w, "index.html", struct{ Applications []models.JobApplication }{Applications: []models.JobApplication{app}})
}

func UpdateApplication(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		log.Printf("Invalid application ID in UpdateApplication: %s", idStr)
		http.Error(w, "Invalid application ID", http.StatusBadRequest)
		return
	}
	if r.Method == "POST" {
		var app models.JobApplication
		app.Company = r.FormValue("company")
		app.Position = r.FormValue("position")
		app.ApplicationDate = r.FormValue("application_date")
		app.Status = r.FormValue("status")

		if err := models.UpdateApplication(id, app); err != nil {
			log.Printf("Error updating application ID %d: %v", id, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	app, err := models.GetApplicationByID(id)
	if err != nil {
		log.Printf("Error fetching application ID %d for edit: %v", id, err)
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}
	renderTemplate(w, "edit.html", app)
}

func DeleteApplication(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID format in DeleteApplication: %s, error: %v", idStr, err)
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	if err := models.DeleteApplication(id); err != nil {
		log.Printf("Error deleting application with ID %d: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func FilterApplications(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received filter request with URL: %s", r.URL.String())
	company := r.URL.Query().Get("company")
	status := r.URL.Query().Get("status")
	log.Printf("Filtering applications with company=%q, status=%q", company, status)

	// Fetch filtered applications
	apps, err := models.FilterApplications(company, status)
	if err != nil {
		log.Printf("Error filtering applications: %v", err)
		http.Error(w, "Internal server error: failed to fetch applications", http.StatusInternalServerError)
		return
	}
	log.Printf("Found %d applications", len(apps))

	// Validate applications
	var validApps []models.JobApplication
	for _, app := range apps {
		if app.ID > 0 {
			log.Printf("Filtered application: ID=%d, Company=%s, Status=%s", app.ID, app.Company, app.Status)
			validApps = append(validApps, app)
		} else {
			log.Printf("Skipping invalid application with ID=%d", app.ID)
		}
	}
	if len(validApps) == 0 {
		log.Printf("No valid applications to render after filtering")
		validApps = []models.JobApplication{}
	}

	// Render the template
	log.Printf("Rendering filtered results with %d valid applications", len(validApps))
	renderTemplate(w, "filtered.html", struct{ Applications []models.JobApplication }{Applications: validApps})
}
