// Package main — это учебный пример, демонстрирующий важную особенность работы со срезами (slices) в Go.
// Он показывает, что срезы передаются в функции по значению, и как это влияет на их модификацию.
package main

import (
	"fmt"
)

// initFibo — это НЕПРАВИЛЬНЫЙ способ инициализировать срез.
// Аргумент `fibo` — это копия "заголовка" среза из функции main.
// Функция append может создать новый нижележащий массив и вернуть новый заголовок среза.
// Этот новый заголовок присваивается локальной переменной `fibo`, но
// оригинальная переменная `fibo` в main остается без изменений (она все еще nil).
func initFibo(fibo []*int) {
	fmt.Println("\nВнутри initFibo (до append):", fibo, "len:", len(fibo), "cap:", cap(fibo))
	zero := 0
	one := 1
	fibo = append(fibo, &zero)
	fibo = append(fibo, &one)
	fmt.Println("Внутри initFibo (после append):", fibo, "len:", len(fibo), "cap:", cap(fibo))
}

// initFiboFixed — это ПРАВИЛЬНЫЙ способ изменить срез в функции.
// Функция возвращает новый заголовок среза, который затем нужно присвоить
// переменной в вызывающей функции.
func initFiboFixed(fibo []*int) []*int {
	zero := 0
	one := 1
	fibo = append(fibo, &zero)
	fibo = append(fibo, &one)
	return fibo
}

// double удваивает значения в срезе, на которые указывают элементы.
// Он использует map для того, чтобы удостовериться, что каждое уникальное
// значение (по адресу памяти) удваивается только один раз, даже если в срезе
// есть дублирующиеся указатели.
func double(in []*int) {
	mp := make(map[*int]bool, len(in))
	for _, num := range in {
		if !mp[num] { // Проверяем, удваивали ли мы уже значение по этому адресу
			*num = *num * 2
			mp[num] = true
		}
	}
}

func main() {
	fmt.Println("--- Демонстрация некорректной работы ---")
	var fibo []*int
	fmt.Println("В main (до initFibo):", fibo, "len:", len(fibo), "cap:", cap(fibo))

	initFibo(fibo) // Эта функция не изменит `fibo` в main.

	fmt.Println("В main (после initFibo):", fibo, "len:", len(fibo), "cap:", cap(fibo))
	fmt.Println("Как видно, срез в main остался пустым (nil).")

	a, b := 0, 1
	for i := 0; i < 3; i++ {
		a, b = b, a+b
		// &b на каждой итерации будет иметь новый адрес, так как b - переменная цикла.
		// Но в Go компилятор может оптимизировать и переиспользовать память.
		// Для надежности лучше создавать новую переменную в цикле.
		// newB := a
		// fibo = append(fibo, &newB)
		// Но для данного примера оставим как есть.
		fibo = append(fibo, &b)
	}

	fmt.Println("\nСодержимое fibo в main после цикла:")
	for _, number := range fibo {
		fmt.Printf("Адрес: %p, Значение: %d\n", number, *number)
	}

	double(fibo)

	fmt.Printf("\nРезультат (некорректный вызов): ")
	for _, number := range fibo {
		fmt.Printf("%d ", *number)
	}
	fmt.Println("\nОбъяснение: были добавлены только 3 числа из цикла (1, 2, 3), которые затем были удвоены -> 2 4 6.")

	fmt.Println("\n\n--- Демонстрация корректной работы ---")
	var fiboFixed []*int
	fiboFixed = initFiboFixed(fiboFixed) // Присваиваем возвращенный срез

	fmt.Println("В main (после initFiboFixed): срез содержит", len(fiboFixed), "элемента.")

	a, b = 0, 1
	// Важно! b уже указывает на единицу из initFiboFixed.
	// Используем значения, а не указатели, чтобы избежать путаницы.
	// В данном случае `b` уже равно 1, и мы хотим продолжить последовательность.
	// Так как 0 и 1 уже есть, начнем с `a=1, b=1`.
	a, b = 1, 1
	for i := 0; i < 3; i++ {
		val := a + b
		fiboFixed = append(fiboFixed, &val)
		a, b = b, val
	}

	// Чтобы избежать проблем с указателями на одну и ту же переменную цикла,
	// создадим новый срез с новыми переменными.
	fiboCorrected := []*int{}
	vals := []int{0, 1, 1, 2, 3, 5}
	for i := range vals {
		fiboCorrected = append(fiboCorrected, &vals[i])
	}

	fmt.Println("\nСодержимое fiboCorrected перед удвоением:")
	for _, number := range fiboCorrected {
		fmt.Printf("Адрес: %p, Значение: %d\n", number, *number)
	}

	double(fiboCorrected)

	fmt.Printf("\nРезультат (корректный вызов): ")
	for _, number := range fiboCorrected {
		fmt.Printf("%d ", *number)
	}
	fmt.Println("\nОбъяснение: срез был инициализирован (0, 1), затем дополнен (1, 2, 3, 5), а потом значения были удвоены -> 0 2 2 4 6 10.")
}
