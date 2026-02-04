// Package main содержит пример решения задачи с собеседования.
// Задача: реализовать функцию, которая параллельно опрашивает несколько реплик базы данных,
// обрабатывает ошибки с логикой повторных попыток (retry) и возвращает первый успешный результат.
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrNotFound — это специальная ошибка, которая означает, что данные не найдены.
// При получении этой ошибки мы не должны делать повторные запросы (retry),
// так как это окончательный ответ от реплики.
var ErrNotFound = errors.New("not found")

// DatabaseHost определяет интерфейс для взаимодействия с хостом базы данных.
// Это позволяет нам использовать как реальные, так и тестовые (mock) реализации.
type DatabaseHost interface {
	DoQuery(ctx context.Context, query string) (string, error)
}

// Response — структура для передачи результата выполнения запроса через канал.
// Она содержит либо сообщение, либо ошибку.
type Response struct {
	Message string
	Err     error
	Host    string // Добавим хост для наглядности в логах
}

const (
	maxAttempts   = 3               // Максимальное количество попыток для одного запроса.
	retryInterval = 500 * time.Millisecond // Интервал между повторными попытками.
	totalTimeout  = 2 * time.Second // Общий таймаут для всей операции DistributedQuery.
)

// DistributedQuery выполняет запрос параллельно к нескольким репликам.
// Она возвращает первый полученный успешный ответ.
// Если все реплики вернули ошибку или истек общий таймаут, функция вернет ошибку.
func DistributedQuery(query string, replicas []DatabaseHost) (string, error) {
	// Создаем контекст с общим таймаутом. Это гарантирует, что функция не будет выполняться вечно.
	ctx, cancel := context.WithTimeout(context.Background(), totalTimeout)
	defer cancel() // Важно вызвать cancel, чтобы освободить ресурсы контекста.

	// Буферизированный канал для результатов. Размер буфера равен количеству реплик,
	// чтобы ни одна горутина не заблокировалась при отправке результата.
	resCh := make(chan Response, len(replicas))
	var wg sync.WaitGroup

	wg.Add(len(replicas))

	// Запускаем по одной горутине на каждую реплику.
	for _, rep := range replicas {
		go func(rep DatabaseHost) {
			defer wg.Done()

			for i := 0; i < maxAttempts; i++ {
				// Перед каждой попыткой проверяем, не был ли отменен контекст (например, по таймауту).
				if ctx.Err() != nil {
					return // Выходим, если операция уже отменена.
				}

				resp, err := rep.DoQuery(ctx, query)

				// Успешный результат или ошибка ErrNotFound - отправляем в канал и выходим.
				if err == nil || errors.Is(err, ErrNotFound) {
					resCh <- Response{Message: resp, Err: err}
					return
				}

				// Для всех остальных ошибок делаем повторную попытку (retry).
				// Используем select, чтобы не блокировать горутину надолго и вовремя среагировать
				// на отмену контекста.
				select {
				case <-time.After(retryInterval):
					// Интервал ожидания прошел, продолжаем цикл для следующей попытки.
					continue
				case <-ctx.Done():
					// Контекст был отменен во время ожидания, выходим.
					return
				}
			}
		}(rep)
	}

	// Запускаем отдельную горутину, которая закроет канал resCh после того,
	// как все воркеры завершат свою работу. Это сигнал о том, что больше результатов не будет.
	go func() {
		wg.Wait()
		close(resCh)
	}()

	// Основной цикл ожидания результатов.
	for {
		select {
		case resp, ok := <-resCh:
			if !ok {
				// Канал закрыт, и мы не получили ни одного успешного ответа.
				// Это означает, что все реплики вернули ошибку (кроме ErrNotFound).
				return "", errors.New("all replicas failed after multiple retries")
			}

			// Получили первый ответ. Если это не ошибка, возвращаем результат.
			if resp.Err == nil {
				fmt.Printf("Success from %s: %s\n", resp.Host, resp.Message)
				cancel() // Отменяем контекст, чтобы остальные горутины прекратили работу.
				return resp.Message, nil
			}

			// Если пришла ошибка ErrNotFound, мы не можем считать ее успехом,
			// но и повторять запрос к этой реплике бессмысленно. Мы просто игнорируем ее
			// и ждем ответов от других реплик.
			if errors.Is(resp.Err, ErrNotFound) {
				fmt.Printf("Result from %s: %s\n", resp.Host, resp.Err)
				// Продолжаем ждать более подходящего ответа.
				continue
			}

		case <-ctx.Done():
			// Сработал общий таймаут.
			return "", fmt.Errorf("query timed out after %s", totalTimeout)
		}
	}
}

