package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Pembayaran struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DesignSelected   string             `json:"design_selected" bson:"design_selected"`
	OrderDescription string             `json:"order_description" bson:"order_description"`
	CardFullname     string             `json:"card_fullname" bson:"card_fullname"`
	CardNumber       string             `json:"card_number" bson:"card_number"`
	CardExpiration   string             `json:"card_expication" bson:"card_expication"`
	CVV              string             `json:"cvv" bson:"cvv"`
	Price            string             `json:"price" bson:"price"`
}
