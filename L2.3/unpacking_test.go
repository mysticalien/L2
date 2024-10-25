package main

import (
	"testing"
)

func TestUnpackString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"a4bc2d5e", "aaaabccddddde", false},
		{"abcd", "abcd", false},
		{"45", "", true},
		{"", "", false},
		{"qwe\\4\\5", "qwe45", false},
		{"qwe\\45", "qwe44444", false},
		{"qwe\\\\5", "qwe\\\\\\\\\\", false},
		{"a\\2b", "a2b", false},
	}

	for _, test := range tests {
		result, err := UnpackString(test.input)

		if test.hasError {
			if err == nil {
				t.Errorf("Ожидалась ошибка для ввода %s, но её не произошло", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Не ожидалась ошибка для ввода %s, но произошла ошибка: %v", test.input, err)
			} else if result != test.expected {
				t.Errorf("Для ввода %s, ожидалось %s, но получено %s", test.input, test.expected, result)
			}
		}
	}
}
