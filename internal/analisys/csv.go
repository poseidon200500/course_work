package analysis

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/poseidon200500/course_work/internal/benchmark"
)

func WriteResultsCSV(results []benchmark.Result, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"storage",
		"scenario",
		"group",
		"description",
		"distribution",

		"insert_duration_ms",
		"materialization_duration_ms",
		"serialization_duration_ms",
		"serialized_bytes",

		"total_inserted",
		"unique_count",

		"before_heap_alloc_mb",
		"before_heap_inuse_mb",
		"before_total_alloc_mb",
		"before_mallocs",
		"before_num_gc",

		"after_insert_heap_alloc_mb",
		"after_insert_heap_inuse_mb",
		"after_insert_total_alloc_mb",
		"after_insert_mallocs",
		"after_insert_num_gc",

		"after_gc_heap_alloc_mb",
		"after_gc_heap_inuse_mb",
		"after_gc_total_alloc_mb",
		"after_gc_mallocs",
		"after_gc_num_gc",

		"heap_alloc_delta_after_insert_mb",
		"heap_alloc_delta_after_gc_mb",
		"heap_inuse_delta_after_insert_mb",
		"heap_inuse_delta_after_gc_mb",
		"total_alloc_delta_mb",
		"mallocs_delta",
		"num_gc_delta",

		"retained_heap_alloc_mb",
		"retained_heap_inuse_mb",
		"bytes_per_inserted",
		"bytes_per_unique",
	}

	if err := writer.Write(header); err != nil {
		return err
	}

	for _, r := range results {
		row := []string{
			r.Name,
			r.Scenario,
			r.Group,
			r.Description,
			r.Distribution,

			formatDurationMs(r.Duration),
			formatDurationMs(r.MaterializationDuration),
			formatDurationMs(r.SerializationDuration),
			strconv.Itoa(r.SerializedBytes),

			strconv.Itoa(r.TotalInserted),
			strconv.Itoa(r.UniqueCount),

			formatFloat(r.BeforeLoad.HeapAllocMB),
			formatFloat(r.BeforeLoad.HeapInuseMB),
			formatFloat(r.BeforeLoad.TotalAllocMB),
			strconv.FormatUint(r.BeforeLoad.Mallocs, 10),
			strconv.FormatUint(uint64(r.BeforeLoad.NumGC), 10),

			formatFloat(r.AfterInsert.HeapAllocMB),
			formatFloat(r.AfterInsert.HeapInuseMB),
			formatFloat(r.AfterInsert.TotalAllocMB),
			strconv.FormatUint(r.AfterInsert.Mallocs, 10),
			strconv.FormatUint(uint64(r.AfterInsert.NumGC), 10),

			formatFloat(r.AfterGC.HeapAllocMB),
			formatFloat(r.AfterGC.HeapInuseMB),
			formatFloat(r.AfterGC.TotalAllocMB),
			strconv.FormatUint(r.AfterGC.Mallocs, 10),
			strconv.FormatUint(uint64(r.AfterGC.NumGC), 10),

			formatFloat(r.HeapAllocDeltaAfterInsertMB),
			formatFloat(r.HeapAllocDeltaAfterGCMB),
			formatFloat(r.HeapInuseDeltaAfterInsertMB),
			formatFloat(r.HeapInuseDeltaAfterGCMB),
			formatFloat(r.TotalAllocDeltaMB),
			strconv.FormatUint(r.MallocsDelta, 10),
			strconv.FormatUint(uint64(r.NumGCDelta), 10),

			formatFloat(r.RetainedHeapAllocMB),
			formatFloat(r.RetainedHeapInuseMB),
			formatFloat(r.BytesPerInserted),
			formatFloat(r.BytesPerUnique),
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	if err := writer.Error(); err != nil {
		return fmt.Errorf("csv writer error: %w", err)
	}

	return nil
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64)
}

func formatDurationMs(d interface{ Microseconds() int64 }) string {
	ms := float64(d.Microseconds()) / 1000.0
	return strconv.FormatFloat(ms, 'f', 6, 64)
}
