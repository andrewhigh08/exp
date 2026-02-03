<h2>Задание: Реализовать библиотеку реализующую key-value хранилище.</h2>
<h3>Требования к реализации.</h3>
<p>1). Хранилище должно быть in-memory.
2). Должны поддерживаться все типы данных, т.е. хранилище должно уметь работать с ключами и значениями типа int (empty interface).<br />
3). Должны поддерживаться следующие операции: вставка, поиск, обновление и удаление (похожее реализовано в sync.Map, CRUD).<br />
4). Хранилище должно быть потокобезопасное (согласованное чтение и запись из разных горутин, синхронизация с sync.RWMutex, concurrency).<br />
5). Хранилище должно предоставлять интерфейс для альтернативной реализации (определить интерфейс для хранилища и методы к нему, чтобы, например, другой разработчик смог реализовать другое хранилище).</p>

```go
package kvstore

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrKeyExists   = errors.New("key already exists")
	ErrKeyNotFound = errors.New("key not found")
)

// Определяем интерфейс для хранилища
type KeyValueStore interface {
	Create(ctx context.Context, key int, value interface{}) error
	Read(ctx context.Context, key int) (interface{}, bool, error)
	Update(ctx context.Context, key int, value interface{}) error
	Delete(ctx context.Context, key int) error
}

// Реализуем тип, который будет поддерживать этот интерфейс
type InMemoryStore struct {
	mu    sync.RWMutex
	store map[int]interface{}
}

// Конструктор для хранилища
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		store: make(map[int]interface{}),
	}
}

// Реализация метода Create (создание) с поддержкой контекста
func (kv *InMemoryStore) Create(ctx context.Context, key int, value interface{}) error {
	select {
	case <-ctx.Done():
		return ctx.Err() // Возвращаем ошибку, если контекст отменён
	default:
		kv.mu.Lock()
		defer kv.mu.Unlock()
		if _, exists := kv.store[key]; exists {
			return ErrKeyExists // Если ключ уже существует, возвращаем ошибку
		}
		kv.store[key] = value
		return nil
	}
}

// Реализация метода Read (чтение) с поддержкой контекста
func (kv *InMemoryStore) Read(ctx context.Context, key int) (interface{}, bool, error) {
	select {
	case <-ctx.Done():
		return nil, false, ctx.Err() // Возвращаем ошибку, если контекст отменён
	default:
		kv.mu.RLock()
		defer kv.mu.RUnlock()
		value, exists := kv.store[key]
		return value, exists, nil
	}
}

// Реализация метода Update (обновление) с поддержкой контекста
func (kv *InMemoryStore) Update(ctx context.Context, key int, value interface{}) error {
	select {
	case <-ctx.Done():
		return ctx.Err() // Возвращаем ошибку, если контекст отменён
	default:
		kv.mu.Lock()
		defer kv.mu.Unlock()
		if _, exists := kv.store[key]; !exists {
			return ErrKeyNotFound // Если ключ не существует, возвращаем ошибку
		}
		kv.store[key] = value
		return nil
	}
}

// Реализация метода Delete (удаление) с поддержкой контекста
func (kv *InMemoryStore) Delete(ctx context.Context, key int) error {
	select {
	case <-ctx.Done():
		return ctx.Err() // Возвращаем ошибку, если контекст отменён
	default:
		kv.mu.Lock()
		defer kv.mu.Unlock()
		if _, exists := kv.store[key]; !exists {
			return ErrKeyNotFound // Если ключ не существует, возвращаем ошибку
		}
		delete(kv.store, key)
		return nil
	}
}
```

```go
package main

import (
    "context"
    "fmt"
    "time"
    "kvstore"
)

func main() {
    store := kvstore.NewInMemoryStore()
    
    // Пример с использованием контекста
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    // Создание нового ключа
    err := store.Create(ctx, 1, "value1")
    if err != nil {
        fmt.Println("Ошибка при создании:", err)
    }
    
    // Пытаемся создать ключ, который уже существует
    err = store.Create(ctx, 1, "value2")
    if err != nil {
        fmt.Println("Ошибка при создании:", err)
    }
    
    // Чтение значения
    value, exists, err := store.Read(ctx, 1)
    if err != nil {
        fmt.Println("Ошибка при чтении:", err)
    } else if exists {
        fmt.Println("Значение:", value)
    } else {
        fmt.Println("Ключ не найден")
    }
    
    // Обновление существующего ключа
    err = store.Update(ctx, 1, "updated_value")
    if err != nil {
        fmt.Println("Ошибка при обновлении:", err)
    }

    // Удаление ключа
    err = store.Delete(ctx, 1)
    if err != nil {
        fmt.Println("Ошибка при удалении:", err)
    }
}
```
