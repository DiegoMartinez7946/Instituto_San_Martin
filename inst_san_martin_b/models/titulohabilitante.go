package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// TituloHabilitante catálogo (colección titulo_habilitante): SI / NO
type TituloHabilitante struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Codigo string             `bson:"codigo" json:"codigo"`
}
