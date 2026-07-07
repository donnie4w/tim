package tc

import (
	"encoding/json"
	"github.com/donnie4w/tim/sys"
	"github.com/donnie4w/tlnet"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"runtime"
	"runtime/metrics"
	"runtime/pprof"
	"time"
)

// SystemMonitor 定义监控信息结构体（多层嵌套）
type SystemMonitor struct {
	Disk        DiskInfo    `json:"disk"`
	CPU         CPUInfo     `json:"cpu"`
	LoadAvg     LoadAvgInfo `json:"load_avg"`
	Memory      MemoryInfo  `json:"memory"`
	GoRuntime   *MetricData `json:"go_runtime"`
	CollectTime string      `json:"collect_time"`
	TaskNumber  TaskNumber  `json:"task_number"`
}

type DiskInfo struct {
	MountPoint  string  `json:"mount_point"`
	TotalGB     float64 `json:"total_gb"`
	UsedGB      float64 `json:"used_gb"`
	FreeGB      float64 `json:"free_gb"`
	UsedPercent float64 `json:"used_percent"`
	Error       string  `json:"error"`
}

type CPUInfo struct {
	ModelName    string  `json:"model_name"`
	LogicalCores int     `json:"logical_cores"`
	UsagePercent float64 `json:"usage_percent"`
	Error        string  `json:"error"`
}

type LoadAvgInfo struct {
	Load1   float64 `json:"load_1"`
	Load5   float64 `json:"load_5"`
	Load15  float64 `json:"load_15"`
	Usage1  float64 `json:"usage_1"`
	Usage5  float64 `json:"usage_5"`
	Usage15 float64 `json:"usage_15"`
	Error   string  `json:"error"`
}

type MemoryInfo struct {
	TotalGB     float64 `json:"total_gb"`
	UsedGB      float64 `json:"used_gb"`
	AvailableGB float64 `json:"available_gb"`
	UsedPercent float64 `json:"used_percent"`
	Error       string  `json:"error"`
}

type GoRuntimeInfo struct {
	Goroutines       int     `json:"goroutines"`
	Threads          int     `json:"threads"`
	GoVersion        string  `json:"go_version"`
	HeapAllocMB      float64 `json:"heap_alloc_mb"`
	TotalAllocMB     float64 `json:"total_alloc_mb"`
	GcCount          uint32  `json:"gc_count"`
	GcTotalMs        float64 `json:"gc_total_ms"`
	GcLastSecondsAgo float64 `json:"gc_last_seconds_ago"`
}

type TaskNumber struct {
	SdkTaskNum int64 `json:"sdk_task_num"`
}

func CollectSystemMonitor() *SystemMonitor {
	r := &SystemMonitor{
		Disk:        getDiskInfo(),
		CPU:         getCPUInfo(),
		LoadAvg:     getLoadAvgInfo(),
		Memory:      getMemInfo(),
		GoRuntime:   GetRuntimeMetrics(),
		CollectTime: time.Now().Format(time.DateTime),
	}
	r.TaskNumber = TaskNumber{SdkTaskNum: sys.Stat.Tx()}
	return r
}

func getDiskInfo() DiskInfo {
	mountPoint := "/"
	if runtime.GOOS == "windows" {
		mountPoint = "C:"
	}

	diskStat, err := disk.Usage(mountPoint)
	if err != nil {
		return DiskInfo{Error: err.Error()}
	}

	return DiskInfo{
		MountPoint:  mountPoint,
		TotalGB:     float64(diskStat.Total) / (1 << 30),
		UsedGB:      float64(diskStat.Used) / (1 << 30),
		FreeGB:      float64(diskStat.Free) / (1 << 30),
		UsedPercent: diskStat.UsedPercent,
	}
}

func getCPUInfo() CPUInfo {
	cpuNum, err := cpu.Counts(true)
	if err != nil {
		return CPUInfo{Error: err.Error()}
	}

	cpuPercent, err := cpu.Percent(1*time.Second, false)
	if err != nil {
		return CPUInfo{Error: err.Error()}
	}

	cpuInfo, err := cpu.Info()
	if err != nil || len(cpuInfo) == 0 {
		return CPUInfo{Error: "获取CPU型号失败"}
	}

	return CPUInfo{
		ModelName:    cpuInfo[0].ModelName,
		LogicalCores: cpuNum,
		UsagePercent: cpuPercent[0],
	}
}

