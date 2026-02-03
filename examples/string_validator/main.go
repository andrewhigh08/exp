// Package main демонстрирует создание строкового валидатора, который проверяет
// соответствие строки набору регулярных выражений, загружаемых из файла.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

// StringValidator хранит скомпилированные регулярные выражения для валидации.
type StringValidator struct {
	patterns []*regexp.Regexp
}

// NewStringValidator — это конструктор для валидатора.
// Он принимает путь к файлу с паттернами и возвращает готовый валидатор или ошибку.
// Такой подход (возврат ошибки вместо паники) является идиоматичным для Go.
func NewStringValidator(filename string) (*StringValidator, error) {
	sv := &StringValidator{}
	err := sv.loadPatterns(filename)
	if err != nil {
		// Если загрузка паттернов не удалась, возвращаем ошибку наверх.
		return nil, fmt.Errorf("не удалось создать валидатор: %w", err)
	}
	return sv, nil
}

// loadPatterns загружает и компилирует регулярные выражения из файла.
func (sv *StringValidator) loadPatterns(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл '%s': %w", filename, err)
	}
	defer file.Close()

	// Использование bufio.Scanner — это эффективный и идиоматичный способ
	// читать файл построчно, который корректно обрабатывает последнюю строку.
	scanner := bufio.NewScanner(file)
	var patterns []*regexp.Regexp
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		// Пропускаем пустые строки
		if line == "" {
			continue
		}

		// MustCompile паникует при ошибке, что хорошо для статических паттернов,
		// но для паттернов из файла лучше использовать Compile и обрабатывать ошибку.
		re, err := regexp.Compile(line)
		if err != nil {
			return fmt.Errorf("не удалось скомпилировать паттерн на строке %d ('%s'): %w", lineNumber, line, err)
		}
		patterns = append(patterns, re)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("ошибка при сканировании файла: %w", err)
	}

	sv.patterns = patterns
	return nil
}

// Validate проверяет, соответствует ли строка ВСЕМ загруженным паттернам.
// Исходная логика была неясной (`mismatchCount <= 3`).
// Новая логика более прямолинейна: строка валидна, если проходит все проверки.
func (sv *StringValidator) Validate(str string) bool {
	// Проходим по всем паттернам.
	for _, p := range sv.patterns {
		// Если строка не соответствует хотя бы одному паттерну, она невалидна.
		if !p.MatchString(str) {
			return false
		}
	}
	// Если строка соответствует всем паттернам.
	return true
}

// createDummyPatternsFile создает временный файл с паттернами для демонстрации.
func createDummyPatternsFile(filename string) error {
	content := `^user_` + "\n" + `\d{3}$` + "\n" + `.*_test$`
	return os.WriteFile(filename, []byte(content), 0644)
}

func main() {
	patternFile := "patterns.cfg"
	// 1. Создаем файл с паттернами для нашего примера.
	if err := createDummyPatternsFile(patternFile); err != nil {
		log.Fatalf("Не удалось создать файл с паттернами: %v", err)
	}
	// Удаляем временный файл в конце.
	defer os.Remove(patternFile)

	fmt.Printf("Загрузка паттернов из файла '%s'...\n", patternFile)
	fmt.Println("Паттерны:\n1. Должно начинаться с 'user_'\n2. Должно содержать 3 цифры\n3. Должно заканчиваться на '_test'")

	// 2. Создаем валидатор.
	validator, err := NewStringValidator(patternFile)
	if err != nil {
		log.Fatalf("Ошибка при создании валидатора: %v", err)
	}

	testCases := []string{
		"user_123_test", // Валидно
		"user_456",      // Невалидно (не заканчивается на _test)
		"admin_123_test",// Невалидно (не начинается с user_)
		"user_12_test",  // Невалидно (не 3 цифры)
	}

	fmt.Println("\n--- Результаты валидации ---")
	for _, tc := range testCases {
		isValid := validator.Validate(tc)
		fmt.Printf("Строка '%-15s' -> Валидна: %t\n", tc, isValid)
	}
}
