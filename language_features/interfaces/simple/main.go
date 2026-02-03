// Package main демонстрирует базовые принципы работы с интерфейсами в Go.
package main

import "fmt"

// --- 1. Определение интерфейсов ---

// Интерфейс в Go — это набор сигнатур методов. Он определяет "контракт" или поведение.
// Любой тип, который реализует все методы интерфейса, неявно удовлетворяет этому интерфейсу.

// Walker определяет поведение "уметь ходить".
type Walker interface {
	Walk() string
}

// Talker определяет поведение "уметь говорить".
type Talker interface {
	Talk() string
}

// TalkWalker — пример вложенного интерфейса.
// Он объединяет в себе методы из других интерфейсов.
// Тип, удовлетворяющий TalkWalker, должен реализовать и Walk(), и Talk().
type TalkWalker interface {
	Walker
	Talker
}

// --- 2. Конкретные реализации ---

// Human — это конкретный тип (struct).
type Human struct {
	Age  int
	Name string
}

// Talk реализует метод из интерфейса Talker.
// Ресивер `(h Human)` — это ресивер-значение. Метод вызывается на копии объекта.
func (h Human) Talk() string {
	return fmt.Sprintf("Меня зовут %s", h.Name)
}

// Walk реализует метод из интерфейса Walker.
func (h Human) Walk() string {
	return fmt.Sprintf("%s идет", h.Name)
}

// Dog - еще один конкретный тип.
type Dog struct {
	Name string
}

// Bark — метод для Dog с ресивером-указателем.
func (d *Dog) Bark() string {
	return "Гав! Меня зовут " + d.Name
}

// Barker - интерфейс для тех, кто умеет лаять.
type Barker interface {
	Bark() string
}

// --- 3. Полиморфизм через интерфейсы ---

// activity — это полиморфная функция. Она может принять любой тип,
// который удовлетворяет интерфейсу TalkWalker, не зная его конкретной реализации.
func activity(tw TalkWalker) {
	fmt.Println("--- Начало активности ---")
	fmt.Println(tw.Talk())
	fmt.Println(tw.Walk())
	fmt.Println("--- Конец активности ---")
}

func main() {
	// Создаем экземпляр конкретного типа Human.
	vasya := Human{33, "Вася"}

	// Мы можем передать `vasya` в функцию `activity`, потому что
	// тип Human реализует все методы интерфейса TalkWalker.
	activity(vasya)

	// --- 4. Ресивер-значение vs. Ресивер-указатель ---
	fmt.Println("\n--- Демонстрация ресиверов ---")

	// Если методы реализованы для ресивера-значения (T), то и значение (T),
	// и указатель на него (*T) удовлетворяют интерфейсу.
	// Компилятор автоматически берет значение или разыменовывает указатель.
	petya := Human{21, "Петя"}
	var petyaWalker Walker = petya // OK
	fmt.Println("Реализация через значение (petya):", petyaWalker.Walk())
	var petyaPtrWalker Walker = &petya // Тоже OK
	fmt.Println("Реализация через указатель (&petya):", petyaPtrWalker.Walk())


	// Рассмотрим другой случай, с ресивером-указателем.
	dog := Dog{"Шарик"}
	var dogPtrBarker Barker = &dog // OK: *Dog имеет метод Bark()
	fmt.Println("Реализация через указатель (&dog):", dogPtrBarker.Bark())

	// var dogBarker Barker = dog // !!! ОШИБКА КОМПИЛЯЦИИ !!!
	// Код выше не скомпилируется, потому что тип `Dog` (значение) НЕ имеет
	// метода Bark(). Метод Bark() есть только у типа `*Dog` (указателя).
	// Компилятор не может автоматически взять адрес `dog`, чтобы удовлетворить интерфейс.
	fmt.Println("Экземпляр Dog (значение) НЕ удовлетворяет интерфейсу Barker, так как метод Bark() определен для *Dog (указателя).")
}
