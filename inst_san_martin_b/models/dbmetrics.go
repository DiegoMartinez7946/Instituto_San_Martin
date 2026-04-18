package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* DBMetricSnapshot muestra histórico ligero en colección dbmetrics (opcional, alimentado al consultar el monitor). */
type DBMetricSnapshot struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt    time.Time          `bson:"createdat" json:"createdAt"`
	PingMs       float64            `bson:"pingms" json:"pingMs"`
	ConnCurrent  int                `bson:"conncurrent" json:"connCurrent"`
	MemResident  int                `bson:"memresidentmb" json:"memResidentMb"`
	ConnTotal    int64              `bson:"conntotalcreated" json:"connTotalCreated"`
	QueriesTotal int64              `bson:"queriesTotal" json:"queriesTotal"`
	UptimeSec    int64              `bson:"uptimesec" json:"uptimeSec"`
}

/* DBMetricsPayload respuesta JSON al panel de administración (derivado de serverStatus + ping). */
type DBMetricsPayload struct {
	CollectedAt time.Time `json:"collectedAt"`
	PingMs      float64   `json:"pingMs"`

	UptimeSec int64 `json:"uptimeSec"`

	Connections struct {
		Current       int   `json:"current"`
		Active        int   `json:"active"`
		Available     int   `json:"available"`
		TotalCreated  int64 `json:"totalCreated"`
		Threaded      int   `json:"threaded,omitempty"`
		ExhaustIsMaster int `json:"exhaustIsMaster,omitempty"`
	} `json:"connections"`

	MemoryMB struct {
		Resident int `json:"resident"`
		Virtual  int `json:"virtual"`
	} `json:"memoryMb"`

	Opcounters struct {
		Query   int64 `json:"query"`
		GetMore int64 `json:"getMore"`
		Insert  int64 `json:"insert"`
		Update  int64 `json:"update"`
		Delete  int64 `json:"delete"`
		Command int64 `json:"command"`
	} `json:"opcounters"`

	Network struct {
		BytesIn      int64 `json:"bytesIn"`
		BytesOut     int64 `json:"bytesOut"`
		NumRequests  int64 `json:"numRequests"`
		Compression  int64 `json:"compression,omitempty"`
	} `json:"network"`

	WiredTigerCacheMB float64 `json:"wiredTigerCacheMb,omitempty"`

	ServerStatusOk bool   `json:"serverStatusOk"`
	ServerStatusMsg string `json:"serverStatusMsg,omitempty"`

	/* Si serverStatus falla por permisos, intentamos dbStats sobre la base de trabajo (readWrite suele alcanzar). */
	PartialMetricsOnly bool `json:"partialMetricsOnly"`
	DbStats            struct {
		DB          string `json:"db,omitempty"`
		Collections int    `json:"collections,omitempty"`
		Objects     int64  `json:"objects,omitempty"`
		DataSize    int64  `json:"dataSize,omitempty"`
		StorageSize int64  `json:"storageSize,omitempty"`
		Indexes     int    `json:"indexes,omitempty"`
		IndexSize   int64  `json:"indexSize,omitempty"`
	} `json:"dbStats,omitempty"`

	Note string `json:"note,omitempty"`
}
