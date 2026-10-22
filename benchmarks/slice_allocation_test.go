// go test -bench . -benchmem
//
// Сравниваем способы заполнить срез из nItems элементов:
//   - append в срез с маленькой начальной ёмкостью (много реаллокаций)
//   - append в срез с заранее выделенной ёмкостью
//   - запись по индексу в срез заранее выделенной длины
package bench

import "testing"

const nItems = 100_000

// sinkInts — запись результата в глобальную переменную не даёт компилятору
// выбросить заполнение среза как неиспользуемое.
var sinkInts []int

func BenchmarkSlice1(b *testing.B) {
	for b.Loop() {
		sl := make([]int, 0, 8)
		for i := range nItems {
			sl = append(sl, i)
		}
		sinkInts = sl
	}
}

func BenchmarkSlice2(b *testing.B) {
	for b.Loop() {
		sl := make([]int, 0, nItems)
		for i := range nItems {
			sl = append(sl, i)
		}
		sinkInts = sl
	}
}

func BenchmarkSlice3(b *testing.B) {
	for b.Loop() {
		sl := make([]int, nItems)
		for i := range nItems {
			sl[i] = i
		}
		sinkInts = sl
	}
}
