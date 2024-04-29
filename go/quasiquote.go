package tabgo

func QuasiquoteList(ast Tab) Tab {
	result := ListToTab(TabList{})
	astList := ToList(ast)
	for i := len(astList) - 1; i >= 0; i-- {
		elt := astList[i]
		if ToBool(IsList(elt)) &&
			len(ToList(elt)) > 0 &&
			ToBool(IsSymbol(ToList(elt)[0])) &&
			ToSymbol(ToList(elt)[0]) == "..unq" {
			result = ListToTab(TabList{
				SymbolToTab("concat"),
				ToList(elt)[1],
				result,
			})
		} else {
			result = ListToTab(TabList{
				SymbolToTab("cons"),
				Quasiquote(elt),
				result,
			})
		}
	}
	return result
}

func Quasiquote(ast Tab) Tab {
	if ToBool(IsList(ast)) &&
		len(ToList(ast)) > 0 &&
		ToBool(IsSymbol(ToList(ast)[0])) &&
		ToSymbol(ToList(ast)[0]) == "unq" {
		return ToList(ast)[1]
	}
	if ToBool(IsList(ast)) {
		return QuasiquoteList(ast)
	}
	if ToBool(IsDict(ast)) || ToBool(IsSymbol(ast)) {
		return ListToTab(TabList{
			SymbolToTab("q"),
			ast,
		})
	}
	return ast
}
