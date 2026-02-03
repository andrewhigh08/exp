// Package main содержит решение задачи о банкомате.
// Задача — разложить запрошенную сумму денег на доступные номиналы банкнот.
package main

import (
	"errors"
	"fmt"
	"sort"
)

// notes содержит доступные номиналы банкнот.
// Важно, что срез отсортирован по убыванию для корректной работы жадного алгоритма.
var notes = []int{
	5000,
	2000,
	1000,
	500,
	100,
	50, // Добавим мелкие купюры для более интересных тестов
	10,
}

var errCannotDispense = errors.New("невозможно выдать запрошенную сумму")
var errInvalidAmount = errors.New("сумма должна быть положительным числом")

// getMoney реализует "жадный" алгоритм для выдачи денег.
// Он принимает сумму и возвращает карту, где ключ - номинал банкноты, а значение - их количество.
//
// Пример:
//   getMoney(5600) -> map[5000:1, 500:1, 100:1], nil
//   getMoney(1234) -> nil, "невозможно выдать запрошенную сумму"
//
// @param {int} value - Запрашиваемая сумма.
// @return {map[int]int} - Карта с количеством банкнот каждого номинала.
// @return {error} - Ошибка, если сумму выдать невозможно.
func getMoney(value int) (result map[int]int, err error) {
	// Проверка на корректность введенной суммы.
	if value <= 0 {
		return nil, errInvalidAmount
	}

	// Для наглядности, убедимся, что наши банкноты всегда отсортированы.
	sort.Sort(sort.Reverse(sort.IntSlice(notes)))

	result = make(map[int]int)
	remaining := value

	// Итерируемся по банкнотам от большей к меньшей.
	for _, note := range notes {
		// Если номинал банкноты больше оставшейся суммы, пропускаем его.
		if note > remaining {
			continue
		}

		// Вычисляем, сколько целых банкнот данного номинала можно выдать.
		count := remaining / note
		if count > 0 {
			result[note] = count
			// Уменьшаем оставшуюся сумму.
			remaining %= note
		}
	}

	// Если в конце осталась какая-то сумма, значит, мы не можем ее выдать.
	if remaining > 0 {
		return nil, errCannotDispense
	}

	return result, nil
}

func main() {
	testCases := []int{
		5600,
		2480,
		1800,
		7770,
		1234, // Невозможно выдать
		10000,
		50,
		0,      // Некорректная сумма
		-100,   // Некорректная сумма
	}

	for _, tc := range testCases {
		fmt.Printf("Запрос: %d\n", tc)
		money, err := getMoney(tc)
		if err != nil {
			fmt.Printf("  Ошибка: %v\n", err)
		} else {
			fmt.Printf("  Результат: %v\n", money)
		}
		fmt.Println("--------------------")
	}
}
