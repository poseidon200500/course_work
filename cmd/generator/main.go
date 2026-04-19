package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/poseidon200500/course_work/internal/benchmark"
	"github.com/poseidon200500/course_work/internal/generator"
)

const (
	defaultSeed  = int64(42)
	defaultZipfS = 1.2
	defaultZipfV = 1.0
)

func main() {
	scenarios := benchmark.DefaultScenarios()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n=== ГЕНЕРАЦИЯ ДАННЫХ ===")
		fmt.Println("1. Сгенерировать по группе сценариев")
		fmt.Println("2. Сгенерировать по конкретному сценарию")
		fmt.Println("3. Сгенерировать все сценарии")
		fmt.Println("4. Создать и сгенерировать кастомный сценарий")
		fmt.Println("0. Выход")

		choice, err := readChoice(reader, "Выберите действие: ")
		if err != nil {
			fmt.Println("Ошибка ввода:", err)
			continue
		}

		switch choice {
		case 1:
			if err := handleGenerateByGroup(reader, scenarios); err != nil {
				fmt.Println("Ошибка:", err)
			}
		case 2:
			if err := handleGenerateByScenario(reader, scenarios); err != nil {
				fmt.Println("Ошибка:", err)
			}
		case 3:
			if err := generateAllScenarios(scenarios); err != nil {
				fmt.Println("Ошибка:", err)
			}
		case 4:
			if err := handleCustomScenario(reader); err != nil {
				fmt.Println("Ошибка:", err)
			}
		case 0:
			fmt.Println("Выход.")
			return
		default:
			fmt.Println("Неизвестная команда.")
		}
	}
}

