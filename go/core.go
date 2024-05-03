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
	return FromNumber(result)
}

func Minus(arguments Tab) Tab {
	args := ToList(arguments)
	if len(args) == 0 {
		return Tab{}
	}
	if len(args) == 1 {
		return FromNumber(-ToNumber(args[0]))
	}
	result := ToNumber(args[0])
	for _, item := range args[1:] {
		result -= ToNumber(item)
	}
	return FromNumber(result)
}

func Multiply(arguments Tab) Tab {
	args := ToList(arguments)
	result := float64(1)
	for _, item := range args {
		result *= ToNumber(item)
	}
	return FromNumber(result)
}

func Divide(arguments Tab) Tab {
	args := ToList(arguments)
	result := ToNumber(args[0])
	for _, item := range args[1:] {
		result /= ToNumber(item)
	}
	return FromNumber(result)
}

func LessThan(arguments Tab) Tab {
	args := ToList(arguments)
	a := ToNumber(args[0])
	b := ToNumber(args[1])
	return FromBool(a < b)
}

func LessThanEqual(arguments Tab) Tab {
	args := ToList(arguments)
	a := ToNumber(args[0])
	b := ToNumber(args[1])
	return FromBool(a <= b)
}

func GreaterThan(arguments Tab) Tab {
	args := ToList(arguments)
	a := ToNumber(args[0])
	b := ToNumber(args[1])
	return FromBool(a > b)
}

func GreaterThanEqual(arguments Tab) Tab {
	args := ToList(arguments)
	a := ToNumber(args[0])
	b := ToNumber(args[1])
	return FromBool(a >= b)
}

func ParseNumber(arguments Tab) Tab {
	str := ToString(ToList(arguments)[0])
	result, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(err)
	}
	return FromNumber(result)
}

func Modulo(arguments Tab) Tab {
	args := ToList(arguments)
	a := ToNumber(args[0])
	b := ToNumber(args[1])
	return FromNumber(float64(int(a) % int(b)))
}

func Pow(arguments Tab) Tab {
	args := ToList(arguments)
	a := ToNumber(args[0])
	b := ToNumber(args[1])
	return FromNumber(math.Pow(a, b))
}

func Round(arguments Tab) Tab {
	args := ToList(arguments)
	a := ToNumber(args[0])
	return FromNumber(float64(int(a)))
}

func Ceil(arguments Tab) Tab {
	args := ToList(arguments)
	a := ToNumber(args[0])
	return FromNumber(float64(int(a) + 1))
}

func Floor(arguments Tab) Tab {
	args := ToList(arguments)
	a := ToNumber(args[0])
	return FromNumber(float64(int(a)))
}

func TIsNumber(arguments Tab) Tab {
	return FromBool(IsNumber(ToList(arguments)[0]))
}

// Strings

func CharAt(arguments Tab) Tab {
	args := ToList(arguments)
	str := ToString(args[0])
	index := int(ToNumber(args[1]))
	return FromString(string(str[index]))
}

func CharCode(arguments Tab) Tab {
	args := ToList(arguments)
	str := ToString(args[0])
	return FromNumber(float64(str[0]))
}

func TIsString(arguments Tab) Tab {
	return FromBool(IsString(ToList(arguments)[0]))
}

func SubStr(arguments Tab) Tab {
	args := ToList(arguments)
	str := ToString(args[0])
	start := int(ToNumber(args[1]))
	strlen := len(str)
	if start > strlen {
		return FromString("")
	}
	var end int
	if len(args) > 2 {
		end = int(ToNumber(args[2]))
		if end < start {
			return FromString("")
		}
		if end > strlen {
			end = strlen
		}
	} else {
		end = strlen
	}
	return FromString(str[start:end])
}

func StrJoin(arguments Tab) Tab {
	args := ToList(arguments)
	str := ToList(args[0])
	if len(str) == 0 {
		return FromString("")
	}
	separator := ToString(args[1])
	value := ToString(str[0])
	for _, item := range str[1:] {
		value += separator + ToString(item)
	}
	return FromString(value)
}

// Todo: should pretty print
func Str(arguments Tab) Tab {
	args := ToList(arguments)
	value := ""
	for _, item := range args {
		value += Print(item, false)
	}
	return FromString(value)
}

func StrStartsWith(arguments Tab) Tab {
	args := ToList(arguments)
	str := ToString(args[0])
	prefix := ToString(args[1])
	return FromBool(strings.HasPrefix(str, prefix))
}

func StrEndsWith(arguments Tab) Tab {
	args := ToList(arguments)
	str := ToString(args[0])
	prefix := ToString(args[1])
	return FromBool(strings.HasSuffix(str, prefix))
}

