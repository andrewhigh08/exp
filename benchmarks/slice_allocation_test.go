// go test -bench . -benchmem ./bench/slice_test.go
package bench

import "testing"

func BenchmarkSlice1(b *testing.B) {

	for i := 0; i < b.N; i++ {
		sl := make([]int, 0, 8)

		for i := 0; i < 100000; i++ {
			sl = append(sl, i)

		}

	}

}

func BenchmarkSlice2(b *testing.B) {

	for i := 0; i < b.N; i++ {
		sl := make([]int, 0, 100000)

		for i := 0; i < 100000; i++ {
			sl = append(sl, i)

		}

	}

}

func BenchmarkSlice3(b *testing.B) {

	for i := 0; i < b.N; i++ {
		sl := make([]int, 100000)

		for i := 0; i < 100000; i++ {
			sl[i] = i
		}

	}

}
