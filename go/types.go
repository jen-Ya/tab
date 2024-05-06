package tabgo

type TabList []Tab
type TabDict map[string]Tab
type TabFunc struct {
	Ast    Tab
	Params Tab
	Env    Tab
}
type TabNativeFunc func(Tab) Tab

type TabVar *Tab

type TabType int

const (
	TabNilType        TabType = iota
	TabListType       TabType = iota
	TabStringType     TabType = iota
	TabSymbolType     TabType = iota
	TabNumberType     TabType = iota
	TabBoolType       TabType = iota
	TabDictType       TabType = iota
	TabTypeType       TabType = iota
	TabFuncType       TabType = iota
	TabNativeFuncType TabType = iota
	TabMacroType      TabType = iota
	TabOtherType      TabType = iota
	TabVarType        TabType = iota
)

func (TT TabType) String() string {
	switch TT {
	case TabNilType:
		return "nil"
	case TabListType:
		return "list"
	case TabStringType:
		return "string"
	case TabSymbolType:
		return "symbol"
	case TabNumberType:
		return "number"
	case TabBoolType:
		return "bool"
	case TabDictType:
		return "dict"
	case TabTypeType:
		return "type"
	case TabFuncType:
		return "func"
	case TabNativeFuncType:
		return "native-func"
	case TabMacroType:
		return "macro"
	case TabOtherType:
		return "other"
	case TabVarType:
		return "var"
	}
	return "unknown"
}

type Tab struct {
	Type     TabType
	Value    interface{}
	Position *TabDict
}

func (T Tab) String() string {
	return Print(T, true)
}

func GetType(T Tab) Tab {
	return FromType(T.Type)
}

func IsList(T Tab) bool {
	return T.Type == TabListType
}

func IsString(T Tab) bool {
	return T.Type == TabStringType
}

func IsSymbol(T Tab) bool {
	return T.Type == TabSymbolType
}

func IsNumber(T Tab) bool {
	return T.Type == TabNumberType
}

func IsBool(T Tab) bool {
	return T.Type == TabBoolType
}

func IsDict(T Tab) bool {
	return T.Type == TabDictType
}

func IsType(T Tab) bool {
	return T.Type == TabTypeType
}

func IsFunc(T Tab) bool {
	return T.Type == TabFuncType
}

func IsNativeFunc(T Tab) bool {
	return T.Type == TabNativeFuncType
}

func IsMacro(T Tab) bool {
	return T.Type == TabMacroType
}

func IsOther(T Tab) bool {
	return T.Type == TabOtherType
}

func IsVar(T Tab) bool {
	return T.Type == TabVarType
}

func IsNil(T Tab) bool {
	return T.Type == TabNilType
}

func FromString(s string) Tab {
	return Tab{Type: TabStringType, Value: &s}
}

func FromSymbol(s string) Tab {
	return Tab{Type: TabSymbolType, Value: &s}
}

func FromNumber(n float64) Tab {
	return Tab{Type: TabNumberType, Value: &n}
}

func FromBool(b bool) Tab {
	return Tab{Type: TabBoolType, Value: &b}
}

func FromList(l TabList) Tab {
	return Tab{Type: TabListType, Value: &l}
}

func FromDict(d TabDict) Tab {
	return Tab{Type: TabDictType, Value: &d}
}

func FromFunc(f TabFunc) Tab {
	return Tab{Type: TabFuncType, Value: &f}
}

func FromMacro(f TabFunc) Tab {
	return Tab{Type: TabMacroType, Value: &f}
}

func FromNativeFunc(f TabNativeFunc) Tab {
	return Tab{Type: TabNativeFuncType, Value: f}
}

func FromType(t TabType) Tab {
	return Tab{Type: TabTypeType, Value: &t}
}

func FromOther(o interface{}) Tab {
	return Tab{Type: TabOtherType, Value: o}
}

func FromVar(o TabVar) Tab {
	return Tab{Type: TabVarType, Value: o}
}

func ToString(T Tab) string {
	return *T.Value.(*string)
}

func ToSymbol(T Tab) string {
	return *T.Value.(*string)
}

func ToDict(T Tab) TabDict {
	return *T.Value.(*TabDict)
}

func ToList(T Tab) TabList {
	return *T.Value.(*TabList)
}

func ToNumber(T Tab) float64 {
	return *T.Value.(*float64)
}

func ToType(T Tab) TabType {
	return *T.Value.(*TabType)
}

func ToBool(T Tab) bool {
	return *T.Value.(*bool)
}

func ToFunc(T Tab) TabFunc {
	return *T.Value.(*TabFunc)
}

func ToNativeFunc(T Tab) TabNativeFunc {
	return T.Value.(TabNativeFunc)
}

func ToMacro(T Tab) TabFunc {
	return *T.Value.(*TabFunc)
}

func ToOther(T Tab) interface{} {
	return T.Value
}

func ToVar(T Tab) TabVar {
	return T.Value.(TabVar)
}

func ArgsToTab(args ...Tab) Tab {
	return FromList(args)
}

func CallTab(fun TabNativeFunc, args ...Tab) Tab {
	return fun(FromList(args))
}
