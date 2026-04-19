package benchmark

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/poseidon200500/course_work/internal/parser"
	"github.com/poseidon200500/course_work/internal/storage"
)

type MemPoint struct {
	HeapAllocMB  float64
	HeapInuseMB  float64
	TotalAllocMB float64
	Mallocs      uint64
	NumGC        uint32
}

type Result struct {
	Name         string
	Scenario     string
	Group        string
	Description  string
	Distribution string

	Duration      time.Duration
	TotalInserted int
	UniqueCount   int

	BeforeLoad  MemPoint
	AfterInsert MemPoint
	AfterGC     MemPoint

	// Дельты оставляем, но интерпретировать их нужно осторожно.
	HeapAllocDeltaAfterInsertMB float64
	HeapAllocDeltaAfterGCMB     float64
	HeapInuseDeltaAfterInsertMB float64
	HeapInuseDeltaAfterGCMB     float64
	TotalAllocDeltaMB           float64
	MallocsDelta                uint64
	NumGCDelta                  uint32

	// Более полезные производные метрики.
	RetainedHeapAllocMB float64
	RetainedHeapInuseMB float64
	BytesPerInserted    float64
	BytesPerUnique      float64

	MaterializationDuration time.Duration
	SerializationDuration   time.Duration
	SerializedBytes         int
}

// RunSingle performs a benchmark run for one storage and one scenario.
// This version tries to isolate runs better by aggressively stabilizing
// GC state before and after the main phases.
func RunSingle(
	storageName string,
	sc Scenario,
	st storage.Storage,
	filename string,
) (Result, error) {
	st.Reset()

	// Максимально стабилизируем память до начала.
	before := forceGCAndRead()

	// На время загрузки можно временно отключить авто-GC,
	// чтобы меньше шумел during-insert.
	oldGCPercent := debug.SetGCPercent(-1)
	defer debug.SetGCPercent(oldGCPercent)

	start := time.Now()

	err := parser.ParseDataStream(filename, func(s string) {
		// Важно: создаём отдельную копию строки,
		// чтобы моделировать "новую входную строку".
		st.Add(strings.Clone(s))
	})
	if err != nil {
		return Result{}, fmt.Errorf("parse stream failed: %w", err)
	}

	duration := time.Since(start)

	afterInsert := readMemPoint()

	// Возвращаем GC обратно
	debug.SetGCPercent(oldGCPercent)

	// Чистим всё, кроме того, что удерживается storage.
	afterGC := forceGCAndRead()

	// Явно удерживаем storage живым до этого момента.
	runtime.KeepAlive(st)

	stats := st.Stats()

	materializationDuration, serializationDuration, serializedBytes, err := measureSerialization(st)
	if err != nil {
		return Result{}, err
	}

	retainedHeapAllocMB := afterGC.HeapAllocMB
	retainedHeapInuseMB := afterGC.HeapInuseMB

	var bytesPerInserted float64
	if stats.TotalInserted > 0 {
		bytesPerInserted = retainedHeapAllocMB * 1024 * 1024 / float64(stats.TotalInserted)
	}

	var bytesPerUnique float64
	if stats.UniqueCount > 0 {
		bytesPerUnique = retainedHeapAllocMB * 1024 * 1024 / float64(stats.UniqueCount)
	}

	result := Result{
		Name:         storageName,
		Scenario:     sc.Name,
		Group:        sc.Group,
		Description:  sc.Description,
		Distribution: string(sc.Distribution),

		Duration:      duration,
		TotalInserted: stats.TotalInserted,
		UniqueCount:   stats.UniqueCount,

		BeforeLoad:  before,
		AfterInsert: afterInsert,
		AfterGC:     afterGC,

		HeapAllocDeltaAfterInsertMB: afterInsert.HeapAllocMB - before.HeapAllocMB,
		HeapAllocDeltaAfterGCMB:     afterGC.HeapAllocMB - before.HeapAllocMB,
		HeapInuseDeltaAfterInsertMB: afterInsert.HeapInuseMB - before.HeapInuseMB,
		HeapInuseDeltaAfterGCMB:     afterGC.HeapInuseMB - before.HeapInuseMB,
		TotalAllocDeltaMB:           afterGC.TotalAllocMB - before.TotalAllocMB,
		MallocsDelta:                afterGC.Mallocs - before.Mallocs,
		NumGCDelta:                  afterGC.NumGC - before.NumGC,

		RetainedHeapAllocMB: retainedHeapAllocMB,
		RetainedHeapInuseMB: retainedHeapInuseMB,
		BytesPerInserted:    bytesPerInserted,
		BytesPerUnique:      bytesPerUnique,

		MaterializationDuration: materializationDuration,
		SerializationDuration:   serializationDuration,
		SerializedBytes:         serializedBytes,
	}

	return result, nil
}

func PrintResults(results []Result) {
	fmt.Println("\n===== RESULTS =====")

	for _, r := range results {
		fmt.Printf(`
[%s | %s]
Group:             %s
Description:       %s
Distribution:      %s

Insert time:       %v
Inserted:          %d
Unique:            %d

BeforeLoad:
  HeapAlloc:       %.2f MB
  HeapInuse:       %.2f MB
  TotalAlloc:      %.2f MB
  Mallocs:         %d
  NumGC:           %d

AfterInsert:
  HeapAlloc:       %.2f MB
  HeapInuse:       %.2f MB
  TotalAlloc:      %.2f MB
  Mallocs:         %d
  NumGC:           %d

AfterGC:
  HeapAlloc:       %.2f MB
  HeapInuse:       %.2f MB
  TotalAlloc:      %.2f MB
  Mallocs:         %d
  NumGC:           %d

Deltas:
  HeapAlloc insert: %.2f MB
  HeapAlloc gc:     %.2f MB
  HeapInuse insert: %.2f MB
  HeapInuse gc:     %.2f MB
  TotalAlloc:       %.2f MB
  Mallocs:          %d
  NumGC:            %d

Retained:
  HeapAlloc:        %.2f MB
  HeapInuse:        %.2f MB
  Bytes/inserted:   %.2f
  Bytes/unique:     %.2f

Serialization:
  To strings:       %v
  Serialize:        %v
  Bytes:            %d
`,
			r.Name, r.Scenario,
			r.Group,
			r.Description,
			r.Distribution,

			r.Duration,
			r.TotalInserted,
			r.UniqueCount,

			r.BeforeLoad.HeapAllocMB,
			r.BeforeLoad.HeapInuseMB,
			r.BeforeLoad.TotalAllocMB,
			r.BeforeLoad.Mallocs,
			r.BeforeLoad.NumGC,

			r.AfterInsert.HeapAllocMB,
			r.AfterInsert.HeapInuseMB,
			r.AfterInsert.TotalAllocMB,
			r.AfterInsert.Mallocs,
			r.AfterInsert.NumGC,

			r.AfterGC.HeapAllocMB,
			r.AfterGC.HeapInuseMB,
			r.AfterGC.TotalAllocMB,
			r.AfterGC.Mallocs,
			r.AfterGC.NumGC,

			r.HeapAllocDeltaAfterInsertMB,
			r.HeapAllocDeltaAfterGCMB,
			r.HeapInuseDeltaAfterInsertMB,
			r.HeapInuseDeltaAfterGCMB,
			r.TotalAllocDeltaMB,
			r.MallocsDelta,
			r.NumGCDelta,

			r.RetainedHeapAllocMB,
			r.RetainedHeapInuseMB,
			r.BytesPerInserted,
			r.BytesPerUnique,

			r.MaterializationDuration,
			r.SerializationDuration,
			r.SerializedBytes,
		)
	}
}
