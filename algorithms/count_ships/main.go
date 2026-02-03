// Package main содержит решение задачи по подсчету "кораблей" (или "островов") на 2D-поле.
// Поле представлено в виде одномерного среза (slice).
package main

import (
	"fmt"
)

// calculateShips считает количество кораблей на поле боя.
// Корабль — это одна или несколько смежных (по горизонтали или вертикали) ячеек со значением 1.
//
// Алгоритм основан на поиске "верхних левых" частей каждого корабля.
// Ячейка считается началом нового корабля, если она содержит "1",
// а ее соседи сверху и слева — "0" (или находятся за пределами поля).
//
// @param {[]int} battleField - поле боя в виде одномерного среза.
// @param {int} width - ширина поля.
// @param {int} height - высота поля (для полноты картины, хотя в данном алгоритме не используется напрямую).
// @return {int} - количество кораблей.
func calculateShips(battleField []int, width int) (int, error) {
	if len(battleField) == 0 {
		return 0, nil
	}
	if len(battleField)%width != 0 {
		return 0, fmt.Errorf("длина поля (%d) не кратна его ширине (%d)", len(battleField), width)
	}

	shipCount := 0
	for i, cell := range battleField {
		// Пропускаем пустые ячейки ("вода")
		if cell == 0 {
			continue
		}

		// Вычисляем координаты ячейки (row, col) для лучшего понимания.
		row := i / width
		col := i % width

		// Проверяем соседа сверху. Если мы в первой строке (row == 0),
		// то соседа сверху нет, что эквивалентно "воде".
		hasTopShip := false
		if row > 0 && battleField[i-width] == 1 {
			hasTopShip = true
		}

		// Проверяем соседа слева. Если мы в первом столбце (col == 0),
		// то соседа слева нет.
		hasLeftShip := false
		if col > 0 && battleField[i-1] == 1 {
			hasLeftShip = true
		}

		// Если у ячейки с "1" нет соседей-кораблей сверху и слева,
		// значит, это "верхняя левая" ячейка нового корабля.
		if !hasTopShip && !hasLeftShip {
			shipCount++
		}
	}

	return shipCount, nil
}

func main() {
	// --- Пример 1: Поле 5x5 ---
	battleField1 := []int{
		1, 0, 0, 1, 1,
		0, 1, 0, 0, 0,
		0, 1, 0, 1, 1,
		0, 1, 0, 0, 0,
		0, 1, 0, 1, 1,
	}
	width1 := 5

	fmt.Println("--- Поле 1 (5x5) ---")
	// Визуализация поля для наглядности
	for i, cell := range battleField1 {
		if i > 0 && i%width1 == 0 {
			fmt.Println()
		}
		fmt.Printf("%d ", cell)
	}
	fmt.Println()

	shipCount1, err := calculateShips(battleField1, width1)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		// Ожидаемый результат: 4 корабля
		// 1. (0,0)
		// 2. (0,3)-(0,4)
		// 3. (1,1)-(2,1)-(3,1)-(4,1)
		// 4. (2,3)-(2,4) и (4,3)-(4,4) - это один большой корабль
		fmt.Printf("Количество кораблей на поле боя 1: %d\n", shipCount1)
	}

	fmt.Println("\n--- Поле 2 (4x3) ---")
	// --- Пример 2: Поле 4x3 ---
	battleField2 := []int{
		1, 1, 0, 0,
		0, 0, 0, 1,
		1, 1, 0, 1,
	}
	width2 := 4
	
	for i, cell := range battleField2 {
		if i > 0 && i%width2 == 0 {
			fmt.Println()
		}
		fmt.Printf("%d ", cell)
	}
	fmt.Println()

	shipCount2, err := calculateShips(battleField2, width2)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		// Ожидаемый результат: 3 корабля
		// 1. (0,0)-(0,1)
		// 2. (1,3)-(2,3)
		// 3. (2,0)-(2,1)
		fmt.Printf("Количество кораблей на поле боя 2: %d\n", shipCount2)
	}
}
