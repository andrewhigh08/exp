// Package main демонстрирует "сложные" и не всегда очевидные аспекты работы с интерфейсами в Go.
// 1. Внутреннее представление интерфейсов (nil-интерфейс vs интерфейс с nil-значением).
// 2. Полиморфизм и безопасное приведение типов (type assertion).
package main

import "fmt"

// --- Определения интерфейсов и типов ---

type Runner interface {
	Run() string
}

type Swimmer interface {
	Swim() string
}

type Flyer interface {
	Fly() string
}

// Ducker — пример вложенного интерфейса, который объединяет несколько других.
type Ducker interface {
	Runner
	Swimmer
	Flyer
}

// Human реализует только интерфейс Runner.
type Human struct {
	Name string
}

func (h *Human) Run() string {
	return fmt.Sprintf("Человек %s бежит", h.Name)
}

func (h *Human) writeCode() {
	fmt.Println("Человек пишет код...")
}

// Duck реализует все три интерфейса (Runner, Swimmer, Flyer) и, следовательно, Ducker.
type Duck struct {
	Name string
}

func (d *Duck) Run() string {
	return fmt.Sprintf("Утка %s бежит", d.Name)
}

func (d *Duck) Swim() string {
	return fmt.Sprintf("Утка %s плывет", d.Name)
}

func (d *Duck) Fly() string {
	return fmt.Sprintf("Утка %s летит", d.Name)
}

// --- Демонстрация 1: Значение интерфейса ---

func interfaceValues() {
	fmt.Println("--- Демонстрация 1: Внутреннее устройство интерфейса ---")
	// Переменная интерфейсного типа состоит из двух компонентов: (тип, значение).
	// `type` — это конкретный тип, который хранится в интерфейсе.
	// `value` — это конкретное значение этого типа.

	// 1. `nil`-интерфейс
	var runner Runner
	fmt.Printf("1. nil-интерфейс: Тип=%T, Значение=%#v\n", runner, runner)
	if runner == nil {
		fmt.Println("   'runner' является nil, так как и тип, и значение в нем — nil.")
	}

	// 2. Интерфейс, содержащий `nil`-указатель
	var unnamedRunner *Human // Это nil-указатель на Human.
	fmt.Printf("2. nil-указатель: Тип=%T, Значение=%#v\n", unnamedRunner, unnamedRunner)

	runner = unnamedRunner // Присваиваем nil-указатель интерфейсу.
	// ВАЖНО: Теперь интерфейс НЕ nil!
	// Его внутреннее представление: (тип: *Human, значение: nil).
	fmt.Printf("3. Интерфейс с nil-указателем: Тип=%T, Значение=%#v\n", runner, runner)
	if runner != nil {
		fmt.Println("   'runner' НЕ является nil, потому что его компонент 'тип' НЕ nil (*Human).")
		// При попытке вызвать метод такого интерфейса произойдет паника,
		// так как мы пытаемся вызвать метод у nil-указателя.
		// runner.Run() // это вызовет панику: panic: runtime error: invalid memory address or nil pointer dereference
	}

	// 3. Интерфейс, содержащий конкретное значение
	namedRunner := &Human{Name: "Джек"}
	runner = namedRunner
	// Внутреннее представление: (тип: *Human, значение: адрес объекта 'namedRunner')
	fmt.Printf("4. Интерфейс с конкретным значением: Тип=%T, Значение=%#v\n", runner, runner)
}

// --- Демонстрация 2: Полиморфизм и Приведение Типов ---

// polymorphism демонстрирует полиморфное поведение: функция принимает любой тип,
// удовлетворяющий интерфейсу Runner, и работает с ним через этот интерфейс.
func polymorphism(r Runner) {
	fmt.Println("Полиморфный вызов:", r.Run())
}

// typeAssertion демонстрирует, как "извлечь" конкретный тип из интерфейса.
func typeAssertion(r Runner) {
	fmt.Println("\n-- Проверка типов для:", r.Run())

	// 1. Приведение типа с проверкой ("comma-ok" idiom)
	// Это безопасный способ проверить, содержит ли интерфейс значение нужного нам типа.
	if human, ok := r.(*Human); ok {
		fmt.Println("   Это Human! Можно вызвать его уникальный метод:")
		human.writeCode()
	} else {
		fmt.Println("   Это не Human.")
	}

	if duck, ok := r.(*Duck); ok {
		fmt.Println("   Это Duck! Можно вызвать его уникальный метод:")
		fmt.Println("  ", duck.Fly())
	} else {
		fmt.Println("   Это не Duck.")
	}

	// 2. `switch` по типу (type switch)
	// Это идиоматичный способ обработать несколько возможных конкретных типов.
	fmt.Println("   Проверка через type switch:")
	switch v := r.(type) {
	case *Human:
		fmt.Printf("   Тип в switch: %T. Имя: %s\n", v, v.Name)
	case *Duck:
		fmt.Printf("   Тип в switch: %T. Имя: %s\n", v, v.Name)
	default:
		fmt.Printf("   Неизвестный тип: %T\n", v)
	}
}

func main() {
	interfaceValues()

	fmt.Println("\n\n--- Демонстрация 2: Полиморфизм и Приведение Типов ---")
	// Создаем экземпляры конкретных типов
	john := &Human{"Джон"}
	donald := &Duck{"Дональд"}

	// Полиморфно передаем их в функции
	polymorphism(john)
	polymorphism(donald)

	// Проверяем, как работает приведение типов для каждого из них
	typeAssertion(john)
	typeAssertion(donald)
}
