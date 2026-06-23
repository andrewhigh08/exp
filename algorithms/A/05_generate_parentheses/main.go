package main

import "fmt"

// genegateBrackets возвращает все правильные скобочные
// последовательности из n пар скобок.
//
// Время: O(4^n / sqrt(n)) — число валидных последовательностей (n-е число
// Каталана), умноженное на O(n) копирования каждой в результат.
// Память: O(n) дополнительной — буфер на 2n байт и глубина рекурсии до 2n
// (без учёта выходного среза, который занимает O(Каталан(n)·n)).
func genegateBrackets(n int) []string {
	res := []string{}
	if n <= 0 {
		// 0 пар -> единственная "пустая" последовательность.
		return append(res, "")
	}

	buf := make([]byte, 0, 2*n) // переиспользуемый буфер на 2n символов

	// open  — количество уже поставленных '('
	// close — количество уже поставленных ')'
	var backtrack func(open, close int)
	backtrack = func(open, close int) {
		if len(buf) == 2*n {
			// баланс уже соблюдён конструктивно -> последовательность валидна
			res = append(res, string(buf)) // string(buf) копирует байты
			return
		}
		// '(' можно ставить, пока не исчерпан лимит открывающих
		if open < n {
			buf = append(buf, '(')
			backtrack(open+1, close)
			buf = buf[:len(buf)-1] // откат
		}
		// ')' можно ставить, только если есть незакрытая открывающая
		if close < open {
			buf = append(buf, ')')
			backtrack(open, close+1)
			buf = buf[:len(buf)-1] // откат
		}
	}

	backtrack(0, 0)
	return res
}

func main() {
	fmt.Println(genegateBrackets(3))
	// [((())) (()()) (())() ()(()) ()()()]
}
