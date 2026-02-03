# Go Experiments

Репозиторий с учебными примерами и экспериментами на Go. Содержит реализации алгоритмов, паттернов конкурентности, паттернов проектирования, особенностей языка и практических примеров.

**Go:** 1.25.5
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

## Паттерны проектирования (`design_patterns/`)

| Паттерн | Директория | Описание |
|---|---|---|
| Adapter | `adapter/` | Адаптация несовместимых интерфейсов (логгер) |
| Decorator | `decorator/` | Кеширование Redis поверх БД |
| Cached Repository | `cached_repository/` | In-memory кеш для репозитория |
| Worker Pool | `worker_pool/` | Распределение задач по воркерам |
| Pipeline | `read_process_write/` | Многостадийная обработка данных |

## Особенности языка (`language_features/`)

| Тема | Директория | Что демонстрирует |
|---|---|---|
| Generics | `generics/` | Параметры типов, constraints (Go 1.18+) |
| Pointers | `pointers/` | Семантика указателей, pass-by-value |
| Slices | `slices/` | Внутреннее устройство слайсов, append |
| Defer | `defer/` | Порядок вызовов defer (LIFO) |
| Fallthrough | `fallthrough/` | Поведение switch/fallthrough |
| Interfaces | `interfaces/` | Реализация интерфейсов |
| — simple | `interfaces/simple/` | Базовое удовлетворение интерфейса |
| — abc | `interfaces/abc/` | Встраивание, type assertions |
| — difficult | `interfaces/difficult/` | nil-интерфейсы vs nil-значения |
| — log_aggregator | `interfaces/log_aggregator/` | Конкурентный пайплайн с интерфейсами |

## Практические примеры (`examples/`)

| Пример | Описание |
|---|---|
| `json_config` | HTTP-сервер с динамической перезагрузкой конфигурации |
| `url_shorter` | Сокращатель URL (`fmt.Stringer`) |
| `cli_spinner` | Анимация спиннера в терминале |
| `string_validator` | Валидация строк через регулярные выражения |

## Бенчмарки (`benchmarks/`)

Сравнение производительности:

- **Аллокация слайсов** — append с малой ёмкостью vs оптимальная ёмкость vs преаллокация
- **Конкатенация строк** — `fmt.Sprintf` vs `strconv.Itoa` vs `strings.Builder`

```bash
cd benchmarks
go test -bench . -benchmem
```

## Кодогенерация (`code_generation/`)

Инструмент `repogen` для автоматической генерации CRUD-репозиториев на основе структур с GORM.

Принцип работы:
1. Парсит структуру с комментарием `//repogen:entity`
2. Генерирует реализацию репозитория (Get, Create, Update, Delete)
3. Создаёт файлы `*_gen.go`

```bash
cd code_generation
go generate ./...
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
| `golang.org/x/sync` | v0.18.0 | `errgroup` для управления горутинами |
| `golang.org/x/tools` | v0.21.0 | AST-парсинг (кодогенерация) |
