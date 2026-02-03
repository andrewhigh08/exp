```go
package main

import "fmt"

func main() {
	a := make([]int, 0, 10) // len=0 cap=10
	a = append(a, 10)       // [10] len=1 cap=10
	a = append(a, 20)       // [10, 20] len=2 cap=10
	b := a[:]               // [10, 20] len=2, cap=2
	// c := a[1:]           // [20] len=1, cap=9
	b = append(b, 30)       // [10,20,40]
	a = append(a, 40)       // [10,20,40]

	fmt.Println(a)
	fmt.Println(b)
}
```

```go
package main

import "fmt"

type A struct {}

func (A) Add(a, b int) int { // func (*A)..
	return a + b
}

func main() {
	c := (*A)(nil).Add(1, 2)
	fmt.Println(c)
}
```

```go
package main

func do_(process func(string), strs []string) {
	for _, str := range strs {
		go func() {
			process(str)
		}()
	}
}

import "sync"

func do(process func(string), strs []string) {
	var wg sync.WaitGroup
	for _, str := range strs {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			process(s)
		}(str)
	}
	wg.Wait()
}

```