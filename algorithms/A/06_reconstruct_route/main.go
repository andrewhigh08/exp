package main

import "fmt"

// Ticket — один билет: откуда и куда.
type Ticket struct {
	From, To string
}

// reconstructRoute восстанавливает единственный неразрывный маршрут без петель и повторов.
// Идея: from->ticket для перехода и множество "to" для поиска старта.
//
// Время: O(N) — построение мап за один проход и линейный обход цепочки из N билетов.
// Память: O(N) — две мапы (byFrom, destinations) и выходной срез route.
func reconstructRoute(tickets []Ticket) []Ticket {
	n := len(tickets)
	if n <= 1 {
		// Пустой вход или один билет — маршрут совпадает со входом.
		// Возвращаем копию, чтобы не отдавать наружу исходный слайс.
		out := make([]Ticket, n)
		copy(out, tickets)
		return out
	}

	// byFrom: быстрый переход из города по ребру маршрута.
	byFrom := make(map[string]Ticket, n)
	// destinations: все города, в которые кто-то приезжает.
	destinations := make(map[string]struct{}, n)

	for _, t := range tickets {
		byFrom[t.From] = t
		destinations[t.To] = struct{}{}
	}

	// Старт — город, который есть среди From, но отсутствует среди To.
	start := ""
	found := false
	for _, t := range tickets {
		if _, isDest := destinations[t.From]; !isDest {
			start = t.From
			found = true
			break
		}
	}
	if !found {
		// Нет однозначного старта (цикл/некорректный вход) — отдаём копию входа.
		out := make([]Ticket, n)
		copy(out, tickets)
		return out
	}

	// Идём from->to, перекладывая билеты в порядке следования.
	route := make([]Ticket, 0, n)
	city := start
	for len(route) < n {
		t, ok := byFrom[city]
		if !ok {
			break // обрыв цепочки — дальше идти некуда
		}
		route = append(route, t)
		city = t.To
	}
	return route
}

func main() {
	tickets := []Ticket{
		{From: "London", To: "Moscow"},
		{From: "NY", To: "London"},
		{From: "Moscow", To: "SPb"},
	}
	route := reconstructRoute(tickets)
	for _, t := range route {
		fmt.Printf("%s -> %s\n", t.From, t.To)
	}
	// Вывод:
	// NY -> London
	// London -> Moscow
	// Moscow -> SPb
}
