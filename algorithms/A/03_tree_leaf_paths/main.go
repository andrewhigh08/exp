package main

import "fmt"

// Node — узел рубрикатора: заголовок и потомки.
type Node struct {
	title    string
	children Tree
}

// Tree — список узлов (срез Node).
type Tree = []Node

// joinPath вручную склеивает накопленные сегменты пути через " => "
// (strings.Join запрещён — собираем строку сами).
func joinPath(path []string) string {
	res := ""
	for i, seg := range path {
		if i > 0 {
			res += " => "
		}
		res += seg
	}
	return res
}

// dfs — обход в пре-порядке (DFS). Накапливаем путь от корня в срезе path.
// Дойдя до листа (нет детей), печатаем полный путь, затем backtrack.
//
// Время: O(N·H) — каждый из N узлов посещается один раз, и на каждом листе путь
// длиной до H (высота дерева) заново склеивается в строку.
// Память: O(H) дополнительной — глубина рекурсии и срез path (без учёта вывода).
func dfs(nodes Tree, path []string) {
	for _, n := range nodes {
		path = append(path, n.title) // добавляем текущий сегмент
		if len(n.children) == 0 {
			fmt.Println(joinPath(path)) // лист — печатаем путь
		} else {
			dfs(n.children, path) // спускаемся глубже
		}
		path = path[:len(path)-1] // backtrack: убираем свой сегмент
	}
}

// printTree — точка входа: печатает все листья с полным путём от корня.
func printTree(rootNodes Tree) {
	dfs(rootNodes, nil)
}

func main() {
	tree := Tree{
		{title: "Вещи", children: Tree{
			{title: "Одежда", children: Tree{
				{title: "Мужская"},
				{title: "Женская"},
			}},
		}},
		{title: "Хобби", children: Tree{
			{title: "Велосипеды", children: Tree{
				{title: "Горные"},
			}},
			{title: "Мангалы"},
		}},
		{title: "Транспорт"},
	}
	printTree(tree)
}
