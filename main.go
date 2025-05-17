package main

import (
	"job-tracker/database" //This is a custom package for database-related operations
	"job-tracker/routes"   // Custom package for the setting up the HTTP routes
	"log"                  //standard libraries for logging
	"net/http"             //standard library http server
	"os/exec"              // executing external commands
	"runtime"              // Detect OS type duirng run time
)

// browser attempts to open the specific URL in the default web browser
func openBrowser(url string) {
	var cmd string
	var args []string

	//checking os and assigning the appropriate command thats needd to open the browser
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin": // macOS
		cmd = "open"
		args = []string{url}
	default: // Linux, BSD, etc.
		cmd = "xdg-open"
		args = []string{url}
	}

	err := exec.Command(cmd, args...).Start()
	if err != nil {
		log.Printf("Failed to open browser: %v", err)
	}
}

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Set up router
	router := routes.SetupRoutes()

	// Start server
	log.Println("Server starting on :8080")
	go func() {
		log.Println("Opening browser at http://localhost:8080")
		openBrowser("http://localhost:8080")
	}()
	log.Fatal(http.ListenAndServe(":8080", router))
}
