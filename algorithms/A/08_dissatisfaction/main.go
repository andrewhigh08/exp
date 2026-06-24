package main

import "fmt"

// quicksortInPlace сортирует срез по возрастанию на месте (in-place),
// без аллокаций сверх стека рекурсии. Схема разбиения Хоара.
func quicksortInPlace(a []int, lo, hi int) {
	for lo < hi {
		// Медиана из трёх как опорный — защита от худшего случая на уже отсортированных данных.
		mid := lo + (hi-lo)/2
		if a[mid] < a[lo] {
			a[mid], a[lo] = a[lo], a[mid]
		}
		if a[hi] < a[lo] {
			a[hi], a[lo] = a[lo], a[hi]
		}
		if a[hi] < a[mid] {
			a[hi], a[mid] = a[mid], a[hi]
		}
		pivot := a[mid]

		i, j := lo-1, hi+1
		for {
			for {
				i++
				if a[i] >= pivot {
					break
				}
			}
			for {
				j--
				if a[j] <= pivot {
					break
				}
			}
			if i >= j {
				break
			}
			a[i], a[j] = a[j], a[i]
		}
		// Рекурсия в меньшую часть, хвостовая итерация по большей — глубина стека O(log N).
		if j-lo < hi-(j+1) {
			quicksortInPlace(a, lo, j)
			lo = j + 1
		} else {
			quicksortInPlace(a, j+1, hi)
			hi = j
		}
	}
}

// lowerBound возвращает индекс первого элемента a[i] >= target (позиция вставки).
// Свой бинарный поиск на полуинтервале [lo, hi); результат в диапазоне [0, len(a)].
func lowerBound(a []int, target int) int {
	lo, hi := 0, len(a)
	for lo < hi {
		mid := lo + (hi-lo)/2
		if a[mid] < target {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}

// totalDissatisfaction — версия без stdlib (свои quicksort и бинпоиск).
// Сортирует goods на месте и для каждой потребности бинпоиском находит позицию
// вставки, сравнивая левого и правого соседа. Версия на пакете sort лежит в ./stdlib.
//
// Время: O((N+M)·log N) — сортировка goods за O(N log N) плюс M бинарных
// поисков по O(log N); N=len(goods), M=len(buyerNeeds).
// Память: O(log N) — стек рекурсии quicksort; сортировка in-place, прочее O(1).
func totalDissatisfaction(goods []int, buyerNeeds []int) int {
	if len(goods) == 0 {
		// Товаров нет — покупать нечего; неудовлетворённость не определена, возвращаем 0.
		return 0
	}

	quicksortInPlace(goods, 0, len(goods)-1)

	total := 0
	for _, need := range buyerNeeds {
		pos := lowerBound(goods, need) // первый goods[pos] >= need
		best := -1                     // минимальная разница, -1 = ещё не задана

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
