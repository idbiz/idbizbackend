package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Pemesanan struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Fullname         string             `json:"fullname" bson:"fullname"`
	Email            string             `json:"email" bson:"email"`
	PhoneNumber      string             `json:"phone_number" bson:"phone_number"`
	Category         DesignCategory     `bson:"category" json:"category"`
	OrderDescription string             `json:"order_description" bson:"order_description"`
	UploadReferences string             `json:"upload_references" bson:"upload_references"`
}
