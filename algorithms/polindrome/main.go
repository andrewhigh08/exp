// Package main содержит реализации функции для проверки, является ли строка палиндромом.
// Палиндром — это строка, которая читается одинаково в обоих направлениях.
package main

import (
	"fmt"
	"strings"
	"unicode"
)

// isPalindromeSimple — это простая, чувствительная к регистру проверка на палиндром.
// Она не игнорирует пробелы и знаки препинания.
// Алгоритм использует два указателя (один в начале, другой в конце) и сравнивает символы, двигаясь к центру.
func isPalindromeSimple(st string) bool {
	// Преобразование строки в срез рун — ключевой шаг для корректной работы с Unicode (например, с кириллицей).
	// Одна кириллическая буква может занимать несколько байт.
	runeSt := []rune(st)
	length := len(runeSt)

	// Итерируемся только до середины строки.
	for i := 0; i < length/2; i++ {
		// Сравниваем i-й символ с начала и i-й символ с конца.
		if runeSt[i] != runeSt[length-1-i] {
			return false // Если символы не совпадают, это не палиндром.
		}
	}
	return true
}

// isPalindromeAdvanced — это более сложная проверка на палиндром.
// Она нечувствительна к регистру и игнорирует все символы, кроме букв.
func isPalindromeAdvanced(st string) bool {
	// Приводим всю строку к нижнему регистру для регистронезависимого сравнения.
	lowerSt := strings.ToLower(st)
	runes := []rune(lowerSt)

	// Используем два указателя: left с начала, right с конца.
	left, right := 0, len(runes)-1

	for left < right {
		// Пропускаем все не-буквенные символы слева.
		if !unicode.IsLetter(runes[left]) {
			left++
			continue
		}
		// Пропускаем все не-буквенные символы справа.
		if !unicode.IsLetter(runes[right]) {
			right--
			continue
		}

		// Сравниваем буквы.
		if runes[left] != runes[right] {
			return false
		}

		// Сдвигаем указатели к центру.
		left++
		right--
	}
	return true
}

func main() {
	testCases := []string{
		"Комок", // Палиндром с заглавной буквой
		"Кабак",
		"казак",
		"шорох",
		"торрот",
		"А роза упала на лапу Азора", // Классический палиндром с пробелами
		"Eva, can I see bees in a cave?", // Английский палиндром со знаками препинания
		"привет",                       // Не палиндром
		"а",                            // Палиндром из одного символа
		"",                             // Пустая строка считается палиндромом
	}

	fmt.Println("--- Простая проверка (isPalindromeSimple) ---")
	for _, tc := range testCases {
		fmt.Printf("Строка: '%-30s' -> Палиндром: %t\n", tc, isPalindromeSimple(tc))
	}

	fmt.Println("\n--- Продвинутая проверка (isPalindromeAdvanced) ---")
	for _, tc := range testCases {
		fmt.Printf("Строка: '%-30s' -> Палиндром: %t\n", tc, isPalindromeAdvanced(tc))
	}
}
