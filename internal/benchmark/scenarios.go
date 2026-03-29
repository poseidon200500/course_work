package benchmark

import (
	"fmt"

	"github.com/poseidon200500/course_work/internal/generator"
	"github.com/poseidon200500/course_work/internal/parser"
	"github.com/poseidon200500/course_work/storage"
)

type Scenario struct {
	Name             string
	Total            int
	WordsPerLine     int
	DuplicatePercent int
	MaxLen           int
}

func DefaultScenarios() []Scenario {
	return []Scenario{
		{
			Name:             "LOW_DUP",
			Total:            1_000_000,
			WordsPerLine:     10,
			DuplicatePercent: 10,
			MaxLen:           8,
		},
		{
			Name:             "MEDIUM_DUP",
			Total:            1_000_000,
			WordsPerLine:     10,
			DuplicatePercent: 40,
			MaxLen:           8,
		},
		{
			Name:             "HIGH_DUP",
			Total:            1_000_000,
			WordsPerLine:     10,
			DuplicatePercent: 80,
			MaxLen:           8,
		},
		{
			Name:             "SHORT_STRINGS",
			Total:            1_000_000,
			WordsPerLine:     10,
			DuplicatePercent: 40,
			MaxLen:           4,
		},
		{
			Name:             "LONG_STRINGS",
			Total:            1_000_000,
			WordsPerLine:     10,
			DuplicatePercent: 40,
			MaxLen:           8,
		},
	}
}

func RunAll(
	storages map[string]func() storage.Storage,
) {

	scenarios := DefaultScenarios()
	var results []Result

	for _, sc := range scenarios {
		fmt.Println("\n=== Scenario:", sc.Name)

		filename := "data_" + sc.Name + ".txt"

		err := generator.GenerateData(
			sc.Total,
			sc.WordsPerLine,
			sc.DuplicatePercent,
			sc.MaxLen,
			filename,
		)
		if err != nil {
			panic(err)
		}

		data, err := parser.ParseData(filename)
		if err != nil {
			panic(err)
		}

		for name, factory := range storages {
			st := factory()

			res := RunSingle(name, sc.Name, st, data)
			results = append(results, res)
		}
	}

	PrintResults(results)
}
