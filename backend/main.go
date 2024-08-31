package main

import (
    "log"
    "net/http"
)

func main() {
    // Serve static files from the "frontend" directory
    fs := http.FileServer(http.Dir("/home/aagam-linux/Desktop/Projects/DSA_Organizer/frontend"))
    http.Handle("/", fs)

    // Initialize the routes
    setupRoutes()

    log.Println("Server started at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}