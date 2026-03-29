package generator

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var (
	vowels          = []rune("aeiou")
	consonants      = []rune("bcdfghjklmnpqrstvwxyz")
	DatastoreFolder = "datastore/"
)

// Генерация псевдослов
func randomWord(maxLen int) string {
	length := rand.Intn(maxLen-1) + 2

	runes := make([]rune, 0, length)

	for len(runes) < length {
		runes = append(runes, consonants[rand.Intn(len(consonants))])
		if len(runes) < length {
			runes = append(runes, vowels[rand.Intn(len(vowels))])
		}
	}

	return string(runes[:length])
}

func GenerateData(total, wordsOnLine int, duplicatePercent int, maxLen int, filename string) error {
	rand.Seed(time.Now().UnixNano())

	dupCount := total * duplicatePercent / 100
	uniqueCount := total - dupCount

	fmt.Println("Total:", total)
	fmt.Println("Unique:", uniqueCount)
	fmt.Println("Duplicates:", dupCount)

	uniqueSet := make(map[string]struct{}, uniqueCount)
	uniqueList := make([]string, 0, uniqueCount)

	// Генерация уникальных слов
	for len(uniqueList) < uniqueCount {
		s := randomWord(maxLen)
		if _, exists := uniqueSet[s]; !exists {
			uniqueSet[s] = struct{}{}
			uniqueList = append(uniqueList, s)
		}
	}

	// Итоговый массив
	result := make([]string, 0, total)
	result = append(result, uniqueList...)

	// Добавляем дубликаты
	for i := 0; i < dupCount; i++ {
		s := uniqueList[rand.Intn(len(uniqueList))]
		result = append(result, s)
	}

	// Перемешивание
	rand.Shuffle(len(result), func(i, j int) {
		result[i], result[j] = result[j], result[i]
	})

	// Запись в файл
	file, err := os.Create(DatastoreFolder + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	lineCounter := 0
	for _, s := range result {
		if lineCounter == wordsOnLine-1 {
			_, err = writer.WriteString(s + "\n")
			lineCounter = 0
		} else {
			_, err = writer.WriteString(s + ",")
			lineCounter++
		}
		if err != nil {
			return err
		}

	}

	return writer.Flush()
}
