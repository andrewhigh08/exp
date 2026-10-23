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
| `bizone` | ETL-пайплайн обработки данных | Fan-out, `errgroup`, параллельная обработка |
| `polindrome` | Проверка палиндрома | Unicode, работа со строками |
| `666` | Семантика слайсов | Указатели, поведение append |
| `opechatka` | Конвертер раскладки клавиатуры | Транслитерация, strings.Builder |
| `parkovka` | Поиск парковки по координатам за O(1) | Hash map, struct как ключ, подводные камни float-ключей |
| `simplify_path` | Упрощение Unix-пути | Стек |
| `rle` | Run-Length Encoding | Сжатие строк |
| `count_ships` | Подсчёт кораблей на сетке | Подсчёт верхних левых клеток за один проход |
| `piramid` | Сумма пирамиды нечётных чисел | Математика |
| `bankomat` | Выдача денег банкоматом | Жадный алгоритм |
| `distributed_query` | Параллельный запрос к репликам БД (первый успешный ответ) | `context.WithTimeout`, retry, `ErrNotFound` |

### Задачи с собеседований (`algorithms/A/`)

У каждой задачи рядом с `main.go` лежит свой `README.md` с условием и разбором решения.

| Директория | Задача | Ключевые концепции |
|---|---|---|
| `01_move_zeroes` | Сдвиг нулей в конец с сохранением порядка | Два указателя, in-place, O(1) памяти |
| `02_bignum_add` | Сложение больших чисел, заданных массивами цифр | Поразрядное сложение с переносом |
| `03_tree_leaf_paths` | Листья рубрикатора с полным путём («Вещи => Одежда => Мужская») | DFS, backtracking |
| `04_compare_versions` | Сравнение версий `vX.Y.Z` | Посегментный разбор без `strconv`, валидация |
| `05_generate_parentheses` | Все правильные скобочные последовательности из n пар | Backtracking, числа Каталана |
| `06_reconstruct_route` | Восстановление маршрута по набору билетов | Hash map, поиск стартового города |
| `07_topo_sort_deps` | Порядок импорта библиотек по зависимостям | Топологическая сортировка, DFS |
| `08_dissatisfaction` | Суммарная неудовлетворённость покупателей (`stdlib/` — решение на пакете `sort`) | Свои quicksort и lower bound, бинарный поиск |
| `09_balanced_team` | Самая «ровная» команда из K отсортированных списков | K указателей, минимальный покрывающий диапазон |
| `10_filter_sellers_by_city` | Фильтр продавцов по списку городов | Множество на `map[string]struct{}` |
| `11_common_items` | Товары, которые есть у каждого продавца | Пересечение множеств |

## Конкурентность (`concurrency/`)

| Директория | Описание | Примитивы |
|---|---|---|
| `sync_basics` | Основы синхронизации | `sync.WaitGroup` |
| `semaphore` | Ограничение конкурентности | Буферизированный канал как семафор |
| `errgroup/` | Группы горутин с ошибками | `errgroup.Group` |
| `errgroup_with_channels` | Errgroup + каналы | `errgroup`, `SetLimit`, каналы |
| `producer_consumer` | Производитель-потребитель | Context, каналы, горутины |
| `pub_sub` | Publish-Subscribe | Fan-out, `sync.RWMutex`, владение каналом (`<-chan`) |
| `sync_channels` | Генератор на каналах | CSP, каналы |
| `result_channel_pattern` | Паттерн Result через канал | Структуры с ошибками |
| `channels_wg_context` | Graceful shutdown | Каналы + WaitGroup + Context |
| `once_with_map` | Уникальные элементы | `sync.Mutex`, дедупликация |
| `maps/reads_writes` | Конкурентное чтение/запись | `sync.RWMutex` |
| `maps/writes` | Конкурентная запись | `sync.Mutex` |
| `max_procs` | GOMAXPROCS | Планировщик, недетерминизм, GOMAXPROCS в Go 1.25+ |
| `spb_ekt_msk` | Fan-in паттерн | Несколько горутин → один канал |
| `worker_pool` | Пул воркеров: параллельная проверка URL | Каналы jobs/results, `sync.WaitGroup`, `net/http` |
| `processor` | Асинхронный обработчик задач (пакет с тестами) | Safe channel close, `sync.RWMutex`, `sync.Once` |

