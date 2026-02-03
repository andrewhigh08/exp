// Package main демонстрирует реализацию конвейерного паттерна обработки данных (Pipeline).
// Этот паттерн позволяет выстраивать цепочку из нескольких шагов обработки,
// где выход одного шага является входом для следующего.
//
// Особенность данной реализации — каждый шаг обработки (Processor) может
// изменять количество элементов данных (один элемент может превратиться в несколько или быть отфильтрован).
package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

// Data — структура данных, которую мы обрабатываем.
type Data struct {
	ID      int
	Payload string
}

// Reader — интерфейс для источника данных.
type Reader interface {
	Read() []*Data
}

// Processor — интерфейс для одного шага обработки.
// Может преобразовать один элемент `Data` в ноль, один или несколько новых элементов.
type Processor interface {
	Process(d *Data) ([]*Data, error)
}

// Writer — интерфейс для приемника обработанных данных.
type Writer interface {
	Write(data []*Data)
}

// Manager — интерфейс, управляющий всем процессом.
type Manager interface {
	Manage()
}

// DataManager — реализация Manager.
type DataManager struct {
	reader     Reader
	processors []Processor
	writer     Writer
}

// NewDataManager — конструктор для DataManager.
func NewDataManager(reader Reader, processors []Processor, writer Writer) *DataManager {
	return &DataManager{
		reader:     reader,
		processors: processors,
		writer:     writer,
	}
}

// Manage управляет потоком данных: читает, конкурентно обрабатывает и записывает.
func (dm *DataManager) Manage() {
	initialData := dm.reader.Read()
	log.Printf("Прочитано %d элементов из источника.", len(initialData))

	var finalResults []*Data
	var finalMu sync.Mutex // Мьютекс для безопасного добавления в общий срез результатов
	var eg errgroup.Group

	// Обрабатываем каждый элемент из начального набора в отдельной горутине.
	for _, item := range initialData {
		item := item // Создаем локальную копию для безопасного использования в замыкании.
		eg.Go(func() error {
			// `currentData` представляет собой набор данных на входе для цепочки процессоров.
			// Начинаем с одного элемента.
			currentData := []*Data{item}

			// Последовательно пропускаем данные через все процессоры.
			for _, processor := range dm.processors {
				// `nextData` будет содержать результаты работы текущего процессора.
				var nextData []*Data
				for _, dataItem := range currentData {
					processed, err := processor.Process(dataItem)
					if err != nil {
						// Если процессор вернул ошибку, пропускаем этот элемент
						// и не передаем его дальше по цепочке.
						log.Printf("Ошибка обработки элемента ID %d: %v. Элемент пропущен.", dataItem.ID, err)
						continue // Пропускаем только `dataItem`, а не весь `item`
					}
					nextData = append(nextData, processed...)
				}
				// Результат этого шага становится входом для следующего.
				currentData = nextData

				// Если на каком-то шаге все данные были отфильтрованы,
				// нет смысла продолжать обработку.
				if len(currentData) == 0 {
					break
				}
			}

			// Если после всех процессоров остались данные, добавляем их в общий результат.
			if len(currentData) > 0 {
				finalMu.Lock()
				finalResults = append(finalResults, currentData...)
				finalMu.Unlock()
			}
			return nil
		})
	}

	// Ожидаем завершения всех горутин. errgroup вернет первую возникшую ошибку.
	if err := eg.Wait(); err != nil {
		log.Printf("Произошла критическая ошибка в одной из горутин: %v", err)
		return
	}

	// Записываем все собранные результаты одним пакетом.
	if len(finalResults) > 0 {
		dm.writer.Write(finalResults)
	} else {
		log.Println("Нет данных для записи после обработки.")
	}
}

// --- Mock-реализации для демонстрации ---

type mockReader struct{}

func (r *mockReader) Read() []*Data {
	return []*Data{
		{ID: 1, Payload: "hello"},
		{ID: 2, Payload: "world"},
		{ID: 3, Payload: "error"}, // Этот элемент вызовет ошибку
	}
}

type duplicatorProcessor struct{}

// Process дублирует каждый элемент.
func (p *duplicatorProcessor) Process(d *Data) ([]*Data, error) {
	log.Printf("Дубликатор: обрабатывается ID %d", d.ID)
	// Имитация ошибки для определенного элемента
	if d.Payload == "error" {
		return nil, errors.New("некорректный payload")
	}
	// Возвращаем два новых элемента
	return []*Data{
		{ID: d.ID, Payload: d.Payload + " (копия 1)"},
		{ID: d.ID, Payload: d.Payload + " (копия 2)"},
	}, nil
}

type upperCaseProcessor struct{}

// Process преобразует Payload в верхний регистр.
func (p *upperCaseProcessor) Process(d *Data) ([]*Data, error) {
	log.Printf("Верхний регистр: обрабатывается ID %d", d.ID)
	d.Payload = strings.ToUpper(d.Payload)
	// Возвращаем один измененный элемент
	return []*Data{d}, nil
}

type mockWriter struct {
	mu   sync.Mutex
	data []*Data
}

func (w *mockWriter) Write(data []*Data) {
	w.mu.Lock()
	defer w.mu.Unlock()
	log.Printf("Запись %d элементов...", len(data))
	w.data = append(w.data, data...)
}

func main() {
	reader := &mockReader{}
	writer := &mockWriter{}
	processors := []Processor{
		&duplicatorProcessor{},
		&upperCaseProcessor{},
	}

	manager := NewDataManager(reader, processors, writer)
	manager.Manage()

	fmt.Println("\n--- Итоговые данные в Writer ---")
	for _, d := range writer.data {
		fmt.Printf("ID: %d, Payload: %s\n", d.ID, d.Payload)
	}
}