func StrLen(arguments Tab) Tab {
	return FromNumber(float64(len(ToString(ToList(arguments)[0]))))
}

func StrReplaceAll(arguments Tab) Tab {
	args := ToList(arguments)
	str := ToString(args[0])
	old := ToString(args[1])
	new := ToString(args[2])
	return FromString(strings.ReplaceAll(str, old, new))
}

func StrConcat(arguments Tab) Tab {
	args := ToList(arguments)
	result := ""
	for _, item := range args {
		result += ToString(item)
	}
	return FromString(result)
}

// Lists

func List(arguments Tab) Tab {
	return arguments
}

func Count(arguments Tab) Tab {
	arg := ToList(arguments)[0]
	if IsString(arg) {
		return FromNumber(float64(len(ToString(arg))))
	}
	return FromNumber(float64(len(ToList(arg))))
}

func Cons(arguments Tab) Tab {
	args := ToList(arguments)
	result := append(TabList{args[0]}, ToList(args[1])...)
	return FromList(result)
}

func Concat(arguments Tab) Tab {
	args := ToList(arguments)
	var result TabList
	for _, item := range args {
		result = append(result, ToList(item)...)
	}
	return FromList(result)
}

func Nth(arguments Tab) Tab {
	args := ToList(arguments)
	index := int(ToNumber(args[1]))
	arg := args[0]
	// if arg.IsString().ToBool() {
	// 	return FromString(string(arg.ToString()[index]))
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
	return FromList(list[1:])
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
	return FromList(values[start:end])
}

func TIsList(arguments Tab) Tab {
	arg := ToList(arguments)[0]
	return FromBool(IsList(arg))
}

// Dicts

func Dict(arguments Tab) Tab {
	args := ToList(arguments)
	result := map[string]Tab{}
	for i := 0; i < len(args); i += 2 {
		result[ToString(args[i])] = args[i+1]
	}
	return FromDict(result)
}

func TIsDict(arguments Tab) Tab {
	return FromBool(IsDict(ToList(arguments)[0]))
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
	return FromBool(ok)
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
	return FromDict(dict)
}

func Keys(arguments Tab) Tab {
	dict := ToDict(ToList(arguments)[0])
	var result TabList
	for key := range dict {
		result = append(result, FromString(key))
	}
	return FromList(result)
}

func Vals(arguments Tab) Tab {
	dict := ToDict(ToList(arguments)[0])
	var result TabList
	for _, value := range dict {
		result = append(result, value)
	}
	return FromList(result)
}

func Entries(arguments Tab) Tab {
	dict := ToDict(ToList(arguments)[0])
	var result TabList
	for key, value := range dict {
		result = append(result, FromList(TabList{FromString(key), value}))
	}
	return FromList(result)
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
	return FromDict(result)
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
	return FromDict(result)
}

// Other

func Equals(arguments Tab) Tab {
	args := ToList(arguments)
	a := args[0]
	b := args[1]
	aType := ToType(GetType(a))
	bType := ToType(GetType(b))
	if aType != bType {
		return FromBool(false)
	}
	switch aType {
	case TabNumberType:
		return FromBool(ToNumber(a) == ToNumber(b))
	case TabStringType:
		return FromBool(ToString(a) == ToString(b))
	case TabSymbolType:
		return FromBool(ToSymbol(a) == ToSymbol(b))
	case TabListType:
		aList := ToList(a)
		bList := ToList(b)
		if len(aList) != len(bList) {
			return FromBool(false)
		}
		for i := 0; i < len(aList); i++ {
			if !ToBool(Equals(FromList(TabList{aList[i], bList[i]}))) {
				return FromBool(false)
			}
		}
		return FromBool(true)
	case TabDictType:
		aDict := ToDict(a)
		bDict := ToDict(b)
		if len(aDict) != len(bDict) {
			return FromBool(false)
		}
		for key, value := range aDict {
			if !ToBool(Equals(FromList(TabList{value, bDict[key]}))) {
				return FromBool(false)
			}
		}
		return FromBool(true)
	case TabFuncType:
		return FromBool(&a == &b)
	case TabNativeFuncType:
		// TODO: cannot compare funcs, what to do?
		return FromBool(&a == &b)
	case TabMacroType:
		return FromBool(&a == &b)
	case TabNilType:
		// Nils are always equal
		return FromBool(true)
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
	return FromBool(IsNil(ToList(arguments)[0]))
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
	return FromDict(*ast.Position)
}

func HasAstPosition(arguments Tab) Tab {
	args := ToList(arguments)
	ast := args[0]
	return FromBool(ast.Position != nil)
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
	return FromString(string(bytes))
}

func Dirname(arguments Tab) Tab {
	args := ToList(arguments)
	path := ToString(args[0])
	return FromString(filepath.Dir(path))
}

func Basename(arguments Tab) Tab {
	args := ToList(arguments)
	path := ToString(args[0])
	return FromString(filepath.Base(path))
}

func PathJoin(arguments Tab) Tab {
	args := ToList(arguments)
	var parts []string
	for _, item := range args {
		parts = append(parts, ToString(item))
	}
	path := filepath.Join(parts...)
	return FromString(path)
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
	return FromString(path)
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
		result = append(result, FromString(file.Name()))
	}
	return FromList(result)
}