## Паттерны проектирования (`design_patterns/`)

| Паттерн | Директория | Описание |
|---|---|---|
| Adapter | `adapter/` | Адаптация несовместимых интерфейсов (логгер) |
| Decorator | `decorator/` | Кеширование (имитация Redis) поверх БД |
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
| — simple | `interfaces/simple/` | Неявная реализация, встраивание интерфейсов, ресивер-значение vs указатель |
| — abc | `interfaces/abc/` | Type assertion между интерфейсами, статический vs динамический тип |
| — difficult | `interfaces/difficult/` | nil-интерфейс vs интерфейс с nil-указателем, type assertion, type switch (конспект — `info.txt`) |
| — typed_nil | `interfaces/typed_nil/` | `v == nil` не ловит typed nil (`*T` в interface) |
| — log_aggregator | `interfaces/log_aggregator/` | Конкурентный пайплайн с интерфейсами |

## Практические примеры (`examples/`)

| Пример | Описание |
|---|---|
| `json_config` | HTTP-сервер на `:8080` с динамической перезагрузкой конфигурации (работает до Ctrl+C) |
| `numeronym` | Нумеронимы (`kubernetes` → `k8s`) через `fmt.Stringer` |
| `cli_spinner` | Анимация спиннера в терминале |
| `string_validator` | Валидация строк через регулярные выражения |
| `sql_example` | SQL: таблицы video/author и подсчёт видео автора по имени (подзапрос с `IN`) |

`json_config` ищет `config.json` в текущем каталоге, поэтому из корня репозитория путь передаётся явно:

```bash
go run ./examples/json_config -config examples/json_config/config.json
```

## Собеседования (`interview/questions/`)

Вопросы с собеседований и ответы на них; Go-сниппеты с ожидаемым выводом проверены на Go 1.27.

| Файл | Темы |
|---|---|
| `bizone.md` | Слайсы и append с общим массивом, метод по nil-указателю, горутины в цикле |
| `CyberProtected.md` | Задание: потокобезопасное in-memory KV-хранилище (CRUD, `sync.RWMutex`, context) |
| `kuper.md` | Задание: gRPC-метод «добавить товар и оформить корзину» |
| `questions-Crosstech.md` | iota, слайсы и массивы, defer, nil в интерфейсе, мапы, мультиплексирование gRPC |
| `sber.md` | Слайсы, defer, WaitGroup, гонка на мапе, слияние отсортированных слайсов |
| `sbs.md` | 15 задач: указатели, переменная цикла, мапы, ресиверы методов, каналы, defer, гонки, goto |
| `защищенныеТелекоммуникации.md` | Слайсы, JSON и неэкспортируемые поля, каналы, GOMAXPROCS(1) и вытеснение горутин |

## Бенчмарки (`benchmarks/`)

Сравнение производительности:

- **Аллокация слайсов** — append с малой ёмкостью vs оптимальная ёмкость vs преаллокация
- **Сборка строк** — `+=` + `fmt.Sprintf` vs `+=` + `strconv.Itoa` vs `strings.Builder`

Бенчмарки написаны через `for b.Loop()` (Go 1.24+).

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

Сборка, проверки и тесты корневого модуля (`code_generation/` — отдельный модуль, см. выше):

```bash
go build ./...
go vet ./...
go test -race ./...
```

## Зависимости

| Модуль | Версия | Назначение |
|---|---|---|
| `golang.org/x/sync` | v0.24.0 | `errgroup` для управления горутинами |
| `golang.org/x/tools` | v0.52.0 | Обход AST и форматирование кода в кодогенераторе (модуль `code_generation/`) |
| `gorm.io/gorm` | v1.31.2 | GORM v2 для сгенерированных репозиториев; генератор берёт из него разбор тегов и именование колонок (модуль `code_generation/`) |
