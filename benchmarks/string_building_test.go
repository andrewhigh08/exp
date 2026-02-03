// go test -bench . -benchmem ./bench/string_test.go
package bench

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func BenchmarkStr(b *testing.B) {

	for i := 0; i < b.N; i++ {
		st := ""

		for i := 0; i < 100; i++ {
			st += fmt.Sprintf("%d", i)
		}

		st = ""

		for i := 0; i < 100; i++ {
			st += fmt.Sprintf("%d", i)
		}

	}

}

func BenchmarkConv(b *testing.B) {

	for i := 0; i < b.N; i++ {
		st := ""

		for i := 0; i < 100; i++ {
			st += strconv.Itoa(i)
		}

		st = ""

		for i := 0; i < 100; i++ {
			st += strconv.Itoa(i)
		}
	}

}

func BenchmarkStringbilder(b *testing.B) {

	for i := 0; i < b.N; i++ {
		var builder strings.Builder

		for i := 0; i < 100; i++ {
			builder.WriteString(strconv.Itoa(i))
		}

		var builder2 strings.Builder

		for i := 0; i < 100; i++ {
			builder2.WriteString(strconv.Itoa(i))
		}

	}
}

