package models

import (
	"fmt"
	"job-tracker/database"
	"log"
)

// JobApplication struct represents the job application structure.
type JobApplication struct {
	ID              int    `json:"id"`
	Company         string `json:"company"`
	Position        string `json:"position"`
	ApplicationDate string `json:"application_date"`
	Status          string `json:"status"`
}

// CreateJobApplication inserts a new job application into the database
func CreateJobApplication(app JobApplication) (int, error) {
	if app.Company == "" || app.Position == "" || app.Status == "" {
		return 0, fmt.Errorf("company, position, and status are required")
	}
	var id int
	err := database.DB.QueryRow(
		"INSERT INTO job_applications (company, position, application_date, status) VALUES (?, ?, ?, ?) RETURNING id",
		app.Company, app.Position, app.ApplicationDate, app.Status,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetAllApplications retrieves all job applications from the database
func GetAllApplications() ([]JobApplication, error) {
	rows, err := database.DB.Query("SELECT id, company, position, application_date, status FROM job_applications")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []JobApplication
	for rows.Next() {
		var app JobApplication
		err := rows.Scan(&app.ID, &app.Company, &app.Position, &app.ApplicationDate, &app.Status)
		if err != nil {
			return nil, err
		}
		applications = append(applications, app)
	}
	return applications, nil
}

// GetApplicationByID retrieves a specific job application by ID
func GetApplicationByID(id int) (JobApplication, error) {
	var app JobApplication
	err := database.DB.QueryRow(
		"SELECT id, company, position, application_date, status FROM job_applications WHERE id = ?",
		id,
	).Scan(&app.ID, &app.Company, &app.Position, &app.ApplicationDate, &app.Status)
	return app, err
}

// UpdateApplication updates an existing job application in the database
func UpdateApplication(id int, app JobApplication) error {
	if app.Company == "" || app.Position == "" || app.Status == "" {
		return fmt.Errorf("company, position, and status are required")
	}
	_, err := database.DB.Exec(
		"UPDATE job_applications SET company = ?, position = ?, application_date = ?, status = ? WHERE id = ?",
		app.Company, app.Position, app.ApplicationDate, app.Status, id,
	)
	return err
}

// DeleteApplication deletes a job application from the database
func DeleteApplication(id int) error {
	if err := database.DB.Ping(); err != nil {
		fmt.Printf("Database connection lost: %v\n", err)
		return err
	}
	fmt.Printf("Executing DELETE for ID: %d\n", id)
	result, err := database.DB.Exec("DELETE FROM job_applications WHERE id = ?", id)
	if err != nil {
		fmt.Printf("Error deleting application with ID %d: %v\n", id, err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error checking rows affected for ID %d: %v\n", id, err)
		return err
	}
	fmt.Printf("Deleted %d rows for ID: %d\n", rowsAffected, id)
	if rowsAffected == 0 {
		fmt.Printf("No rows deleted for ID %d (possibly ID does not exist)\n", id)
		return fmt.Errorf("no application found with ID %d", id)
	}
	return nil
}

// FilterApplications filters job applications based on company and status
func FilterApplications(company, status string) ([]JobApplication, error) {
	log.Printf("Executing FilterApplications with company=%q, status=%q", company, status)
	if err := database.DB.Ping(); err != nil {
		log.Printf("Database connection lost: %v", err)
		return nil, err
	}

	query := "SELECT id, company, position, application_date, status FROM job_applications WHERE 1=1"
	var args []interface{}
	var paramCount = 1

	if company != "" {
		query += fmt.Sprintf(" AND lower(company) LIKE lower(?)")
		args = append(args, "%"+company+"%")
		paramCount++
	}
	if status != "" {
		query += fmt.Sprintf(" AND status = ?")
		args = append(args, status)
		paramCount++
	}

	log.Printf("Executing query: %s with args: %v", query, args)
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var applications []JobApplication
	for rows.Next() {
		var app JobApplication
		err := rows.Scan(&app.ID, &app.Company, &app.Position, &app.ApplicationDate, &app.Status)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		applications = append(applications, app)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v", err)
		return nil, err
	}
	log.Printf("Successfully retrieved %d applications", len(applications))
	return applications, nil
}
