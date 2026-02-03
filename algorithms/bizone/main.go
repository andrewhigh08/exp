// Package main демонстрирует реализацию конвейера обработки данных (ETL/pipeline pattern).
// Он состоит из чтения, параллельной обработки и записи данных с использованием интерфейсов.
package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
)

// Data - структура данных, которую мы обрабатываем в нашем конвейере.
type Data struct {
	ID      int
	Payload map[string]interface{}
}

// Reader определяет интерфейс для источника данных.
type Reader interface {
	Read() []*Data
}

// Processor определяет интерфейс для одного шага обработки данных.
// Каждый процессор принимает данные, обрабатывает их и возвращает измененные данные.
type Processor interface {
	Process(d Data) (*Data, error)
}

// Writer определяет интерфейс для приемника обработанных данных.
type Writer interface {
	Write(data []*Data)
}

// Manager управляет всем процессом конвейера.
type Manager interface {
	Manage()
}

// managerImpl - конкретная реализация интерфейса Manager.
type managerImpl struct {
	reader     Reader
	processors []Processor
	writer     Writer
}

// NewManager - конструктор для создания нового Manager.
func NewManager(reader Reader, processors []Processor, writer Writer) Manager {
	return &managerImpl{
		reader:     reader,
		processors: processors,
		writer:     writer,
	}
}

// Manage - основной метод, который управляет процессом чтения, обработки и записи.
func (m *managerImpl) Manage() {
	// Шаг 1: Чтение исходных данных.
	dataList := m.reader.Read()
	log.Printf("Прочитано %d записей.", len(dataList))

	// Канал для сбора обработанных данных от параллельных воркеров.
	// Буфер канала равен количеству данных, чтобы воркеры не блокировались.
	dataChan := make(chan *Data, len(dataList))

	// errgroup используется для управления группой горутин и их ошибками.
	// Он позволяет легко дождаться завершения всех горутин и получить первую возникшую ошибку.
	var g errgroup.Group

	// Шаг 2: Параллельная обработка каждой записи.
	for _, data := range dataList {
		// Создаем локальную копию переменной `data` для безопасного использования в замыкании (closure).
		// Это предотвращает гонку данных, когда несколько горутин могут получить указатель на одну и ту же переменную цикла.
		d := data
		g.Go(func() error {
			tempData := d
			// Последовательно применяем все процессоры к одной записи.
			for _, processor := range m.processors {
				var err error
				tempData, err = processor.Process(*tempData)
				if err != nil {
					// Если любой из процессоров возвращает ошибку, вся группа горутин будет отменена.
					return fmt.Errorf("ошибка обработки данных с ID %d: %w", d.ID, err)
				}
			}
			// Отправляем успешно обработанные данные в канал.
			dataChan <- tempData
			return nil
		})
	}

	// Ожидаем завершения всех горутин в группе.
	// Если хотя бы одна горутина вернула ошибку, wg.Wait() вернет эту ошибку.
	if err := g.Wait(); err != nil {
		log.Printf("Произошла ошибка во время обработки: %v. Процесс остановлен.", err)
		close(dataChan) // Закрываем канал, чтобы цикл сбора данных завершился.
		return          // Прекращаем выполнение.
	}
	// Закрываем канал после того, как все горутины успешно завершились.
	// Это сигнал для цикла ниже, что больше данных не будет.
	close(dataChan)

	// Шаг 3: Сбор всех обработанных данных из канала.
	var processedData []*Data
	for d := range dataChan {
		processedData = append(processedData, d)
	}

	log.Printf("Успешно обработано %d записей.", len(processedData))

	// Шаг 4: Запись обработанных данных.
	if len(processedData) > 0 {
		m.writer.Write(processedData)
	} else {
		log.Println("Нет данных для записи.")
	}
}

// --- Mock-реализации для демонстрации работы ---

// mockReader имитирует чтение данных из источника.
type mockReader struct{}

func (r *mockReader) Read() []*Data {
	return []*Data{
		{ID: 1, Payload: map[string]interface{}{"value": 10}},
		{ID: 2, Payload: map[string]interface{}{"value": 20}},
		{ID: 3, Payload: map[string]interface{}{"value": 30}},
	}
}

// addTimestampProcessor имитирует процессор, который добавляет временную метку.
type addTimestampProcessor struct{}

func (p *addTimestampProcessor) Process(d Data) (*Data, error) {
	d.Payload["timestamp"] = time.Now().Unix()
	log.Printf("Процессор 1: Добавлена временная метка для ID %d", d.ID)
	// Имитация небольшой задержки
	time.Sleep(100 * time.Millisecond)
	return &d, nil
}

// stringifyValueProcessor имитирует процессор, который преобразует числовое значение в строку.
type stringifyValueProcessor struct{}

func (p *stringifyValueProcessor) Process(d Data) (*Data, error) {
	if val, ok := d.Payload["value"].(int); ok {
		d.Payload["value_str"] = strconv.Itoa(val)
		log.Printf("Процессор 2: Преобразовано значение для ID %d", d.ID)
	}
	// Этот процессор никогда не возвращает ошибку
	return &d, nil
}

// mockWriter имитирует запись данных в приемник (например, базу данных или лог).
type mockWriter struct{}

func (w *mockWriter) Write(data []*Data) {
	log.Printf("--- Начало записи %d обработанных записей ---", len(data))
	for _, d := range data {
		fmt.Printf("Запись: ID=%d, Payload=%v\n", d.ID, d.Payload)
	}
	log.Println("--- Конец записи ---")
}

func main() {
	log.Println("Запуск конвейера обработки данных...")

	// Создаем экземпляры наших mock-компонентов.
	reader := &mockReader{}
	writer := &mockWriter{}
	processors := []Processor{
		&addTimestampProcessor{},
		&stringifyValueProcessor{},
	}

	// Создаем и запускаем менеджер.
	manager := NewManager(reader, processors, writer)
	manager.Manage()

	log.Println("Конвейер завершил работу.")
}
