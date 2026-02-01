package monitor

// import (
// 	"context"
// 	"os"
// 	"time"

// 	"github.com/shirou/gopsutil/v3/cpu"
// 	"github.com/shirou/gopsutil/v3/disk"
// 	"github.com/shirou/gopsutil/v3/mem"
// 	"github.com/shirou/gopsutil/v3/net"
// 	"github.com/shirou/gopsutil/v3/process"
// )

// type ResourceMetrics struct {
// 	CPUUsagePercent   []float64               `json:"cpu_usage_percent"`
// 	Ram               *mem.VirtualMemoryStat  `json:"ram_usage"`
// 	DiskUsage         *disk.UsageStat         `json:"disk_usage"`
// 	Network           []net.IOCountersStat    `json:"network"`
// 	ProcessCPUPercent float64                 `json:"process_cpu_percent"`
// 	ProcessMemoryMB   *process.MemoryInfoStat `json:"process_memory_mb"`
// 	Errors            []error                 `json:"errors"`
// }

// type MonirtorService struct {
// 	callBack func(*ResourceMetrics)
// 	cancel   context.CancelFunc
// }

// type ResourceHandlerI interface {
// 	Start(interval time.Duration, callBack func(*ResourceMetrics))
// 	Stop()
// }

// func NewMonitorService() *MonirtorService {
// 	return &MonirtorService{}
// }

// func (M *MonirtorService) Start(interval time.Duration) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	M.cancel = cancel
// 	ticker := time.NewTicker(interval)
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 		case <-ticker.C:
// 			RM := fetchMetrics()
// 			M.callBack(RM)
// 		}
// 	}
// }

// func (M *MonirtorService) Stop() {
// 	M.cancel()
// }

// // no need to export
// //
// //	FetchMetrics collects current system and process metrics.
// func fetchMetrics() *ResourceMetrics {
// 	var R ResourceMetrics
// 	addErr := func(err error) {
// 		if err != nil {
// 			R.Errors = append(R.Errors, err)
// 		}
// 	}
// 	var err error
// 	R.CPUUsagePercent, err = cpu.Percent(0, false)
// 	addErr(err)

// 	R.Ram, err = mem.VirtualMemory()
// 	addErr(err)

// 	R.DiskUsage, err = disk.Usage("/")
// 	addErr(err)

// 	R.Network, err = net.IOCounters(false)
// 	addErr(err)

// 	proc, err := process.NewProcess(int32(os.Getpid()))
// 	addErr(err)

// 	if proc != nil {
// 		R.ProcessCPUPercent, err = proc.CPUPercent()
// 		addErr(err)

// 		R.ProcessMemoryMB, err = proc.MemoryInfo()
// 		addErr(err)
// 	}
// 	return &R
// }
