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
		return t.FromString(color.New(attrs...).SprintFunc()(arg))
	}
}

func init() {
	Export = t.FromDict(t.TabDict{
		"cyan": t.FromNativeFunc(WrapColor(color.FgCyan)),
		"red":  t.FromNativeFunc(WrapColor(color.FgRed)),
	})
}
