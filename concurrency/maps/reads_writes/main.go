// Package main демонстрирует использование `sync.RWMutex` для синхронизации
// конкурентных операций чтения и записи в `map`.
//
// ПРОБЛЕМА:
// Если множество горутин читают данные, а одна горутина их изменяет,
// использование обычного `sync.Mutex` становится неэффективным. `Mutex` заблокирует
// всех читателей, даже если они могли бы безопасно читать данные одновременно.
//
// РЕШЕНИЕ:
// `sync.RWMutex` (мьютекс чтения-записи) решает эту проблему. Он позволяет:
// - Любому количеству читателей одновременно получать доступ к данным (через `RLock`).
// - Только одному писателю получать эксклюзивный доступ (через `Lock`).
//
// Пока писатель удерживает блокировку, все читатели ждут.
// Пока хотя бы один читатель удерживает блокировку, писатель ждет.
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	const numWriters = 2
	const numReaders = 10
	const initialDataSize = 100

	storage := make(map[int]int, initialDataSize)
	var wg sync.WaitGroup
	var mu sync.RWMutex // Используем RWMutex

	// --- Фаза 1: Первичное заполнение карты ---
	wg.Add(initialDataSize)
	fmt.Printf("Запуск %d горутин для первоначального заполнения карты...\n", initialDataSize)
	for i := 0; i < initialDataSize; i++ {
		go func(key int) {
			defer wg.Done()
			mu.Lock() // Эксклюзивная блокировка для записи
			storage[key] = key * key
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	fmt.Printf("Карта заполнена. Размер: %d\n\n", len(storage))

	// --- Фаза 2: Симуляция конкурентного доступа ---
	// Запускаем много читателей и несколько писателей одновременно.
	fmt.Printf("Запуск %d читателей и %d писателей...\n", numReaders, numWriters)
	wg.Add(numReaders + numWriters)

	// Горутины-читатели
	for i := 0; i < numReaders; i++ {
		go func(workerID int) {
			defer wg.Done()
			// Имитация многократных чтений
			for j := 0; j < 5; j++ {
				mu.RLock() // Блокировка на чтение (разделяемая)
				// Несколько горутин могут одновременно находиться здесь.
				randomKey := rand.Intn(initialDataSize)
				value := storage[randomKey]
				fmt.Printf("Читатель #%d: прочитал значение %d по ключу %d\n", workerID, value, randomKey)
				mu.RUnlock() // Освобождаем блокировку на чтение.
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	// Горутины-писатели
	for i := 0; i < numWriters; i++ {
		go func(workerID int) {
			defer wg.Done()
			// Имитация редких записей
			time.Sleep(20 * time.Millisecond)

			mu.Lock() // Эксклюзивная блокировка на запись
			// В этот момент все читатели (и другие писатели) будут ждать.
			randomKey := rand.Intn(initialDataSize)
			newValue := -workerID
			fmt.Printf(">> Писатель #%d: устанавливает значение %d по ключу %d <<\n", workerID, newValue, randomKey)
			storage[randomKey] = newValue
			mu.Unlock() // Освобождаем эксклюзивную блокировку.

		}(i)
	}

	wg.Wait()
	fmt.Println("\nВсе операции завершены.")
	fmt.Printf("Итоговый размер карты: %d\n", len(storage))
}
