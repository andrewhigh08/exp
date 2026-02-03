// Package main реализует потокобезопасный менеджер для сообщений по паттерну "Издатель-подписчик" (Publisher-Subscriber).
//
// Паттерн позволяет компонентам (издателям) отправлять сообщения в именованные "топики",
// не зная, кто их получит. Другие компоненты (подписчики) могут подписываться на эти
// топики, чтобы получать копии всех отправленных в них сообщений (Fan-Out).
package main

import (
	"log"
	"sync"
	"time"
)

// PubSubManager управляет подписками и рассылкой сообщений.
type PubSubManager struct {
	// mu защищает доступ к `topics`. RWMutex выбран потому, что публикаций
	// (чтение списка подписчиков) обычно гораздо больше, чем изменений в подписках.
	mu sync.RWMutex
	// topics хранит для каждого ID топика срез каналов его подписчиков.
	topics map[string][]chan any
}

// NewPubSubManager создает новый экземпляр менеджера.
func NewPubSubManager() *PubSubManager {
	return &PubSubManager{
		topics: make(map[string][]chan any),
	}
}

// Publish отправляет сообщение всем подписчикам указанного топика.
// Рассылка происходит по принципу Fan-Out.
func (p *PubSubManager) Publish(topicID string, msg any) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Проверяем, есть ли подписчики на данный топик.
	if subscribers, found := p.topics[topicID]; found {
		// Клонируем срез подписчиков, чтобы не блокировать мьютекс надолго.
		// Это быстрая операция, после которой можно отпустить мьютекс.
		subsCopy := make([]chan any, len(subscribers))
		copy(subsCopy, subscribers)

		go func() {
			// Отправляем сообщение всем подписчикам в отдельной горутине.
			for _, subChan := range subsCopy {
				// Используем неблокирующую отправку, чтобы медленный или неактивный
				// подписчик не мог заблокировать рассылку для остальных.
				select {
				case subChan <- msg:
				default:
					// Если канал подписчика переполнен или заблокирован,
					// мы просто пропускаем отправку ему этого сообщения.
					log.Printf("Канал подписчика для топика '%s' заблокирован. Сообщение пропущено.", topicID)
				}
			}
		}()
	}
}

// Subscribe подписывает нового клиента на топик и возвращает канал для получения сообщений.
func (p *PubSubManager) Subscribe(topicID string) chan any {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Создаем канал для нового подписчика.
	// Буферизация помогает справиться с кратковременными пиками сообщений.
	ch := make(chan any, 10)

	// Добавляем канал в список подписчиков топика.
	p.topics[topicID] = append(p.topics[topicID], ch)

	return ch
}

// Unsubscribe отписывает клиента от топика.
// subChan должен быть типа `chan any`, чтобы его можно было закрыть.
func (p *PubSubManager) Unsubscribe(topicID string, subChan chan any) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if subscribers, found := p.topics[topicID]; found {
		// Создаем новый срез, исключая из него отписавшийся канал.
		newSubscribers := make([]chan any, 0, len(subscribers)-1)
		for _, sub := range subscribers {
			if sub != subChan {
				newSubscribers = append(newSubscribers, sub)
			}
		}
		// Обновляем список подписчиков.
		p.topics[topicID] = newSubscribers
		// Закрываем канал, чтобы потребитель знал, что подписка прекращена.
		close(subChan)
	}
}

// Close завершает работу менеджера, отписывая всех подписчиков и закрывая их каналы.
func (p *PubSubManager) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for topicID, subscribers := range p.topics {
		for _, subChan := range subscribers {
			close(subChan)
		}
		// Очищаем карту топиков.
		delete(p.topics, topicID)
	}
}

func main() {
	m := NewPubSubManager()
	defer m.Close() // Гарантируем корректное завершение работы.

	// Подписчик 1
	sub1Chan := m.Subscribe("news")
	go func() {
		for msg := range sub1Chan {
			log.Printf("Подписчик 1 получил: %v", msg)
		}
		log.Println("Подписчик 1: канал закрыт.")
	}()

	// Подписчик 2
	sub2Chan := m.Subscribe("news")
	go func() {
		for msg := range sub2Chan {
			log.Printf("Подписчик 2 получил: %v", msg)
			time.Sleep(500 * time.Millisecond) // Имитация медленного потребителя
		}
		log.Println("Подписчик 2: канал закрыт.")
	}()

	// Публикуем сообщения
	m.Publish("news", "Привет, мир!")
	m.Publish("news", "Вторая новость")
	m.Publish("other_topic", "Это сообщение никто не получит")

	time.Sleep(1 * time.Second)

	// Отписываем первого подписчика
	log.Println("Отписываем Подписчика 1...")
	m.Unsubscribe("news", sub1Chan)

	// Публикуем еще одно сообщение, его получит только второй подписчик.
	m.Publish("news", "Третья новость для оставшихся")

	time.Sleep(2 * time.Second)
	log.Println("Завершение работы main.")
}
