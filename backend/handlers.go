package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
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

type Algorithm struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SubtopicID  primitive.ObjectID `json:"subtopicId" bson:"subtopicId"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	CodeFileID  primitive.ObjectID `json:"codeFileId" bson:"codeFileId"`
	Code        string             `json:"code"`
}

var (
	cardsCollection      *mongo.Collection
	subtopicsCollection  *mongo.Collection
	algorithmsCollection *mongo.Collection
	client               *mongo.Client
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
	algorithmsCollection = client.Database("DSA_Cards").Collection("Algorithms")
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

func getAlgorithmsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the subtopicId from the query parameters
	subtopicID := r.URL.Query().Get("subtopicId")
	if subtopicID == "" {
		http.Error(w, "subtopicId is required", http.StatusBadRequest)
		return
	}

	// Convert subtopicID to ObjectID
	subtopicObjectID, err := primitive.ObjectIDFromHex(subtopicID)
	if err != nil {
		http.Error(w, "Invalid subtopicId", http.StatusBadRequest)
		return
	}

	// Find algorithms related to the subtopic
	var algorithms []Algorithm
	cursor, err := algorithmsCollection.Find(context.TODO(), bson.M{"subtopicId": subtopicObjectID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var algorithm Algorithm
		cursor.Decode(&algorithm)

		// Fetch the file from GridFS
		bucket, err := gridfs.NewBucket(
			client.Database("DSA_Cards"),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var buf bytes.Buffer
		_, err = bucket.DownloadToStream(algorithm.CodeFileID, &buf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		algorithm.Code = buf.String()
		algorithms = append(algorithms, algorithm)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(algorithms)
}

func uploadAlgorithmHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the form values
	subtopicID := r.FormValue("subtopicId")
	name := r.FormValue("name")
	description := r.FormValue("description")

	// Get the file from the form
	file, header, err := r.FormFile("code")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Convert subtopicID to ObjectID
	subtopicObjectID, err := primitive.ObjectIDFromHex(subtopicID)
	if err != nil {
		http.Error(w, "Invalid subtopicId", http.StatusBadRequest)
		return
	}

	// Create a GridFS bucket
	bucket, err := gridfs.NewBucket(
		client.Database("DSA_Cards"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Upload the file to GridFS
	uploadStream, err := bucket.OpenUploadStream(header.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadStream.Close()

	_, err = io.Copy(uploadStream, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new algorithm
	newAlgorithm := Algorithm{
		SubtopicID:  subtopicObjectID,
		Name:        name,
		Description: description,
		CodeFileID:  uploadStream.FileID.(primitive.ObjectID),
	}

	// Insert the algorithm into the database
	result, err := algorithmsCollection.InsertOne(context.TODO(), newAlgorithm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the ID of the new algorithm
	newAlgorithm.ID = result.InsertedID.(primitive.ObjectID)

	// Return the new algorithm as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newAlgorithm)
}

func getAlgorithmHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get the algorithmId from the URL parameters
	vars := mux.Vars(r)
	algorithmID := vars["algorithmId"]

	// Convert algorithmID to ObjectID
	algorithmObjectID, err := primitive.ObjectIDFromHex(algorithmID)
	if err != nil {
		http.Error(w, "Invalid algorithmId", http.StatusBadRequest)
		return
	}

	// Find the algorithm by ID
	var algorithm Algorithm
	err = algorithmsCollection.FindOne(context.TODO(), bson.M{"_id": algorithmObjectID}).Decode(&algorithm)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Algorithm not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the file from GridFS
	bucket, err := gridfs.NewBucket(
		client.Database("DSA_Cards"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	_, err = bucket.DownloadToStream(algorithm.CodeFileID, &buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	algorithm.Code = buf.String()

	json.NewEncoder(w).Encode(algorithm)
}
