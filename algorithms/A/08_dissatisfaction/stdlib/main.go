package main

import (
	"fmt"
	"sort"
)

// totalDissatisfaction — версия на stdlib (та же задача, что в родительском
// каталоге, но без своих quicksort/бинпоиска):
//   - sort.Ints       — сортировка на месте;
//   - sort.SearchInts — позиция вставки: первый индекс, где goods[i] >= need,
//     то есть в точности lowerBound из ручной версии.
//
// Время: O((N+M)·log N) — сортировка goods за O(N log N) плюс M бинарных
// поисков по O(log N); N=len(goods), M=len(buyerNeeds).
// Память: O(1) дополнительной — sort.Ints сортирует на месте.
func totalDissatisfaction(goods []int, buyerNeeds []int) int {
	if len(goods) == 0 {
		// Товаров нет — покупать нечего; неудовлетворённость не определена, возвращаем 0.
		return 0
	}

	sort.Ints(goods)

	total := 0
	for _, need := range buyerNeeds {
		pos := sort.SearchInts(goods, need) // первый goods[pos] >= need
		best := -1                          // минимальная разница, -1 = ещё не задана

		// Правый сосед: первый элемент >= need (разница >= 0).
		if pos < len(goods) {
			d := goods[pos] - need
			if best < 0 || d < best {
				best = d
			}
		}
		// Левый сосед: последний элемент < need (разница > 0).
		if pos-1 >= 0 {
			d := need - goods[pos-1]
			if best < 0 || d < best {
				best = d
			}
		}
		total += best
	}
	return total
}

func main() {
	// Пример из условия: goods=[8,3,5], needs=[5,6] -> 1.
	fmt.Println(totalDissatisfaction([]int{8, 3, 5}, []int{5, 6})) // 1
}
