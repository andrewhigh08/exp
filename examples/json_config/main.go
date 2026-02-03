// Package main демонстрирует создание HTTP-сервера, который:
// 1. Динамически (на лету) перезагружает конфигурацию из JSON-файла.
// 2. По запросу на эндпоинт `/ping` конкурентно опрашивает список серверов из конфига.
//
// В коде исправлены критические состояния гонки и применены идиоматичные подходы.
package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Config определяет структуру нашего JSON-конфига.
// Использование структуры вместо `map[string]interface{}` является более безопасным
// и идиоматичным подходом, так как обеспечивает строгую типизацию.
type Config struct {
	Servers []string `json:"servers"`
}

// App — основная структура нашего приложения.
// Она инкапсулирует зависимости: текущую конфигурацию и мьютекс для ее защиты.
type App struct {
	config Config
	mu     sync.RWMutex // RWMutex идеален для конфига: много читателей, редкие писатели.
}

// loadConfig периодически читает и обновляет конфигурацию приложения.
// Эта функция должна запускаться в отдельной горутине.
func (a *App) loadConfig(path string) {
	for {
		// Читаем файл
		data, err := os.ReadFile(path)
		if err != nil {
			log.Printf("Ошибка чтения файла конфигурации '%s': %v", path, err)
			time.Sleep(5 * time.Second) // В случае ошибки повторяем не так часто
			continue
		}

		var newConfig Config
		if err := json.Unmarshal(data, &newConfig); err != nil {
			log.Printf("Ошибка парсинга JSON из файла '%s': %v", path, err)
			time.Sleep(5 * time.Second)
			continue
		}

		// Блокируем мьютекс на запись, чтобы безопасно обновить конфигурацию.
		a.mu.Lock()
		a.config = newConfig
		a.mu.Unlock()

		log.Println("Конфигурация успешно обновлена.")
		time.Sleep(5 * time.Second) // Перезагружаем каждые 5 секунд
	}
}

// pingHandler — это обработчик для эндпоинта /ping.
func (a *App) pingHandler(w http.ResponseWriter, r *http.Request) {
	// Блокируем мьютекс на чтение, чтобы безопасно получить копию списка серверов.
	a.mu.RLock()
	servers := make([]string, len(a.config.Servers))
	copy(servers, a.config.Servers)
	a.mu.RUnlock()

	// responseMap будет содержать результаты опроса.
	responseMap := make(map[string]string)
	// Для защиты responseMap от конкурентной записи из горутин нужен отдельный мьютекс.
	var responseMu sync.Mutex
	var wg sync.WaitGroup

	log.Printf("Начинаю опрос %d серверов...", len(servers))

	for _, serverURL := range servers {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			// Выполняем GET-запрос.
			resp, err := http.Get(url)
			var status string
			if err != nil {
				status = "ERROR: " + err.Error()
			} else {
				defer resp.Body.Close()
				status = resp.Status
			}

			// Защищаем запись в responseMap с помощью мьютекса.
			responseMu.Lock()
			responseMap[url] = status
			responseMu.Unlock()

		}(serverURL)
	}

	// Ожидаем завершения всех запросов.
	wg.Wait()
	log.Println("Опрос завершен.")

	// Отправляем результат клиенту в формате JSON.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseMap)
}

func main() {
	// Определяем флаг для пути к файлу конфигурации.
	configPath := flag.String("config", "config.json", "путь к файлу config.json")
	flag.Parse()

	// Создаем экземпляр нашего приложения.
	app := &App{
		config: Config{},
	}

	// Запускаем горутину для динамической перезагрузки конфига.
	go app.loadConfig(*configPath)

	// Регистрируем обработчик эндпоинта.
	http.HandleFunc("/ping", app.pingHandler)

	log.Println("Сервер запущен на порту :8080")
	log.Printf("Для проверки откройте в браузере http://localhost:8080/ping")
	// Запускаем HTTP-сервер. log.Fatal остановит программу, если сервер не сможет запуститься.
	log.Fatal(http.ListenAndServe(":8080", nil))
}
