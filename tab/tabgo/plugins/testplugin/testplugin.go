package main

import (
	t "jen-ya.de/tabgo"
)

var Export t.Tab

func doSomething(arguments t.Tab) t.Tab {
	args := t.ToList(arguments)
	return t.CallTab(t.Plus, args[0], args[0], args[1], args[1])
}

func init() {
	Export = t.DictToTab(t.TabDict{
		"do-something": t.NativeFuncToTab(doSomething),
	})
}
