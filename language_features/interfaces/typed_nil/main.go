// Package main — классический собеседовательный кейс:
// nil-интерфейс (nil, nil) vs интерфейс с typed nil (*T, nil).
//
// Проверка `v == nil` ловит только первый случай. Второй выглядит как <nil>
// при печати, но сравнение даёт false; вызов метода часто паникует.
package main

import "fmt"

type (
	iface interface {
		Echo()
	}

	Inst struct {
		val int
	}
)

func (i *Inst) Echo() {
	// Ресивер может быть (*Inst)(nil). Без проверки i.val паникует.
	if i == nil {
		fmt.Println("<nil Inst>")
		return
	}
	fmt.Println(i.val)
}

func f(v iface) error {
	// Интерфейс = (type, data). v == nil только если оба компонента nil.
	if v == nil {
		return fmt.Errorf("empty param")
	}

	// bar *Inst == nil, но в iface это (*Inst, nil) — typed nil.
	// v == nil здесь false, поэтому ловим type assert.
	if inst, ok := v.(*Inst); ok && inst == nil {
		return fmt.Errorf("empty param")
	}

	v.Echo()
	return nil
}

func main() {
	var foo = &Inst{} // (*Inst, &Inst{val:0}) — не nil
	var bar *Inst     // bar == nil, тип *Inst
	var err error

	{
		// v = (*Inst, &Inst{val:0})
		// v == nil → false
		// Echo печатает 0
		err = f(foo)     // nil
		fmt.Println(err) // <nil>
	}

	{
		// v = (*Inst, nil) — typed nil, не (nil, nil)
		// v == nil → false; ловим через type assert
		err = f(bar)     // "empty param"
		fmt.Println(err) // empty param
	}
}
