// Package main демонстрирует основы дженериков (обобщенного программирования) в Go,
// введенных в версии 1.18. Дженерики позволяют писать функции и типы, которые
// могут работать с любым типом из заданного набора. С Go 1.27 собственные
// параметры типа могут объявлять и методы (generic-методы, см. демо 6).
package main

import (
	"fmt"
)

// Number — это интерфейс, который используется как "ограничение" (constraint) для дженериков.
// Он определяет набор типов, с которыми может работать обобщенная функция.
// `~int64` означает "любой тип, чей базовый (underlying) тип - это int64".
type Number interface {
	~int64 | ~float64 // Поддерживает int64, float64 и любые типы на их основе (например, CustomInt).
}

// CustomInt — это пользовательский тип на основе int64.
// Он удовлетворяет ограничению `~int64` из интерфейса `Number`.
type CustomInt int64

func (ci CustomInt) IsPositive() bool {
	return ci > 0
}

// Numbers — это обобщенный (generic) тип.
// Он представляет собой срез элементов типа `T`, где `T` должен удовлетворять ограничению `Number`.
type Numbers[T Number] []T

// Map — generic-метод (Go 1.27+): кроме параметра типа `T`, пришедшего от получателя,
// у него есть собственный параметр типа `U`. До Go 1.27 методы не могли объявлять
// свои параметры типа, и такое преобразование приходилось писать свободной функцией:
// func mapNumbers[T Number, U any](ns Numbers[T], f func(T) U) []U.
func (ns Numbers[T]) Map[U any](f func(T) U) []U {
	out := make([]U, 0, len(ns))
	for _, n := range ns {
		out = append(out, f(n))
	}
	return out
}

// --- Обобщенные (Generic) Функции ---

// sum — простая обобщенная функция.
// [V int64 | float64] — это "параметры типа" (type parameters).
// `V` может быть либо `int64`, либо `float64` — ровно эти типы: без `~` пользовательский
// CustomInt (базовый тип — int64) сюда не подойдет, см. демо 5.
func sum[V int64 | float64](numbers []V) V {
	var sum V // `sum` будет иметь тот же тип, что и элементы среза.
	for _, num := range numbers {
		sum += num
	}
	return sum
}

// contains — обобщенная функция с ограничением `comparable`.
// `comparable` — это встроенное ограничение, которому удовлетворяют все типы,
// для которых разрешена операция сравнения `==` и `!=` (числа, строки, указатели, структуры и т.д.,
// но не срезы, карты или функции).
// С Go 1.20 ему удовлетворяют и интерфейсы, в том числе `any` (раньше — нет). Но для них `==`
// безопасен не всегда: если динамические типы совпадают и несравнимы, будет паника во время
// выполнения — contains([]any{[]int{1}}, any([]int{1})) компилируется, но падает с
// "runtime error: comparing uncomparable type []int".
func contains[T comparable](elements []T, searchEl T) bool {
	for _, el := range elements {
		if searchEl == el {
			return true
		}
	}
	return false
}

// sumUnionInterface — это версия `sum`, использующая кастомный интерфейс `Number` в качестве ограничения.
// Именованное ограничение удобно переиспользовать, а благодаря `~` в `Number` эта версия гибче `sum`:
// она принимает и пользовательские типы на основе int64/float64 (см. демо 5).
func sumUnionInterface[V Number](numbers []V) V {
	var sum V
	for _, num := range numbers {
		sum += num
	}
	return sum
}

// show — обобщенная функция с ограничением `any`.
// `any` — это встроенный псевдоним для `interface{}`. Он означает, что `T` может быть абсолютно любым типом.
func show[T any](entities ...T) {
	fmt.Printf("Тип: %T, Значения: %v\n", entities, entities)
}

// --- Демонстрационные функции ---

