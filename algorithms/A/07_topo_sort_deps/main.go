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

// topoSort печатает узлы так, что каждый узел идёт ПОСЛЕ всех своих зависимостей.
// Алгоритм: DFS с пост-порядком. Состояние узла кодируется двумя множествами:
//   - inProgress: узел на стеке рекурсии (защита от повторного входа и детектор цикла);
//   - done: узел уже напечатан.
//
// Время: O(V+E) — каждый узел и ребро обходятся ровно по разу.
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
			fmt.Printf("обнаружен цикл на узле %q\n", node)
			return
		}
		inProgress[node] = true
		// Отсутствующий ключ -> nil-slice -> range пуст: узел-сирота просто пропустит цикл.
		// Порядок обхода зависимостей не важен: любой валидный топопорядок допустим.
		for _, d := range deps[node] {
			visit(d)
		}

		inProgress[node] = false
		done[node] = true
		fmt.Println(node) // пост-порядок: печать ПОСЛЕ всех зависимостей
	}
	// Стартуем со всех узлов. Порядок обхода зависит от итерации map (Go её
	// рандомизирует), поэтому вывод между запусками может отличаться — но всегда
	// остаётся валидным топологическим порядком.
	for _, n := range collectNodes(deps) {
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
	// Пример валидного вывода (порядок между запусками может отличаться из-за
	// рандомизации обхода map, но каждый узел всегда печатается ПОСЛЕ всех своих
	// зависимостей):
	//   core, linux, nvcc, gpu, tensorflow, mylib, requests, mylib2
}
