package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
Task:
	1. Отсортировать строки в файле по аналогии с консольной утилитой sort (man sort — смотрим описание и основные параметры):
	на входе подается файл из несортированными строками, на выходе — файл с отсортированными.
		Реализовать поддержку утилитой следующих ключей:
			-k — указание колонки для сортировки (слова в строке могут выступать в качестве колонок, по умолчанию разделитель — пробел);
			-n — сортировать по числовому значению;
			-r — сортировать в обратном порядке;
			-u — не выводить повторяющиеся строки.
Дополнительно:
	Реализовать поддержку утилитой следующих ключей:
		-M — сортировать по названию месяца;
		-b — игнорировать хвостовые пробелы;
		-c — проверять отсортированы ли данные;
		-h — сортировать по числовому значению с учетом суффиксов.
*/

// parseCommandLine обрабатывает флаги командной строки и сохраняет их значения в структуре SortOptions.
// -k определяет номер колонки для сортировки (по умолчанию первая).
// -n включает числовую сортировку, чтобы строки с числами сортировались как числа, а не строки.
// -r включает сортировку в обратном порядке.
// -u включает уникальность строк, удаляя дубликаты.
// -o позволяет указать файл для записи результата.

type SortOptions struct {
	column       int
	numeric      bool
	reverseOrder bool
	uniqueOnly   bool
	outputFile   string
}

type SortableLines struct {
	lines   []string
	options SortOptions
}

func (s SortableLines) Len() int           { return len(s.lines) }
func (s SortableLines) Swap(i, j int)      { s.lines[i], s.lines[j] = s.lines[j], s.lines[i] }
func (s SortableLines) Less(i, j int) bool { return lineComparison(s.lines[i], s.lines[j], s.options) }

func main() {
	opts := parseCommandLine()
	filePath := flag.Arg(0)

	lines, err := loadLines(filePath)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Ошибка загрузки файла: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}

	sortedLines, err := sortLines(lines, opts)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Ошибка сортировки: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}

	if opts.outputFile != "" {
		err = saveLines(opts.outputFile, sortedLines)
	} else {
		err = saveLines(filePath, sortedLines)
	}
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Ошибка записи: %v\n", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}

func parseCommandLine() SortOptions {
	var opts SortOptions
	flag.IntVar(&opts.column, "k", 1, "Указать колонку для сортировки")
	flag.BoolVar(&opts.numeric, "n", false, "Числовая сортировка")
	flag.BoolVar(&opts.reverseOrder, "r", false, "Обратный порядок сортировки")
	flag.BoolVar(&opts.uniqueOnly, "u", false, "Только уникальные строки")
	flag.StringVar(&opts.outputFile, "o", "", "Файл для записи результата")
	flag.Parse()
	return opts
}

func loadLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func saveLines(filePath string, lines []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}

func lineComparison(line1, line2 string, opts SortOptions) bool {
	parts1 := strings.Fields(line1)
	parts2 := strings.Fields(line2)

	if opts.column > len(parts1) || opts.column > len(parts2) {
		return line1 < line2
	}

	var comparison int
	if opts.numeric {
		num1, err1 := strconv.Atoi(parts1[opts.column-1])
		num2, err2 := strconv.Atoi(parts2[opts.column-1])
		if err1 != nil || err2 != nil {
			comparison = strings.Compare(parts1[opts.column-1], parts2[opts.column-1])
		} else {
			comparison = num1 - num2
		}
	} else {
		comparison = strings.Compare(parts1[opts.column-1], parts2[opts.column-1])
	}

	return comparison < 0
}

func uniqueLines(lines []string) []string {
	existing := make(map[string]struct{})
	var result []string
	for _, line := range lines {
		if _, found := existing[line]; !found {
			existing[line] = struct{}{}
			result = append(result, line)
		}
	}
	return result
}

func reverseOrder(lines []string) {
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}
}

func sortLines(lines []string, opts SortOptions) ([]string, error) {
	if opts.uniqueOnly {
		lines = uniqueLines(lines)
	}

	sortable := SortableLines{lines: lines, options: opts}
	sort.Sort(sortable)

	if opts.reverseOrder {
		reverseOrder(lines)
	}

	return lines, nil
}
