package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Note is the model that governs all notes objects retrived or inserted into the DB
type Note struct {
	ID         primitive.ObjectID `bson:"_id"`
	Text       string             `json:"text"`
	Title      string             `json:"title"`
	Created_at time.Time          `json:"created_at"`
	Updated_at time.Time          `json:"updated_at"`
	Note_id    string             `json:"note_id"`
}