func demoSum() {
	fmt.Println("--- 1. Обобщенная функция `sum` ---")
	floats := []float64{1.1, 2.2, 3.3}
	ints := []int64{1, 2, 3}

	fmt.Println("Сумма float64:", sum(floats))
	fmt.Println("Сумма int64:", sum(ints))
	// Можно явно указать тип, но компилятор обычно выводит его сам.
	fmt.Println("Сумма int64 (явно):", sum[int64](ints))
}

func demoContains() {
	fmt.Println("\n--- 2. Обобщенная функция `contains` с ограничением `comparable` ---")
	type Person struct {
		Name string
		Age  int64
	}

	ints := []int64{1, 2, 3, 4, 5}
	fmt.Println("Содержит ли срез {1,2,3,4,5} число 4?:", contains(ints, 4))

	strings := []string{"Вася", "Дима", "Катя"}
	fmt.Println("Содержит ли срез {Вася, Дима, Катя} строку 'Катя'?:", contains(strings, "Катя"))
	fmt.Println("Содержит ли срез {Вася, Дима, Катя} строку 'Саша'?:", contains(strings, "Саша"))

	// Структуры также являются `comparable`, если все их поля `comparable`.
	people := []Person{
		{Name: "Вася", Age: 20},
		{Name: "Даша", Age: 23},
	}
	fmt.Println("Содержит ли срез Person{Вася, 20}?:", contains(people, Person{Name: "Вася", Age: 20}))
	fmt.Println("Содержит ли срез Person{Вася, 21}?:", contains(people, Person{Name: "Вася", Age: 21}))
}

func demoAny() {
	fmt.Println("\n--- 3. Обобщенная функция `show` с ограничением `any` ---")
	show(1, 2, 3)
	show("test1", "test2")
	show(map[string]int64{"first": 1})
}

func demoUnionInterface() {
	fmt.Println("\n--- 4. Использование интерфейса как ограничения ---")
	ints := Numbers[int64]{1, 2, 3, 4, 5}
	floats := Numbers[float64]{1.0, 2.5, 3.5}

	fmt.Println("Сумма (Numbers[int64]):", sumUnionInterface(ints))
	fmt.Println("Сумма (Numbers[float64]):", sumUnionInterface(floats))
}

func demoTypeApproximation() {
	fmt.Println("\n--- 5. Приближение типа (Type Approximation) с помощью `~` ---")
	// `CustomInt` — это наш кастомный тип, но его базовый тип — `int64`.
	customInts := []CustomInt{10, 20, 30}

	// Благодаря `~int64` в интерфейсе `Number`, наша функция `sumUnionInterface`
	// может работать с `[]CustomInt` напрямую, без преобразования типов.
	// Это мощный механизм для работы с пользовательскими типами.
	fmt.Println("Сумма `[]CustomInt` напрямую:", sumUnionInterface(customInts))

	// А `sum` с ограничением `int64 | float64` (без `~`) такой срез не примет:
	// sum(customInts) // ошибка компиляции: CustomInt does not satisfy int64 | float64 (possibly missing ~ for int64 in int64 | float64)
}

func demoGenericMethods() {
	fmt.Println("\n--- 6. Generic-методы (Go 1.27+) ---")
	ints := Numbers[int64]{1, 2, 3}

	// `U` выводится из аргумента-функции: здесь `U` = string.
	labels := ints.Map(func(n int64) string { return fmt.Sprintf("#%d", n) })
	fmt.Println("Numbers[int64] -> []string:", labels)

	// Параметр типа метода можно указать и явно — так же, как у обобщенной функции.
	halves := ints.Map[float64](func(n int64) float64 { return float64(n) / 2 })
	fmt.Println("Numbers[int64] -> []float64:", halves)

	// Ограничения: у методов интерфейса параметров типа быть не может, а generic-метод
	// не реализует метод интерфейса. Например, Numbers[int64] не удовлетворяет
	// interface{ Map(func(int64) string) []string } — ошибка компиляции
	// "does not implement ... (wrong type for method Map)".
}

func main() {
	demoSum()
	demoContains()
	demoAny()
	demoUnionInterface()
	demoTypeApproximation()
	demoGenericMethods()
}
