задание:
```go
package main

import (
	"context"
)

// ------ internal/service/entity/basket/basket.go
type Basket struct {
	ID     uint64
	UserID uint64
	Items  []BasketItem
	Status string
	Total  uint64
}

// ------ internal/service/entity/basket/item.go
type BasketItem struct {
	BasketID  uint64
	ProductID uint64
	Count     uint64
	Price     uint64
}

// ------ internal/gateway/grpc/basket/dependencies.go
type BasketRepository interface {
	Load(ctx context.Context, userID uint64) (*Basket, error)
	Save(ctx context.Context, b *Basket) error
}

type CheckoutProducer interface {
	SendMessage(ctx context.Context, basket Basket) error
}

// ------ internal/gateway/grpc/basket/server.go
type BasketServer struct {
	repo     BasketRepository
	producer CheckoutProducer
}

// ------ pkg/server/grpc/basket.grpc.pb.go
type BasketServiceServer interface {
	AddItemAndOrder(context.Context, *AddItemRequest) (*EmptyResponse, error)
}

type AddItemRequest struct {
	UserID    uint64
	ProductID uint64
	Price     uint64
	Count     uint64
}

type EmptyResponse struct {
}

// ------ internal/gateway/grpc/basket/add_item.go
/*
- Реализовать метод grpc-сервера, который:
  - добавляет товар в корзину
  - оформляет корзину

- В оформленную корзину (basket.Status=="ordered") изменения вносить нельзя.
  При изменении состава корзины надо пересчитывать basket.Total=sum(count*price)

- Для оформления корзины необходимо:
  - сменить её статус на ordered
  - сообщить возможным потребителям с помощью сообщений в брокер.

- Все элементы в корзине должны быть уникальны по ключу ProductID
*/

func (bs *BasketServer) AddItemAndOrder(ctx context.Context, req *AddItemRequest) (*EmptyResponse, error) {
	basket, err := bs.repo.Load(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// place your code
	
	err = bs.repo.Save(ctx, basket)
	if err != nil {
		return nil, err
	}

	return &EmptyResponse{}, nil
}

```
решение:
```go
const maxBasketTotal = 1000

func (bs *BasketServer) AddItemAndOrder(ctx context.Context, req *AddItemRequest) (*EmptyResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}

	// Проверяем, что req.UserID не равен 0
	if req.UserID <= 0 {
		return nil, fmt.Errorf("userID is invalid: %d", req.UserID)
	}

	// Загружаем корзину по UserID
	basket, err := bs.repo.Load(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if basket == nil {
		return nil, fmt.Errorf("basket is nil")
	}

	// Проверяем, что basket.ID не равен 0
	if basket.ID <= 0 {
		return nil, fmt.Errorf("basket ID is invalid: %d", basket.ID)
	}

	// Проверяем, не оформлена ли корзина
	if basket.Status == "ordered" {
		return nil, fmt.Errorf("basket is already ordered and cannot be modified")
	}

	// Проверяем, что req.Count и req.Price не равны 0
	if req.Count <= 0 || req.Price <= 0 {
		return nil, fmt.Errorf("count or price is invalid: count=%d, price=%d", req.Count, req.Price)
	}

	// Флаг для проверки, есть ли товар уже в корзине
	itemFound := false

	// Пересчитываем сумму только за добавленный/обновленный элемент
	totalDelta := req.Count * req.Price

	// Итерируем по существующим элементам в корзине
	for i := range basket.Items {
		if basket.Items[i].ProductID == req.ProductID {
			// Если товар уже в корзине, увеличиваем количество и обновляем цену
			totalDelta = (basket.Items[i].Count+req.Count)*req.Price - (basket.Items[i].Count * basket.Items[i].Price)
			basket.Items[i].Count += req.Count
			basket.Items[i].Price = req.Price
			itemFound = true
			break
		}
	}

	// Если товар не найден, добавляем его в корзину
	if !itemFound {
		newItem := BasketItem{
			BasketID:  basket.ID,
			ProductID: req.ProductID,
			Count:     req.Count,
			Price:     req.Price,
		}
		basket.Items = append(basket.Items, newItem)
	}

	// Обновляем общую сумму корзины
	basket.Total = basket.Total + totalDelta

	// Проверяем, что basket.Total не превышает максимально допустимую сумму корзины
	if basket.Total > maxBasketTotal {
		return nil, fmt.Errorf("basket.Total %d exceeds maximum allowed value %d", basket.Total, maxBasketTotal)
	}

	// Обновляем статус корзины на "ordered"
	basket.Status = "ordered"

	// Сохраняем изменения в репозитории
	err = bs.repo.Save(ctx, basket)
	if err != nil {
		return nil, err
	}

	// Проверяем, что bs.producer не равен нулю
	if bs.producer != nil {
		// Отправляем сообщение через CheckoutProducer
		err = bs.producer.SendMessage(ctx, *basket)
		if err != nil {
			return nil, err
		}
	} else {
		log.Println("Warning: producer is nil, no message sent to CheckoutProducer for basket ID", basket.ID)
	}

	return &EmptyResponse{}, nil
}
```
