package main

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type Resp struct {
	Response []byte
	// Error error — в errgroup ошибку обычно возвращают отдельно,
	// но можно оставить, если это часть бизнес-логики.
}

func main() {
	// Пример вызова
	MyChanGroup(context.Background(), []string{"192.168.0.1", "127.0.0.1", "google.com"})
}

func MyChanGroup(ctx context.Context, addrs []string) error {
	// 1. Создаем errgroup с контекстом.
	// Если любая горутина вернет error != nil, ctxGroup отменится для всех остальных.
	g, ctxGroup := errgroup.WithContext(ctx)

	// 2. Буферизированный канал (оптимизация)
	ch := make(chan Resp, len(addrs))

	// Имитация клиента
	clMock := func(ctx context.Context, addr string) (Resp, error) {
		// Здесь реальная логика, которая уважает ctx
		return Resp{Response: []byte("data from " + addr)}, nil
	}

	g.SetLimit(10) // Максимум 10 активных горутин одновременно

	for _, addr := range addrs {
		addr := addr // Shadowing (для версий Go < 1.22)

		// 3. g.Go запускает горутину. Не нужно Add/Done.
		g.Go(func() error {
			// Используем ctxGroup! Если соседняя горутина упадет, этот контекст закроется.
			resp, err := clMock(ctxGroup, addr)
			if err != nil {
				return err // Это вызовет cancel() для всех остальных
			}

			select {
			case ch <- resp:
				return nil
			case <-ctxGroup.Done():
				return ctxGroup.Err() // Не зависаем при записи
			}
		})
	}

	// 4. Горутина для закрытия канала
	go func() {
		// Ждем завершения всех горутин (успешного или с ошибкой)
		_ = g.Wait()
		close(ch)
	}()

	// 5. Читаем результаты
	for resp := range ch {
		fmt.Printf("Received: %s\n", resp.Response)
	}

	// 6. Проверяем, была ли ошибка в группе
	if err := g.Wait(); err != nil {
		fmt.Printf("Finished with error: %v\n", err)
		return err
	}

	fmt.Println("Finished successfully")
	return nil
}
