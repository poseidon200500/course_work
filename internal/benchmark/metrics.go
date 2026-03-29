package benchmark

import (
	"fmt"
	"runtime"
	"time"

	"github.com/poseidon200500/course_work/storage"
)

type Result struct {
	Name          string
	Scenario      string
	Duration      time.Duration
	TotalInserted int
	UniqueCount   int
	AllocMB       float64
	TotalAllocMB  float64
	NumGC         uint32
}

func measureMemory() runtime.MemStats {
	var m runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m)
	return m
}

func RunSingle(
	name string,
	scenario string,
	st storage.Storage,
	data []string,
) Result {

	startMem := measureMemory()
	start := time.Now()

	for _, s := range data {
		st.Add(s)
	}

	duration := time.Since(start)
	endMem := measureMemory()

	stats := st.Stats()

	return Result{
		Name:          name,
		Scenario:      scenario,
		Duration:      duration,
		TotalInserted: stats.TotalInserted,
		UniqueCount:   stats.UniqueCount,
		AllocMB:       float64(endMem.Alloc-startMem.Alloc) / 1024 / 1024,
		TotalAllocMB:  float64(endMem.TotalAlloc-startMem.TotalAlloc) / 1024 / 1024,
		NumGC:         endMem.NumGC - startMem.NumGC,
	}
}

func PrintResults(results []Result) {
	fmt.Println("\n===== RESULTS =====")

	for _, r := range results {
		fmt.Printf(`
[%s | %s]
Time:           %v
Inserted:       %d
Unique:         %d
Alloc:          %.2f MB
TotalAlloc:     %.2f MB
GC runs:        %d
`, r.Name, r.Scenario, r.Duration, r.TotalInserted, r.UniqueCount, r.AllocMB, r.TotalAllocMB, r.NumGC)
	}
}
