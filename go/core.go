package tabgo

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"strconv"
	"strings"
	"time"
)

// Math

func Plus(arguments Tab) Tab {
	args := ToList(arguments)
	result := float64(0)
	for _, item := range args {
		result += ToNumber(item)
	}
	return NumberToTab(result)
}

func Minus(arguments Tab) Tab {
	args := ToList(arguments)
	if len(args) == 0 {
		return Tab{}
	}
	if len(args) == 1 {
		return NumberToTab(-ToNumber(args[0]))
	}
	result := ToNumber(args[0])
	for _, item := range args[1:] {
		result -= ToNumber(item)
	}
	return NumberToTab(result)
}

func Multiply(arguments Tab) Tab {
	args := ToList(arguments)
	result := float64(1)
	for _, item := range args {
		result *= ToNumber(item)
	}
	return NumberToTab(result)
}

func Divide(arguments Tab) Tab {
	args := ToList(arguments)
	result := ToNumber(args[0])
	for _, item := range args[1:] {
		result /= ToNumber(item)
	}
	return NumberToTab(result)
}

func LessThan(arguments Tab) Tab {
	args := ToList(arguments)
	a := ToNumber(args[0])
	b := ToNumber(args[1])
	return BoolToTab(a < b)
}

func LessThanEqual(arguments Tab) Tab {
	args := ToList(arguments)
	a := ToNumber(args[0])
	b := ToNumber(args[1])
	return BoolToTab(a <= b)
}

func GreaterThan(arguments Tab) Tab {
	args := ToList(arguments)
	a := ToNumber(args[0])
	b := ToNumber(args[1])
	return BoolToTab(a > b)
}

func GreaterThanEqual(arguments Tab) Tab {
	args := ToList(arguments)
	a := ToNumber(args[0])
	b := ToNumber(args[1])
	return BoolToTab(a >= b)
}

func ParseNumber(arguments Tab) Tab {
	str := ToString(ToList(arguments)[0])
	result, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(err)
	}
	return NumberToTab(result)
}

func Modulo(arguments Tab) Tab {
	args := ToList(arguments)
	a := ToNumber(args[0])
	b := ToNumber(args[1])
	return NumberToTab(float64(int(a) % int(b)))
}

func Pow(arguments Tab) Tab {
	args := ToList(arguments)
	a := ToNumber(args[0])
	b := ToNumber(args[1])
	return NumberToTab(math.Pow(a, b))
}

func Round(arguments Tab) Tab {
	args := ToList(arguments)
	a := ToNumber(args[0])
	return NumberToTab(float64(int(a)))
}

func Ceil(arguments Tab) Tab {
	args := ToList(arguments)
	a := ToNumber(args[0])
	return NumberToTab(float64(int(a) + 1))
}

func Floor(arguments Tab) Tab {
	args := ToList(arguments)
	a := ToNumber(args[0])
	return NumberToTab(float64(int(a)))
}

func TIsNumber(arguments Tab) Tab {
	return IsNumber(ToList(arguments)[0])
}

// Strings

func CharAt(arguments Tab) Tab {
	args := ToList(arguments)
	str := ToString(args[0])
	index := int(ToNumber(args[1]))
	return StringToTab(string(str[index]))
}

func CharCode(arguments Tab) Tab {
	args := ToList(arguments)
	str := ToString(args[0])
	return NumberToTab(float64(str[0]))
}

func TIsString(arguments Tab) Tab {
	return IsString(ToList(arguments)[0])
}

func SubStr(arguments Tab) Tab {
	args := ToList(arguments)
	str := ToString(args[0])
	start := int(ToNumber(args[1]))
	strlen := len(str)
	if start > strlen {
		return StringToTab("")
	}
	var end int
	if len(args) > 2 {
		end = int(ToNumber(args[2]))
		if end < start {
			return StringToTab("")
		}
		if end > strlen {
			end = strlen
		}
	} else {
		end = strlen
	}
	return StringToTab(str[start:end])
}

func StrJoin(arguments Tab) Tab {
	args := ToList(arguments)
	str := ToList(args[0])
	if len(str) == 0 {
		return StringToTab("")
	}
	separator := ToString(args[1])
	value := ToString(str[0])
	for _, item := range str[1:] {
		value += separator + ToString(item)
	}
	return StringToTab(value)
}

