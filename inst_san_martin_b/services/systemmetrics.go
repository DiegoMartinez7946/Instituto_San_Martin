package services

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"sync"
	"time"

	"github.com/benjacifre10/san_martin_b/db"
	"github.com/benjacifre10/san_martin_b/models"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

var (
	systemMetricsSnapshotMu   sync.Mutex
	systemMetricsLastInsert     time.Time
	systemMetricsInsertMinGap   = 60 * time.Second
)

func diskMonitorPath() string {
	wd, err := os.Getwd()
	if err != nil || wd == "" {
		if runtime.GOOS == "windows" {
			return `C:\`
		}
		return "/"
	}
	if runtime.GOOS == "windows" {
		vol := filepath.VolumeName(wd)
		if vol != "" {
			return vol + `\`
		}
		return `C:\`
	}
	return "/"
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

/* GetSystemMetricsService CPU/RAM/disco del host + proceso Go + top RSS. */
func GetSystemMetricsService() (models.SystemMetricsPayload, string, int) {
	out := models.SystemMetricsPayload{
		CollectedAt: time.Now().UTC(),
		Note: "CPU host: muestra ~250 ms de muestreo. «Top procesos» ordena por memoria RSS (no por CPU); el proceso de esta API aparece como «backend Go». En Linux/macOS también se muestra load average.",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 22*time.Second)
	defer cancel()

	if hi, err := host.InfoWithContext(ctx); err == nil && hi != nil {
		out.Host.Hostname = hi.Hostname
		out.Host.UptimeSec = hi.Uptime
		out.Host.OS = hi.OS
		out.Host.Platform = hi.Platform
		out.Host.KernelVer = hi.KernelVersion
	}

	if avg, err := load.AvgWithContext(ctx); err == nil && avg != nil {
		out.LoadAvg.Load1 = avg.Load1
		out.LoadAvg.Load5 = avg.Load5
		out.LoadAvg.Load15 = avg.Load15
	}

	cpuPercents, err := cpu.PercentWithContext(ctx, 250*time.Millisecond, false)
	if err == nil && len(cpuPercents) > 0 {
		out.CpuPercent = cpuPercents[0]
	}

	if vm, err := mem.VirtualMemoryWithContext(ctx); err == nil && vm != nil {
		out.Memory.TotalBytes = vm.Total
		out.Memory.UsedBytes = vm.Used
		out.Memory.UsedPercent = vm.UsedPercent
		out.Memory.AvailableBytes = vm.Available
	}

	if sw, err := mem.SwapMemoryWithContext(ctx); err == nil && sw != nil {
		out.Swap.TotalBytes = sw.Total
		out.Swap.UsedPercent = sw.UsedPercent
	}

	dpath := diskMonitorPath()
	if du, err := disk.UsageWithContext(ctx, dpath); err == nil && du != nil {
		out.Disk.Path = dpath
		out.Disk.TotalBytes = du.Total
		out.Disk.UsedBytes = du.Used
		out.Disk.UsedPercent = du.UsedPercent
	}

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	out.GoRuntime.Goroutines = runtime.NumGoroutine()
	out.GoRuntime.HeapAllocBytes = ms.HeapAlloc
	out.GoRuntime.HeapSysBytes = ms.HeapSys
	out.GoRuntime.NumGC = ms.NumGC
	out.GoRuntime.LastGCPauseNs = ms.PauseNs[(ms.NumGC+255)%256]

	pid := int32(os.Getpid())
	out.ProcessSelf.PID = pid
	if self, err := process.NewProcessWithContext(ctx, pid); err == nil {
		out.ProcessSelf.Name, _ = self.NameWithContext(ctx)
		out.ProcessSelf.CPUPercent, _ = self.CPUPercentWithContext(ctx)
		if mi, err := self.MemoryInfoWithContext(ctx); err == nil && mi != nil {
			out.ProcessSelf.RSSBytes = mi.RSS
		}
	}

	procs, err := process.ProcessesWithContext(ctx)
	if err == nil && len(procs) > 0 {
		type row struct {
			p   *process.Process
			rss uint64
		}
		rows := make([]row, 0, len(procs))
		for _, p := range procs {
			if ctx.Err() != nil {
				break
			}
			mi, err := p.MemoryInfoWithContext(ctx)
			if err != nil || mi == nil {
				continue
			}
			rows = append(rows, row{p: p, rss: mi.RSS})
		}
		sort.Slice(rows, func(i, j int) bool { return rows[i].rss > rows[j].rss })
		topN := minInt(8, len(rows))
		out.TopProcessesByRSS = make([]models.TopProcessMetric, 0, topN)
		for i := 0; i < topN; i++ {
			p := rows[i].p
			name, _ := p.NameWithContext(ctx)
			cpuP, _ := p.CPUPercentWithContext(ctx)
			out.TopProcessesByRSS = append(out.TopProcessesByRSS, models.TopProcessMetric{
				PID:        p.Pid,
				Name:       name,
				RSSBytes:   rows[i].rss,
				CPUPercent: cpuP,
			})
		}
	}

	maybePersistSystemSnapshot(out)
	return out, "", 200
}

func maybePersistSystemSnapshot(p models.SystemMetricsPayload) {
	systemMetricsSnapshotMu.Lock()
	defer systemMetricsSnapshotMu.Unlock()
	if time.Since(systemMetricsLastInsert) < systemMetricsInsertMinGap {
		return
	}
	systemMetricsLastInsert = time.Now()

	snap := models.SystemMetricSnapshot{
		CreatedAt:       p.CollectedAt,
		CpuPercent:      p.CpuPercent,
		MemUsedPercent:  p.Memory.UsedPercent,
		DiskUsedPercent: p.Disk.UsedPercent,
		Goroutines:      p.GoRuntime.Goroutines,
		HeapAllocBytes:  p.GoRuntime.HeapAllocBytes,
	}
	if err := db.InsertSystemMetricSnapshot(snap); err != nil {
		log.Println("InsertSystemMetricSnapshot:", err)
	}
}
