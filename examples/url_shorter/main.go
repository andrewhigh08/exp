// Package main демонстрирует реализацию интерфейса `fmt.Stringer`
// для создания кастомного строкового представления для пользовательского типа.
package main

import (
	"fmt"
)

// Abbreviator — это пользовательский тип на основе строки.
// Мы определяем для него собственный метод `String()`.
type Abbreviator string

// String реализует интерфейс `fmt.Stringer` для типа Abbreviator.
// Когда значение этого типа передается в функцию пакета fmt (например, Println),
// для его отображения будет автоматически вызван этот метод.
//
// Логика: "kubernetes" -> "k" + "8" (длина - 2) + "s" -> "k8s".
func (s Abbreviator) String() string {
	// Преобразуем в срез рун для корректной работы с многобайтными символами (например, кириллицей).
	runes := []rune(s)
	length := len(runes)

	// Если строка слишком короткая для аббревиатуры, возвращаем ее как есть.
	if length <= 2 {
		return string(s)
	}

	// Формируем аббревиатуру.
	// Использование fmt.Sprintf более читаемо и идиоматично, чем ручная конкатенация.
	return fmt.Sprintf("%c%d%c", runes[0], length-2, runes[length-1])
}

func main() {
	testCases := []Abbreviator{
		"kubernetes",
		"internationalization",
		"localization",
		"hi",        // Короткая строка
		"адаптация", // Пример с кириллицей
	}

	fmt.Println("--- Демонстрация fmt.Stringer с кастомным типом ---")
	for _, str := range testCases {
		// При вызове fmt.Println(str), Go автоматически обнаруживает,
		// что тип `Abbreviator` имеет метод `String() string`, и вызывает его.
		fmt.Printf("Исходная строка: '%s', результат: %s\n", str, str)
	}
}
