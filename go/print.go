package tabgo

import (
	"fmt"
	"strings"
)

type TabPrinter func(Tab, bool) string

func printEscape(str string) string {
	// escape newlines, tabs, quotes, and backslashes
	// TODO: can this be combined with Escape?
	// Probably no, because the parser does not care about backslashes or quotes
	str = strings.Replace(str, "\\", "\\\\", -1)
	str = strings.Replace(str, "\"", "\\\"", -1)
	str = strings.Replace(str, "\n", "\\n", -1)
	str = strings.Replace(str, "\t", "\\t", -1)
	return str
}

func PrintList(ast Tab, readable bool) string {
	var parts []string
	for _, val := range ToList(ast) {
		parts = append(parts, Print(val, readable))
	}
	return "(" + strings.Join(parts, " ") + ")"
}

func PrintDict(ast Tab, readable bool) string {
	parts := []string{"dict"}
	for key, val := range ToDict(ast) {
		parts = append(parts, "\""+key+"\" "+Print(val, readable))
	}
	return "(" + strings.Join(parts, " ") + ")"
}

var Printers map[TabType]TabPrinter

func init() {
	Printers = map[TabType]TabPrinter{
		TabListType: PrintList,
		TabSymbolType: func(ast Tab, readable bool) string {
			return ToSymbol(ast)
		},
		TabStringType: func(ast Tab, readable bool) string {
			if readable {
				return "\"" + printEscape(ToString(ast)) + "\""
			}
			return ToString(ast)
		},
		TabNumberType: func(ast Tab, readable bool) string {
			return fmt.Sprint(ToNumber(ast))
		},
		TabBoolType: func(ast Tab, readable bool) string {
			if ToBool(ast) {
				return "true"
			}
			return "false"
		},
		TabDictType: PrintDict,
		TabTypeType: func(ast Tab, readable bool) string {
			typestr := ToType(ast).String()
			return fmt.Sprintf("#<type:%s>", typestr)
		},
		TabFuncType: func(ast Tab, readable bool) string {
			return Print(ToFunc(ast).Ast, readable)
		},
		TabMacroType: func(ast Tab, readable bool) string {
			return "#<macro>"
		},
		TabNativeFuncType: func(ast Tab, readable bool) string {
			return "#<nativefunc>"
		},
		TabOtherType: func(ast Tab, readable bool) string {
			return "#<other>"
		},
		TabNilType: func(ast Tab, readable bool) string {
			return "nil"
		},
		TabVarType: func(ast Tab, readable bool) string {
			return "(var " + Print(*ToVar(ast), readable) + ")"
		},
	}
}

func AddPrinter(typ TabType, printer TabPrinter) {
	Printers[typ] = printer
}

func Print(ast Tab, readable bool) string {
	typ := ast.Type
	if printer, ok := Printers[typ]; ok {
		return printer(ast, readable)
	}
	return "#<unknown>"
}