func getLoadAvgInfo() LoadAvgInfo {
	loadAvg, err := load.Avg()
	if err != nil {
		return LoadAvgInfo{Error: err.Error()}
	}

	cpuCore, err := cpu.Counts(true)
	if err != nil {
		return LoadAvgInfo{Error: err.Error()}
	}

	return LoadAvgInfo{
		Load1:   loadAvg.Load1,
		Load5:   loadAvg.Load5,
		Load15:  loadAvg.Load15,
		Usage1:  (loadAvg.Load1 / float64(cpuCore)) * 100,
		Usage5:  (loadAvg.Load5 / float64(cpuCore)) * 100,
		Usage15: (loadAvg.Load15 / float64(cpuCore)) * 100,
	}
}

func getMemInfo() MemoryInfo {
	memStat, err := mem.VirtualMemory()
	if err != nil {
		return MemoryInfo{Error: err.Error()}
	}

	return MemoryInfo{
		TotalGB:     float64(memStat.Total) / (1 << 30),
		UsedGB:      float64(memStat.Used) / (1 << 30),
		AvailableGB: float64(memStat.Available) / (1 << 30),
		UsedPercent: memStat.UsedPercent,
	}
}

type MetricData struct {
	GoVersion  string `json:"go_version"`
	Goroutines int    `json:"goroutines"`
	Threads    int    `json:"threads"`

	HeapAllocMB float64 `json:"heap_alloc_mb"`
	HeapSysMB   float64 `json:"heap_sys_mb"`
	//TotalAllocMB float64 `json:"total_alloc_mb"`
	HeapObjects  uint64  `json:"heap_objects"`
	HeapGoalMB   float64 `json:"heap_goal_mb"`
	HeapStacksMB float64 `json:"heap_stacks_mb"`

	GcCountAutomatic uint64 `json:"gc_count_automatic"`
	GcCountForced    uint64 `json:"gc_count_forced"`
	GcCountTotal     uint64 `json:"gc_count_total"`
	//PauseTotalMs     float64 `json:"pause_total_ms"`

	AllocBytes   uint64 `json:"alloc_bytes"`
	AllocObjects uint64 `json:"alloc_objects"`
	FreeBytes    uint64 `json:"free_bytes"`
	FreeObjects  uint64 `json:"free_objects"`
	TinyAllocs   uint64 `json:"tiny_allocs"`

	MemFreeMB         float64 `json:"mem_free_mb"`
	MemReleasedMB     float64 `json:"mem_released_mb"`
	MemUnusedMB       float64 `json:"mem_unused_mb"`
	MetaMCacheInUseMB float64 `json:"meta_mcache_inuse_mb"`
	MetaMSpanInUseMB  float64 `json:"meta_mspan_inuse_mb"`
	MemOtherMB        float64 `json:"mem_other_mb"`

	SchedLatenciesMs float64 `json:"sched_latencies_ms"`
}

