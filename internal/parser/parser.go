package parser

import (
	"bufio"
	"os"
	"strings"
)

func ParseData(filepath string) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// если файл большой — увеличим буфер
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	var result []string

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		words := strings.Split(line, ",")
		for _, w := range words {
			if w != "" {
				result = append(result, w)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
