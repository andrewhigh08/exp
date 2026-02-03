// Package main содержит реализации алгоритма сжатия данных RLE (Run-Length Encoding).
// RLE — это простой алгоритм, который заменяет последовательности повторяющихся
// символов на один символ и количество его повторений.
// Например, строка "AAABBC" сжимается в "3A2B1C".
package main

import (
	"fmt"
	"strconv"
	"strings"
)

// rleInefficient демонстрирует неэффективный, но прямолинейный подход к RLE.
// ПРОБЛЕМА: Использование оператора `+=` для конкатенации строк в цикле.
// В Go строки неизменяемы (immutable). Каждая операция `res += ...` создает
// совершенно новую строку в памяти и копирует в нее старое и новое содержимое.
// Это приводит к большому количеству лишних аллокаций памяти и имеет
// квадратичную сложность O(n^2) от длины результирующей строки.
func rleInefficient(inStr string) string {
	if len(inStr) == 0 {
		return ""
	}

	runes := []rune(inStr)
	var (
		result string
		count  = 1
		// Начинаем с первого символа.
		prevChar = runes[0]
	)

	// Идем со второго символа, сравнивая его с предыдущим.
	for i := 1; i < len(runes); i++ {
		if runes[i] == prevChar {
			// Символ тот же, увеличиваем счетчик.
			count++
		} else {
			// Символ изменился. Записываем результат для предыдущей серии.
			result += strconv.Itoa(count) + string(prevChar)
			// Сбрасываем счетчик и обновляем предыдущий символ.
			count = 1
			prevChar = runes[i]
		}
	}

	// Важно не забыть записать результат для последней серии символов после выхода из цикла.
	result += strconv.Itoa(count) + string(prevChar)

	return result
}

// rleEfficient демонстрирует эффективный и идиоматичный подход к RLE с использованием `strings.Builder`.
// РЕШЕНИЕ: `strings.Builder` использует внутренний байтовый буфер, который может расти
// без необходимости каждый раз переаллоцировать всю строку. Это сводит сложность
// операции построения строки к амортизированной O(n) от ее длины.
func rleEfficient(inStr string) string {
	if len(inStr) == 0 {
		return ""
	}

	runes := []rune(inStr)
	var result strings.Builder
	// Оптимизация: предварительно выделяем память, чтобы уменьшить количество реаллокаций.
	result.Grow(len(runes)) // Приблизительная оценка; в худшем случае результат будет 2*len(runes).

	count := 1
	prevChar := runes[0]

	for i := 1; i < len(runes); i++ {
		if runes[i] == prevChar {
			count++
		} else {
			result.WriteString(strconv.Itoa(count))
			result.WriteRune(prevChar)
			count = 1
			prevChar = runes[i]
		}
	}

	// Записываем последнюю серию.
	result.WriteString(strconv.Itoa(count))
	result.WriteRune(prevChar)

	return result.String()
}

func main() {
	testCases := []string{
		"AAAbbc",
		"WWWWWWWWWWWWBWWWWWWWWWWWWBBBWWWWWWWWWWWWWWWWWWWWWWWWBWWWWWWWWWWWWWW",
		"abc",
		"AAAAA",
		"",
	}

	for _, tc := range testCases {
		fmt.Printf("--- Тест: '%s' ---\n", tc)
		inefficientResult := rleInefficient(tc)
		efficientResult := rleEfficient(tc)
		fmt.Printf("Неэффективный: %s\n", inefficientResult)
		fmt.Printf("Эффективный:   %s\n", efficientResult)
		fmt.Printf("Результаты совпадают: %t\n", inefficientResult == efficientResult)
		fmt.Println()
	}
}
