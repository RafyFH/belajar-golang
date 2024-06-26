package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Ticket struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ConcertName string             `json:"concert_name"`
	Artist      string             `json:"artist"`
	Price       float64            `json:"price"`
	Available   bool               `json:"available"`
}

type Booking struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TicketID primitive.ObjectID `json:"ticket_id"`
	Name     string             `json:"name"`
	Email    string             `json:"email"`
	Quantity int                `json:"quantity"`
	Total    float64            `json:"total"`
}
