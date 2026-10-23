package processor

import (
	"errors"
	"sync"
)

// Экспортируемые ошибки библиотеки
var (
	ErrQueueFull = errors.New("queue is full")
	ErrClosed    = errors.New("processor is closed")
)

// Task — интерфейс задачи
type Task interface {
	Do()
}

// Processor — интерфейс асинхронного обработчика
type Processor interface {
	Do(t Task) error
	Close()
}

type processor struct {
	tasks chan Task
	wg    sync.WaitGroup
	once  sync.Once

	mu       sync.RWMutex
	isClosed bool
}

// NewProcessor создает обработчик.
// n — число одновременно выполняемых задач (воркеров)
// x — размер очереди (буфера)
func NewProcessor(n, x int) Processor {
	p := &processor{
		tasks: make(chan Task, x),
	}

	p.wg.Add(n)
	for range n { // Поддерживается в Go 1.22+
		go func() {
			defer p.wg.Done()
			for t := range p.tasks {
				t.Do()
			}
		}()
	}

	return p
}

func (p *processor) Do(t Task) error {
	// Держим блокировку на чтение (RLock) на время проверки флага isClosed
	// и отправки в канал. Это защищает от гонки: Close() выставляет isClosed
	// под Lock и только после этого закрывает канал p.tasks, поэтому если мы
	// увидели isClosed == false, канал не закроют, пока мы не отпустим RLock.
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.isClosed {
		return ErrClosed
	}

	select {
	case p.tasks <- t:
		return nil
	default: // Буфер полон — сразу возвращаем ошибку
		return ErrQueueFull
	}
}

func (p *processor) Close() {
	// sync.Once защищает от повторного вызова close(p.tasks)
	p.once.Do(func() {
		p.mu.Lock()
		p.isClosed = true
		p.mu.Unlock()

		close(p.tasks) // Безопасно закрываем канал
		p.wg.Wait()    // Ждем завершения всех воркеров
	})
}
