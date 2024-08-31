package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Card struct {
    ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name string             `json:"name"`
}

type Subtopic struct {
    ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    TopicID     primitive.ObjectID `json:"topicId" bson:"topicId"`
    Name        string             `json:"name"`
    Description string             `json:"description"`
}


var (
    cardsCollection *mongo.Collection
    subtopicsCollection *mongo.Collection
    client          *mongo.Client
)

func init() {
    // Load environment variables from .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    // Get MongoDB URL from environment variables
    mongoURL := os.Getenv("MONGODB_URL")
    if mongoURL == "" {
        log.Fatalf("MONGODB_URL not set in .env file")
    }

    // Set client options
    clientOptions := options.Client().ApplyURI(mongoURL)

    // Connect to MongoDB
    client, err = mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    // Check the connection
    err = client.Ping(context.TODO(), nil)
    if err != nil {
        log.Fatalf("Failed to ping MongoDB: %v", err)
    }

    // Get a handle for the cards collection
    cardsCollection = client.Database("DSA_Cards").Collection("Topics")
    subtopicsCollection = client.Database("DSA_Cards").Collection("Subtopics")
}

func getCardsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var cards []Card
    cursor, err := cardsCollection.Find(context.TODO(), bson.D{})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.TODO())

    for cursor.Next(context.TODO()) {
        var card Card
        cursor.Decode(&card)
        cards = append(cards, card)
    }

    if err := cursor.Err(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(cards)
}

func addCardHandler(w http.ResponseWriter, r *http.Request) {
    var newCard Card
    json.NewDecoder(r.Body).Decode(&newCard)

    result, err := cardsCollection.InsertOne(context.TODO(), newCard)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Set the ID of the new card
    newCard.ID = result.InsertedID.(primitive.ObjectID)

    // Return the updated list of cards
    getCardsHandler(w, r)
}

func getSubtopicsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var subtopics []Subtopic
    cursor, err := subtopicsCollection.Find(context.TODO(), bson.D{})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer cursor.Close(context.TODO())

    for cursor.Next(context.TODO()) {
        var subtopic Subtopic
        cursor.Decode(&subtopic)
        subtopics = append(subtopics, subtopic)
    }

    if err := cursor.Err(); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(subtopics)
}

func addSubtopicHandler(w http.ResponseWriter, r *http.Request) {
    var newSubtopic Subtopic
    json.NewDecoder(r.Body).Decode(&newSubtopic)

    result, err := subtopicsCollection.InsertOne(context.TODO(), newSubtopic)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Set the ID of the new subtopic
    newSubtopic.ID = result.InsertedID.(primitive.ObjectID)

    // Return the updated list of subtopics
    getSubtopicsHandler(w, r)
}