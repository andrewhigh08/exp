// go test -bench . -benchmem
//
// Сравниваем способы сборки строки из чисел:
//   - += с fmt.Sprintf
//   - += с strconv.Itoa
//   - strings.Builder + strconv.Itoa
//
// Важно: первые два бенчмарка включают стоимость конкатенации через +=
// (квадратичные аллокации), а не только конвертацию int→string.
package bench

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

const nNums = 100

// sinkString — запись результата в глобальную переменную не даёт компилятору
// выбросить сборку строки как неиспользуемую.
var sinkString string

func BenchmarkSprintfPlus(b *testing.B) {
	for b.Loop() {
		st := ""
		for j := range nNums {
			st += fmt.Sprintf("%d", j)
		}
		sinkString = st
	}
}

func BenchmarkItoaPlus(b *testing.B) {
	for b.Loop() {
		st := ""
		for j := range nNums {
			st += strconv.Itoa(j)
		}
		sinkString = st
	}
}

func BenchmarkBuilderItoa(b *testing.B) {
	for b.Loop() {
		var builder strings.Builder
		builder.Grow(nNums * 2) // грубая оценка; убирает лишние реаллокации
		for j := range nNums {
			builder.WriteString(strconv.Itoa(j))
		}
		sinkString = builder.String()
	}
}
