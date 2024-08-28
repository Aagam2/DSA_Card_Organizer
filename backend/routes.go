package main

import "net/http"

func setupRoutes() {
    http.HandleFunc("/api/cards", getCardsHandler)
    http.HandleFunc("/api/cards/add", addCardHandler)
    // Add more routes as needed
}
