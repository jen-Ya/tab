package tabgo

import "fmt"

type TabEquals func(Tab, Tab) bool

var TabEqualsMap map[TabType]TabEquals

func init() {
	TabEqualsMap = map[TabType]TabEquals{
		TabNumberType: func(a Tab, b Tab) bool {
			return ToNumber(a) == ToNumber(b)
		},
		TabStringType: func(a Tab, b Tab) bool {
			return ToString(a) == ToString(b)
		},
		TabSymbolType: func(a Tab, b Tab) bool {
			return ToSymbol(a) == ToSymbol(b)
		},
		TabBoolType: func(a Tab, b Tab) bool {
			return ToBool(a) == ToBool(b)
		},
		TabListType: func(a Tab, b Tab) bool {
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
		},
		TabDictType: func(a Tab, b Tab) bool {
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
		},
		TabFuncType: func(a Tab, b Tab) bool {
			return a == b
		},
		TabNativeFuncType: func(a Tab, b Tab) bool {
			return a == b
		},
		TabMacroType: func(a Tab, b Tab) bool {
			return a == b
		},
		// vars are equal if they are the same object
		// so (= (var 42) (var 42)) is false
		TabVarType: func(a Tab, b Tab) bool {
			return a == b
		},
		// Nils are always equal
		TabNilType: func(a Tab, b Tab) bool {
			return true
		},
		TabTypeType: func(a Tab, b Tab) bool {
			return ToType(a) == ToType(b)
		},
	}
}

func Equals(a Tab, b Tab) bool {
	aType := a.Type
	bType := b.Type
	if aType != bType {
		return false
	}
	if equals, ok := TabEqualsMap[aType]; ok {
		return equals(a, b)
	}
	panic(fmt.Sprintf("Equals not implemented for type %s", aType))
}

func TEquals(arguments Tab) Tab {
	args := ToList(arguments)
	a := args[0]
	b := args[1]
	return FromBool(Equals(a, b))
}