// Todo: should pretty print
func Str(arguments Tab) Tab {
	args := ToList(arguments)
	value := ""
	for _, item := range args {
		value += Print(item, false)
	}
	return StringToTab(value)
}

func StrStartsWith(arguments Tab) Tab {
	args := ToList(arguments)
	str := ToString(args[0])
	prefix := ToString(args[1])
	return BoolToTab(strings.HasPrefix(str, prefix))
}

func StrEndsWith(arguments Tab) Tab {
	args := ToList(arguments)
	str := ToString(args[0])
	prefix := ToString(args[1])
	return BoolToTab(strings.HasSuffix(str, prefix))
}

func StrLen(arguments Tab) Tab {
	return NumberToTab(float64(len(ToString(ToList(arguments)[0]))))
}

func StrReplaceAll(arguments Tab) Tab {
	args := ToList(arguments)
	str := ToString(args[0])
	old := ToString(args[1])
	new := ToString(args[2])
	return StringToTab(strings.ReplaceAll(str, old, new))
}

func StrConcat(arguments Tab) Tab {
	args := ToList(arguments)
	result := ""
	for _, item := range args {
		result += ToString(item)
	}
	return StringToTab(result)
}

// Lists

func List(arguments Tab) Tab {
	return ListToTab(ToList(arguments))
}

func Count(arguments Tab) Tab {
	arg := ToList(arguments)[0]
	if ToBool(IsString(arg)) {
		return NumberToTab(float64(len(ToString(arg))))
	}
	return NumberToTab(float64(len(ToList(arg))))
}

func Cons(arguments Tab) Tab {
	args := ToList(arguments)
	result := append(TabList{args[0]}, ToList(args[1])...)
	return ListToTab(result)
}

func Concat(arguments Tab) Tab {
	args := ToList(arguments)
	var result TabList
	for _, item := range args {
		result = append(result, ToList(item)...)
	}
	return ListToTab(result)
}

func Nth(arguments Tab) Tab {
	args := ToList(arguments)
	index := int(ToNumber(args[1]))
	arg := args[0]
	// if arg.IsString().ToBool() {
	// 	return StringToTab(string(arg.ToString()[index]))
	// }
	return ToList(arg)[index]
}

func First(arguments Tab) Tab {
	// TODO: this is safe if list is empty or nil, but maybe it shouldn't be?
	// e, g.: return ToList(ToList(arguments)[0])[0]
	args := ToList(arguments)
	if args[0].Type == TabNilType {
		return Tab{}
	}
	list := ToList(args[0])
	if len(list) == 0 {
		return Tab{}
	}
	return list[0]
}

func Last(arguments Tab) Tab {
	values := ToList(ToList(arguments)[0])
	return values[len(values)-1]
}

func Rest(arguments Tab) Tab {
	// TODO: this is safe if list is empty or nil, but maybe it shouldn't be?
	args := ToList(arguments)
	if args[0].Type == TabNilType {
		return Tab{}
	}
	list := ToList(args[0])
	if len(list) == 0 {
		return Tab{}
	}
	return ListToTab(list[1:])
}

func Slice(arguments Tab) Tab {
	args := ToList(arguments)
	values := ToList(args[0])
	start := int(ToNumber(args[1]))
	var end int
	if len(args) > 2 {
		end = int(ToNumber(args[2]))
	} else {
		end = len(values)
	}
	return ListToTab(values[start:end])
}

func TIsList(arguments Tab) Tab {
	return IsList(ToList(arguments)[0])
}

// Dicts

func Dict(arguments Tab) Tab {
	args := ToList(arguments)
	result := map[string]Tab{}
	for i := 0; i < len(args); i += 2 {
		result[ToString(args[i])] = args[i+1]
	}
	return DictToTab(result)
}

func TIsDict(arguments Tab) Tab {
	return IsDict(ToList(arguments)[0])
}

func Get(arguments Tab) Tab {
	args := ToList(arguments)
	dict := ToDict(args[0])
	key := ToString(args[1])
	return dict[key]
}

func Has(arguments Tab) Tab {
	args := ToList(arguments)
	dict := ToDict(args[0])
	key := ToString(args[1])
	_, ok := dict[key]
	return BoolToTab(ok)
}

