package tabgo

func QuasiquoteList(ast Tab) Tab {
	result := FromList(TabList{})
	astList := ToList(ast)
	for i := len(astList) - 1; i >= 0; i-- {
		elt := astList[i]
		if IsList(elt) &&
			len(ToList(elt)) > 0 &&
			IsSymbol(ToList(elt)[0]) &&
			ToSymbol(ToList(elt)[0]) == "..unq" {
			result = FromList(TabList{
				FromSymbol("concat"),
				ToList(elt)[1],
				result,
			})
		} else {
			result = FromList(TabList{
				FromSymbol("cons"),
				Quasiquote(elt),
				result,
			})
		}
	}
	return result
}

func Quasiquote(ast Tab) Tab {
	if IsList(ast) &&
		len(ToList(ast)) > 0 &&
		IsSymbol(ToList(ast)[0]) &&
		ToSymbol(ToList(ast)[0]) == "unq" {
		return ToList(ast)[1]
	}
	if IsList(ast) {
		return QuasiquoteList(ast)
	}
	if IsDict(ast) || IsSymbol(ast) {
		return FromList(TabList{
			FromSymbol("q"),
			ast,
		})
	}
	return ast
}
