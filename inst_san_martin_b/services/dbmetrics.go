package services

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/benjacifre10/san_martin_b/config"
	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	dbMetricsSnapshotMu   sync.Mutex
	dbMetricsLastInsert   time.Time
	dbMetricsInsertMinGap = 60 * time.Second
)

func friendlyServerStatusDenied(raw string) string {
	s := strings.TrimSpace(raw)
	low := strings.ToLower(s)
	if strings.Contains(low, "not allowed") || strings.Contains(low, "unauthorized") {
		return "Permisos: el usuario de conexión no puede ejecutar serverStatus (métricas a nivel cluster). " +
			"En MongoDB Atlas → Database Access → editar el usuario → añadir rol «clusterMonitor». " +
			"Detalle: " + s
	}
	return s
}

func intFromMap(m bson.M, key string) int {
	v, ok := m[key]
	if !ok || v == nil {
		return 0
	}
	switch t := v.(type) {
	case int32:
		return int(t)
	case int64:
		return int(t)
	case int:
		return t
	case float64:
		return int(t)
	default:
		return 0
	}
}

func int64FromMap(m bson.M, key string) int64 {
	v, ok := m[key]
	if !ok || v == nil {
		return 0
	}
	switch t := v.(type) {
	case int32:
		return int64(t)
	case int64:
		return t
	case int:
		return int64(t)
	case float64:
		return int64(t)
	default:
		return 0
	}
}

func float64FromMap(m bson.M, key string) float64 {
	v, ok := m[key]
	if !ok || v == nil {
		return 0
	}
	switch t := v.(type) {
	case float64:
		return t
	case int32:
		return float64(t)
	case int64:
		return float64(t)
	default:
		return 0
	}
}

func applyDbStatsFallback(ctx context.Context, out *models.DBMetricsPayload) {
	db := config.WorkingDatabase()
	var raw bson.M
	err := db.RunCommand(ctx, bson.D{{Key: "dbStats", Value: 1}, {Key: "scale", Value: 1}}).Decode(&raw)
	if err != nil {
		return
	}
	ok := false
	switch v := raw["ok"].(type) {
	case bool:
		ok = v
	case int32:
		ok = v == 1
	case int64:
		ok = v == 1
	case float64:
		ok = int(v) == 1
	}
	if !ok {
		return
	}
	out.PartialMetricsOnly = true
	if s, _ := raw["db"].(string); s != "" {
		out.DbStats.DB = s
	}
	out.DbStats.Collections = intFromMap(raw, "collections")
	out.DbStats.Objects = int64FromMap(raw, "objects")
	out.DbStats.DataSize = int64FromMap(raw, "dataSize")
	out.DbStats.StorageSize = int64FromMap(raw, "storageSize")
	out.DbStats.Indexes = intFromMap(raw, "indexes")
	out.DbStats.IndexSize = int64FromMap(raw, "indexSize")
}

func nestedMap(m bson.M, key string) bson.M {
	v, ok := m[key]
	if !ok || v == nil {
		return nil
	}
	switch t := v.(type) {
	case bson.M:
		return t
	case map[string]interface{}:
		out := bson.M{}
		for k, val := range t {
			out[k] = val
		}
		return out
	default:
		return nil
	}
}