func Set(arguments Tab) Tab {
	args := ToList(arguments)
	dict := make(TabDict)
	key := ToString(args[1])
	value := args[2]
	for key, value := range ToDict(args[0]) {
		dict[key] = value
	}
	dict[key] = value
	return DictToTab(dict)
}

func Keys(arguments Tab) Tab {
	dict := ToDict(ToList(arguments)[0])
	var result TabList
	for key := range dict {
		result = append(result, StringToTab(key))
	}
	return ListToTab(result)
}

func Vals(arguments Tab) Tab {
	dict := ToDict(ToList(arguments)[0])
	var result TabList
	for _, value := range dict {
		result = append(result, value)
	}
	return ListToTab(result)
}

func Entries(arguments Tab) Tab {
	dict := ToDict(ToList(arguments)[0])
	var result TabList
	for key, value := range dict {
		result = append(result, ListToTab(TabList{StringToTab(key), value}))
	}
	return ListToTab(result)
}

func Assoc(arguments Tab) Tab {
	args := ToList(arguments)
	dict := ToDict(args[0])
	kvs := args[1:]
	result := make(TabDict)
	for key, value := range dict {
		result[key] = value
	}
	for i := 0; i < len(kvs); i += 2 {
		key := ToString(kvs[i])
		value := kvs[i+1]
		result[key] = value
	}
	return DictToTab(result)
}

func Dissoc(arguments Tab) Tab {
	args := ToList(arguments)
	dict := ToDict(args[0])
	deleteSet := make(map[string]bool)
	for _, key := range args[1:] {
		deleteSet[ToString(key)] = true
	}
	result := make(TabDict)
	for key, value := range dict {
		if deleteSet[key] {
			continue
		}
		result[key] = value
	}
	return DictToTab(result)
}

// Other

func Equals(arguments Tab) Tab {
	args := ToList(arguments)
	a := args[0]
	b := args[1]
	aType := ToType(GetType(a))
	bType := ToType(GetType(b))
	if aType != bType {
		return BoolToTab(false)
	}
	switch aType {
	case TabNumberType:
		return BoolToTab(ToNumber(a) == ToNumber(b))
	case TabStringType:
		return BoolToTab(ToString(a) == ToString(b))
	case TabSymbolType:
		return BoolToTab(ToSymbol(a) == ToSymbol(b))
	case TabListType:
		aList := ToList(a)
		bList := ToList(b)
		if len(aList) != len(bList) {
			return BoolToTab(false)
		}
		for i := 0; i < len(aList); i++ {
			if !ToBool(Equals(ListToTab(TabList{aList[i], bList[i]}))) {
				return BoolToTab(false)
			}
		}
		return BoolToTab(true)
	case TabDictType:
		aDict := ToDict(a)
		bDict := ToDict(b)
		if len(aDict) != len(bDict) {
			return BoolToTab(false)
		}
		for key, value := range aDict {
			if !ToBool(Equals(ListToTab(TabList{value, bDict[key]}))) {
				return BoolToTab(false)
			}
		}
		return BoolToTab(true)
	case TabFuncType:
		return BoolToTab(&a == &b)
	case TabNativeFuncType:
		// TODO: cannot compare funcs, what to do?
		return BoolToTab(&a == &b)
	case TabMacroType:
		return BoolToTab(&a == &b)
	case TabNilType:
		// Nils are always equal
		return BoolToTab(true)
	default:
		panic(fmt.Sprint("Unknown type ", aType))
	}
}

func TPrint(arguments Tab) Tab {
	parts := []string{}
	for _, item := range ToList(arguments) {
		parts = append(parts, Print(item, true))
	}
	str := strings.Join(parts, " ")
	fmt.Println(str)
	return Tab{}
}

func Println(arguments Tab) Tab {
	parts := []string{}
	for _, item := range ToList(arguments) {
		parts = append(parts, Print(item, false))
	}
	str := strings.Join(parts, " ")
	fmt.Println(str)
	return Tab{}
}

func Exit(arguments Tab) Tab {
	code := int(ToNumber(ToList(arguments)[0]))
	os.Exit(code)
	return Tab{}
}

func TIsNil(arguments Tab) Tab {
	return IsNil(ToList(arguments)[0])
}