// Time

func TimeMs(arguments Tab) Tab {
	return FromNumber(float64(time.Now().UnixMilli()))
}

// Symbols

func Symbol(arguments Tab) Tab {
	return FromSymbol(ToString(ToList(arguments)[0]))
}

func TIsSymbol(arguments Tab) Tab {
	return FromBool(IsSymbol(ToList(arguments)[0]))
}

// Funcs

func TIsFunc(arguments Tab) Tab {
	arg := ToList(arguments)[0]
	return FromBool(IsFunc(arg) || IsNativeFunc(arg))
}

// Bools

func TIsBool(arguments Tab) Tab {
	return FromBool(IsBool(ToList(arguments)[0]))
}

// Plugins

func LoadPlugin(arguments Tab) Tab {
	args := ToList(arguments)
	filename := ToString(args[0])
	p, err := plugin.Open(filename)
	if err != nil {
		panic(FromString("Failed to open plugin: " + err.Error()))
	}
	value, err := p.Lookup("Export")
	if err != nil {
		panic(FromString("Failed to lookup Export"))
	}
	tab, ok := value.(*Tab)
	if !ok {
		panic(FromString("Plugin Export is not a Tab"))
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
	return FromString(string(out))
}

// Vars

func Var(arguments Tab) Tab {
	arg := ToList(arguments)[0]
	return FromVar(&arg)
}

func TIsVar(arguments Tab) Tab {
	return FromBool(IsVar(ToList(arguments)[0]))
}

func VarGet(arguments Tab) Tab {
	args := ToList(arguments)
	tvar := args[0]
	return *ToVar(tvar)
}

func VarSet(arguments Tab) Tab {
	args := ToList(arguments)
	tvar := args[0]
	value := args[1]
	*ToVar(tvar) = value
	return value
}

func AddCore(env Tab) Tab {
	// Math
	EnvSet(env, FromSymbol("is-number"), FromNativeFunc(TIsNumber))
	EnvSet(env, FromSymbol("+"), FromNativeFunc(Plus))
	EnvSet(env, FromSymbol("-"), FromNativeFunc(Minus))
	EnvSet(env, FromSymbol("*"), FromNativeFunc(Multiply))
	EnvSet(env, FromSymbol("/"), FromNativeFunc(Divide))
	EnvSet(env, FromSymbol(">"), FromNativeFunc(GreaterThan))
	EnvSet(env, FromSymbol(">="), FromNativeFunc(GreaterThanEqual))
	EnvSet(env, FromSymbol("<"), FromNativeFunc(LessThan))
	EnvSet(env, FromSymbol("<="), FromNativeFunc(LessThanEqual))
	EnvSet(env, FromSymbol("%"), FromNativeFunc(Modulo))
	EnvSet(env, FromSymbol("pow"), FromNativeFunc(Pow))
	EnvSet(env, FromSymbol("round"), FromNativeFunc(Round))
	EnvSet(env, FromSymbol("ceil"), FromNativeFunc(Ceil))
	EnvSet(env, FromSymbol("floor"), FromNativeFunc(Floor))
	EnvSet(env, FromSymbol("parse-number"), FromNativeFunc(ParseNumber))

	// Strings
	EnvSet(env, FromSymbol("str"), FromNativeFunc(Str))
	EnvSet(env, FromSymbol("is-string"), FromNativeFunc(TIsString))
	EnvSet(env, FromSymbol("char-at"), FromNativeFunc(CharAt))
	EnvSet(env, FromSymbol("char-code"), FromNativeFunc(CharCode))
	EnvSet(env, FromSymbol("sub-str"), FromNativeFunc(SubStr))
	EnvSet(env, FromSymbol("str-join"), FromNativeFunc(StrJoin))
	EnvSet(env, FromSymbol("str-starts-with"), FromNativeFunc(StrStartsWith))
	EnvSet(env, FromSymbol("str-ends-with"), FromNativeFunc(StrEndsWith))
	EnvSet(env, FromSymbol("str-len"), FromNativeFunc(StrLen))
	EnvSet(env, FromSymbol("str-replace-all"), FromNativeFunc(StrReplaceAll))
	EnvSet(env, FromSymbol("str-concat"), FromNativeFunc(StrConcat))

	// Lists
	EnvSet(env, FromSymbol("list"), FromNativeFunc(List))
	EnvSet(env, FromSymbol("is-list"), FromNativeFunc(TIsList))
	EnvSet(env, FromSymbol("count"), FromNativeFunc(Count))
	EnvSet(env, FromSymbol("cons"), FromNativeFunc(Cons))
	EnvSet(env, FromSymbol("concat"), FromNativeFunc(Concat))
	EnvSet(env, FromSymbol("nth"), FromNativeFunc(Nth))
	EnvSet(env, FromSymbol("first"), FromNativeFunc(First))
	EnvSet(env, FromSymbol("last"), FromNativeFunc(Last))
	EnvSet(env, FromSymbol("slice"), FromNativeFunc(Slice))
	EnvSet(env, FromSymbol("rest"), FromNativeFunc(Rest))

	// Dicts
	EnvSet(env, FromSymbol("dict"), FromNativeFunc(Dict))
	EnvSet(env, FromSymbol("is-dict"), FromNativeFunc(TIsDict))
	EnvSet(env, FromSymbol("get"), FromNativeFunc(Get))
	EnvSet(env, FromSymbol("has"), FromNativeFunc(Has))
	EnvSet(env, FromSymbol("set"), FromNativeFunc(Set))
	EnvSet(env, FromSymbol("keys"), FromNativeFunc(Keys))
	EnvSet(env, FromSymbol("vals"), FromNativeFunc(Vals))
	EnvSet(env, FromSymbol("entries"), FromNativeFunc(Entries))
	EnvSet(env, FromSymbol("assoc"), FromNativeFunc(Assoc))
	EnvSet(env, FromSymbol("dissoc"), FromNativeFunc(Dissoc))

	// Meta
	EnvSet(env, FromSymbol("tokenize"), FromNativeFunc(TTokenize))
	EnvSet(env, FromSymbol("parse"), FromNativeFunc(TParse))
	EnvSet(env, FromSymbol("read-string"), FromNativeFunc(ReadString))
	EnvSet(env, FromSymbol("get-ast-position"), FromNativeFunc(GetAstPosition))
	EnvSet(env, FromSymbol("env-get"), FromNativeFunc(TEnvGet))
	EnvSet(env, FromSymbol("env-set"), FromNativeFunc(TEnvSet))
	EnvSet(env, FromSymbol("env-new"), FromNativeFunc(EnvNew))
	EnvSet(env, FromSymbol("load-plugin"), FromNativeFunc(LoadPlugin))

	// Time
	EnvSet(env, FromSymbol("time-ms"), FromNativeFunc(TimeMs))

	// Symbol
	EnvSet(env, FromSymbol("symbol"), FromNativeFunc(Symbol))
	EnvSet(env, FromSymbol("is-symbol"), FromNativeFunc(TIsSymbol))

	// Funcs
	EnvSet(env, FromSymbol("is-func"), FromNativeFunc(TIsFunc))

	// Bools
	EnvSet(env, FromSymbol("is-boolean"), FromNativeFunc(TIsBool))

	// Files
	EnvSet(env, FromSymbol("file-read"), FromNativeFunc(FileRead))
	EnvSet(env, FromSymbol("dirname"), FromNativeFunc(Dirname))
	EnvSet(env, FromSymbol("basename"), FromNativeFunc(Basename))
	EnvSet(env, FromSymbol("path-join"), FromNativeFunc(PathJoin))
	EnvSet(env, FromSymbol("path-resolve"), FromNativeFunc(PathResolve))
	EnvSet(env, FromSymbol("read-dir"), FromNativeFunc(ReadDir))

	// Other
	EnvSet(env, FromSymbol("="), FromNativeFunc(Equals))
	EnvSet(env, FromSymbol("print"), FromNativeFunc(TPrint))
	EnvSet(env, FromSymbol("println"), FromNativeFunc(Println))
	EnvSet(env, FromSymbol("exit"), FromNativeFunc(Exit))
	EnvSet(env, FromSymbol("exec"), FromNativeFunc(Exec))
	EnvSet(env, FromSymbol("is-nil"), FromNativeFunc(TIsNil))

	// Vars
	EnvSet(env, FromSymbol("var"), FromNativeFunc(Var))
	EnvSet(env, FromSymbol("is-var"), FromNativeFunc(TIsVar))
	EnvSet(env, FromSymbol("deref"), FromNativeFunc(VarGet))
	EnvSet(env, FromSymbol("reset"), FromNativeFunc(VarSet))

	return env
}
