package main

import (
    "log"
    "net/http"
)

func main() {
    r := setupRoutes()
    // Serve static files from the "static" directory
    staticFileDirectory := http.Dir("/home/aagam-linux/Desktop/Projects/DSA_Organizer/frontend")
    staticFileHandler := http.StripPrefix("/", http.FileServer(staticFileDirectory))

    // Handle static files and API routes
    r.PathPrefix("/").Handler(staticFileHandler)
    log.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}