func TTokenize(arguments Tab) Tab {
	args := ToList(arguments)
	text := ToString(args[0])
	keepComments := false
	filename := ""
	if len(args) > 1 {
		keepComments = ToBool(args[1])
	}
	if len(args) > 2 {
		filename = ToString(args[2])
	}
	tokens, err := Tokenize(text, keepComments, filename)
	if err != nil {
		panic(err)
	}
	return tokens
}

func TParse(arguments Tab) Tab {
	args := ToList(arguments)
	tokens := args[0]
	ast, err := Parse(tokens)
	if err != nil {
		panic(err)
	}
	return ast
}

func ReadString(arguments Tab) Tab {
	args := ToList(arguments)
	string := ToString(args[0])
	keepComments := false
	filepath := ""
	if len(args) > 1 {
		if val, ok := ToDict(args[1])["keepComments"]; ok {
			keepComments = ToBool(val)
		}
		if val, ok := ToDict(args[1])["filename"]; ok {
			filepath = ToString(val)
		}
	}
	return Read(string, keepComments, filepath)
}

func GetAstPosition(arguments Tab) Tab {
	args := ToList(arguments)
	ast := args[0]
	if ast.Position == nil {
		return Tab{}
	}
	return DictToTab(*ast.Position)
}

func HasAstPosition(arguments Tab) Tab {
	args := ToList(arguments)
	ast := args[0]
	return BoolToTab(ast.Position != nil)
}

func SetAstPosition(arguments Tab) Tab {
	args := ToList(arguments)
	ast := args[0]
	position := ToDict(args[1])
	ast.Position = &position
	return ast
}

func TEnvGet(arguments Tab) Tab {
	args := ToList(arguments)
	env := args[0]
	key := args[1]
	return EnvGet(env, key)
}

func TEnvSet(arguments Tab) Tab {
	args := ToList(arguments)
	env := args[0]
	key := args[1]
	value := args[2]
	return EnvSet(env, key, value)
}

func EnvNew(arguments Tab) Tab {
	args := ToList(arguments)
	parent := args[0]
	return Env(parent)
}

// Files

func FileRead(arguments Tab) Tab {
	args := ToList(arguments)
	filepath := ToString(args[0])
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	return StringToTab(string(bytes))
}

func Dirname(arguments Tab) Tab {
	args := ToList(arguments)
	path := ToString(args[0])
	return StringToTab(filepath.Dir(path))
}

func Basename(arguments Tab) Tab {
	args := ToList(arguments)
	path := ToString(args[0])
	return StringToTab(filepath.Base(path))
}

func PathJoin(arguments Tab) Tab {
	args := ToList(arguments)
	var parts []string
	for _, item := range args {
		parts = append(parts, ToString(item))
	}
	path := filepath.Join(parts...)
	return StringToTab(path)
}

func PathResolve(arguments Tab) Tab {
	args := ToList(arguments)
	var parts []string
	for _, item := range args {
		parts = append(parts, ToString(item))
	}
	path, err := filepath.Abs(filepath.Join(parts...))
	if err != nil {
		panic(err)
	}
	return StringToTab(path)
}

func ReadDir(arguments Tab) Tab {
	args := ToList(arguments)
	path := ToString(args[0])
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	var result TabList
	for _, file := range files {
		result = append(result, StringToTab(file.Name()))
	}
	return ListToTab(result)
}

// Time

func TimeMs(arguments Tab) Tab {
	return NumberToTab(float64(time.Now().UnixMilli()))
}

// Symbols

func Symbol(arguments Tab) Tab {
	return SymbolToTab(ToString(ToList(arguments)[0]))
}

func TIsSymbol(arguments Tab) Tab {
	return IsSymbol(ToList(arguments)[0])
}

// Funcs

func TIsFunc(arguments Tab) Tab {
	return BoolToTab(ToBool(IsFunc(ToList(arguments)[0])) || ToBool(IsNativeFunc(ToList(arguments)[0])))
}

// Bools

func TIsBool(arguments Tab) Tab {
	return IsBool(ToList(arguments)[0])
}

// Plugins

func LoadPlugin(arguments Tab) Tab {
	args := ToList(arguments)
	filename := ToString(args[0])
	p, err := plugin.Open(filename)
	if err != nil {
		panic(StringToTab("Failed to open plugin: " + err.Error()))
	}
	value, err := p.Lookup("Export")
	if err != nil {
		panic(StringToTab("Failed to lookup Export"))
	}
	tab, ok := value.(*Tab)
	if !ok {
		panic(StringToTab("Plugin Export is not a Tab"))
	}
	return *tab
}

