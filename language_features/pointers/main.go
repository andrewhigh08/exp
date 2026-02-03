// Package main демонстрирует ключевые аспекты работы с указателями в Go.
// В частности, он показывает, что в Go всё, включая указатели, передается по значению.
package main

import "fmt"

type Person struct {
	name string
	age  uint8
}

// --- Демонстрация 1: Неправильная попытка изменить указатель ---

// changeLocalPointer получает КОПИЮ указателя на Person.
func changeLocalPointer(person *Person) {
	fmt.Printf("  Внутри changeLocalPointer: полученный указатель указывает на -> %v\n", person)
	// Эта строка переназначает ЛОКАЛЬНУЮ переменную `person`.
	// Теперь она указывает на совершенно новый объект Person в памяти.
	// Оригинальный указатель в функции main остается без изменений.
	person = &Person{
		name: "Владимир",
		age:  25,
	}
	fmt.Printf("  Внутри changeLocalPointer: локальный указатель теперь указывает на -> %v\n", person)
}

// --- Демонстрация 2: Правильное изменение данных по указателю ---

// modifyUnderlyingData получает КОПИЮ указателя, но использует ее для
// доступа и изменения ОРИГИНАЛЬНОГО объекта Person.
func modifyUnderlyingData(person *Person) {
	// Мы не меняем сам указатель, а меняем поля в структуре, на которую он указывает.
	person.name = "Владимир"
	person.age = 25
}

// --- Демонстрация 3: Изменение самого оригинального указателя ---

// changeOriginalPointer — это продвинутый способ, который позволяет изменить
// указатель в вызывающей функции. Для этого используется указатель на указатель (**Person).
func changeOriginalPointer(person **Person) {
	// Разыменовывая `person` один раз (`*person`), мы получаем доступ
	// к оригинальному указателю из main и можем его переназначить.
	*person = &Person{
		name: "Алексей",
		age:  40,
	}
}

func main() {
	// --- Сценарий 1 ---
	fmt.Println("--- 1. Попытка изменить указатель напрямую (не сработает) ---")
	person1 := &Person{name: "Иван", age: 30}
	fmt.Printf("До вызова:  %s, %d, адрес: %p\n", person1.name, person1.age, person1)
	changeLocalPointer(person1)
	fmt.Printf("После вызова: %s, %d, адрес: %p (не изменился)\n", person1.name, person1.age, person1)

	// --- Сценарий 2 ---
	fmt.Println("\n--- 2. Изменение данных по указателю (идиоматичный способ) ---")
	person2 := &Person{name: "Иван", age: 30}
	fmt.Printf("До вызова:  %s, %d\n", person2.name, person2.age)
	modifyUnderlyingData(person2)
	fmt.Printf("После вызова: %s, %d (значения изменились)\n", person2.name, person2.age)

	// --- Сценарий 3 ---
	fmt.Println("\n--- 3. Изменение самого указателя через двойной указатель ---")
	person3 := &Person{name: "Иван", age: 30}
	fmt.Printf("До вызова:  %s, %d, адрес: %p\n", person3.name, person3.age, person3)
	// Передаем адрес нашего указателя `person3`
	changeOriginalPointer(&person3)
	fmt.Printf("После вызова: %s, %d, адрес: %p (указатель теперь другой)\n", person3.name, person3.age, person3)
}
