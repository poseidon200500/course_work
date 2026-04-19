package generator

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

var (
	DatastoreFolder = "datastore/"
	vowels          = []rune("aeiou")
	consonants      = []rune("bcdfghjklmnpqrstvwxyz")
)

type DistributionType string

const (
	DistributionUniform DistributionType = "uniform"
	DistributionZipf    DistributionType = "zipf"
)

type Config struct {
	Total            int
	WordsOnLine      int
	DuplicatePercent int
	MaxLen           int
	Filename         string

	Seed          int64
	Deterministic bool
	Distribution  DistributionType

	// параметры Zipf
	ZipfS float64
	ZipfV float64
}

func newRand(cfg Config) *rand.Rand {
	if cfg.Deterministic {
		return rand.New(rand.NewSource(cfg.Seed))
	}

	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

// Генерация псевдослов
func randomWord(r *rand.Rand, maxLen int) string {
	if maxLen < 2 {
		return ""
	}

	length := r.Intn(maxLen-1) + 2

	runes := make([]rune, 0, length)

	for len(runes) < length {
		runes = append(runes, consonants[r.Intn(len(consonants))])
		if len(runes) < length {
			runes = append(runes, vowels[r.Intn(len(vowels))])
		}
	}

	return string(runes[:length])
}

func GenerateDataWithConfig(cfg Config) error {
	if cfg.Total <= 0 {
		return fmt.Errorf("total must be > 0")
	}
	if cfg.WordsOnLine <= 0 {
		return fmt.Errorf("wordsOnLine must be > 0")
	}
	if cfg.DuplicatePercent < 0 || cfg.DuplicatePercent > 100 {
		return fmt.Errorf("duplicatePercent must be in [0, 100]")
	}
	if cfg.MaxLen < 2 {
		return fmt.Errorf("maxLen must be >= 2")
	}
	if cfg.Filename == "" {
		return fmt.Errorf("filename is empty")
	}
	if cfg.Distribution == "" {
		cfg.Distribution = DistributionUniform
	}

	r := newRand(cfg)

	dupCount := cfg.Total * cfg.DuplicatePercent / 100
	uniqueCount := cfg.Total - dupCount

	uniqueSet := make(map[string]struct{}, uniqueCount)
	uniqueList := make([]string, 0, uniqueCount)

	for len(uniqueList) < uniqueCount {
		s := randomWord(r, cfg.MaxLen)
		if _, exists := uniqueSet[s]; !exists {
			uniqueSet[s] = struct{}{}
			uniqueList = append(uniqueList, s)
		}
	}

	result := make([]string, 0, cfg.Total)
	result = append(result, uniqueList...)

	switch cfg.Distribution {
	case DistributionUniform:
		fillDuplicatesUniform(r, &result, uniqueList, dupCount)
	case DistributionZipf:
		if err := fillDuplicatesZipf(r, &result, uniqueList, dupCount, cfg.ZipfS, cfg.ZipfV); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown distribution: %s", cfg.Distribution)
	}

	r.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	file, err := os.Create(DatastoreFolder + cfg.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for i, s := range result {
		_, err = writer.WriteString(s)
		if err != nil {
			return err
		}

		isLast := i == len(result)-1
		endOfLine := (i+1)%cfg.WordsOnLine == 0

		switch {
		case isLast:
		case endOfLine:
			_, err = writer.WriteString("\n")
		default:
			_, err = writer.WriteString(",")
		}

		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func fillDuplicatesUniform(r *rand.Rand, result *[]string, uniqueList []string, dupCount int) {
	for i := 0; i < dupCount; i++ {
		s := uniqueList[r.Intn(len(uniqueList))]
		*result = append(*result, s)
	}
}

func fillDuplicatesZipf(r *rand.Rand, result *[]string, uniqueList []string, dupCount int, s, v float64) error {
	if len(uniqueList) == 0 {
		return fmt.Errorf("uniqueList is empty")
	}
	if s <= 1.0 {
		return fmt.Errorf("zipf parameter s must be > 1.0")
	}
	if v < 1.0 {
		return fmt.Errorf("zipf parameter v must be >= 1.0")
	}

	imax := uint64(len(uniqueList) - 1)
	zipf := rand.NewZipf(r, s, v, imax)
	if zipf == nil {
		return fmt.Errorf("failed to create zipf generator")
	}

	for i := 0; i < dupCount; i++ {
		idx := int(zipf.Uint64())
		if idx < 0 || idx >= len(uniqueList) {
			return fmt.Errorf("zipf index out of range: %d", idx)
		}
		*result = append(*result, uniqueList[idx])
	}

	return nil
}

// Очистка папки DatastoreFolder
func ClearDatastore() error {
	// если папки нет — создаём
	if _, err := os.Stat(DatastoreFolder); os.IsNotExist(err) {
		return os.MkdirAll(DatastoreFolder, os.ModePerm)
	}

	entries, err := os.ReadDir(DatastoreFolder)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(DatastoreFolder, entry.Name())

		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	}

	return nil
}
