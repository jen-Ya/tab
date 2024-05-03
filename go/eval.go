package tabgo

import (
	"fmt"
)

func EvalAst(ast Tab, env Tab) Tab {
	// fmt.Println("EvalAst:", Print(ast))
	switch ast.Type {
	case TabSymbolType:
		// todo: special symbols
		switch ToSymbol(ast) {
		case ".env":
			return env
		case ".filename":
			return ToDict(CallTab(GetAstPosition, ast))["filename"]
		}
		return EnvGet(env, ast)
	case TabListType:
		results := TabList{}
		for _, item := range ToList(ast) {
			results = append(results, Eval(item, env))
		}
		return FromList(results)
	case TabDictType:
		results := TabDict{}
		for key, item := range ToDict(ast) {
			results[key] = Eval(item, env)
		}
		return FromDict(results)
	}
	return ast
}

func Macroexpand(ast Tab, env Tab) Tab {
	for IsList(ast) &&
		len(ToList(ast)) > 0 &&
		IsSymbol(ToList(ast)[0]) &&
		IsMacro(EnvGet(env, ToList(ast)[0])) {
		astList := ToList(ast)
		macro := ToMacro(EnvGet(env, astList[0]))
		// if first param is '.caller', provide caller function also
		_env := Env(macro.Env)
		EnvSet(_env, FromSymbol(".caller"), astList[0])
		args := astList[1:]
		// TODO: .caller ?
		for i, param := range ToList(macro.Params) {
			if ToSymbol(param) == ".." {
				EnvSet(_env, ToList(macro.Params)[i+1], FromList(args[i:]))
				break
			}
			if len(args) <= i {
				break
			}
			EnvSet(_env, param, args[i])
		}
		ast = Eval(macro.Ast, _env)
	}
	return ast
}

func AstPositionToString(ast Tab) string {
	if ast.Position == nil {
		return "#<unknown>"
	}
	position := *ast.Position
	return fmt.Sprintf(
		"%s:%.0f:%.0f",
		ToString(position["filename"]),
		ToNumber(position["startLine"])+1,
		ToNumber(position["startChar"])+1,
	)
}

