package tabgo

func Read(code string, keepComments bool, filepath string) Tab {
	// fmt.Println("tokenizing")
	tokens, err := Tokenize(code, keepComments, filepath)
	if err != nil {
		panic(err)
	}
	// fmt.Println("parsing")
	ast, err := Parse(tokens)
	if err != nil {
		panic(err)
	}
	return ast
}
