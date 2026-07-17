// go test -bench . -benchmem
//
// Сравниваем способы сборки строки из чисел:
//  - += с fmt.Sprintf
//  - += с strconv.Itoa
//  - strings.Builder + strconv.Itoa
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

func BenchmarkSprintfPlus(b *testing.B) {
	var sink string
	for i := 0; i < b.N; i++ {
		st := ""
		for j := 0; j < nNums; j++ {
			st += fmt.Sprintf("%d", j)
		}
		sink = st
	}
	_ = sink
}

func BenchmarkItoaPlus(b *testing.B) {
	var sink string
	for i := 0; i < b.N; i++ {
		st := ""
		for j := 0; j < nNums; j++ {
			st += strconv.Itoa(j)
		}
		sink = st
	}
	_ = sink
}

func BenchmarkBuilderItoa(b *testing.B) {
	var sink string
	for i := 0; i < b.N; i++ {
		var builder strings.Builder
		builder.Grow(nNums * 2) // грубая оценка; убирает лишние реаллокации
		for j := 0; j < nNums; j++ {
			builder.WriteString(strconv.Itoa(j))
		}
		sink = builder.String()
	}
	_ = sink
}
