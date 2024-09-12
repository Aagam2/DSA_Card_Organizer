package main

import (
    "github.com/gorilla/mux"
)

func setupRoutes() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/api/cards", getCardsHandler).Methods("GET")
    r.HandleFunc("/api/cards/add", addCardHandler).Methods("POST")
    r.HandleFunc("/api/subtopics", getSubtopicsHandler).Methods("GET")
    r.HandleFunc("/api/subtopics/add", addSubtopicHandler).Methods("POST")
    r.HandleFunc("/api/algorithms", getAlgorithmsHandler).Methods("GET")
    r.HandleFunc("/api/algorithms/add", uploadAlgorithmHandler).Methods("POST")
    r.HandleFunc("/api/algorithms/{algorithmId}", getAlgorithmHandler).Methods("GET")
    r.HandleFunc("/api/notes/{algorithmId}", getNotesHandler).Methods("GET")
    r.HandleFunc("/api/notes/{algorithmId}", saveNotesHandler).Methods("POST")
    return r
}