func Eval(ast Tab, env Tab) (evaled Tab) {
	// Traces panic, useful for debugging
	// Unfortunately, it also traces recovered panics
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		if tab, ok := r.(Tab); ok {
	// 			fmt.Printf("%s: %s\n", AstPositionToString(ast), Print(tab, true))
	// 		} else {
	// 			fmt.Printf("%s: %s\n", AstPositionToString(ast), r)
	// 		}
	// 		panic(r)
	// 	}
	// }()
	for {
		// fmt.Println("Eval:", Print(ast))
		if !IsList(ast) {
			return EvalAst(ast, env)
		}
		if len(ToList(ast)) == 0 {
			return ast
		}
		ast = Macroexpand(ast, env)
		if !IsList(ast) {
			// todo: evalast
			return EvalAst(ast, env)
		}
		// if first element is symbol
		if IsSymbol(ToList(ast)[0]) {
			// fmt.Println("Symbol:", ToList(ast)[0].ToSymbol())
			// swich on symbol
			switch ToSymbol(ToList(ast)[0]) {

			case "let":
				list := ToList(ast)
				v := Eval(list[2], env)
				EnvSet(env, list[1], v)
				return v

			case "eval":
				list := ToList(ast)[1:]
				result := EvalAst(
					FromList(list),
					env,
				)
				ast = ToList(result)[0]
				if len(ToList(result)) > 1 {
					env = ToList(result)[1]
				}
				continue
			// Allow multiple expressions to be evaluated in sequence
			// Also maybe it would be enough to implement as an immediatly invoked anonymous function
			case "with":
				env = Env(env)
				// key value key value ... expression
				list := ToList(ToList(ast)[1])
				for i := 0; i < len(list); i += 2 {
					key := list[i]
					value := Eval(list[i+1], env)
					EnvSet(env, key, value)
				}
				ast = ToList(ast)[2]
				continue

			case "do":
				list := ToList(ast)
				for i := 1; i < len(list)-1; i++ {
					Eval(list[i], env)
				}
				ast = list[len(list)-1]
				continue

			case "if":
				list := ToList(ast)
				cond := Eval(list[1], env)
				// condition is not nil and not boolean false
				truthy := !IsNil(cond) && (!IsBool(cond) || ToBool(cond))
				if truthy {
					// truthy
					ast = list[2]
					continue
				} else if len(list) > 3 {
					// falsy -> else
					ast = list[3]
					continue
				} else {
					// falsy -> no else
					return Tab{}
				}

			case "q":
				return ToList(ast)[1]

			// TODO: quasiquoteexpand
			case "qq":
				// fmt.Println("QQ FORM DETECTED")
				ast = Quasiquote(ToList(ast)[1])
				// fmt.Println("QQ EXPANDED:", Print(ast))
				continue

			case "qqexpand":
				return Quasiquote(ToList(ast)[1])

			case "apply":
				list := ToList(ast)
				// eval first element as function to call
				first := Eval(list[1], env)
				applyArgs := ToList(EvalAst(FromList(list[2:]), env))
				concats := applyArgs[0 : len(applyArgs)-1]
				last := ToList(applyArgs[len(applyArgs)-1])
				funcArgs := append(concats, last...)
				var f TabFunc
				if IsNativeFunc(first) {
					return ToNativeFunc(first)(FromList(funcArgs))
				} else if IsFunc(first) {
					f = ToFunc(first)
				} else if IsMacro(first) {
					f = ToMacro(first)
				} else {
					panic(fmt.Sprintf("Cannot call non-function: %s", Print(ast, true)))
				}

				env = Env(f.Env)
				for i, param := range ToList(f.Params) {
					if ToSymbol(param) == ".." {
						EnvSet(env, ToList(f.Params)[i+1], FromList(funcArgs[i:]))
						break
					}
					if len(funcArgs) <= i {
						break
					}
					EnvSet(env, param, funcArgs[i])
				}
				ast = f.Ast
				continue

			case "lambda":
				params := ToList(ast)[1]
				if !IsList(params) {
					params = FromList(TabList{params})
				}
				return FromFunc(TabFunc{
					Ast:    ToList(ast)[2],
					Params: params,
					Env:    env,
				})

			case "macrof":
				params := ToList(ast)[1]
				if !IsList(params) {
					params = FromList(TabList{params})
				}
				return FromMacro(TabFunc{
					Ast:    ToList(ast)[2],
					Params: params,
					Env:    env,
				})
			case "macroexpand":
				ast = Macroexpand(ToList(ast)[1], env)
				return ast
			case "try":
				astList := ToList(ast)
				defer func() {
					if r := recover(); r != nil {
						symbol := ToList(astList[2])[1]
						_env := Env(env)
						// check if paniced with type Tab
						if tab, ok := r.(Tab); ok {
							// if so, set .error to tab
							EnvSet(_env, symbol, tab)
						} else {
							// otherwise, set .error to string
							EnvSet(_env, symbol, FromString(fmt.Sprintf("%s", r)))
						}
						evaled = Eval(ToList(astList[2])[2], _env)
					}
				}()
				return Eval(astList[1], env)
			case "throw":
				panic(Eval(ToList(ast)[1], env))
			}

		}

		ulist := ToList(ast)
		first := Eval(ulist[0], env)
		switch ToType(GetType(first)) {
		case TabFuncType:
			fun := ToFunc(first)
			args := ToList(EvalAst(FromList(ulist[1:]), env))
			env = Env(fun.Env)
			for i, param := range ToList(fun.Params) {
				if IsNil(param) {
					continue
				}
				if ToSymbol(param) == ".." {
					EnvSet(env, ToList(fun.Params)[i+1], FromList(args[i:]))
					break
				}
				if len(args) <= i {
					break
				}
				EnvSet(env, param, args[i])
			}
			ast = ToFunc(first).Ast
			continue
		case TabNativeFuncType:
			args := EvalAst(FromList(ulist[1:]), env)
			return ToNativeFunc(first)(args)
		default:
			panic(fmt.Sprintf("Cannot call non-function of type %s at %s: %s", first.Type.String(), AstPositionToString(ast), Print(ast, true)))
		}

	}
}
