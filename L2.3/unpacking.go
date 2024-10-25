package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

/*
Task:
	Создать Go-функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы/руны.
	Например:
		- "a4bc2d5e" => "aaaabccddddde"
		- "abcd" => "abcd"
		- "45" => "" (некорректная строка)
		- "" => ""
Дополнительно:
	1. Реализовать поддержку escape-последовательностей.
		Например:
			qwe\4\5 => qwe45 (*)
			qwe\45 => qwe44444 (*)
			qwe\\5 => qwe\\\\\ (*)
	2. В случае если была передана некорректная строка, функция должна возвращать ошибку.
	3. Написать unit-тесты.
*/

func UnpackString(input string) (string, error) {
	var result strings.Builder
	var prevRune rune
	isEscaped := false

	for i, rune := range input {
		if unicode.IsDigit(rune) && !isEscaped {
			if i == 0 || prevRune == 0 {
				return "", errors.New("\"\" (некорректная строка)")
			}
			repeat, _ := strconv.Atoi(string(rune))
			result.WriteString(strings.Repeat(string(prevRune), repeat))
			prevRune = 0
		} else {
			if !isEscaped && rune == '\\' {
				isEscaped = true
			} else {
				if prevRune != 0 {
					result.WriteRune(prevRune)
				}
				prevRune = rune
				isEscaped = false
			}
		}
	}

	if prevRune != 0 {
		result.WriteRune(prevRune)
	}

	return result.String(), nil
}

func main() {
	examples := []string{"a4bc2d5e", "abcd", "45", "qwe\\4\\5", "qwe\\45", "qwe\\\\5", ""}
	for _, example := range examples {
		unpacked, err := UnpackString(example)
		if err != nil {
			fmt.Printf("Распакованная строка %s: %s\n", example, err)
		} else {
			fmt.Printf("Распакованная строка %s: %s\n", example, unpacked)
		}
	}
}
