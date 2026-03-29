package main

import (
	"github.com/poseidon200500/course_work/internal/benchmark"
	"github.com/poseidon200500/course_work/internal/storage"
)

func main() {

	storages := map[string]func() storage.Storage{
		"BASE":      func() storage.Storage { return storage.NewBaseStorage() },
		"INTERN":    func() storage.Storage { return storage.NewInternStorage() },
		"UNIQUE_V1": func() storage.Storage { return storage.NewUniqueStorageV1() },
		"UNIQUE_V2": func() storage.Storage { return storage.NewUniqueStorageV2() },
	}

	benchmark.RunAll(storages)
}
