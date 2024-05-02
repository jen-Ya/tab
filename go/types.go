package tabgo

type TabList []Tab
type TabDict map[string]Tab
type TabFunc struct {
	Ast    Tab
	Params Tab
	Env    Tab
}
type TabNativeFunc func(Tab) Tab

type TabVar struct {
	Pointer *Tab
}

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
	TAbVarType        TabType = iota
)

type Tab struct {
	Type     TabType
	Value    interface{}
	Position *TabDict
}

func (T Tab) String() string {
	return Print(T, true)
}

func GetType(T Tab) Tab {
	return TypeToTab(T.Type)
}

func IsList(T Tab) Tab {
	return BoolToTab(T.Type == TabListType)
}

func IsString(T Tab) Tab {
	return BoolToTab(T.Type == TabStringType)
}

func IsSymbol(T Tab) Tab {
	return BoolToTab(T.Type == TabSymbolType)
}

func IsNumber(T Tab) Tab {
	return BoolToTab(T.Type == TabNumberType)
}

func IsBool(T Tab) Tab {
	return BoolToTab(T.Type == TabBoolType)
}

func IsDict(T Tab) Tab {
	return BoolToTab(T.Type == TabDictType)
}

func IsType(T Tab) Tab {
	return BoolToTab(T.Type == TabTypeType)
}

func IsFunc(T Tab) Tab {
	return BoolToTab(T.Type == TabFuncType)
}

func IsNativeFunc(T Tab) Tab {
	return BoolToTab(T.Type == TabNativeFuncType)
}

func IsMacro(T Tab) Tab {
	return BoolToTab(T.Type == TabMacroType)
}

func IsOther(T Tab) Tab {
	return BoolToTab(T.Type == TabOtherType)
}

func IsVar(T Tab) Tab {
	return BoolToTab(T.Type == TAbVarType)
}

func IsNil(T Tab) Tab {
	return BoolToTab(T.Type == TabNilType)
}

func StringToTab(s string) Tab {
	return Tab{Type: TabStringType, Value: &s}
}

func SymbolToTab(s string) Tab {
	return Tab{Type: TabSymbolType, Value: &s}
}

func NumberToTab(n float64) Tab {
	return Tab{Type: TabNumberType, Value: &n}
}

func BoolToTab(b bool) Tab {
	return Tab{Type: TabBoolType, Value: &b}
}

func ListToTab(l TabList) Tab {
	return Tab{Type: TabListType, Value: &l}
}

func DictToTab(d TabDict) Tab {
	return Tab{Type: TabDictType, Value: &d}
}

func FuncToTab(f TabFunc) Tab {
	return Tab{Type: TabFuncType, Value: &f}
}

func MacroToTab(f TabFunc) Tab {
	return Tab{Type: TabMacroType, Value: &f}
}

func NativeFuncToTab(f TabNativeFunc) Tab {
	return Tab{Type: TabNativeFuncType, Value: f}
}

func TypeToTab(t TabType) Tab {
	return Tab{Type: TabTypeType, Value: &t}
}

func OtherToTab(o interface{}) Tab {
	return Tab{Type: TabOtherType, Value: o}
}

func VarToTab(o TabVar) Tab {
	return Tab{Type: TAbVarType, Value: o}
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
	return ListToTab(args)
}

func CallTab(fun TabNativeFunc, args ...Tab) Tab {
	return fun(ListToTab(args))
}