/* GetDBMetricsService ping + serverStatus y opcionalmente persiste snapshot. */
func GetDBMetricsService() (models.DBMetricsPayload, string, int) {
	atlasHelp := " En MongoDB Atlas: Database Access → tu usuario → Edit → Built-in Role → agregar «clusterMonitor» (solo métricas, lectura)."
	out := models.DBMetricsPayload{
		CollectedAt: time.Now().UTC(),
		Note:        "Los contadores de operaciones y red son acumulados desde el último reinicio del proceso mongod. El uso de CPU del servidor no forma parte de serverStatus; en MongoDB Atlas podés ver CPU en el panel de métricas del cluster." + atlasHelp,
	}

	client := config.MongoConnection
	ctx, cancel := context.WithTimeout(context.Background(), 18*time.Second)
	defer cancel()

	t0 := time.Now()
	if err := client.Ping(ctx, nil); err != nil {
		out.ServerStatusOk = false
		out.ServerStatusMsg = err.Error()
		return out, "No se pudo hacer ping a MongoDB", 503
	}
	out.PingMs = float64(time.Since(t0).Microseconds()) / 1000.0

	var raw bson.M
	err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "serverStatus", Value: 1}}).Decode(&raw)
	if err != nil {
		out.ServerStatusOk = false
		out.ServerStatusMsg = friendlyServerStatusDenied(err.Error())
		applyDbStatsFallback(ctx, &out)
		maybePersistSnapshot(out)
		return out, "", 200
	}

	okStatus := false
	switch v := raw["ok"].(type) {
	case bool:
		okStatus = v
	case float64:
		okStatus = int(v) == 1
	case int32:
		okStatus = v == 1
	case int64:
		okStatus = v == 1
	default:
		if raw["uptime"] != nil || nestedMap(raw, "connections") != nil {
			okStatus = true
		}
	}
	if !okStatus {
		if v, has := raw["errmsg"].(string); has {
			out.ServerStatusMsg = friendlyServerStatusDenied(v)
		} else {
			out.ServerStatusMsg = friendlyServerStatusDenied("serverStatus devolvió ok=false")
		}
		out.ServerStatusOk = false
		applyDbStatsFallback(ctx, &out)
		maybePersistSnapshot(out)
		return out, "", 200
	}

	out.ServerStatusOk = true
	out.UptimeSec = int64FromMap(raw, "uptime")

	conn := nestedMap(raw, "connections")
	if conn != nil {
		out.Connections.Current = intFromMap(conn, "current")
		out.Connections.Active = intFromMap(conn, "active")
		out.Connections.Available = intFromMap(conn, "available")
		out.Connections.TotalCreated = int64FromMap(conn, "totalCreated")
		out.Connections.Threaded = intFromMap(conn, "threaded")
		out.Connections.ExhaustIsMaster = intFromMap(conn, "exhaustIsMaster")
	}

	mem := nestedMap(raw, "mem")
	if mem != nil {
		out.MemoryMB.Resident = intFromMap(mem, "resident")
		out.MemoryMB.Virtual = intFromMap(mem, "virtual")
	}

	op := nestedMap(raw, "opcounters")
	if op != nil {
		out.Opcounters.Query = int64FromMap(op, "query")
		out.Opcounters.GetMore = int64FromMap(op, "getmore")
		out.Opcounters.Insert = int64FromMap(op, "insert")
		out.Opcounters.Update = int64FromMap(op, "update")
		out.Opcounters.Delete = int64FromMap(op, "delete")
		out.Opcounters.Command = int64FromMap(op, "command")
	}

	net := nestedMap(raw, "network")
	if net != nil {
		out.Network.BytesIn = int64FromMap(net, "bytesIn")
		out.Network.BytesOut = int64FromMap(net, "bytesOut")
		out.Network.NumRequests = int64FromMap(net, "numRequests")
		out.Network.Compression = int64FromMap(net, "compression")
	}

	wt := nestedMap(raw, "wiredTiger")
	if wt != nil {
		cache := nestedMap(wt, "cache")
		if cache != nil {
			// "bytes currently in the cache" como número aproximado en MB
			bytesInCache := float64FromMap(cache, "bytes currently in the cache")
			if bytesInCache > 0 {
				out.WiredTigerCacheMB = bytesInCache / (1024 * 1024)
			}
		}
	}

	maybePersistSnapshot(out)
	return out, "", 200
}

func maybePersistSnapshot(p models.DBMetricsPayload) {
	dbMetricsSnapshotMu.Lock()
	defer dbMetricsSnapshotMu.Unlock()
	if time.Since(dbMetricsLastInsert) < dbMetricsInsertMinGap {
		return
	}
	dbMetricsLastInsert = time.Now()

	snap := models.DBMetricSnapshot{
		CreatedAt:    p.CollectedAt,
		PingMs:       p.PingMs,
		ConnCurrent:  p.Connections.Current,
		MemResident:  p.MemoryMB.Resident,
		ConnTotal:    p.Connections.TotalCreated,
		QueriesTotal: p.Opcounters.Query,
		UptimeSec:    p.UptimeSec,
	}
	if err := db.InsertDBMetricSnapshot(snap); err != nil {
		log.Println("InsertDBMetricSnapshot:", err)
	}
}
