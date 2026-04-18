package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/* SystemMetricSnapshot histórico ligero en colección systemmetrics (opcional). */
type SystemMetricSnapshot struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CreatedAt      time.Time          `bson:"createdat" json:"createdAt"`
	CpuPercent     float64            `bson:"cpupercent" json:"cpuPercent"`
	MemUsedPercent float64            `bson:"memusedpercent" json:"memUsedPercent"`
	DiskUsedPercent float64           `bson:"diskusedpercent" json:"diskUsedPercent"`
	Goroutines     int                `bson:"goroutines" json:"goroutines"`
	HeapAllocBytes uint64             `bson:"heapallocbytes" json:"heapAllocBytes"`
}

/* TopProcessMetric filas para identificar procesos que más RAM consumen (instantáneo). */
type TopProcessMetric struct {
	PID      int32   `json:"pid"`
	Name     string  `json:"name"`
	RSSBytes uint64  `json:"rssBytes"`
	CPUPercent float64 `json:"cpuPercent,omitempty"`
}

/* SystemMetricsPayload métricas del host + proceso Go API backend (solo panel admin). */
type SystemMetricsPayload struct {
	CollectedAt time.Time `json:"collectedAt"`

	Host struct {
		Hostname  string `json:"hostname"`
		UptimeSec uint64 `json:"uptimeSec"`
		OS        string `json:"os"`
		Platform  string `json:"platform"`
		KernelVer string `json:"kernelVer,omitempty"`
	} `json:"host"`

	CpuPercent float64 `json:"cpuPercent"`

	LoadAvg struct {
		Load1  float64 `json:"load1,omitempty"`
		Load5  float64 `json:"load5,omitempty"`
		Load15 float64 `json:"load15,omitempty"`
	} `json:"loadAvg"`

	Memory struct {
		TotalBytes     uint64  `json:"totalBytes"`
		UsedBytes      uint64  `json:"usedBytes"`
		UsedPercent    float64 `json:"usedPercent"`
		AvailableBytes uint64  `json:"availableBytes"`
	} `json:"memory"`

	Swap struct {
		TotalBytes  uint64  `json:"totalBytes"`
		UsedPercent float64 `json:"usedPercent"`
	} `json:"swap"`

	Disk struct {
		Path        string  `json:"path"`
		TotalBytes  uint64  `json:"totalBytes"`
		UsedBytes   uint64  `json:"usedBytes"`
		UsedPercent float64 `json:"usedPercent"`
	} `json:"disk"`

	GoRuntime struct {
		Goroutines     int    `json:"goroutines"`
		HeapAllocBytes uint64 `json:"heapAllocBytes"`
		HeapSysBytes   uint64 `json:"heapSysBytes"`
		NumGC          uint32 `json:"numGc"`
		LastGCPauseNs  uint64 `json:"lastGcPauseNs"`
	} `json:"goRuntime"`

	ProcessSelf struct {
		PID        int32   `json:"pid"`
		CPUPercent float64 `json:"cpuPercent"`
		RSSBytes   uint64  `json:"rssBytes"`
		Name       string  `json:"name,omitempty"`
	} `json:"processSelf"`

	/* Ordenados por RSS descendente (no implica mayor CPU). */
	TopProcessesByRSS []TopProcessMetric `json:"topProcessesByRss,omitempty"`

	Note string `json:"note,omitempty"`
}
