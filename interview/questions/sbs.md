1.
```go
package main

import "fmt"

func main() {

	intRef := new(int)

	*intRef = 10

	ref := intRef

	*ref++

	fmt.Println(*intRef)
}
```
```go
package main

import "fmt"

func main() {

	intRef := new(int)
	fmt.Println(intRef) // 0x12345678901

	*intRef = 10
	fmt.Println(*intRef) // 0

	ref := intRef
	fmt.Println(ref) // 0x12345678901

	*ref++
	fmt.Println(*intRef) // 11

}
```
2.
```go
package main

import "fmt"

func main() {

	var refInSlice []*int

	for i := 0; i < 5; i++ {
		refInSlice = append(refInSlice, &i)
	}

	for _, ref := range refInSlice {
		fmt.Println(*ref)
	}

}
```
```go
package main

import "fmt"

func main() {

	var refInSlice []*int // nil

	fmt.Println(refInSlice) // []

	for i := 0; i < 5; i++ {
		refInSlice = append(refInSlice, &i) // кладем в мапу 5 адресов i (0x01234567890)
	}

	for _, ref := range refInSlice {
		fmt.Println(*ref) // разименовываем *ref, т.е. получаем значения из адресов 0 1 2 3 4
	}

}
```
3.
```go
package main

import "fmt"

func main() {

	mapToTest := map[string]uint8{
		"John":  11,
		"Дима":  27,
	}

/* 	// Удалите запись с John. Но прежде, убедитесь, что она есть!
	if // ... Code  := mapToTest["John"]; // ... Code {
		// ... Code
	} */

	fmt.Println(mapToTest) // Вывестись должен только Дима

}
```
```go
package main

import "fmt"

func main() {

	mapToTest := map[string]uint8{
		"John": 11,
		"Дима": 27,
	}

	/* 	// Удалите запись с John. Но прежде, убедитесь, что она есть!
	   	if // ... Code  := mapToTest["John"]; // ... Code {
	   		// ... Code
	   	}
	*/
	if _, ok := mapToTest["John"]; ok {
		delete(mapToTest, "John") // delete(m, k)  удалить элемент m[k] из карты m
	}

	fmt.Println(mapToTest) // Вывестись должен только Дима

}

```
4.
```go
package main

import "fmt"

func main() {

	var declare map[string]int

	initialize := make(map[string]int)

	// Что выведет принт?

	fmt.Println(declare == nil, initialize == nil) // Запишите: ... , ...

	return

	// Исправьте код

	declare["first"] = 1

	initializeMapInMap := make(map[string]map[string]int)

	// Заполните initializeMapInMap одной записью
	// ... code

	fmt.Println(initializeMapInMap)

}
```
```go
package main

import "fmt"

func main() {

	var declare map[string]int // ниловая мапа
	fmt.Println(declare)       // map[]

	initialize := make(map[string]int) // инициализация значениями по умолчанию [""] 0
	fmt.Println(initialize)            // map[''] 0

	// Что выведет принт?

	fmt.Println(declare == nil, initialize == nil) // Запишите: ... , ... true, false

	// return

	// Исправьте код

	// declare["first"] = 1 // ошибка нельзя писать в ниловую мапу

	initializeMapInMap := make(map[string]map[string]int)

	// Заполните initializeMapInMap одной записью
	// ... code
	initializeMapInMap["name"] = initialize

	fmt.Println(initializeMapInMap)

}
```
5.
```go
package main

import "fmt"

// В чем разница при объявлении метода принимающего копию своей структуры и ссылку на неё?
// Что выведет fmt.Println в каждом случае?

type MethodStruct struct {
	Name string
}

func (m MethodStruct) GetCopyFromName() string {
	return m.Name
}

func (m MethodStruct) ChangeCopyFromName(s string) {
	m.Name = s
}

func (m *MethodStruct) GetNameFromRef() string {
	return m.Name
}

func (m *MethodStruct) ChangeNameFromRef(s string) {
	m.Name = s
}

func main() {

	var getName = MethodStruct{
		Name: "John",
	}

	fmt.Println("33.. Copy name: ", getName.GetCopyFromName())
	fmt.Println("34.. Ref  name: ", getName.GetNameFromRef())

	getName.ChangeCopyFromName("Dima")

	fmt.Println("38.. Copy name: ", getName.GetCopyFromName())
	fmt.Println("39.. Ref  name: ", getName.GetNameFromRef())

	getName.ChangeNameFromRef("Dima")

	fmt.Println("43.. Copy name: ", getName.GetCopyFromName())
	fmt.Println("44.. Ref  name: ", getName.GetNameFromRef())

}
```
```go
package main

import "fmt"

// В чем разница при объявлении метода принимающего копию своей структуры и ссылку на неё?
// Что выведет fmt.Println в каждом случае?

type MethodStruct struct {
	Name string
}

func (m MethodStruct) GetCopyFromName() string {
	return m.Name
}

func (m MethodStruct) ChangeCopyFromName(s string) {
	m.Name = s
}

func (m *MethodStruct) GetNameFromRef() string {
	return m.Name
}

func (m *MethodStruct) ChangeNameFromRef(s string) {
	m.Name = s
}

func main() {

	var getName = MethodStruct{
		Name: "John",
	}

	fmt.Println("33.. Copy name: ", getName.GetCopyFromName())
	fmt.Println("34.. Ref  name: ", getName.GetNameFromRef())

	getName.ChangeCopyFromName("Dima")

	fmt.Println("38.. Copy name: ", getName.GetCopyFromName()) // тут ошибся, распечатает john
	fmt.Println("39.. Ref  name: ", getName.GetNameFromRef())

	getName.ChangeNameFromRef("Dima")

	fmt.Println("43.. Copy name: ", getName.GetCopyFromName())
	fmt.Println("44.. Ref  name: ", getName.GetNameFromRef())

}
```
6.
```go
package main

import (
	"fmt"
)

func main() {

	message := make(chan int)

	message <- 11

	go PrintChanData(message)

}

func PrintChanData(c chan int) {

	data := <-c

	fmt.Println(data)
}
```
```go
package main

import (
	"fmt"
)

func main() {

	message := make(chan int)

	go PrintChanData(message)

	message <- 11

	var input string
	fmt.Scanln(&input)
}

func PrintChanData(c chan int) {

	data := <-c

	fmt.Println(data)
}
```
7.
```go
package main

import "fmt"

// Необходимо, используя оператор «…», передать через аргумент функции Sum() N чисел и вернуть результат их суммы.

func main() {
	fmt.Println(Sum())
}

func Sum() int {
	return 0
}
```
```go
package main

import "fmt"

// Необходимо, используя оператор «…», передать через аргумент функции Sum() N чисел и вернуть результат их суммы.

func main() {
	fmt.Println(Sum(0, 1, 2, 3))
}

func Sum(input ...int) int {
	sum := 0
	for _, v := range input {
		sum += v
	}
	return sum
}
```
8.
```go
package main

import (
	"fmt"
)

func testDefer(name string) string {

	fmt.Println("в функции testDefer")

	return name
}

func main() {

	defer fmt.Println(testDefer("функция вернула значение"))

	fmt.Println("конец main функции")
}
```
```go
package main

import (
	"fmt"
)

func testDefer(name string) string {

	fmt.Println("в функции testDefer")

	return name
}

func main() {

	defer fmt.Println(testDefer("функция вернула значение"))

	fmt.Println("конец main функции")
}
```
9.
```go
package main

import "fmt"

type cache struct {
	data map[int]int
}

func main() {
	cch := cache{
		data: make(map[int]int),
	}

	defer fmt.Println(cch.get(10)) // какое выведется число?

	defer cch.store(10, 11).update(10, 22)
}

func (c cache) store(key, val int) cache {
	c.data[key] = val
	return c
}

func (c cache) update(key, val int) {
	if _, ok := c.data[key]; ok {
		c.data[key] = val
	}
}

func (c cache) get(key int) int {
	val := c.data[key]
	return val
}
```
```go
package main

import "fmt"

type cache struct {
	data map[int]int
}

func main() {
	cch := cache{
		data: make(map[int]int),
	}

	defer fmt.Println(cch.get(10)) // какое выведется число?  выведется 0, тк мы меняем копию

	defer cch.store(10, 11).update(10, 22)
}

func (c cache) store(key, val int) cache {
	c.data[key] = val
	return c
}

func (c cache) update(key, val int) {
	if _, ok := c.data[key]; ok {
		c.data[key] = val
	}
}

func (c cache) get(key int) int {
	val := c.data[key]
	return val
}
```
10.
```go
package main

import "fmt"

// Что такое состояние гонки и как решить эту проблему?

func main() {

	var counter int

	for i := 0; i < 1000; i++ {
		go func() {
			counter++
		}()
	}

	fmt.Println(counter)
}
```
```go
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Что такое состояние гонки и как решить эту проблему?

func main_() {

	var (
		counter int64
		wg      sync.WaitGroup
	)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}
	wg.Wait()

	fmt.Println(counter)
}

func main() {

	var (
		counter int
		wg      sync.WaitGroup
		mu      sync.Mutex
	)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()

	fmt.Println(counter)
}
```
11.
```go
package main

import "fmt"

func main() {

	in := make(chan int)

	go func(out <-chan int) {
		for i := 0; i <= 4; i++ {

			fmt.Println("начало цикла")

			out <- i

			fmt.Println("конец цикла")

		}

		fmt.Println("горутина завершилась")
	}(in)

	for i := range in {
		fmt.Println("\tданные из канала: ", i)
	}
}
```
```go
package main

import "fmt"

func main() {

	in := make(chan int)

	go func(out chan int) { // канал только для вычитывания out <-chan int
		for i := 0; i <= 4; i++ {

			fmt.Println("начало цикла")

			out <- i

			fmt.Println("конец цикла")

		}

		fmt.Println("горутина завершилась")
		close(out)
	}(in)

	for i := range in {
		fmt.Println("\tданные из канала: ", i)
	}
}
```
12.
```go
package main

import (
	"fmt"
	"sync"
)

// Необходимо привести код к корректному виду
// Код должен отрабатывать верно: без ошибок; без потенциальных багов
// Можно предложить несколько вариантов решения

func main() {
	ch := make(chan int)

	mu := sync.Mutex{}

	wg := sync.WaitGroup{}

	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func() {
			wg.Done()
			mu.Lock()
			ch <- i
			mu.Unlock()
		}()
	}

	wg.Wait()

	for {
		select {
		case v := <-ch:
			fmt.Println(v)

		}
	}

}
```
```go
package main

import (
	"fmt"
	"sync"
)

// Необходимо привести код к корректному виду
// Код должен отрабатывать верно: без ошибок; без потенциальных багов
// Можно предложить несколько вариантов решения

func main() {
	ch := make(chan int)
	wg := sync.WaitGroup{}

	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(j int) { // с версии 1.21 можно так не делать еще есть вариант с i := i
			defer wg.Done()
			ch <- j
		}(i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for v := range ch {
		fmt.Println(v)
	}
}

func main_() {
	ch := make(chan int, 10) // добавим буфер для уменьшения блокировок
	wg := sync.WaitGroup{}

	wg.Add(10)

	for i := 0; i < 10; i++ {
		go func(j int) { // с версии 1.21 можно так не делать еще есть вариант с i := i
			defer wg.Done()
			ch <- j
		}(i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for {
		select {
		case v, ok := <-ch:
			if !ok {
				fmt.Println("канал закрыт")
				return
			}
			fmt.Println(v)
		}
	}

}
```
13.
```go
package main

import "fmt"

// Что выведется в трех разных случаях и почему? Случаи:
// - вывод на старте без изменений;
// - если раскомментировать append на 12 строке;
// - если раскомментировать append на 25 строке.

func sliceConverter(arr []int) {

	// arr = append(arr, 1)

	for i := range arr {
		arr[i] = 5
	}

	fmt.Println(arr)
}

func main() {

	mainArr := []int{1, 2, 3, 4, 5}

	// mainArr = append(mainArr, 1)

	sliceConverter(mainArr)

	fmt.Println(mainArr)

}
```
```go
package main

import "fmt"

// Что выведется в трех разных случаях и почему? Случаи:
// - вывод на старте без изменений; [5 5 5 5 5] \n [5 5 5 5 5]
// - если раскомментировать append на 12 строке; [5 5 5 5 5 5] \n [1 2 3 4 5]
// - если раскомментировать append на 25 строке. [5 5 5 5 5 5 5]  /n [5 5 5 5 5 5 5] тут ошибся ввторым выведется [5 5 5 5 5 5]

func sliceConverter(arr []int) {

	arr = append(arr, 1) // произойдет переаллокация слайса arr будет ссылаться на новый массив len=6, cap=10
	//25 строчка, в функцию подается слайс len=6, cap=10   len становится 7, изменится исходный массив, который подается в функцию (mainArr)

	for i := range arr { // здесь
		arr[i] = 5
	}

	fmt.Println(arr)
}

func main() {

	mainArr := []int{1, 2, 3, 4, 5} // len=5, cap=5

	mainArr = append(mainArr, 1) // произойдет переаллокация слайса arr будет ссылаться на новый массив len=6, cap=10

	sliceConverter(mainArr)

	fmt.Println(mainArr)

}
```
14. Можно менять то что в defer и передвигать его, в результате должно вывестись 222
```go
package main

import "fmt"

func main() {
	fmt.Println(SomeTest(0, 0))
}

func SomeTest(x, y int) int {

	if x == 0 || y == 0 {
		goto get_x_y
	}

	defer func() {
		x = 111
		y = 111
	}()

get_x_y:
	x = 11
	y = 11

	return x + y
}
```
```go
func SomeTest(x, y int) int {

	if x == 0 || y == 0 {
		goto get_x_y
	}

get_x_y:
	x = 11
	y = 11

	defer fmt.Println(func() int {
		x = 111
		y = 111
		return 0
	}())

	return x + y

}
```
15.
```go
package main

import "fmt"

func main() {

	i := 0

metka:

	if i < 10 {
		i++
		goto metka
	}

	fmt.Println(i)
}
```