// --- Mock-реализация для демонстрации ---

// mockHost имитирует хост базы данных.
type mockHost struct {
	name         string
	flaky        bool // Если true, хост будет возвращать ошибки.
	notFound     bool // Если true, хост вернет ошибку ErrNotFound.
	slow         bool // Если true, хост будет отвечать медленно.
	flakyCounter int
}

// DoQuery реализует интерфейс DatabaseHost для mockHost.
func (h *mockHost) DoQuery(ctx context.Context, query string) (string, error) {
	// Имитация долгого запроса
	if h.slow {
		select {
		case <-time.After(1 * time.Second):
			// Продолжаем выполнение
		case <-ctx.Done():
			return "", ctx.Err() // Возвращаем ошибку, если контекст отменен
		}
	}

	if h.notFound {
		return "", ErrNotFound
	}

	if h.flaky {
		h.flakyCounter++
		// Допустим, хост отвечает успешно только с третьей попытки.
		if h.flakyCounter < 3 {
			return "", errors.New("temporary connection error")
		}
	}

	return fmt.Sprintf("result from %s", h.name), nil
}

func main() {
	fmt.Println("--- Сценарий 1: Одна из реплик отвечает успешно ---")
	replicas1 := []DatabaseHost{
		&mockHost{name: "Replica 1 (flaky)", flaky: true},
		&mockHost{name: "Replica 2 (ok)"},
		&mockHost{name: "Replica 3 (slow)", slow: true},
	}
	result, err := DistributedQuery("SELECT * FROM users", replicas1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Final Result: %s\n", result)
	}
	// Ожидаемый результат: "result from Replica 2 (ok)"


	fmt.Println("\n--- Сценарий 2: Все реплики возвращают ошибку ---")
	replicas2 := []DatabaseHost{
		&mockHost{name: "Replica 1 (flaky)", flaky: true},
		&mockHost{name: "Replica 2 (flaky)", flaky: true},
	}
	result, err = DistributedQuery("SELECT * FROM users", replicas2)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	// Ожидаемый результат: "all replicas failed after multiple retries"


	fmt.Println("\n--- Сценарий 3: Таймаут ---")
	replicas3 := []DatabaseHost{
		&mockHost{name: "Replica 1 (very slow)", slow: true},
		&mockHost{name: "Replica 2 (very slow)", slow: true},
	}
	// Установим таймаут меньше, чем время ответа реплик
	// (Для этого сценария, можно было бы передать кастомный таймаут в функцию,
	// но для простоты примера оставим глобальный)
	result, err = DistributedQuery("SELECT * FROM users", replicas3)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	// Ожидаемый результат: "query timed out after 2s"


	fmt.Println("\n--- Сценарий 4: Одна реплика не находит данные, другая успешна ---")
	replicas4 := []DatabaseHost{
		&mockHost{name: "Replica 1 (not found)", notFound: true},
		&mockHost{name: "Replica 2 (ok)"},
	}
	result, err = DistributedQuery("SELECT * FROM users WHERE id=123", replicas4)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Final Result: %s\n", result)
	}
	// Ожидаемый результат: "result from Replica 2 (ok)"
}
