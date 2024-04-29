package tabgo

import (
	"fmt"
	"strings"
)

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

func Print(ast Tab, readable bool) string {
	switch ToType(GetType(ast)) {
	case TabListType:
		return PrintList(ast, readable)
	case TabSymbolType:
		return ToSymbol(ast)
	case TabStringType:
		if readable {
			return "\"" + printEscape(ToString(ast)) + "\""
		}
		return ToString(ast)
	case TabNumberType:
		return fmt.Sprint(ToNumber(ast))
	case TabBoolType:
		if ToBool(ast) {
			return "true"
		}
		return "false"
	case TabDictType:
		return PrintDict(ast, readable)
	case TabTypeType:
		// TODO: implement
		return "#<type>"
	case TabFuncType:
		// TODO: implement?
		return Print(ToFunc(ast).Ast, readable)
	case TabMacroType:
		// TODO: implement?
		return "#<macro>"
	case TabNativeFuncType:
		// TODO: implement?
		return "#<nativefunc>"
	case TabOtherType:
		return "#<other>"
	case TabNilType:
		return "nil"
	}
	return "#<unknown>"
}

// const printList = (mal, start, stop, readable = true) => {
// 	return start + mal.map((val) => printString(val, readable)).join(' ') + stop;
// };

// const printHashmap = (mal, readable = true) => {
// 	let keys = Object.keys(mal);
// 	const printKey = key => '"' + key + '"' + ' ' + printString(mal[key], readable);
// 	return '(' + [
// 		'dict',
// 		...keys.map(printKey),
// 	].join(' ') + ')';
// };

// export const printString = (mal, readable = true) => {
// 	if(mal instanceof TalNil) {
// 		return 'nil';
// 	}
// 	if(mal === undefined) {
// 		return 'undefined';
// 	}
// 	if(
// 		[
// 			'boolean',
// 			'number',
// 		].indexOf(typeof mal) > -1
// 	) {
// 		return mal.toString();
// 	}
// 	if(typeof mal === 'string') {
// 		if(readable) {
// 			return '"' + mal
// 				.replace(/\\/g, '\\\\')
// 				.replace(/"/g, '\\"')
// 				.replace(/\n/g, '\\n')
// 				.replace(/\t/g, '\\t')
// 				+ '"';
// 		}
// 		return mal;
// 	}
// 	if(mal instanceof Array) {
// 		return printList(mal, '(', ')', readable);
// 	}
// 	if(mal instanceof TalMap) {
// 		return printHashmap(mal, readable);
// 	}
// 	if(mal instanceof TalFunc) {
// 		return '#<lambda>';
// 	}
// 	if(mal instanceof Function) {
// 		return '#<function>';
// 	}
// 	if(mal instanceof TalVar) {
// 		return `(var ${ printString(mal.value) })`;
// 	}
// 	if(mal instanceof TalSymbol) {
// 		return mal.valueOf();
// 	}
// 	return `#js<${ inspect(mal) }>`;
// };
