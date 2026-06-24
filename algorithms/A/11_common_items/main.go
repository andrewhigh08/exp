package main

import "fmt"

// Seller — продавец и его товары.
type Seller struct {
	sellerID int
	items    []string
}

// commonItems возвращает товары, которые есть у КАЖДОГО продавца (пересечение
// списков). Порядок результата совпадает с порядком товаров первого продавца.
//
// Время: O(G·I) — построение множеств товаров (суммарно по всем продавцам) плюс
// проверка каждого товара первого продавца по остальным G-1 множествам за O(1);
// G=len(goods), I — среднее число товаров у продавца.
// Память: O(G·I) — множества товаров продавцов.
func commonItems(goods []Seller) []string {
	if len(goods) == 0 {
		return nil
	}

	// Множества товаров всех продавцов, кроме первого.
	sets := make([]map[string]struct{}, len(goods))
	for i := 1; i < len(goods); i++ {
		s := make(map[string]struct{}, len(goods[i].items))
		for _, item := range goods[i].items {
			s[item] = struct{}{}
		}
		sets[i] = s
	}

	// Идём по товарам первого продавца в их порядке; берём те, что есть у всех.
	var res []string
	seen := make(map[string]struct{}) // защита от повторов в списке первого продавца
	for _, item := range goods[0].items {
		if _, dup := seen[item]; dup {
			continue
		}
		inAll := true
		for i := 1; i < len(goods); i++ {
			if _, ok := sets[i][item]; !ok {
				inAll = false
				break
			}
		}
		if inAll {
			seen[item] = struct{}{}
			res = append(res, item)
		}
	}
	return res
}

func main() {
	goods := []Seller{
		{sellerID: 1, items: []string{"варежки", "шуба", "валенки", "стол", "шапка", "шарф", "кофта", "рубашка"}},
		{sellerID: 2, items: []string{"шкаф", "тумба", "шапка", "стол", "ложка", "кофта"}},
		{sellerID: 3, items: []string{"вилка", "кофта", "тарелка", "шарф", "шуба", "стол", "кастрюля"}},
		{sellerID: 4, items: []string{"ковер", "тумба", "кофта", "пылесос", "валенки", "ложка", "стол"}},
	}
	fmt.Println(commonItems(goods)) // [стол кофта]

	single := []Seller{
		{sellerID: 4, items: []string{"ковер", "тумба", "кофта", "пылесос", "валенки", "ложка", "стол"}},
	}
	fmt.Println(commonItems(single)) // [ковер тумба кофта пылесос валенки ложка стол]
}
