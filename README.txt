Job Tracker Website
A simple web application to track job applications, built with Go, SQLite, and Gorilla Mux.
Prerequisites

Go 1.21 or later (download here)
A writable directory for the SQLite database

Setup

Extract the zip file to a directory (e.g., job-tracker).
Navigate to the project directory:cd job-tracker


Install dependencies:go mod tidy


Build the application:go build -o job-tracker


Run the application:./job-tracker  # macOS/Linux
job-tracker.exe  # Windows



Usage

The server runs at http://localhost:8080.
Use the web interface to add, edit, delete, and filter job applications.
The SQLite database (job_tracker.db) is created automatically in the project directory.

Troubleshooting

Port in use: If port 8080 is busy, edit main.go to use another port (e.g., :8081).
Database errors: Ensure the project directory is writable. Delete job_tracker.db to reset the database.
Dependencies fail: Run go get github.com/gorilla/mux@v1.8.1 and go get modernc.org/sqlite@v1.33.1 with an internet connection.

