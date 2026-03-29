package main

import (
	"fmt"

	"github.com/poseidon200500/course_work/generator"
)

func main() {
	err := generator.GenerateData(
		1_000, // количество строк
		10,    // количество строк на строке
		40,    // % дубликатов
		8,     // макс длина строки
		"data.txt",
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Done")
}