func handleGenerateByGroup(reader *bufio.Reader, scenarios []benchmark.Scenario) error {
	groups := benchmark.GroupScenarios(scenarios)
	groupOrder := benchmark.SortedGroupNames(groups)

	if len(groupOrder) == 0 {
		return fmt.Errorf("группы сценариев не найдены")
	}

	fmt.Println("\n=== ВЫБОР ГРУППЫ ===")
	for i, groupName := range groupOrder {
		fmt.Printf("%d. %s\n", i+1, benchmark.FormatGroupName(groupName))
		printGroupDescriptions(groups[groupName])
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

	confirm, err := readYesNo(reader, "Сгенерировать все сценарии этой группы? (y/n): ")
	if err != nil {
		return err
	}
	if !confirm {
		fmt.Println("Отменено.")
		return nil
	}

	return generateScenarioList(selectedScenarios)
}

func handleGenerateByScenario(reader *bufio.Reader, scenarios []benchmark.Scenario) error {
	selected, err := selectScenario(reader, scenarios)
	if err != nil {
		return err
	}
	if selected == nil {
		return nil
	}

	fmt.Println("\nВыбран сценарий:")
	printScenarioDetails(*selected)

	confirm, err := readYesNo(reader, "Сгенерировать этот сценарий? (y/n): ")
	if err != nil {
		return err
	}
	if !confirm {
		fmt.Println("Отменено.")
		return nil
	}

	return generateScenarioList([]benchmark.Scenario{*selected})
}

func generateAllScenarios(scenarios []benchmark.Scenario) error {
	fmt.Println("\nБудут сгенерированы все сценарии:")
	for _, sc := range scenarios {
		fmt.Printf("- [%s] %s — %s\n", benchmark.FormatGroupName(sc.Group), sc.Name, sc.Description)
	}

	return generateScenarioList(scenarios)
}

func handleCustomScenario(reader *bufio.Reader) error {
	fmt.Println("\n=== СОЗДАНИЕ КАСТОМНОГО СЦЕНАРИЯ ===")

	name, err := readLine(reader, "Имя сценария: ")
	if err != nil {
		return err
	}
	name = sanitizeScenarioName(name)
	if name == "" {
		return fmt.Errorf("имя сценария не может быть пустым")
	}

	description, err := readLine(reader, "Описание сценария: ")
	if err != nil {
		return err
	}

	group, err := readLine(reader, "Группа сценария (по умолчанию custom): ")
	if err != nil {
		return err
	}
	if strings.TrimSpace(group) == "" {
		group = "custom"
	}

	total, err := readIntPrompt(reader, "Количество строк: ")
	if err != nil {
		return err
	}

	wordsOnLine, err := readIntPrompt(reader, "Количество слов в строке файла: ")
	if err != nil {
		return err
	}

	duplicatePercent, err := readIntPrompt(reader, "Процент дубликатов (0-100): ")
	if err != nil {
		return err
	}

	maxLen, err := readIntPrompt(reader, "Максимальная длина слова: ")
	if err != nil {
		return err
	}

	distribution, err := readDistribution(reader)
	if err != nil {
		return err
	}

	sc := benchmark.Scenario{
		Name:             name,
		Total:            total,
		WordsPerLine:     wordsOnLine,
		DuplicatePercent: duplicatePercent,
		MaxLen:           maxLen,
		Group:            group,
		Distribution:     distribution,
		Description:      description,
	}

	fmt.Println("\nСоздан сценарий:")
	printScenarioDetails(sc)

	confirm, err := readYesNo(reader, "Сгенерировать этот сценарий? (y/n): ")
	if err != nil {
		return err
	}
	if !confirm {
		fmt.Println("Отменено.")
		return nil
	}

	return generateScenarioList([]benchmark.Scenario{sc})
}

func generateScenarioList(scenarios []benchmark.Scenario) error {
	if len(scenarios) == 0 {
		return fmt.Errorf("список сценариев пуст")
	}

	for _, sc := range scenarios {
		if err := regenerateScenario(sc); err != nil {
			return err
		}
	}

	fmt.Println("\nГенерация завершена.")
	return nil
}

func regenerateScenario(sc benchmark.Scenario) error {
	filename := scenarioFilename(sc)
	fullPath := filepath.Join(generator.DatastoreFolder, filename)

	if _, err := os.Stat(fullPath); err == nil {
		fmt.Printf("Удаление старого файла сценария %s: %s\n", sc.Name, fullPath)
		if err := os.Remove(fullPath); err != nil {
			return fmt.Errorf("не удалось удалить старый файл сценария %s: %w", sc.Name, err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("ошибка проверки файла сценария %s: %w", sc.Name, err)
	}

	cfg := buildConfig(sc)

	fmt.Printf("\nГенерация сценария %s...\n", sc.Name)
	printScenarioDetails(sc)
	fmt.Printf("Файл: %s\n", fullPath)

	if err := generator.GenerateDataWithConfig(cfg); err != nil {
		return fmt.Errorf("ошибка генерации сценария %s: %w", sc.Name, err)
	}

	fmt.Printf("Сценарий %s успешно сгенерирован.\n", sc.Name)
	return nil
}

func buildConfig(sc benchmark.Scenario) generator.Config {
	return generator.Config{
		Total:            sc.Total,
		WordsOnLine:      sc.WordsPerLine,
		DuplicatePercent: sc.DuplicatePercent,
		MaxLen:           sc.MaxLen,
		Filename:         scenarioFilename(sc),
		Deterministic:    true,
		Seed:             defaultSeed,
		Distribution:     sc.Distribution,
		ZipfS:            defaultZipfS,
		ZipfV:            defaultZipfV,
	}
}

func scenarioFilename(sc benchmark.Scenario) string {
	return fmt.Sprintf("data_%s.txt", sc.Name)
}

func selectScenario(reader *bufio.Reader, scenarios []benchmark.Scenario) (*benchmark.Scenario, error) {
	fmt.Println("\n=== ВЫБОР СЦЕНАРИЯ ===")
	for i, sc := range scenarios {
		fmt.Printf("%d. [%s] %s — %s\n",
			i+1,
			benchmark.FormatGroupName(sc.Group),
			sc.Name,
			sc.Description,
		)
	}
	fmt.Println("0. Назад")

	choice, err := readChoice(reader, "Выберите сценарий: ")
	if err != nil {
		return nil, err
	}
	if choice == 0 {
		return nil, nil
	}
	if choice < 1 || choice > len(scenarios) {
		return nil, fmt.Errorf("некорректный номер сценария")
	}

	selected := scenarios[choice-1]
	return &selected, nil
}

func printScenarioDetails(sc benchmark.Scenario) {
	fmt.Printf("Имя: %s\n", sc.Name)
	fmt.Printf("Группа: %s\n", benchmark.FormatGroupName(sc.Group))
	fmt.Printf("Описание: %s\n", sc.Description)
	fmt.Printf("Всего строк: %d\n", sc.Total)
	fmt.Printf("Слов в строке файла: %d\n", sc.WordsPerLine)
	fmt.Printf("Процент дубликатов: %d\n", sc.DuplicatePercent)
	fmt.Printf("Максимальная длина слова: %d\n", sc.MaxLen)
	fmt.Printf("Распределение: %s\n", sc.Distribution)
}

func printGroupDescriptions(scenarios []benchmark.Scenario) {
	for _, sc := range scenarios {
		fmt.Printf("   - %s: %s\n", sc.Name, sc.Description)
	}
}

func readChoice(reader *bufio.Reader, prompt string) (int, error) {
	text, err := readLine(reader, prompt)
	if err != nil {
		return 0, err
	}

	value, err := strconv.Atoi(text)
	if err != nil {
		return 0, fmt.Errorf("ожидалось число")
	}

	return value, nil
}

func readIntPrompt(reader *bufio.Reader, prompt string) (int, error) {
	return readChoice(reader, prompt)
}

func readLine(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

func readYesNo(reader *bufio.Reader, prompt string) (bool, error) {
	text, err := readLine(reader, prompt)
	if err != nil {
		return false, err
	}

	text = strings.ToLower(strings.TrimSpace(text))
	return text == "y" || text == "yes" || text == "д" || text == "да", nil
}

func readDistribution(reader *bufio.Reader) (generator.DistributionType, error) {
	fmt.Println("Выберите распределение:")
	fmt.Println("1. Uniform")
	fmt.Println("2. Zipf")

	choice, err := readChoice(reader, "Ваш выбор: ")
	if err != nil {
		return "", err
	}

	switch choice {
	case 1:
		return generator.DistributionUniform, nil
	case 2:
		return generator.DistributionZipf, nil
	default:
		return "", fmt.Errorf("некорректный выбор распределения")
	}
}

func sanitizeScenarioName(name string) string {
	name = strings.TrimSpace(strings.ToUpper(name))
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	return name
}
