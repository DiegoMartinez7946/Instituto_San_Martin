package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// ModalidadDoc catálogo (colección modalidad): TITULAR, PROVISIONAL, SUPLENTE
type ModalidadDoc struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Codigo string             `bson:"codigo" json:"codigo"`
}
