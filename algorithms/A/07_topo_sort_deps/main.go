package main

import "fmt"

// collectNodes собирает полное множество всех узлов графа: и ключи мапы,
// и те имена, что встречаются ТОЛЬКО как зависимости (gpu, core, requests).
// Без этого «листовые» зависимости-сироты никогда бы не попали в обход как стартовые.
func collectNodes(deps map[string][]string) []string {
	seen := make(map[string]bool)
	var nodes []string
	add := func(name string) {
		if !seen[name] {
			seen[name] = true
			nodes = append(nodes, name)
		}
	}
	for node, list := range deps {
		add(node)
		for _, d := range list {
			add(d)
		}
	}
	return nodes
}

// insertionSort — собственная сортировка строк (пакет sort запрещён).
// Нужна лишь для ДЕТЕРМИНИРОВАННОГО вывода: на корректность топосорта не влияет,
// любой валидный порядок допустим. Сложность O(n^2) на небольших списках.
func insertionSort(a []string) {
	for i := 1; i < len(a); i++ {
		key := a[i]
		j := i - 1
		for j >= 0 && lessStr(key, a[j]) {
			a[j+1] = a[j]
			j--
		}
		a[j+1] = key
	}
}

// lessStr — лексикографическое сравнение строк байт за байтом (пакет strings запрещён).
func lessStr(x, y string) bool {
	n := len(x)
	if len(y) < n {
		n = len(y)
	}
	for i := 0; i < n; i++ {
		if x[i] != y[i] {
			return x[i] < y[i]
		}
	}
	// общий префикс совпал — короче считается меньше
	return len(x) < len(y)
}

// topoSort печатает узлы так, что каждый узел идёт ПОСЛЕ всех своих зависимостей.
// Алгоритм: DFS с пост-порядком. Состояние узла кодируется двумя множествами:
//   - inProgress: узел на стеке рекурсии (защита от повторного входа и детектор цикла);
//   - done: узел уже напечатан.
//
// Время: O(V+E) на сам обход (каждый узел и ребро — по разу), но ручные
// insertionSort ради детерминированного вывода добавляют до O(V^2 + sum(deg^2)).
// Память: O(V) — мапы done/inProgress, множество узлов и стек рекурсии до O(V).
func topoSort(deps map[string][]string) {
	done := make(map[string]bool)       // узел полностью обработан и напечатан
	inProgress := make(map[string]bool) // узел в текущей ветке рекурсии

	var visit func(node string)
	visit = func(node string) {
		if done[node] {
			return // уже напечатан — пропускаем
		}
		if inProgress[node] {
			// По условию циклов нет. Но честная реализация их детектит:
			// повторный вход в узел, который ещё на стеке, — это цикл.
			fmt.Printf("обнаружен цикл на узле %q\n", node)
			return
		}
		inProgress[node] = true

		// Обходим зависимости в отсортированном порядке ради детерминизма.
		// Отсутствующий ключ -> nil-slice -> range пуст: узел-сирота просто пропустит цикл.
		ds := deps[node]
		if len(ds) > 1 {
			cp := make([]string, len(ds))
			copy(cp, ds) // не мутируем исходную мапу
			insertionSort(cp)
			ds = cp
		}
		for _, d := range ds {
			visit(d)
		}

		inProgress[node] = false
		done[node] = true
		fmt.Println(node) // пост-порядок: печать ПОСЛЕ всех зависимостей
	}

	// Стартуем со всех узлов в детерминированном порядке.
	nodes := collectNodes(deps)
	insertionSort(nodes)
	for _, n := range nodes {
		visit(n)
	}
}

func main() {
	deps := map[string][]string{
		"tensorflow": {"nvcc", "gpu", "linux"},
		"nvcc":       {"linux"},
		"linux":      {"core"},
		"mylib":      {"tensorflow"},
		"mylib2":     {"requests"},
	}
	topoSort(deps)
	// Фактический детерминированный вывод этой реализации (проверено go run):
	// core
	// gpu
	// linux
	// nvcc
	// tensorflow
	// mylib
	// requests
	// mylib2
	// Каждый узел напечатан ПОСЛЕ всех своих зависимостей — это один из валидных
	// топологических порядков (в общем случае порядок не единственный).
}
