package main

import "fmt"

// minInt — нижняя граница int для безопасного старта максимума.
const minInt = -int(^uint(0)>>1) - 1

// Result — выбранные уровни по каждой функции.
type Result map[string]int

// smallestRangeByRole выбирает по одному уровню из каждого отсортированного
// по неубыванию слайса так, чтобы (max-min) выбранных был минимален.
// Возвращает (ответ, true) либо (nil, false), если выбор невозможен
// (пустой вход или пустой слайс у какой-то роли).
//
// Время: O(K · sum(Ni)) — на каждом сдвиге указателя минимум ищется линейно за
// O(K); всего сдвигов не больше суммы длин (K ролей, Ni — длина i-й). С min-кучей
// было бы O(sum(Ni)·log K).
// Память: O(K) — roles/arrs/idx/bestIdx и результирующая мапа.
func smallestRangeByRole(levels map[string][]int) (Result, bool) {
	K := len(levels)
	if K == 0 {
		return nil, false
	}

	// Фиксируем порядок ролей: индексные слайсы дешевле обращений по ключу map.
	roles := make([]string, 0, K)
	for r, arr := range levels {
		if len(arr) == 0 {
			return nil, false // хотя бы одна роль без уровней -> ответа нет
		}
		roles = append(roles, r)
	}
	arrs := make([][]int, K)
	for i := range roles {
		arrs[i] = levels[roles[i]]
	}

	idx := make([]int, K) // указатели в каждом массиве

	// curMax — максимум среди текущих выбранных; поддерживаем инкрементально.
	curMax := minInt
	for i := 0; i < K; i++ {
		if v := arrs[i][0]; v > curMax {
			curMax = v
		}
	}

	bestSpread := int(^uint(0) >> 1) // +inf
	bestIdx := make([]int, K)

	for {
		// Минимум среди текущих элементов — O(K). Он определяет ширину окна
		// и какой указатель двигать.
		minPos := 0
		minVal := arrs[0][idx[0]]
		for i := 1; i < K; i++ {
			if v := arrs[i][idx[i]]; v < minVal {
				minVal, minPos = v, i
			}
		}

		if s := curMax - minVal; s < bestSpread {
			bestSpread = s
			copy(bestIdx, idx)
			if s == 0 {
				break // разброс 0 — лучше не бывает
			}
		}

		// Двигать имеет смысл ТОЛЬКО указатель минимума: уменьшить (max-min)
		// способно лишь увеличение текущего минимума.
		idx[minPos]++
		if idx[minPos] == len(arrs[minPos]) {
			break // массив исчерпан — окно больше не накроет все K ролей
		}
		// Новый элемент >= прежнего (массив отсортирован) — max только растёт.
		if nv := arrs[minPos][idx[minPos]]; nv > curMax {
			curMax = nv
		}
	}

	res := make(Result, K)
	for i := 0; i < K; i++ {
		res[roles[i]] = arrs[i][bestIdx[i]]
	}
	return res, true
}

func main() {
	t1 := map[string][]int{
		"backend": {1, 2, 2, 3}, "frontend": {1, 3}, "qa": {3, 4, 4}, "design": {2, 3},
	}
	r1, ok1 := smallestRangeByRole(t1)
	fmt.Println("test1:", r1, ok1) // -> все 3, spread 0

	t2 := map[string][]int{
		"backend": {5}, "frontend": {3, 6, 7, 10}, "qa": {3, 9, 11, 18}, "design": {20},
	}
	r2, ok2 := smallestRangeByRole(t2)
	fmt.Println("test2:", r2, ok2) // -> backend:5 frontend:6 qa:9 design:20

	r3, ok3 := smallestRangeByRole(map[string][]int{"x": {7}})
	fmt.Println("test3:", r3, ok3) // один массив
	r4, ok4 := smallestRangeByRole(map[string][]int{"a": {1}, "b": {}})
	fmt.Println("test4:", r4, ok4) // пустой слайс роли -> false
	r5, ok5 := smallestRangeByRole(map[string][]int{})
	fmt.Println("test5:", r5, ok5) // пустой вход -> false
}
