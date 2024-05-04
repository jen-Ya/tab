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

func wrapDo(astlist []Tab) Tab {
	if len(astlist) > 1 {
		return FromList(append([]Tab{FromSymbol("do")}, astlist...))
	}
	return astlist[0]
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
				args := ToList(ast)[1:]
				// key value key value ... expression
				keyvals := ToList(args[0])
				for i := 0; i < len(keyvals); i += 2 {
					key := keyvals[i]
					value := Eval(keyvals[i+1], env)
					EnvSet(env, key, value)
				}
				exps := args[1:]
				for i := 0; i < len(exps)-1; i++ {
					Eval(exps[i], env)
				}
				ast = exps[len(exps)-1]
				continue

			case "do":
				args := ToList(ast)[1:]
				for i := 0; i < len(args)-1; i++ {
					Eval(args[i], env)
				}
				ast = args[len(args)-1]
				continue

			case "if":
				args := ToList(ast)[1:]
				cond := Eval(args[0], env)
				// condition is not nil and not boolean false
				truthy := !IsNil(cond) && (!IsBool(cond) || ToBool(cond))
				if truthy {
					// truthy
					ast = args[1]
					continue
				} else if len(args) > 2 {
					// falsy -> else
					ast = args[2]
					continue
				} else {
					// falsy -> no else
					return Tab{}
				}

			case "q":
				return ToList(ast)[1]

			// TODO: quasiquoteexpand
			case "qq":
				ast = Quasiquote(ToList(ast)[1])
				continue

			case "qqexpand":
				return Quasiquote(ToList(ast)[1])

			case "apply":
				args := ToList(ast)[1:]
				// eval first element as function to call
				first := Eval(args[0], env)
				applyArgs := ToList(EvalAst(FromList(args[1:]), env))
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

			case "f":
				args := ToList(ast)[1:]
				params := args[0]
				if !IsList(params) {
					params = FromList(TabList{params})
				}
				body := wrapDo(args[1:])
				return FromFunc(TabFunc{
					Ast:    body,
					Params: params,
					Env:    env,
				})

			case "macrof":
				args := ToList(ast)[1:]
				params := args[0]
				if !IsList(params) {
					params = FromList(TabList{params})
				}
				body := wrapDo((args[1:]))
				return FromMacro(TabFunc{
					Ast:    body,
					Params: params,
					Env:    env,
				})
			case "macroexpand":
				ast = Macroexpand(ToList(ast)[1], env)
				return ast
			case "try":
				args := ToList(ast)[1:]
				defer func() {
					if r := recover(); r != nil {
						symbol := ToList(args[1])[1]
						_env := Env(env)
						// check if paniced with type Tab
						if tab, ok := r.(Tab); ok {
							// if so, set .error to tab
							EnvSet(_env, symbol, tab)
						} else {
							// otherwise, set .error to string
							EnvSet(_env, symbol, FromString(fmt.Sprintf("%s", r)))
						}
						evaled = Eval(ToList(args[1])[2], _env)
					}
				}()
				return Eval(args[0], env)
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