func Q(arguments Tab) Tab {
	return ToList(arguments)[0]
}

func Exec(arguments Tab) Tab {
	args := ToList(arguments)
	command := ToString(args[0])
	cmd := exec.Command(command)
	if len(args) > 1 {
		for _, item := range ToList(args[1]) {
			cmd.Args = append(cmd.Args, ToString(item))
		}
	}
	if len(args) > 2 {
		opts := ToDict(args[2])
		if dir, ok := opts["dir"]; ok {
			cmd.Dir = ToString(dir)
		}

	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	return StringToTab(string(out))
}

// Vars

func Var(arguments Tab) Tab {
	arg := ToList(arguments)[0]
	return VarToTab(TabVar{
		Pointer: &arg,
	})
}

func TIsVar(arguments Tab) Tab {
	return IsVar(ToList(arguments)[0])
}

func VarGet(arguments Tab) Tab {
	args := ToList(arguments)
	tvar := args[0]
	return *ToVar(tvar).Pointer
}

func VarSet(arguments Tab) Tab {
	args := ToList(arguments)
	tvar := args[0]
	value := args[1]
	*ToVar(tvar).Pointer = value
	return value
}

func AddCore(env Tab) Tab {
	// Math
	EnvSet(env, SymbolToTab("is-number"), NativeFuncToTab(TIsNumber))
	EnvSet(env, SymbolToTab("+"), NativeFuncToTab(Plus))
	EnvSet(env, SymbolToTab("-"), NativeFuncToTab(Minus))
	EnvSet(env, SymbolToTab("*"), NativeFuncToTab(Multiply))
	EnvSet(env, SymbolToTab("/"), NativeFuncToTab(Divide))
	EnvSet(env, SymbolToTab(">"), NativeFuncToTab(GreaterThan))
	EnvSet(env, SymbolToTab(">="), NativeFuncToTab(GreaterThanEqual))
	EnvSet(env, SymbolToTab("<"), NativeFuncToTab(LessThan))
	EnvSet(env, SymbolToTab("<="), NativeFuncToTab(LessThanEqual))
	EnvSet(env, SymbolToTab("%"), NativeFuncToTab(Modulo))
	EnvSet(env, SymbolToTab("pow"), NativeFuncToTab(Pow))
	EnvSet(env, SymbolToTab("round"), NativeFuncToTab(Round))
	EnvSet(env, SymbolToTab("ceil"), NativeFuncToTab(Ceil))
	EnvSet(env, SymbolToTab("floor"), NativeFuncToTab(Floor))
	EnvSet(env, SymbolToTab("parse-number"), NativeFuncToTab(ParseNumber))

	// Strings
	EnvSet(env, SymbolToTab("str"), NativeFuncToTab(Str))
	EnvSet(env, SymbolToTab("is-string"), NativeFuncToTab(TIsString))
	EnvSet(env, SymbolToTab("char-at"), NativeFuncToTab(CharAt))
	EnvSet(env, SymbolToTab("char-code"), NativeFuncToTab(CharCode))
	EnvSet(env, SymbolToTab("sub-str"), NativeFuncToTab(SubStr))
	EnvSet(env, SymbolToTab("str-join"), NativeFuncToTab(StrJoin))
	EnvSet(env, SymbolToTab("str-starts-with"), NativeFuncToTab(StrStartsWith))
	EnvSet(env, SymbolToTab("str-ends-with"), NativeFuncToTab(StrEndsWith))
	EnvSet(env, SymbolToTab("str-len"), NativeFuncToTab(StrLen))
	EnvSet(env, SymbolToTab("str-replace-all"), NativeFuncToTab(StrReplaceAll))
	EnvSet(env, SymbolToTab("str-concat"), NativeFuncToTab(StrConcat))

	// Lists
	EnvSet(env, SymbolToTab("list"), NativeFuncToTab(List))
	EnvSet(env, SymbolToTab("is-list"), NativeFuncToTab(TIsList))
	EnvSet(env, SymbolToTab("count"), NativeFuncToTab(Count))
	EnvSet(env, SymbolToTab("cons"), NativeFuncToTab(Cons))
	EnvSet(env, SymbolToTab("concat"), NativeFuncToTab(Concat))
	EnvSet(env, SymbolToTab("nth"), NativeFuncToTab(Nth))
	EnvSet(env, SymbolToTab("first"), NativeFuncToTab(First))
	EnvSet(env, SymbolToTab("last"), NativeFuncToTab(Last))
	EnvSet(env, SymbolToTab("slice"), NativeFuncToTab(Slice))
	EnvSet(env, SymbolToTab("rest"), NativeFuncToTab(Rest))

	// Dicts
	EnvSet(env, SymbolToTab("dict"), NativeFuncToTab(Dict))
	EnvSet(env, SymbolToTab("is-dict"), NativeFuncToTab(TIsDict))
	EnvSet(env, SymbolToTab("get"), NativeFuncToTab(Get))
	EnvSet(env, SymbolToTab("has"), NativeFuncToTab(Has))
	EnvSet(env, SymbolToTab("set"), NativeFuncToTab(Set))
	EnvSet(env, SymbolToTab("keys"), NativeFuncToTab(Keys))
	EnvSet(env, SymbolToTab("vals"), NativeFuncToTab(Vals))
	EnvSet(env, SymbolToTab("entries"), NativeFuncToTab(Entries))
	EnvSet(env, SymbolToTab("assoc"), NativeFuncToTab(Assoc))
	EnvSet(env, SymbolToTab("dissoc"), NativeFuncToTab(Dissoc))

	// Meta
	EnvSet(env, SymbolToTab("tokenize"), NativeFuncToTab(TTokenize))
	EnvSet(env, SymbolToTab("parse"), NativeFuncToTab(TParse))
	EnvSet(env, SymbolToTab("read-string"), NativeFuncToTab(ReadString))
	EnvSet(env, SymbolToTab("get-ast-position"), NativeFuncToTab(GetAstPosition))
	EnvSet(env, SymbolToTab("env-get"), NativeFuncToTab(TEnvGet))
	EnvSet(env, SymbolToTab("env-set"), NativeFuncToTab(TEnvSet))
	EnvSet(env, SymbolToTab("env-new"), NativeFuncToTab(EnvNew))
	EnvSet(env, SymbolToTab("load-plugin"), NativeFuncToTab(LoadPlugin))

	// Time
	EnvSet(env, SymbolToTab("time-ms"), NativeFuncToTab(TimeMs))

	// Symbol
	EnvSet(env, SymbolToTab("symbol"), NativeFuncToTab(Symbol))
	EnvSet(env, SymbolToTab("is-symbol"), NativeFuncToTab(TIsSymbol))

	// Funcs
	EnvSet(env, SymbolToTab("is-func"), NativeFuncToTab(TIsFunc))

	// Bools
	EnvSet(env, SymbolToTab("is-boolean"), NativeFuncToTab(TIsBool))

	// Files
	EnvSet(env, SymbolToTab("file-read"), NativeFuncToTab(FileRead))
	EnvSet(env, SymbolToTab("dirname"), NativeFuncToTab(Dirname))
	EnvSet(env, SymbolToTab("basename"), NativeFuncToTab(Basename))
	EnvSet(env, SymbolToTab("path-join"), NativeFuncToTab(PathJoin))
	EnvSet(env, SymbolToTab("path-resolve"), NativeFuncToTab(PathResolve))
	EnvSet(env, SymbolToTab("read-dir"), NativeFuncToTab(ReadDir))

	// Other
	EnvSet(env, SymbolToTab("="), NativeFuncToTab(Equals))
	EnvSet(env, SymbolToTab("print"), NativeFuncToTab(TPrint))
	EnvSet(env, SymbolToTab("println"), NativeFuncToTab(Println))
	EnvSet(env, SymbolToTab("exit"), NativeFuncToTab(Exit))
	EnvSet(env, SymbolToTab("exec"), NativeFuncToTab(Exec))
	EnvSet(env, SymbolToTab("is-nil"), NativeFuncToTab(TIsNil))

	// Vars
	EnvSet(env, SymbolToTab("var"), NativeFuncToTab(Var))
	EnvSet(env, SymbolToTab("is-var"), NativeFuncToTab(TIsVar))
	EnvSet(env, SymbolToTab("deref"), NativeFuncToTab(VarGet))
	EnvSet(env, SymbolToTab("reset"), NativeFuncToTab(VarSet))

	return env
}
