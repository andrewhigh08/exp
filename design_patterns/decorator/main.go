// Package main демонстрирует реализацию структурного паттерна "Декоратор".
//
// Паттерн "Декоратор" (Decorator) позволяет динамически добавлять объектам
// новую функциональность, оборачивая их в "обертки". Декораторы предоставляют
// ту же самую функциональность, что и оборачиваемые объекты, плюс что-то свое.
//
// Ключевая идея: Декоратор и оборачиваемый объект реализуют один и тот же интерфейс.
//
// Компоненты паттерна:
// 1. Component: Общий интерфейс для всех объектов. (`DB` в нашем примере)
// 2. ConcreteComponent: Базовая реализация, которую мы хотим "украсить". (`PostgresDB`)
// 3. Decorator: Абстрактный класс или структура, которая содержит ссылку на
//    объект Component и реализует его интерфейс.
// 4. ConcreteDecorator: Конкретная реализация декоратора, добавляющая свою логику. (`RedisCacheDecorator`)
package main

import (
	"fmt"
	"sync"
	"time"
)

// DB — это общий интерфейс Component.
type DB interface {
	Query(query string) string
}

// --- Конкретный компонент ---

// PostgresDB — это ConcreteComponent, базовая реализация.
type PostgresDB struct{}

func (db *PostgresDB) Query(query string) string {
	// Имитация долгого запроса к реальной базе данных.
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Выполняю запрос к PostgreSQL...")
	return "Результат из PostgreSQL для запроса: " + query
}

// --- Конкретный декоратор ---

// RedisCacheDecorator — это ConcreteDecorator. Он добавляет кэширование.
type RedisCacheDecorator struct {
	// Декоратор "оборачивает" другой объект, который тоже реализует интерфейс DB.
	// Это может быть как базовый PostgresDB, так и другой декоратор.
	DB DB

	// Дополнительное состояние и функциональность.
	Cache map[string]string // Имитация кеша Redis
	mu    sync.RWMutex      // Мьютекс для потокобезопасного доступа к кешу.
}

// NewRedisCacheDecorator — конструктор для нашего декоратора.
func NewRedisCacheDecorator(db DB) *RedisCacheDecorator {
	return &RedisCacheDecorator{
		DB:    db,
		Cache: make(map[string]string),
	}
}

// Query — реализация метода интерфейса DB. Здесь и происходит "декорирование".
func (r *RedisCacheDecorator) Query(query string) string {
	// 1. Добавленная логика: проверяем наличие в кеше.
	r.mu.RLock()
	if cachedResult, ok := r.Cache[query]; ok {
		r.mu.RUnlock()
		fmt.Println("Результат найден в Redis кеше!")
		return cachedResult
	}
	r.mu.RUnlock()

	// 2. Если в кеше нет, вызываем метод оборачиваемого объекта.
	fmt.Println("В кеше не найдено, обращаемся к базе данных...")
	result := r.DB.Query(query)

	// 3. Еще одна добавленная логика: сохраняем результат в кеш.
	fmt.Println("Сохраняем результат в кеш...")
	r.mu.Lock()
	r.Cache[query] = result
	r.mu.Unlock()

	return result
}

func main() {
	// 1. Создаем базовый объект (ConcreteComponent).
	db := &PostgresDB{}

	// 2. "Украшаем" (оборачиваем) его декоратором кеширования.
	// `cachedDB` теперь имеет тот же интерфейс, что и `db`, но с дополнительной логикой.
	cachedDB := NewRedisCacheDecorator(db)

	fmt.Println("--- Первый запрос (ожидается обращение к БД) ---")
	result1 := cachedDB.Query("SELECT * FROM users WHERE id = 1")
	fmt.Printf("Результат: %s\n\n", result1)

	fmt.Println("--- Второй, идентичный запрос (ожидается результат из кеша) ---")
	result2 := cachedDB.Query("SELECT * FROM users WHERE id = 1")
	fmt.Printf("Результат: %s\n\n", result2)

	// Можно создавать цепочки декораторов. Например, добавить декоратор для логирования:
	// loggedAndCachedDB := NewLoggingDecorator(cachedDB)
	// loggedAndCachedDB.Query("SELECT * FROM products")
}
