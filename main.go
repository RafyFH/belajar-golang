package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"booking-api/tiketing"
)

var client *mongo.Client
var ticketsCollection *mongo.Collection
var bookingsCollection *mongo.Collection

func main() {
	// Set up MongoDB client
	clientOptions := options.Client().ApplyURI("mongodb+srv://rafyfakhrizal299:Dynmt3Du4rr@cluster0.empkic1.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")

	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Set up collections
	ticketsCollection = client.Database("ticketing").Collection("tickets")
	bookingsCollection = client.Database("ticketing").Collection("bookings")

	// Seed data to MongoDB
	// seedData()

	// Set up router
	router := mux.NewRouter()

	// Ticket Routes
	router.HandleFunc("/api/tickets", getTickets).Methods("GET")
	router.HandleFunc("/api/tickets/{id}", getTicketByID).Methods("GET")

	router.HandleFunc("/api/auth", tiketing.SignIn).Methods("POST")

	// Booking Routes with Authentication
	router.HandleFunc("/api/bookings", tiketing.SignIn).Methods("POST")
	router.HandleFunc("/api/bookings", tiketing.IsAuthenticated(tiketing.BookTicket)).Methods("POST")
	router.HandleFunc("/api/bookings/{id}", tiketing.IsAuthenticated(tiketing.GetBooking)).Methods("GET")
	router.HandleFunc("/api/bookings/{id}", tiketing.IsAuthenticated(tiketing.CancelBooking)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func seedData() {
	// Implementasi seeding data ke MongoDB
}

func getTickets(w http.ResponseWriter, r *http.Request) {
	// Context for the database operations
	ctx := context.Background()

	// Find options to limit results
	findOptions := options.Find()

	// Slice to store the decoded documents
	var tickets []bson.M

	// Find all documents
	cur, err := ticketsCollection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	// Iterate through the cursor and decode each document
	for cur.Next(ctx) {
		var ticket bson.M
		err := cur.Decode(&ticket)
		if err != nil {
			log.Fatal(err)
		}
		tickets = append(tickets, ticket)
	}

	// If error while iterating
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Convert to JSON
	jsonBytes, err := json.Marshal(tickets)
	if err != nil {
		log.Fatal(err)
	}

	// Set Content-Type header and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// Function to fetch a ticket by ID
func getTicketByID(w http.ResponseWriter, r *http.Request) {
	// Context for the database operations
	ctx := context.Background()

	// Get the ticket ID from the request parameters
	params := mux.Vars(r)
	ticketID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		log.Fatal(err)
	}

	// Filter to find the specific ticket by ID
	filter := bson.D{{"_id", ticketID}}

	// Variable to store the decoded document
	var ticket bson.M

	// Find the document
	err = ticketsCollection.FindOne(ctx, filter).Decode(&ticket)
	if err != nil {
		log.Fatal(err)
	}

	// Convert to JSON
	jsonBytes, err := json.Marshal(ticket)
	if err != nil {
		log.Fatal(err)
	}

	// Set Content-Type header and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
