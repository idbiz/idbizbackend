package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Portofolio struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Category    DesignCategory     `json:"category,omitempty" bson:"category,omitempty"`
	DesignTitle string             `json:"design_title,omitempty" bson:"design_title,omitempty"`
	DesignDesc  string             `json:"design_desc,omitempty" bson:"design_desc,omitempty"`
	DesignImage string             `json:"design_image,omitempty" bson:"design_image,omitempty"`
}
