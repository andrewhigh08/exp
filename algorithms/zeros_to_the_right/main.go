// Package main содержит решение задачи по перемещению всех нулевых элементов среза в его конец.
// При этом относительный порядок ненулевых элементов должен быть сохранен.
package main

import "fmt"

// moveZerosToEndNewSlice решает задачу путем создания нового среза.
//
// Алгоритм:
// 1. Создает новый пустой срез `result`.
// 2. Проходит по исходному срезу, копируя все ненулевые элементы в `result`.
// 3. Подсчитывает количество нулей.
// 4. Добавляет в конец `result` необходимое количество нулей.
//
// Плюсы: Простой и понятный код.
// Минусы: Требует дополнительной памяти O(n) для создания нового среза.
func moveZerosToEndNewSlice(input []int) []int {
	length := len(input)
	// Создаем новый срез с предвыделенной емкостью для эффективности.
	result := make([]int, 0, length)
	countZero := 0

	for _, value := range input {
		if value != 0 {
			result = append(result, value)
		} else {
			countZero++
		}
	}

	// Добавляем нули в конец.
	for i := 0; i < countZero; i++ {
		result = append(result, 0)
	}

	return result
}

// moveZerosToEndInPlace решает задачу "на месте" (in-place) без выделения дополнительной памяти.
//
// Алгоритм (метод двух указателей):
// 1. Используем указатель `insertPos` (или "снежный ком"), который указывает на позицию,
//    куда следует поместить следующий ненулевой элемент.
// 2. Итерируемся по срезу. Когда встречаем ненулевой элемент,
//    мы помещаем его в позицию `insertPos` и сдвигаем `insertPos`.
// 3. После первого прохода все ненулевые элементы будут собраны в начале среза
//    в правильном порядке.
// 4. Заполняем оставшуюся часть среза (с `insertPos` до конца) нулями.
//
// Плюсы: Эффективность по памяти, сложность O(1). Это предпочтительное решение на собеседованиях.
// Минусы: Модифицирует исходный срез.
func moveZerosToEndInPlace(input []int) {
	insertPos := 0

	// Перемещаем все ненулевые элементы в начало.
	for _, value := range input {
		if value != 0 {
			input[insertPos] = value
			insertPos++
		}
	}

	// Заполняем оставшуюся часть среза нулями.
	for i := insertPos; i < len(input); i++ {
		input[i] = 0
	}
}

func main() {
	testCases := [][]int{
		{0, 1, 2, 3, 1, 2, 9, 2, 3, 4, 6, 0, 0, 12, 34, 34},
		{0, 0, 0, 1, 2, 3},
		{1, 2, 3, 0, 0, 0},
		{1, 2, 3},
		{0, 0, 0},
		{},
		{4, 0, 2, 0, 1, 0, 3},
	}

	fmt.Println("--- Сравнение двух методов ---")
	for i, originalSlice := range testCases {
		fmt.Printf("\n--- Тест #%d ---\n", i+1)
		fmt.Printf("Исходный срез: %v\n", originalSlice)

		// --- Метод 1: Создание нового среза ---
		// Копируем, чтобы не изменять исходный тестовый пример
		sliceForMethod1 := make([]int, len(originalSlice))
		copy(sliceForMethod1, originalSlice)
		result1 := moveZerosToEndNewSlice(sliceForMethod1)
		fmt.Printf("Результат (новый срез): %v\n", result1)

		// --- Метод 2: "На месте" (in-place) ---
		sliceForMethod2 := make([]int, len(originalSlice))
		copy(sliceForMethod2, originalSlice)
		moveZerosToEndInPlace(sliceForMethod2)
		fmt.Printf("Результат (in-place):    %v\n", sliceForMethod2)
	}
}
