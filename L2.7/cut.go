package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

/*
Task:
	Реализовать утилиту аналог консольной команды cut (man cut).
	Утилита должна принимать строки через STDIN, разбивать по разделителю (TAB) на колонки и выводить запрошенные.
	Реализовать поддержку утилитой следующих ключей:
		-f — "fields": выбрать поля (колонки);
		-d — "delimiter": использовать другой разделитель;
		-s — "separated": только строки с разделителем.
*/

var (
	filePath      string
	fields        string
	delimiter     string
	onlySeparated bool
)

func init() {
	flag.StringVar(&fields, "f", "", "Specify fields (columns) to cut, e.g. '1,3'")
	flag.StringVar(&delimiter, "d", "\t", "Specify a custom delimiter (default: TAB)")
	flag.BoolVar(&onlySeparated, "s", false, "Only print lines with the delimiter")
	flag.StringVar(&filePath, "file", "", "Path to input file")
}

func main() {
	flag.Parse()
	if err := customCut(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func customCut() error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	fieldIndices := parseFields(fields)
	for scanner.Scan() {
		line := scanner.Text()
		if onlySeparated && !strings.Contains(line, delimiter) {
			continue
		}
		columns := strings.Split(line, delimiter)
		printSelectedFields(columns, fieldIndices)
	}
	return scanner.Err()
}

func parseFields(fields string) []int {
	if fields == "" {
		return nil
	}
	var indices []int
	for _, f := range strings.Split(fields, ",") {
		var index int
		if _, err := fmt.Sscanf(f, "%d", &index); err == nil {
			indices = append(indices, index-1)
		}
	}
	return indices
}

func printSelectedFields(columns []string, indices []int) {
	var output []string
	for _, idx := range indices {
		if idx >= 0 && idx < len(columns) {
			output = append(output, columns[idx])
		}
	}
	if len(output) > 0 {
		fmt.Println(strings.Join(output, "\t"))
	}
}
