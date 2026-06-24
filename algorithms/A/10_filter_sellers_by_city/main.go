package main

import "fmt"

// filterSellersByCity возвращает продавцов, чей город входит в список искомых.
// Множество искомых городов даёт проверку «город нужен?» за O(1), поэтому
// достаточно одного прохода по продавцам (наивный вложенный перебор дал бы O(S·C)).
//
// Время: O(S + C) — построение множества из C городов плюс проход по S продавцам;
// S=len(sellers), C=len(cities).
// Память: O(S + C) — множество городов и результирующая мапа.
func filterSellersByCity(sellers map[int]string, cities []string) map[int]string {
	want := make(map[string]struct{}, len(cities))
	for _, c := range cities {
		want[c] = struct{}{}
	}

	res := make(map[int]string)
	for id, city := range sellers {
		if _, ok := want[city]; ok {
			res[id] = city
		}
	}
	return res
}

func main() {
	sellers := map[int]string{
		1: "Москва",
		2: "Самара",
		3: "Самара",
		4: "Тула",
		5: "Ростов",
		6: "Казань",
		7: "Курган",
		8: "Пенза",
	}
	citiesToFind := []string{"Самара", "Казань", "Тула"}

	fmt.Println(filterSellersByCity(sellers, citiesToFind))
	// map[2:Самара 3:Самара 4:Тула 6:Казань]
}
