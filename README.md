# Go Experiments

Репозиторий с учебными примерами и экспериментами на Go. Содержит реализации алгоритмов, паттернов конкурентности, паттернов проектирования, особенностей языка и практических примеров.

**Go:** 1.27.2
**Модуль:** `github.com/andrewhigh08/exp`

## Структура проекта

```
exp/
├── algorithms/          # Алгоритмы и структуры данных
├── benchmarks/          # Бенчмарки производительности
├── code_generation/     # Кодогенератор репозиториев (repogen)
├── concurrency/         # Паттерны конкурентности
├── design_patterns/     # Паттерны проектирования
├── examples/            # Практические примеры
├── interview/           # Вопросы с собеседований
└── language_features/   # Особенности языка Go
```

## Алгоритмы (`algorithms/`)

| Директория | Описание | Ключевые концепции |
|---|---|---|
| `zeros_to_the_right` | Перемещение нулей в конец слайса | Два указателя, in-place |
| `bizone` | ETL-пайплайн обработки данных | Fan-out, параллельная обработка |
| `polindrome` | Проверка палиндрома | Unicode, работа со строками |
| `666` | Семантика слайсов | Указатели, поведение append |
| `opechatka` | Конвертер раскладки клавиатуры | Транслитерация, strings.Builder |
| `parkovka` | Поиск парковочного места за O(1) | Hash map |
| `simplify_path` | Упрощение Unix-пути | Стек |
| `rle` | Run-Length Encoding | Сжатие строк |
| `count_ships` | Подсчёт кораблей на сетке | Обход 2D-массива |
| `piramid` | Сумма пирамиды нечётных чисел | Математика |
| `bankomat` | Выдача денег банкоматом | Жадный алгоритм |
| `distributed_query` | Распределённые запросы | Агрегация данных, обработка ошибок |

## Конкурентность (`concurrency/`)

| Директория | Описание | Примитивы |
|---|---|---|
| `sync_basics` | Основы синхронизации | `sync.WaitGroup` |
| `semaphore` | Ограничение конкурентности | Буферизированный канал как семафор |
| `errgroup/` | Группы горутин с ошибками | `errgroup.Group` |
| `errgroup_with_channels` | Errgroup + каналы | `errgroup`, `SetLimit`, каналы |
| `producer_consumer` | Производитель-потребитель | Context, каналы, горутины |
| `pub_sub` | Publish-Subscribe | Fan-out, `sync.RWMutex` |
| `sync_channels` | Генератор на каналах | CSP, каналы |
| `result_channel_pattern` | Паттерн Result через канал | Структуры с ошибками |
| `channels_wg_context` | Graceful shutdown | Каналы + WaitGroup + Context |
| `once_with_map` | Уникальные элементы | `sync.Mutex`, дедупликация |
| `maps/reads_writes` | Конкурентное чтение/запись | `sync.RWMutex` |
| `maps/writes` | Конкурентная запись | `sync.Mutex` |
| `max_procs` | GOMAXPROCS | Планировщик, недетерминизм |
| `spb_ekt_msk` | Fan-in паттерн | Несколько горутин → один канал |
| `processor` | Асинхронный обработчик задач | Safe channel close, `sync.RWMutex`, `sync.Once` |

## Паттерны проектирования (`design_patterns/`)

| Паттерн | Директория | Описание |
|---|---|---|
| Adapter | `adapter/` | Адаптация несовместимых интерфейсов (логгер) |
| Decorator | `decorator/` | Кеширование Redis поверх БД |
| Cached Repository | `cached_repo/` | In-memory кеш для репозитория |
| Worker Pool | `worker_pool/` | Распределение задач по воркерам |
| Pipeline | `read_process_write/` | Многостадийная обработка данных |

## Особенности языка (`language_features/`)

| Тема | Директория | Что демонстрирует |
|---|---|---|
| Generics | `generics/` | Параметры типов, constraints (`~T`, `comparable`, `any`), generic-методы (Go 1.27+) |
| Pointers | `pointers/` | Семантика указателей, pass-by-value |
| Slices | `slices/` | Внутреннее устройство слайсов, append |
| Defer | `defer/` | Момент вычисления аргументов отложенного вызова (ресивер-значение vs указатель vs замыкание) |
| Fallthrough | `fallthrough/` | Поведение switch/fallthrough |
| Interfaces | `interfaces/` | Реализация интерфейсов |
| — simple | `interfaces/simple/` | Базовое удовлетворение интерфейса |
| — abc | `interfaces/abc/` | Встраивание, type assertions |
| — difficult | `interfaces/difficult/` | nil-интерфейсы vs nil-значения |
| — typed_nil | `interfaces/typed_nil/` | `v == nil` не ловит typed nil (`*T` в interface) |
| — log_aggregator | `interfaces/log_aggregator/` | Конкурентный пайплайн с интерфейсами |

## Практические примеры (`examples/`)

| Пример | Описание |
|---|---|
| `json_config` | HTTP-сервер с динамической перезагрузкой конфигурации |
| `numeronym` | Нумеронимы (`kubernetes` → `k8s`) через `fmt.Stringer` |
| `cli_spinner` | Анимация спиннера в терминале |
| `string_validator` | Валидация строк через регулярные выражения |

## Бенчмарки (`benchmarks/`)

Сравнение производительности:

- **Аллокация слайсов** — append с малой ёмкостью vs оптимальная ёмкость vs преаллокация
- **Сборка строк** — `+=` + `fmt.Sprintf` vs `+=` + `strconv.Itoa` vs `strings.Builder`

```bash
cd benchmarks
go test -bench . -benchmem
```

## Кодогенерация (`code_generation/`)

Инструмент `repogen` (отдельный модуль) генерирует CRUD-репозитории поверх GORM v2 (`gorm.io/gorm`) для структур с комментарием `//repogen:entity`.

Принцип работы:
1. `go generate` выполняет директиву `//go:generate go tool repogen` из `gen.go`. Генератор подключён к модулю tool-директивой (`tool repogen/cmd/repogen` в `go.mod`, Go 1.24+), поэтому `go install` и `PATH` не нужны
2. Парсит файл и находит структуры с `//repogen:entity`. Первичный ключ — поле с тегом `gorm:"primaryKey"` (или устаревшим `primary_key`), без тега — поле `ID`; колонка — из `column:` или по именованию GORM (`UserID` → `user_id`)
3. Генерирует репозиторий (Get, Create, Update, Delete) в `<файл>_gen.go` с пометкой `// Code generated by repogen. DO NOT EDIT.`; сгенерированный `gen_gen.go` лежит в репозитории

```bash
cd code_generation
go generate ./...              # пересоздаёт gen_gen.go
go build ./... && go vet ./...
```

## Запуск примеров

Каждый пример — самостоятельный `main.go`, который можно запустить напрямую:

```bash
# Алгоритм
go run algorithms/zeros_to_the_right/main.go

# Конкурентность
go run concurrency/producer_consumer/main.go

# Паттерн проектирования
go run design_patterns/decorator/main.go

# Особенность языка
go run language_features/generics/main.go
```

## Зависимости

| Модуль | Версия | Назначение |
|---|---|---|
| `golang.org/x/sync` | v0.24.0 | `errgroup` для управления горутинами |
| `golang.org/x/tools` | v0.52.0 | Обход AST и форматирование кода в кодогенераторе (модуль `code_generation/`) |
| `gorm.io/gorm` | v1.31.2 | GORM v2 для сгенерированных репозиториев; генератор берёт из него разбор тегов и именование колонок (модуль `code_generation/`) |
