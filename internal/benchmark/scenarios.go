package benchmark

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/poseidon200500/course_work/internal/generator"
	"github.com/poseidon200500/course_work/internal/parser"
	"github.com/poseidon200500/course_work/internal/storage"
)

const (
	defaultSeed  = int64(42)
	defaultZipfS = 1.2
	defaultZipfV = 1.0
)

type Scenario struct {
	Name             string
	Total            int
	WordsPerLine     int
	DuplicatePercent int
	MaxLen           int
	Group            string
	Distribution     generator.DistributionType
	Description      string
}

func DefaultScenarios() []Scenario {
	return []Scenario{
		// Группа 1. Влияние процента дубликатов (равномерное распределение)
		{
			Name:             "DUP_10_UNIFORM",
			Group:            "duplicate_ratio",
			Description:      "Низкая доля дубликатов, равномерное распределение",
			Total:            1_000_000,
			WordsPerLine:     10,
			DuplicatePercent: 10,
			MaxLen:           8,
			Distribution:     generator.DistributionUniform,
		},
		{
			Name:             "DUP_40_UNIFORM",
			Group:            "duplicate_ratio",
			Description:      "Средняя доля дубликатов, равномерное распределение",
			Total:            1_000_000,
			WordsPerLine:     10,
			DuplicatePercent: 40,
			MaxLen:           8,
			Distribution:     generator.DistributionUniform,
		},
		{
			Name:             "DUP_80_UNIFORM",
			Group:            "duplicate_ratio",
			Description:      "Высокая доля дубликатов, равномерное распределение",
			Total:            1_000_000,
			WordsPerLine:     10,
			DuplicatePercent: 80,
			MaxLen:           8,
			Distribution:     generator.DistributionUniform,
		},

		// Группа 2. Влияние длины строк (равномерное распределение)
		{
			Name:             "LEN_4_UNIFORM",
			Group:            "string_length",
			Description:      "Короткие строки, равномерное распределение",
			Total:            10_000,
			WordsPerLine:     10,
			DuplicatePercent: 40,
			MaxLen:           4,
			Distribution:     generator.DistributionUniform,
		},
		{
			Name:             "LEN_8_UNIFORM",
			Group:            "string_length",
			Description:      "Средние строки, равномерное распределение",
			Total:            1_000_000,
			WordsPerLine:     10,
			DuplicatePercent: 40,
			MaxLen:           8,
			Distribution:     generator.DistributionUniform,
		},
		{
			Name:             "LEN_12_UNIFORM",
			Group:            "string_length",
			Description:      "Более длинные строки, равномерное распределение",
			Total:            1_000_000,
			WordsPerLine:     10,
			DuplicatePercent: 40,
			MaxLen:           12,
			Distribution:     generator.DistributionUniform,
		},

		// Группа 3. Влияние размера хранилища (равномерное распределение)
		{
			Name:             "SIZE_100K_UNIFORM",
			Group:            "dataset_size",
			Description:      "Малое хранилище, равномерное распределение",
			Total:            100_000,
			WordsPerLine:     10,
			DuplicatePercent: 40,
			MaxLen:           8,
			Distribution:     generator.DistributionUniform,
		},
		{
			Name:             "SIZE_1M_UNIFORM",
			Group:            "dataset_size",
			Description:      "Среднее хранилище, равномерное распределение",
			Total:            1_000_000,
			WordsPerLine:     10,
			DuplicatePercent: 40,
			MaxLen:           8,
			Distribution:     generator.DistributionUniform,
		},
		{
			Name:             "SIZE_5M_UNIFORM",
			Group:            "dataset_size",
			Description:      "Большое хранилище, равномерное распределение",
			Total:            5_000_000,
			WordsPerLine:     10,
			DuplicatePercent: 40,
			MaxLen:           8,
			Distribution:     generator.DistributionUniform,
		},

		// Группа 4. Сравнение распределений при одном и том же наборе параметров
		{
			Name:             "DIST_UNIFORM_40",
			Group:            "distribution_type",
			Description:      "Равномерное распределение при 40% дубликатов",
			Total:            1_000_000,
			WordsPerLine:     10,
			DuplicatePercent: 40,
			MaxLen:           8,
			Distribution:     generator.DistributionUniform,
		},
		{
			Name:             "DIST_ZIPF_40",
			Group:            "distribution_type",
			Description:      "Распределение Ципфа при 40% дубликатов",
			Total:            1_000_000,
			WordsPerLine:     10,
			DuplicatePercent: 40,
			MaxLen:           8,
			Distribution:     generator.DistributionZipf,
		},
		// Группа 5. Сценарии, в которых unique имеет больше шансов показать преимущество
		{
			Name:             "UNIQUE_FRIENDLY_32_95_UNIFORM",
			Group:            "unique_friendly",
			Description:      "Длинные строки, 95% дубликатов, равномерное распределение",
			Total:            1_000_000,
			WordsPerLine:     10,
			DuplicatePercent: 95,
			MaxLen:           32,
			Distribution:     generator.DistributionUniform,
		},
		{
			Name:             "UNIQUE_FRIENDLY_64_95_UNIFORM",
			Group:            "unique_friendly",
			Description:      "Очень длинные строки, 95% дубликатов, равномерное распределение",
			Total:            1_000_000,
			WordsPerLine:     10,
			DuplicatePercent: 95,
			MaxLen:           64,
			Distribution:     generator.DistributionUniform,
		},
		{
			Name:             "UNIQUE_FRIENDLY_32_95_ZIPF",
			Group:            "unique_friendly",
			Description:      "Длинные строки, 95% дубликатов, распределение Ципфа",
			Total:            1_000_000,
			WordsPerLine:     10,
			DuplicatePercent: 95,
			MaxLen:           32,
			Distribution:     generator.DistributionZipf,
		},
		{
			Name:             "UNIQUE_FRIENDLY_64_99_ZIPF",
			Group:            "unique_friendly",
			Description:      "Очень длинные строки, 99% дубликатов, распределение Ципфа",
			Total:            1_000_000,
			WordsPerLine:     10,
			DuplicatePercent: 99,
			MaxLen:           64,
			Distribution:     generator.DistributionZipf,
		},
	}
}

