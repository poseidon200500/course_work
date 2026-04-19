package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	analysis "github.com/poseidon200500/course_work/internal/analisys"
	"github.com/poseidon200500/course_work/internal/benchmark"
	"github.com/poseidon200500/course_work/internal/storage"
)

// func main() {
// 	storages := map[string]func() storage.Storage{
// 		"BASE":      func() storage.Storage { return storage.NewBaseStorage() },
// 		"INTERN":    func() storage.Storage { return storage.NewInternStorage() },
// 		"UNIQUE_V1": func() storage.Storage { return storage.NewUniqueStorageV1() },
// 		"UNIQUE_V2": func() storage.Storage { return storage.NewUniqueStorageV2() },
// 	}

// 	scenarios := benchmark.DefaultScenarios()
// 	reader := bufio.NewReader(os.Stdin)

// 	for {
// 		fmt.Println("\n=== БЕНЧМАРКИ ХРАНИЛИЩ ===")
// 		fmt.Println("1. Запустить по группе сценариев")
// 		fmt.Println("2. Запустить по конкретному сценарию")
// 		fmt.Println("3. Запустить все сценарии")
// 		fmt.Println("0. Выход")

// 		choice, err := readChoice(reader, "Выберите действие: ")
// 		if err != nil {
// 			fmt.Println("Ошибка ввода:", err)
// 			continue
// 		}

// 		switch choice {
// 		case 1:
// 			if err := handleBenchmarkByGroup(reader, scenarios, storages); err != nil {
// 				fmt.Println("Ошибка:", err)
// 			}
// 		case 2:
// 			if err := handleBenchmarkByScenario(reader, scenarios, storages); err != nil {
// 				fmt.Println("Ошибка:", err)
// 			}
// 		case 3:
// 			if _, err := benchmark.RunSelected(storages, scenarios); err != nil {
// 				fmt.Println("Ошибка:", err)
// 			}
// 		case 0:
// 			fmt.Println("Выход.")
// 			return
// 		default:
// 			fmt.Println("Неизвестная команда.")
// 		}
// 	}
// }

func main() {
	makeCSV()
}

func handleBenchmarkByGroup(reader *bufio.Reader, scenarios []benchmark.Scenario, storages map[string]func() storage.Storage) error {
	groups := benchmark.GroupScenarios(scenarios)
	groupOrder := benchmark.SortedGroupNames(groups)

	if len(groupOrder) == 0 {
		return fmt.Errorf("группы сценариев не найдены")
	}

	fmt.Println("\n=== ВЫБОР ГРУППЫ ===")
	for i, groupName := range groupOrder {
		fmt.Printf("%d. %s\n", i+1, benchmark.FormatGroupName(groupName))
		for _, sc := range groups[groupName] {
			fmt.Printf("   - %s: %s\n", sc.Name, sc.Description)
		}
	}
	fmt.Println("0. Назад")

	choice, err := readChoice(reader, "Выберите группу: ")
	if err != nil {
		return err
	}

	if choice == 0 {
		return nil
	}
	if choice < 1 || choice > len(groupOrder) {
		return fmt.Errorf("некорректный номер группы")
	}

	selectedGroup := groupOrder[choice-1]
	selectedScenarios := groups[selectedGroup]

	fmt.Printf("\nВыбрана группа: %s\n", benchmark.FormatGroupName(selectedGroup))
	for _, sc := range selectedScenarios {
		fmt.Printf("- %s: %s\n", sc.Name, sc.Description)
	}

	_, err = benchmark.RunSelected(storages, selectedScenarios)
	return err
}

func handleBenchmarkByScenario(reader *bufio.Reader, scenarios []benchmark.Scenario, storages map[string]func() storage.Storage) error {
	fmt.Println("\n=== ВЫБОР СЦЕНАРИЯ ===")
	for i, sc := range scenarios {
		fmt.Printf(
			"%d. [%s] %s — %s\n",
			i+1,
			benchmark.FormatGroupName(sc.Group),
			sc.Name,
			sc.Description,
		)
	}
	fmt.Println("0. Назад")

	choice, err := readChoice(reader, "Выберите сценарий: ")
	if err != nil {
		return err
	}

	if choice == 0 {
		return nil
	}
	if choice < 1 || choice > len(scenarios) {
		return fmt.Errorf("некорректный номер сценария")
	}

	selected := scenarios[choice-1]

	fmt.Printf("\nВыбран сценарий:\n")
	printScenarioDetails(selected)
	_, err = benchmark.RunSelected(storages, []benchmark.Scenario{selected})
	return err
}

func printScenarioDetails(sc benchmark.Scenario) {
	fmt.Printf("Имя: %s\n", sc.Name)
	fmt.Printf("Группа: %s\n", benchmark.FormatGroupName(sc.Group))
	fmt.Printf("Описание: %s\n", sc.Description)
	fmt.Printf("Всего строк: %d\n", sc.Total)
	fmt.Printf("Строк на строку файла: %d\n", sc.WordsPerLine)
	fmt.Printf("Процент дубликатов: %d\n", sc.DuplicatePercent)
	fmt.Printf("Макс. длина слова: %d\n", sc.MaxLen)
	fmt.Printf("Распределение: %s\n", sc.Distribution)
}

func readChoice(reader *bufio.Reader, prompt string) (int, error) {
	fmt.Print(prompt)

	text, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	text = strings.TrimSpace(text)
	value, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("ожидалось число")
	}

	return value, nil
}

func makeCSV() {
	storages := map[string]func() storage.Storage{
		"BASE":      func() storage.Storage { return storage.NewBaseStorage() },
		"INTERN":    func() storage.Storage { return storage.NewInternStorage() },
		"UNIQUE_V1": func() storage.Storage { return storage.NewUniqueStorageV1() },
		"UNIQUE_V2": func() storage.Storage { return storage.NewUniqueStorageV2() },
	}

	scenarios := benchmark.DefaultScenarios()

	results, err := benchmark.RunSelected(storages, scenarios)
	if err != nil {
		panic(err)
	}

	if err := analysis.WriteResultsCSV(results, "benchmark_results.csv"); err != nil {
		panic(err)
	}

	fmt.Println("CSV и графики успешно сохранены.")
}
