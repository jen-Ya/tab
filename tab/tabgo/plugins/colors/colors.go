package main

import (
	"github.com/fatih/color"
	t "jen-ya.de/tabgo"
)

var Export t.Tab

func WrapColor(attrs ...color.Attribute) t.TabNativeFunc {
	return func(arguments t.Tab) t.Tab {
		args := t.ToList(arguments)
		arg := t.ToString(args[0])
		return t.StringToTab(color.New(attrs...).SprintFunc()(arg))
	}
}

func init() {
	Export = t.DictToTab(t.TabDict{
		"cyan": t.NativeFuncToTab(WrapColor(color.FgCyan)),
		"red":  t.NativeFuncToTab(WrapColor(color.FgRed)),
	})
}
