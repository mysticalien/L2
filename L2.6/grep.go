package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

/*
Task:
	Реализовать утилиту фильтрации по аналогии с консольной утилитой (man grep — смотрим описание и основные параметры).
	Реализовать поддержку утилитой следующих ключей:
		-A - "after": печатать +N строк после совпадения;
		-B - "before": печатать +N строк до совпадения;
		-C - "context": (A+B) печатать ±N строк вокруг совпадения;
		-c - "count": количество строк;
		-i - "ignore-case": игнорировать регистр;
		-v - "invert": вместо совпадения, исключать;
		-F - "fixed": точное совпадение со строкой, не паттерн;
		-n - "line num": напечатать номер строки.
*/

// Опции для утилиты grep
type GrepOptions struct {
	afterContext  int
	beforeContext int
	context       int
	countOnly     bool
	ignoreCase    bool
	invertMatch   bool
	fixedMatch    bool
	printLineNum  bool
}

// Парсинг опций командной строки
func parseGrepOptions() GrepOptions {
	var opts GrepOptions
	flag.IntVar(&opts.afterContext, "A", 0, "Печатать +N строк после совпадения")
	flag.IntVar(&opts.beforeContext, "B", 0, "Печатать +N строк до совпадения")
	flag.IntVar(&opts.context, "C", 0, "Печатать ±N строк вокруг совпадения")
	flag.BoolVar(&opts.countOnly, "c", false, "Выводить только количество совпадений")
	flag.BoolVar(&opts.ignoreCase, "i", false, "Игнорировать регистр")
	flag.BoolVar(&opts.invertMatch, "v", false, "Инвертировать совпадение")
	flag.BoolVar(&opts.fixedMatch, "F", false, "Точное совпадение строки")
	flag.BoolVar(&opts.printLineNum, "n", false, "Печатать номер строки")
	flag.Parse()

	// Если указан флаг -C, устанавливаем значения -A и -B
	if opts.context > 0 {
		opts.beforeContext = opts.context
		opts.afterContext = opts.context
	}

	return opts
}

// Основная функция для обработки поиска
func grep(lines []string, pattern string, opts GrepOptions) {
	if opts.ignoreCase {
		pattern = strings.ToLower(pattern)
	}

	var regex *regexp.Regexp
	var err error
	if !opts.fixedMatch {
		regex, err = regexp.Compile(pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка компиляции регулярного выражения: %v\n", err)
			return
		}
	}

	lineCount := len(lines)
	var count int

	for i, line := range lines {
		lineToMatch := line
		if opts.ignoreCase {
			lineToMatch = strings.ToLower(line)
		}

		match := (opts.fixedMatch && lineToMatch == pattern) || (!opts.fixedMatch && regex.MatchString(lineToMatch))
		if opts.invertMatch {
			match = !match
		}

		if match {
			count++
			if opts.countOnly {
				continue
			}

			start := max(0, i-opts.beforeContext)
			end := min(lineCount, i+opts.afterContext+1)

			for j := start; j < end; j++ {
				if opts.printLineNum {
					fmt.Printf("%d:%s\n", j+1, lines[j])
				} else {
					fmt.Println(lines[j])
				}
			}
		}
	}

	if opts.countOnly {
		fmt.Println(count)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	opts := parseGrepOptions()
	pattern := flag.Arg(0)
	filePath := flag.Arg(1)

	if filePath == "" {
		fmt.Fprintln(os.Stderr, "Укажите путь к файлу для поиска.")
		os.Exit(1)
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка открытия файла: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения файла: %v\n", err)
		os.Exit(1)
	}

	grep(lines, pattern, opts)
}
