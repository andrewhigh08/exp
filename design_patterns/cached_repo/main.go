// Package main демонстрирует реализацию паттерна "Декоратор" (Decorator)
// для добавления функциональности кэширования к существующему репозиторию данных.
//
// Паттерн позволяет динамически добавлять объектам новую функциональность,
// оборачивая их в полезные "обертки".
//
// Здесь `CachedRepository` является декоратором для любого объекта, реализующего
// интерфейс `Repository`, добавляя ему слой in-memory кэширования.
package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// Repository определяет общий интерфейс для доступа к данным.
// Это может быть база данных, внешний API и т.д.
type Repository interface {
	Get(key string) (string, error)
	MGet(keys ...string) ([]string, error)
	Set(key, value string) error
	Del(key string) error
}

// --- Декоратор: Кэширующий репозиторий ---

// CachedRepository — это декоратор, который добавляет кэширование.
// Он реализует тот же интерфейс `Repository`, что и оборачиваемый объект.
type CachedRepository struct {
	repo  Repository        // Оборачиваемый репозиторий (например, БД)
	cache map[string]string // In-memory кэш
	mu    sync.RWMutex      // Мьютекс для потокобезопасного доступа к кэшу
}

// NewCachedRepository создает новый экземпляр кэширующего репозитория.
func NewCachedRepository(repo Repository) *CachedRepository {
	return &CachedRepository{
		repo:  repo,
		cache: make(map[string]string),
	}
}

// Get реализует стратегию "Cache-Aside".
// 1. Попробовать найти в кэше.
// 2. Если в кэше нет -> загрузить из основного репозитория.
// 3. Поместить загруженное значение в кэш.
// 4. Вернуть значение.
func (c *CachedRepository) Get(key string) (string, error) {
	// Сначала проверяем кэш с блокировкой на чтение (RLock),
	// чтобы не мешать другим читателям.
	c.mu.RLock()
	if value, ok := c.cache[key]; ok {
		c.mu.RUnlock()
		fmt.Printf("[CACHE HIT] Get key: %s\n", key)
		return value, nil
	}
	// Важно отпустить блокировку чтения перед тем, как делать что-то еще.
	c.mu.RUnlock()

	fmt.Printf("[CACHE MISS] Get key: %s -> fetching from DB\n", key)
	// Если в кэше нет, загружаем из основного репозитория.
	value, err := c.repo.Get(key)
	if err != nil {
		return "", err
	}

	// Сохраняем значение в кэше с эксклюзивной блокировкой на запись.
	c.mu.Lock()
	c.cache[key] = value
	c.mu.Unlock()

	return value, nil
}

// MGet выполняет пакетное получение данных.
// Он эффективно находит ключи, которых нет в кэше, и запрашивает только их.
func (c *CachedRepository) MGet(keys ...string) ([]string, error) {
	results := make([]string, len(keys))
	missingKeys := make([]string, 0)
	// Создаем карту для быстрого поиска индекса ключа, чтобы избежать вложенного цикла.
	keyIndexMap := make(map[string]int, len(keys))
	for i, key := range keys {
		keyIndexMap[key] = i
	}

	c.mu.RLock()
	for _, key := range keys {
		if value, ok := c.cache[key]; ok {
			fmt.Printf("[CACHE HIT] MGet key: %s\n", key)
			results[keyIndexMap[key]] = value
		} else {
			fmt.Printf("[CACHE MISS] MGet key: %s\n", key)
			missingKeys = append(missingKeys, key)
		}
	}
	c.mu.RUnlock()

	if len(missingKeys) > 0 {
		fmt.Printf("MGet fetching %d missing keys from DB: %v\n", len(missingKeys), missingKeys)
		missingValues, err := c.repo.MGet(missingKeys...)
		if err != nil {
			return nil, err
		}

		c.mu.Lock()
		for i, value := range missingValues {
			key := missingKeys[i]
			c.cache[key] = value
			results[keyIndexMap[key]] = value
		}
		c.mu.Unlock()
	}

	return results, nil
}

// Set реализует стратегию "Write-Through" (с некоторыми упрощениями).
// Сначала обновляем кэш, затем основное хранилище.
func (c *CachedRepository) Set(key, value string) error {
	fmt.Printf("Set key: %s. Updating cache and DB.\n", key)
	c.mu.Lock()
	c.cache[key] = value
	c.mu.Unlock()

	// Передаем вызов дальше, в основной репозиторий.
	return c.repo.Set(key, value)
}

// Del реализует стратегию "Write-Through" для удаления.
// Сначала удаляем из кэша, затем из основного хранилища.
func (c *CachedRepository) Del(key string) error {
	fmt.Printf("Del key: %s. Deleting from cache and DB.\n", key)
	c.mu.Lock()
	delete(c.cache, key)
	c.mu.Unlock()

	return c.repo.Del(key)
}

// --- Mock-реализация для демонстрации ---

// mockDBRepository имитирует реальный репозиторий (например, базу данных)
// с искусственной задержкой для наглядности работы кэша.
type mockDBRepository struct {
	data map[string]string
	mu   sync.Mutex
}

func newMockDB() *mockDBRepository {
	return &mockDBRepository{
		data: map[string]string{
			"user:1": "John",
			"user:2": "Jane",
		},
	}
}

func (db *mockDBRepository) Get(key string) (string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	time.Sleep(100 * time.Millisecond) // Имитация задержки БД
	if val, ok := db.data[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("key not found")
}

func (db *mockDBRepository) MGet(keys ...string) ([]string, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	time.Sleep(200 * time.Millisecond) // Пакетная операция тоже занимает время
	results := make([]string, len(keys))
	for i, key := range keys {
		results[i] = db.data[key]
	}
	return results, nil
}

func (db *mockDBRepository) Set(key, value string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	time.Sleep(50 * time.Millisecond)
	db.data[key] = value
	return nil
}

func (db *mockDBRepository) Del(key string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	time.Sleep(50 * time.Millisecond)
	delete(db.data, key)
	return nil
}

func main() {
	// 1. Создаем основной репозиторий (наша "база данных").
	dbRepo := newMockDB()
	// 2. Создаем кэширующий декоратор, оборачивая основной репозиторий.
	cachedRepo := NewCachedRepository(dbRepo)

	fmt.Println("--- Первый запрос Get ---")
	val, _ := cachedRepo.Get("user:1")
	fmt.Printf("Получено значение: %s\n\n", val)

	fmt.Println("--- Второй запрос Get (должен быть быстрее из-за кэша) ---")
	val, _ = cachedRepo.Get("user:1")
	fmt.Printf("Получено значение: %s\n\n", val)

	fmt.Println("--- Запрос MGet ---")
	vals, _ := cachedRepo.MGet("user:1", "user:2", "user:3")
	fmt.Printf("Получены значения: %s\n\n", strings.Join(vals, ", "))

	fmt.Println("--- Второй запрос MGet (user:1 и user:2 из кэша) ---")
	vals, _ = cachedRepo.MGet("user:1", "user:2", "user:3")
	fmt.Printf("Получены значения: %s\n\n", strings.Join(vals, ", "))

	fmt.Println("--- Запрос Set ---")
	_ = cachedRepo.Set("user:4", "Alice")
	val, _ = cachedRepo.Get("user:4")
	fmt.Printf("Проверка после Set: %s\n\n", val)

	fmt.Println("--- Запрос Del ---")
	_ = cachedRepo.Del("user:1")
	_, err := cachedRepo.Get("user:1") // Должен быть промах кэша и ошибка БД
	fmt.Printf("Проверка после Del: %v\n", err)
}