func RunAll(storages map[string]func() storage.Storage) {
	scenarios := DefaultScenarios()
	var results []Result

	// очищаем папку перед запуском
	err := generator.ClearDatastore()
	if err != nil {
		panic(err)
	}

	for _, sc := range scenarios {
		fmt.Println("\n=== Scenario:", sc.Name)

		filename := fmt.Sprintf("data_%s.txt", sc.Name)

		cfg := generator.Config{
			Total:            sc.Total,
			WordsOnLine:      sc.WordsPerLine,
			DuplicatePercent: sc.DuplicatePercent,
			MaxLen:           sc.MaxLen,
			Filename:         filename,
			Deterministic:    true,
			Seed:             42,
			Distribution:     sc.Distribution,
			ZipfS:            1.2,
			ZipfV:            1.0,
		}

		fmt.Println("Generating file:", filename)
		err := generator.GenerateDataWithConfig(cfg)
		if err != nil {
			panic(err)
		}

		fullPath := filepath.Join(generator.DatastoreFolder, filename)

		for name, factory := range storages {
			st := factory()

			res, err := RunSingle(name, sc, st, fullPath)
			if err != nil {
				fmt.Printf("run single failed for %s/%s: %v\n", name, sc.Name, err)
				continue
			}

			results = append(results, res)
		}
	}

	PrintResults(results)
}

func RunSelected(storages map[string]func() storage.Storage, scenarios []Scenario) ([]Result, error) {
	if len(scenarios) == 0 {
		return nil, fmt.Errorf("нет сценариев для запуска")
	}

	var results []Result

	fmt.Println("\nПодготовка данных...")
	for _, sc := range scenarios {
		if err := EnsureScenarioData(sc); err != nil {
			return nil, fmt.Errorf("failed to ensure scenario data for %s: %w", sc.Name, err)
		}
	}

	fmt.Println("\nЗапуск бенчмарков...")
	for _, sc := range scenarios {
		fmt.Printf("\n=== Scenario: %s ===\n", sc.Name)

		filename := ScenarioFullPath(sc)

		for name, factory := range storages {
			st := factory()

			res, err := RunSingle(name, sc, st, filename)
			if err != nil {
				return nil, fmt.Errorf("failed to run single benchmark %s/%s: %w", name, sc.Name, err)
			}

			results = append(results, res)
		}
	}

	PrintResults(results)
	return results, nil
}

func QuickScenarios() []Scenario {
	return []Scenario{
		{
			Name:             "QUICK_DUP_40",
			Group:            "quick",
			Description:      "Быстрый локальный тест",
			Total:            50_000,
			WordsPerLine:     10,
			DuplicatePercent: 40,
			MaxLen:           8,
		},
	}
}

func ScenarioFilename(sc Scenario) string {
	return fmt.Sprintf("data_%s.txt", sc.Name)
}

func ScenarioFullPath(sc Scenario) string {
	return filepath.Join(generator.DatastoreFolder, ScenarioFilename(sc))
}

func ScenarioFileExists(sc Scenario) bool {
	_, err := os.Stat(ScenarioFullPath(sc))
	return err == nil
}

func EnsureScenarioData(sc Scenario) error {
	if ScenarioFileExists(sc) {
		return nil
	}

	cfg := generator.Config{
		Total:            sc.Total,
		WordsOnLine:      sc.WordsPerLine,
		DuplicatePercent: sc.DuplicatePercent,
		MaxLen:           sc.MaxLen,
		Filename:         ScenarioFilename(sc),
		Deterministic:    true,
		Seed:             defaultSeed,
		Distribution:     sc.Distribution,
		ZipfS:            defaultZipfS,
		ZipfV:            defaultZipfV,
	}

	fmt.Printf("Generating missing file for scenario %s...\n", sc.Name)
	return generator.GenerateDataWithConfig(cfg)
}

func LoadScenarioData(sc Scenario) ([]string, error) {
	return parser.ParseData(ScenarioFullPath(sc))
}

func SortedGroupNames(groups map[string][]Scenario) []string {
	order := []string{
		"duplicate_ratio",
		"string_length",
		"dataset_size",
		"distribution_type",
		"quick",
		"custom",
	}

	var result []string
	seen := make(map[string]struct{})

	for _, name := range order {
		if _, ok := groups[name]; ok {
			result = append(result, name)
			seen[name] = struct{}{}
		}
	}

	for name := range groups {
		if _, ok := seen[name]; !ok {
			result = append(result, name)
		}
	}

	return result
}

func FormatGroupName(group string) string {
	switch group {
	case "duplicate_ratio":
		return "Влияние процента дубликатов"
	case "string_length":
		return "Влияние длины строк"
	case "dataset_size":
		return "Влияние размера хранилища"
	case "distribution_type":
		return "Влияние типа распределения"
	case "quick":
		return "Быстрые тесты"
	case "custom":
		return "Кастомные сценарии"
	case "unique_friendly":
		return "Сценарии, благоприятные для unique"
	default:
		return group
	}
}

func GroupScenarios(scenarios []Scenario) map[string][]Scenario {
	result := make(map[string][]Scenario)
	for _, sc := range scenarios {
		result[sc.Group] = append(result[sc.Group], sc)
	}
	return result
}
