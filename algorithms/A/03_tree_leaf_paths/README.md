# Конечные узлы рубрикатора с полным путём

## Условие
На сайте есть рубрикатор:
```
Вещи
    Одежда
        Мужская
        Женская
Хобби
    Велосипеды
        Горные
    Мангалы
Транспорт
```
Надо распечатать конечные узлы рубрикатора (у которых нет детей) с полным путём.

Ожидаемый вывод:
```
Вещи => Одежда => Мужская
Вещи => Одежда => Женская
Хобби => Велосипеды => Горные
Хобби => Мангалы
Транспорт
```

Соответственно, есть класс `Node` такого вида:
```ts
class Node {
    text: string
    children: Node[]
}
```

## Входные параметры
Массив нод первого уровня `[Вещи, Хобби, Транспорт]`.

## Шаблон для Go
```go
package main

import "fmt"

type Node struct {
    title    string
    children Tree
}

type Tree = []Node

var input = Tree{
    Node{
        title: "Вещи",
        children: Tree{
            Node{
                title: "Одежда",
                children: Tree{
                    Node{title: "Мужская"},
                    Node{title: "Женская"},
                },
            },
        },
    },
    Node{
        title: "Хобби",
        children: Tree{
            Node{
                title: "Велосипеды",
                children: Tree{
                    Node{title: "Горные"},
                },
            },
            Node{title: "Мангалы"},
        },
    },
    Node{title: "Транспорт"},
}

func main() {
    printTree(input)
}

func printTree(rootNodes Tree) {
    // code
}
```