func GetRuntimeMetrics() *MetricData {
	data := &MetricData{}
	data.GoVersion = runtime.Version()
	data.Goroutines = runtime.NumGoroutine()
	data.Threads = pprof.Lookup("threadcreate").Count()

	sampleNames := []string{
		"/gc/cycles/automatic:gc-cycles",
		"/gc/cycles/forced:gc-cycles",
		"/gc/heap/allocs:bytes",
		"/gc/heap/allocs:objects",
		"/gc/heap/frees:bytes",
		"/gc/heap/frees:objects",
		"/gc/heap/goal:bytes",
		"/gc/heap/objects:objects",
		"/gc/heap/tiny/allocs:objects",
		"/memory/classes/heap/free:bytes",
		"/memory/classes/heap/released:bytes",
		"/memory/classes/heap/unused:bytes",
		"/memory/classes/heap/objects:bytes",
		"/memory/classes/heap/stacks:bytes",
		"/memory/classes/metadata/mcache/inuse:bytes",
		"/memory/classes/metadata/mspan/inuse:bytes",
		"/memory/classes/other:bytes",
		"/memory/classes/total:bytes",
		"/gc/pauses:seconds",
		"/sched/latencies:seconds",
	}

	samples := make([]metrics.Sample, len(sampleNames))
	for i, name := range sampleNames {
		samples[i].Name = name
	}

	metrics.Read(samples)

	getUint64 := func(name string) uint64 {
		for _, s := range samples {
			if s.Name == name {
				return s.Value.Uint64()
			}
		}
		return 0
	}

	getHistogram := func(name string) *metrics.Float64Histogram {
		for _, s := range samples {
			if s.Name == name {
				return s.Value.Float64Histogram()
			}
		}
		return nil
	}

	data.GcCountAutomatic = getUint64("/gc/cycles/automatic:gc-cycles")
	data.GcCountForced = getUint64("/gc/cycles/forced:gc-cycles")
	data.GcCountTotal = data.GcCountAutomatic + data.GcCountForced

	data.AllocBytes = getUint64("/gc/heap/allocs:bytes")
	data.AllocObjects = getUint64("/gc/heap/allocs:objects")
	data.FreeBytes = getUint64("/gc/heap/frees:bytes")
	data.FreeObjects = getUint64("/gc/heap/frees:objects")
	data.TinyAllocs = getUint64("/gc/heap/tiny/allocs:objects")

	data.HeapObjects = getUint64("/gc/heap/objects:objects")
	data.HeapGoalMB = float64(getUint64("/gc/heap/goal:bytes")) / 1024 / 1024
	data.HeapStacksMB = float64(getUint64("/memory/classes/heap/stacks:bytes")) / 1024 / 1024

	data.MemFreeMB = float64(getUint64("/memory/classes/heap/free:bytes")) / 1024 / 1024
	data.MemReleasedMB = float64(getUint64("/memory/classes/heap/released:bytes")) / 1024 / 1024
	data.MemUnusedMB = float64(getUint64("/memory/classes/heap/unused:bytes")) / 1024 / 1024
	data.MetaMCacheInUseMB = float64(getUint64("/memory/classes/metadata/mcache/inuse:bytes")) / 1024 / 1024
	data.MetaMSpanInUseMB = float64(getUint64("/memory/classes/metadata/mspan/inuse:bytes")) / 1024 / 1024
	data.MemOtherMB = float64(getUint64("/memory/classes/other:bytes")) / 1024 / 1024
	data.HeapSysMB = float64(getUint64("/memory/classes/total:bytes")) / 1024 / 1024

	heapObjects := float64(getUint64("/memory/classes/heap/objects:bytes"))
	heapUnused := float64(getUint64("/memory/classes/heap/unused:bytes"))
	if heapObjects > heapUnused {
		data.HeapAllocMB = float64(heapObjects-heapUnused) / 1024 / 1024
	}

	//var m runtime.MemStats
	//runtime.ReadMemStats(&m)

	//data.TotalAllocMB = float64(m.TotalAlloc) / 1024 / 1024
	//data.HeapAllocMB = float64(m.HeapAlloc) / 1024 / 1024
	//data.PauseTotalMs = float64(m.PauseTotalNs) / 1_000_000

	if h := getHistogram("/sched/latencies:seconds"); h != nil {
		data.SchedLatenciesMs = histogramMean(h) * 1000
	}

	return data
}

func histogramMean(h *metrics.Float64Histogram) float64 {
	if h == nil || len(h.Counts) == 0 || len(h.Buckets) < 2 {
		return 0
	}

	var totalCount uint64
	var weightedSum float64

	for i := range h.Counts {
		count := h.Counts[i]
		if count == 0 {
			continue
		}
		bucketValue := h.Buckets[i+1]
		weightedSum += float64(count) * bucketValue
		totalCount += count
	}

	if totalCount == 0 {
		return 0
	}
	return weightedSum / float64(totalCount)
}

func dashboardData(hc *tlnet.HttpContext) {
	if bean := CollectSystemMonitor(); bean != nil {
		jsonData, err := json.MarshalIndent(bean, "", "  ")
		if err == nil {
			hc.ResponseBytes(0, jsonData)
		}
	}
}
