package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Task представляет задачу с URL для скачивания/проверки
type Task struct {
	URL string
}

// Result содержит результат проверки URL
type Result struct {
	URL        string
	StatusCode int
	Error      error
	Duration   time.Duration
}

// worker — функция, которая читает из канала jobs и делает HTTP-запрос по URL
func worker(id int, jobs <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	// Настроим HTTP-клиент с таймаутом
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for j := range jobs {
		fmt.Printf("Воркер %d: начал обработку %s\n", id, j.URL)

		start := time.Now()
		resp, err := client.Get(j.URL)
		duration := time.Since(start)

		result := Result{
			URL:      j.URL,
			Duration: duration,
			Error:    err,
		}

		if err == nil {
			result.StatusCode = resp.StatusCode
			resp.Body.Close() // Обязательно закрываем тело ответа
		}

		fmt.Printf("Воркер %d: закончил обработку %s\n", id, j.URL)
		results <- result
	}
}

func main() {
	// Список URL-ов для проверки
	urls := []string{
		"https://golang.org",
		"https://google.com",
		"https://github.com",
		"https://stackoverflow.com",
		"https://pkg.go.dev",
		"https://invalid-url.com.example", // Пример нерабочего URL
	}

	numJobs := len(urls)
	const numWorkers = 3

	// Создаем буферизованные каналы для задач и результатов
	jobs := make(chan Task, numJobs)
	results := make(chan Result, numJobs)

	// WaitGroup для синхронизации завершения всех воркеров
	var wg sync.WaitGroup

	// Запускаем воркеров
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Отправляем задачи в канал.
	for _, u := range urls {
		jobs <- Task{URL: u}
	}
	// Закрываем канал задач: больше задач не поступит.
	close(jobs)

	// Ожидаем завершения всех воркеров в отдельной горутине.
	// Когда все закончат работу, закрываем канал результатов.
	go func() {
		wg.Wait()
		close(results)
	}()

	fmt.Println("\n--- Вывод результатов ---")
	// Читаем результаты по мере их поступления
	for res := range results {
		if res.Error != nil {
			fmt.Printf("❌ ОШИБКА  \t %s: %v (заняло %v)\n", res.URL, res.Error, res.Duration)
		} else {
			fmt.Printf("✅ %d \t %s (заняло %v)\n", res.StatusCode, res.URL, res.Duration)
		}
	}

	fmt.Println("Все URL обработаны.")
}
