package tiketing

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"booking-api/models"
)

var bookingsCollection *mongo.Collection
var ticketsCollection *mongo.Collection

func BookTicket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var booking models.Booking
	_ = json.NewDecoder(r.Body).Decode(&booking)

	// Insert booking into MongoDB
	result, err := bookingsCollection.InsertOne(context.Background(), booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	booking.ID = result.InsertedID.(primitive.ObjectID)
	// Update ticket availability in MongoDB
	err = updateTicketAvailability(booking.TicketID, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(booking)
}

func GetBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	bookingID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	var booking models.Booking
	err = bookingsCollection.FindOne(context.Background(), bson.M{"_id": bookingID}).Decode(&booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(booking)
}

func CancelBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	bookingID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	var booking models.Booking
	err = bookingsCollection.FindOneAndDelete(context.Background(), bson.M{"_id": bookingID}).Decode(&booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update ticket availability in MongoDB
	err = updateTicketAvailability(booking.TicketID, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Booking canceled successfully"})
}

func updateTicketAvailability(ticketID primitive.ObjectID, available bool) error {
	_, err := ticketsCollection.UpdateOne(context.Background(), bson.M{"_id": ticketID}, bson.M{"$set": bson.M{"available": available}})
	return err
}
