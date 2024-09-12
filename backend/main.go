package main

import (
    "log"
    "net/http"
)

func main() {
    r := setupRoutes()

    // Serve static files from the "static" directory
    staticFileDirectory := http.Dir("/home/aagam-linux/Desktop/Projects/DSA_Organizer/frontend/static")
    staticFileHandler := http.StripPrefix("/static/", http.FileServer(staticFileDirectory))

    // Serve HTML files from the "templates" directory
    templatesFileDirectory := http.Dir("/home/aagam-linux/Desktop/Projects/DSA_Organizer/frontend/templates")
    templatesFileHandler := http.StripPrefix("/", http.FileServer(templatesFileDirectory))

    // Handle static files
    r.PathPrefix("/static/").Handler(staticFileHandler)

    // Handle HTML templates
    r.PathPrefix("/").Handler(templatesFileHandler)

    log.Println("Server is running on http://localhost:8080/")
    log.Fatal(http.ListenAndServe(":8080", r))
}
