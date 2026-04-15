// Package models provides ...
package models

/* ArgentinaTimeDisplay documento en colección "time" (referencia horaria Argentina, IANA). */
type ArgentinaTimeDisplay struct {
	Hora   string `bson:"hora" json:"hora"`
	Fecha  string `bson:"fecha" json:"fecha"`
	Dia    string `bson:"dia" json:"dia"`
	Mes    string `bson:"mes" json:"mes"`
	Ano       string `bson:"año" json:"año"`
	UpdatedAt string `bson:"updatedat,omitempty" json:"updatedAt,omitempty"`
}
