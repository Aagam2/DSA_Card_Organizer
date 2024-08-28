package main

import (
    "encoding/json"
    "net/http"
)

type Card struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

var cards []Card

func getCardsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(cards)
}

func addCardHandler(w http.ResponseWriter, r *http.Request) {
    var newCard Card
    json.NewDecoder(r.Body).Decode(&newCard)
    cards = append(cards, newCard)
    w.WriteHeader(http.StatusCreated)
}
