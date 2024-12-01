package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Feedback struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Category FeedbackCategory   `bson:"category" json:"category"`
	Comments string             `json:"comments" bson:"comments"`
	Image    string             `json:"image" bson:"image"`
}
