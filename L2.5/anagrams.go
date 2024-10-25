package main

import (
	"fmt"
	"sort"
	"strings"
)

/*
Task:
	Написать функцию поиска всех множеств анаграмм по словарю.
	Например:
		'пятак', 'пятка' и 'тяпка' — принадлежат одному множеству;
		'листок', 'слиток' и 'столик' — другому.
	Требования
		1. Входные данные для функции: ссылка на массив, каждый элемент которого — слово на русском языке в кодировке utf8.
		2. Выходные данные: ссылка на мапу множеств анаграмм.
		3. Ключ — первое встретившееся в словаре слово из множества. Значение — ссылка на массив, каждый элемент которого, слово из множества.
		4. Массив должен быть отсортирован по возрастанию.
		5. Множества из одного элемента не должны попасть в результат.
		6. Все слова должны быть приведены к нижнему регистру.
		7. В результате каждое слово должно встречаться только один раз.
*/

func findAnagrams(words []string) *map[string][]string {
	anagramGroups := make(map[string][]string)
	seen := make(map[string]string)

	for _, word := range words {
		lowerWord := strings.ToLower(word)
		sortedWord := sortString(lowerWord)

		// Если это первая анаграмма данного множества, запоминаем оригинальное слово как ключ
		if _, exists := seen[sortedWord]; !exists {
			seen[sortedWord] = lowerWord
		}
		key := seen[sortedWord]
		anagramGroups[key] = append(anagramGroups[key], lowerWord)
	}

	for key, group := range anagramGroups {
		if len(group) < 2 {
			delete(anagramGroups, key)
		} else {
			sort.Strings(group)
		}
	}

	return &anagramGroups
}

func sortString(s string) string {
	chars := strings.Split(s, "")
	sort.Strings(chars)
	return strings.Join(chars, "")
}

func main() {
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "fg"}
	anagrams := findAnagrams(words)
	for key, group := range *anagrams {
		fmt.Println(key, ": ", group)
	}
}
