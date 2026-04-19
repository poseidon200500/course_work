package benchmark

import (
	"encoding/json"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/poseidon200500/course_work/internal/storage"
)

const (
	gcCycles = 5
	gcPause  = 50 * time.Millisecond
)

func measureSerialization(st storage.Storage) (time.Duration, time.Duration, int, error) {
	startMaterialize := time.Now()
	data := st.GetAll()
	materializationDuration := time.Since(startMaterialize)

	startSerialization := time.Now()
	raw, err := json.Marshal(data)
	serializationDuration := time.Since(startSerialization)
	if err != nil {
		return 0, 0, 0, err
	}

	return materializationDuration, serializationDuration, len(raw), nil
}

func toMB(v uint64) float64 {
	return float64(v) / 1024.0 / 1024.0
}

func readMemPoint() MemPoint {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return MemPoint{
		HeapAllocMB:  toMB(m.HeapAlloc),
		HeapInuseMB:  toMB(m.HeapInuse),
		TotalAllocMB: toMB(m.TotalAlloc),
		Mallocs:      m.Mallocs,
		NumGC:        m.NumGC,
	}
}

// stabilizeMemoryState aggressively runs GC several times and asks runtime
// to return as much memory to OS as possible.
func stabilizeMemoryState() {
	for i := 0; i < gcCycles; i++ {
		runtime.GC()
		debug.FreeOSMemory()
		time.Sleep(gcPause)
	}
}

func forceGCAndRead() MemPoint {
	stabilizeMemoryState()
	return readMemPoint()
}
