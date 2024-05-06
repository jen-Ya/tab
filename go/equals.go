package tabgo

import "fmt"

func Equals(a Tab, b Tab) bool {
	aType := a.Type
	bType := b.Type
	if aType != bType {
		return false
	}
	switch aType {
	// numbers, strings, symbols, and bools are equal if their values are equal
	case TabNumberType:
		return ToNumber(a) == ToNumber(b)
	case TabStringType:
		return ToString(a) == ToString(b)
	case TabSymbolType:
		return ToSymbol(a) == ToSymbol(b)
	case TabBoolType:
		return ToBool(a) == ToBool(b)
	// lists and dicts are equal if their elements are equal
	case TabListType:
		aList := ToList(a)
		bList := ToList(b)
		if len(aList) != len(bList) {
			return false
		}
		for i := 0; i < len(aList); i++ {
			if !Equals(aList[i], bList[i]) {
				return false
			}
		}
		return true
	case TabDictType:
		aDict := ToDict(a)
		bDict := ToDict(b)
		if len(aDict) != len(bDict) {
			return false
		}
		for key, value := range aDict {
			if !Equals(value, bDict[key]) {
				return false
			}
		}
		return true
	// funcs are equal if they are the same object
	case TabFuncType:
		return a == b
	case TabNativeFuncType:
		return a == b
	case TabMacroType:
		return a == b
	// vars are equal if they are the same object
	// so (= (var 42) (var 42)) is false
	case TabVarType:
		return a == b
	case TabNilType:
		// Nils are always equal
		return true
	default:
		panic(fmt.Sprint("Unknown type ", aType))
	}
}

func TEquals(arguments Tab) Tab {
	args := ToList(arguments)
	a := args[0]
	b := args[1]
	return FromBool(Equals(a, b))
}
