package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DesignCategory struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Category string             `json:"category" bson:"category"`
}
