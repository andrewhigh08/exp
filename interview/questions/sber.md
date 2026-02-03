
```go
package main

import (
	"fmt"
	"sync"
	"time"
)

func mod(a []int) {
	a = append(a, 125)// len 5 cap 8

	for i := range a {
		a[i] = 5 // a [5 5 5 5 5]
	}

	fmt.Println(a) // 5 5 5 5 5
}

func main() {
	sl := []int{1, 2, 3, 4} // len 4 cap 4
	mod(sl)
	fmt.Println(sl) // 1 2 3 4
}
```

```go
package main

import (
"fmt"
)

func main() {
	i := 0
	defer fmt.Println(i) // 0
	i++
	return
}
```

```go
package main

import (
"fmt"
"sync"
)

func main() {
	wg := new(sync.WaitGroup)

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			fmt.Println(i)
			wg.Done()
		}()
	}

	wg.Wait()
}
```

```go
package main

import (
"fmt"
"sync"
"time"
)

func main() {
	var wg sync.WaitGroup{}

	wg.Add(1)
	go func() {
		time.Sleep(time.Second * 2)
		fmt.Println("1")
		wg.Done()
	}()

	go func() {
		fmt.Println("2")
	}()

	wg.Wait()
	fmt.Println("3")
}
```

```go
package main

import (
"fmt"
"time"
)

var m map[string]int

func main() {
	m = make(map[string]int)
	go f1()
	go f2()

	time.Sleep(time.Second * 10)

	fmt.Printf("print map = %+v", m)["f1" :99999]
}

func f1() {
	for i := 0; i < 100000; i++ {
		m["f1"]++
	}
}

func f2() {
	for i := 0; i < 100000; i++ {
		m["f2"]++
	}
}
```

Написать функцию, которая на вход принимает 2 отсортированных слайса с типом Int
в ответ выдает 1 отсортированный слайс, содержащий данные из входящих слайсов
Оценить сложность алгоритма:
```go
//  j
sl1[1, 3, 6, 9, 12] //len = 5
//  k
sl2[-1, 2, 4, 5, 6, 8] //len = 6

// res[-1, 1, 2, 3, 4, 5, 6, 6, 8, 9, 12]   len = 11

func merge(sl1 []int, sl2 []int) []int {
	l1, l2 := len(sl1), len(sl2)
	lSum := l1 + l2
	res := make([]int, 0, lSum)

	for j, k := 0, 0; j < l1 && k < l2;{
		if sl1[j] > ls2[k]{
			res = append(res, ls2[k])
			k++
		} else {
			res = append(res, ls1[j])
			j++
		}
	}

	if j < l1 {
		res = append(res, sl1[j:]...)
	}

	if k < l2 {
		res = append(res, sl2[k:]...)
	}

	return res
}
```

Архитектура - Постановка задачи: есть некая база, в которой хранятся данные о пользователях
есть несколько разных микросервисов, которые эти данные читают и пишут в разных комбинациях..
После увеличения количества микросервисов и запросов БД не справляется с нагрузкой -
предложить варианты масштабирования и решения проблем..





