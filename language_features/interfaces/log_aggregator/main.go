// Package main демонстрирует реализацию конкурентного конвейера (pipeline)
// для агрегации и обработки логов из различных источников.
package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"
)

// LogMessage — структура нашего лога.
type LogMessage struct {
	Timestamp time.Time
	Level     string
	Message   string
}

// --- Интерфейсы компонентов конвейера ---

// LogReader читает лог-сообщения из источника.
// Метод ReadLog должен возвращать `io.EOF`, когда сообщения заканчиваются.
type LogReader interface {
	ReadLog() (*LogMessage, error)
}

// LogTransformer преобразует лог-сообщение.
type LogTransformer interface {
	Transform(*LogMessage) (*LogMessage, error)
}

// LogStorage сохраняет обработанное лог-сообщение.
type LogStorage interface {
	StoreLog(*LogMessage) error
}

// LogManager запускает и управляет процессом агрегации.
type LogManager interface {
	Aggregate()
}

// --- Реализация менеджера ---

// LogAggregator — реализация LogManager.
type LogAggregator struct {
	reader       LogReader
	transformers []LogTransformer // Теперь это срез для поддержки цепочки трансформаций
	storage      LogStorage
	numWorkers   int // Количество воркеров для параллельной обработки
}

// NewLogAggregator — конструктор для LogAggregator.
func NewLogAggregator(reader LogReader, transformers []LogTransformer, storage LogStorage, numWorkers int) *LogAggregator {
	return &LogAggregator{
		reader:       reader,
		transformers: transformers,
		storage:      storage,
		numWorkers:   numWorkers,
	}
}

// Aggregate запускает конвейер: читает логи и распределяет их по воркерам для обработки.
func (la *LogAggregator) Aggregate() {
	var wg sync.WaitGroup
	jobs := make(chan *LogMessage, la.numWorkers)

	// 1. Запускаем пул воркеров
	wg.Add(la.numWorkers)
	for i := 0; i < la.numWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			// Воркер читает сообщения из канала `jobs` до тех пор, пока он не будет закрыт.
			for logMsg := range jobs {
				processLog(workerID, logMsg, la.transformers, la.storage)
			}
		}(i)
	}

	// 2. Читаем логи из источника и отправляем их в канал `jobs`
	for {
		logMsg, err := la.reader.ReadLog()
		if err != nil {
			// Если источник иссяк, прекращаем чтение.
			if errors.Is(err, io.EOF) {
				fmt.Println("Источник логов иссяк. Завершение чтения.")
				break
			}
			// Логируем ошибку чтения и продолжаем.
			log.Printf("Ошибка чтения лога: %v\n", err)
			continue
		}
		jobs <- logMsg
	}

	// 3. Закрываем канал `jobs`, чтобы воркеры завершили свою работу после обработки всех сообщений.
	close(jobs)

	// 4. Ожидаем, пока все воркеры полностью завершат работу.
	wg.Wait()
	fmt.Println("Вся обработка завершена.")
}

// processLog выполняет полную цепочку обработки для одного лог-сообщения.
func processLog(workerID int, msg *LogMessage, transformers []LogTransformer, storage LogStorage) {
	fmt.Printf("[Воркер %d] Начал обработку сообщения: %s\n", workerID, msg.Message)
	currentMsg := msg

	// Применяем все трансформации последовательно.
	for _, t := range transformers {
		var err error
		currentMsg, err = t.Transform(currentMsg)
		if err != nil {
			log.Printf("[Воркер %d] Ошибка трансформации лога '%s': %v. Лог пропущен.", workerID, msg.Message, err)
			return // Прерываем обработку этого сообщения.
		}
	}

	// Сохраняем итоговый результат.
	if err := storage.StoreLog(currentMsg); err != nil {
		log.Printf("[Воркер %d] Ошибка сохранения лога '%s': %v.", workerID, msg.Message, err)
	} else {
		fmt.Printf("[Воркер %d] Успешно сохранил лог: %s\n", workerID, currentMsg.Message)
	}
}

// --- Mock-реализации для демонстрации ---

type mockReader struct {
	messages []*LogMessage
	mu       sync.Mutex
}

func (r *mockReader) ReadLog() (*LogMessage, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.messages) == 0 {
		return nil, io.EOF
	}
	msg := r.messages[0]
	r.messages = r.messages[1:]
	return msg, nil
}

type addPrefixTransformer struct{ prefix string }

func (t *addPrefixTransformer) Transform(msg *LogMessage) (*LogMessage, error) {
	msg.Message = t.prefix + msg.Message
	return msg, nil
}

type toUpperTransformer struct{}

func (t *toUpperTransformer) Transform(msg *LogMessage) (*LogMessage, error) {
	if msg.Message == "[PROCESSED] special_error" {
		return nil, errors.New("специальная ошибка трансформации")
	}
	msg.Message = strings.ToUpper(msg.Message)
	return msg, nil
}

type mockStorage struct{}

func (s *mockStorage) StoreLog(msg *LogMessage) error {
	fmt.Printf("---ХРАНИЛИЩЕ: Сохранен лог (Уровень: %s): %s\n", msg.Level, msg.Message)
	return nil
}

func main() {
	// 1. Создаем компоненты конвейера.
	reader := &mockReader{
		messages: []*LogMessage{
			{Timestamp: time.Now(), Level: "INFO", Message: "user logged in"},
			{Timestamp: time.Now(), Level: "WARN", Message: "disk space is low"},
			{Timestamp: time.Now(), Level: "INFO", Message: "special_error"},
			{Timestamp: time.Now(), Level: "DEBUG", Message: "request received"},
		},
	}
	transformers := []LogTransformer{
		&addPrefixTransformer{prefix: "[PROCESSED] "},
		&toUpperTransformer{},
	}
	storage := &mockStorage{}

	// 2. Создаем и запускаем менеджер агрегации с 2 воркерами.
	manager := NewLogAggregator(reader, transformers, storage, 2)
	manager.Aggregate()
}
