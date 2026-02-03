// Package main — это руководство по работе со срезами (slices) в Go.
// Срезы — это мощный, но иногда запутанный механизм. Этот файл демонстрирует
// их внутреннее устройство и ключевые концепции.
package main

import "fmt"

// printSliceInfo — это вспомогательная функция для наглядной демонстрации
// состояния среза: его длины, ёмкости, самих элементов и адреса первого элемента.
func printSliceInfo(name string, s []int) {
	// %p - выводит адрес указателя. &s[0] - адрес первого элемента в базовом массиве.
	// Это позволяет увидеть, ссылаются ли разные срезы на один и тот же массив.
	fmt.Printf("%s: len=%d, cap=%d, %v, addr=%p\n", name, len(s), cap(s), s, &s[0])
}

// --- Демонстрации ---

func demo1_SubSlicing() {
	fmt.Println("--- 1. Создание под-среза (sub-slicing) ---")
	// Срез — это на самом деле "заголовок" (header), который содержит:
	// 1. Указатель на базовый массив (где лежат данные).
	// 2. Длину (len) — количество видимых элементов.
	// 3. Ёмкость (cap) — общее количество элементов от начала среза до конца базового массива.

	original := []int{1, 2, 3, 4, 5, 6}
	printSliceInfo("original", original)

	// Создаем под-срез. `sub` получает НОВЫЙ заголовок, но он указывает
	// на тот же самый базовый массив, что и `original`.
	sub := original[2:4]
	printSliceInfo("sub     ", sub)
	fmt.Println("=> `sub` указывает на тот же базовый массив, что и `original`.\n")

	// Изменение элемента в под-срезе меняет и оригинальный срез!
	fmt.Println("Изменяем sub[0] = 99")
	sub[0] = 99
	printSliceInfo("original", original)
	printSliceInfo("sub     ", sub)
	fmt.Println("=> Изменение в `sub` отразилось в `original`, так как массив у них общий.")
}

func demo2_AppendWithCapacity() {
	fmt.Println("\n--- 2. `append` при наличии свободной ёмкости ---")
	original := []int{1, 2, 3, 4, 5, 6}
	printSliceInfo("original", original)

	// sub имеет len=2, но cap=4, так как в базовом массиве после него есть еще элементы (5, 6).
	sub := original[2:4]
	printSliceInfo("sub     ", sub)

	// Добавляем элемент в `sub`. Так как есть свободная ёмкость,
	// Go не будет создавать новый массив, а запишет значение `100`
	// в следующую ячейку базового массива, на которую указывает `original`.
	fmt.Println("\nВыполняем `sub = append(sub, 100)`")
	sub = append(sub, 100)

	printSliceInfo("original", original) // `original` изменился! Элемент `5` был перезаписан.
	printSliceInfo("sub     ", sub)
	fmt.Println("=> `append` к `sub` изменил базовый массив, что повлияло на `original`!")
}

func demo3_AppendWithReallocation() {
	fmt.Println("\n--- 3. `append`, вызывающий реаллокацию ---")
	original := []int{1, 2, 3, 4}
	printSliceInfo("original", original)

	// `sub` использует весь остаток базового массива. len=2, cap=2.
	sub := original[2:4]
	printSliceInfo("sub     ", sub)

	// Добавляем элемент в `sub`. Так как свободной ёмкости нет (`len == cap`),
	// Go создает СОВЕРШЕННО НОВЫЙ базовый массив (обычно удвоенной ёмкости),
	// копирует в него старые элементы `sub` и добавляет новый.
	fmt.Println("\nВыполняем `sub = append(sub, 200)`")
	sub = append(sub, 200)

	printSliceInfo("original", original) // `original` не изменился!
	printSliceInfo("sub     ", sub)      // `sub` теперь указывает на новый адрес.
	fmt.Println("=> `append` создал новый массив для `sub`. `original` не затронут.")
}

// --- Безопасные и "in-place" функции ---

// doubleSafe создает и возвращает новый срез.
// Это "безопасный" способ, так как он никогда не изменит исходные данные.
func doubleSafe(nums []int) []int {
	res := make([]int, len(nums)) // Создаем новый срез (и новый базовый массив).
	for i, v := range nums {
		res[i] = v * 2
	}
	return res
}

// doubleInPlace изменяет элементы исходного среза "на месте".
// Это эффективно по памяти, но может приводить к побочным эффектам.
func doubleInPlace(nums []int) {
	for i := range nums {
		nums[i] *= 2
	}
}

func demo4_Functions() {
	fmt.Println("\n--- 4. Функции, работающие со срезами ---")
	data := []int{10, 20, 30}
	fmt.Println("Исходный срез:", data)

	safeResult := doubleSafe(data)
	fmt.Println("Результат `doubleSafe`:", safeResult)
	fmt.Println("Исходный срез после `doubleSafe`:", data, "(не изменился)")

	fmt.Println()
	doubleInPlace(data)
	fmt.Println("Исходный срез после `doubleInPlace`:", data, "(изменился)")
}

func main() {
	demo1_SubSlicing()
	demo2_AppendWithCapacity()
	demo3_AppendWithReallocation()
	demo4_Functions()
}
