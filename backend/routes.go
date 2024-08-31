package main

import "net/http"

func setupRoutes() {
    http.HandleFunc("/api/cards", getCardsHandler)
    http.HandleFunc("/api/cards/add", addCardHandler)
    http.HandleFunc("/api/subtopics", getSubtopicsHandler)
    http.HandleFunc("/api/subtopics/add", addSubtopicHandler)
}