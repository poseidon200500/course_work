package main

import (
	"flag"
	"fmt"
	"os"

	analysis "github.com/poseidon200500/course_work/internal/analisys"
	"github.com/poseidon200500/course_work/internal/benchmark"
	"github.com/poseidon200500/course_work/internal/generator"
	"github.com/poseidon200500/course_work/internal/storage"
)

func main() {
	// ===== FLAGS =====
	var (
		name       = flag.String("name", "CUSTOM", "scenario name")
		total      = flag.Int("total", 1_000_000, "total strings")
		dup        = flag.Int("dup", 40, "duplicate percent")
		maxLen     = flag.Int("maxlen", 8, "max string length")
		dist       = flag.String("dist", "uniform", "distribution: uniform|zipf")
		storageArg = flag.String("storage", "all", "storage: base|intern|v1|v2|all")
		output     = flag.String("out", "result.csv", "output csv file")
	)

	flag.Parse()

	// ===== DISTRIBUTION =====
	var distribution generator.DistributionType
	switch *dist {
	case "uniform":
		distribution = generator.DistributionUniform
	case "zipf":
		distribution = generator.DistributionZipf
	default:
		panic("unknown distribution")
	}

	// ===== SCENARIO =====
	sc := benchmark.Scenario{
		Name:             *name,
		Group:            "custom",
		Description:      "CLI scenario",
		Total:            *total,
		WordsPerLine:     10,
		DuplicatePercent: *dup,
		MaxLen:           *maxLen,
		Distribution:     distribution,
	}

	// ===== STORAGE FACTORIES =====
	allStorages := map[string]func() storage.Storage{
		"BASE":      func() storage.Storage { return storage.NewBaseStorage() },
		"INTERN":    func() storage.Storage { return storage.NewInternStorage() },
		"UNIQUE_V1": func() storage.Storage { return storage.NewUniqueStorageV1() },
		"UNIQUE_V2": func() storage.Storage { return storage.NewUniqueStorageV2() },
	}

	storages := make(map[string]func() storage.Storage)

	switch *storageArg {
	case "base":
		storages["BASE"] = allStorages["BASE"]
	case "intern":
		storages["INTERN"] = allStorages["INTERN"]
	case "v1":
		storages["UNIQUE_V1"] = allStorages["UNIQUE_V1"]
	case "v2":
		storages["UNIQUE_V2"] = allStorages["UNIQUE_V2"]
	case "all":
		storages = allStorages
	default:
		panic("unknown storage")
	}

	// ===== DATA PREP =====
	fmt.Println("Preparing data...")
	if err := benchmark.EnsureScenarioData(sc); err != nil {
		panic(err)
	}

	filename := benchmark.ScenarioFullPath(sc)

	// ===== RUN =====
	fmt.Println("Running benchmark...")
	var results []benchmark.Result

	for name, factory := range storages {
		st := factory()

		res, err := benchmark.RunSingle(name, sc, st, filename)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}

		results = append(results, res)
	}

	benchmark.PrintResults(results)

	// ===== SAVE CSV =====
	if err := analysis.WriteResultsCSV(results, *output); err != nil {
		panic(err)
	}

	fmt.Println("Saved to:", *output)

	// важно для GC (чтобы не держать память)
	os.Exit(0)
}
