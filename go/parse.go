package tabgo

import "fmt"

func Parse(tokens Tab) (Tab, error) {
	list := ToList(tokens)
	next := list[0]
	rest := list
	ParseError := func(message string) error {
		if next.Position == nil {
			return fmt.Errorf("ParseError at unkown location\n%s", message)
		}
		position := *next.Position
		return fmt.Errorf("ParserError at %s:%.0f:%.0f\n%s", ToString(position["filename"]), ToNumber(position["startLine"])+1, ToNumber(position["startChar"])+1, message)
	}
	IsPeek := func(expected ...TabToken) bool {
		for i, _ := range expected {
			if TabToken(ToNumber(ToDict(rest[i])["kind"])) != expected[i] {
				return false
			}
		}
		return true
	}
	Consume := func(expected TabToken) (Tab, error) {
		// fmt.Printf("Consuming %s\n", expected.String())
		if !IsPeek(expected) {
			return Tab{}, ParseError(fmt.Sprintf("unexpected token %s %s, expected %s", next.Type, Print(next, true), expected))
		}
		// fmt.Println("consumed", expected, next.GetType().ToString())
		consumed := ToDict(next)["value"]
		consumed.Position = &TabDict{
			"filename":  (*next.Position)["filename"],
			"startLine": (*next.Position)["startLine"],
			"startChar": (*next.Position)["startChar"],
			"endLine":   (*next.Position)["endLine"],
			"endChar":   (*next.Position)["endChar"],
		}
		rest = rest[1:]
		next = rest[0]
		return consumed, nil
	}
	ConsumeAtom := func() (Tab, error) {
		if IsPeek(TabSymbolToken) {
			return Consume(TabSymbolToken)
		}
		if IsPeek(TabStringToken) {
			return Consume(TabStringToken)
		}
		if IsPeek(TabNumberToken) {
			return Consume(TabNumberToken)
		}
		if IsPeek(TabNilToken) {
			return Consume(TabNilToken)
		}
		if IsPeek(TabBooleanToken) {
			return Consume(TabBooleanToken)
		}
		return Tab{}, ParseError(fmt.Sprintf("atom not implemented %.0f / %s", ToNumber(ToDict(next)["kind"]), Print(next, true)))
	}

	var ConsumeList func(explicit bool) (Tab, error)

	ConsumeExpression := func() (Tab, error) {
		// fmt.Println("ConsumeExpression")
		if IsPeek(TabOpenToken) {
			// fmt.Println("Explicit List")
			o, err := Consume(TabOpenToken)
			if err != nil {
				return Tab{}, err
			}
			if IsPeek(TabCloseToken) {
				c, err := Consume(TabCloseToken)
				if err != nil {
					return Tab{}, err
				}
				result := FromList(TabList{})
				result.Position = &TabDict{
					"filename":  (*o.Position)["filename"],
					"startLine": (*o.Position)["startLine"],
					"startChar": (*o.Position)["startChar"],
					"endLine":   (*c.Position)["endLine"],
					"endChar":   (*c.Position)["endChar"],
				}
				return result, nil
			}
			var exp Tab
			if IsPeek(TabIndentToken) {
				_, err := Consume(TabIndentToken)
				if err != nil {
					return Tab{}, err
				}
				exp, err = ConsumeList(true)

				if err != nil {
					return Tab{}, err
				}
				_, err = Consume(TabDedentToken)
				if err != nil {
					return Tab{}, err
				}
			} else {
				var err error
				exp, err = ConsumeList(true)

				if err != nil {
					return Tab{}, err
				}
			}
			c, err := Consume(TabCloseToken)
			if err != nil {
				return Tab{}, err
			}
			exp.Position = &TabDict{
				"filename":  (*o.Position)["filename"],
				"startLine": (*o.Position)["startLine"],
				"startChar": (*o.Position)["startChar"],
				"endLine":   (*c.Position)["endLine"],
				"endChar":   (*c.Position)["endChar"],
			}
			return exp, nil
		}
		_atom, err := ConsumeAtom()
		if err != nil {
			return Tab{}, err
		}
		return _atom, nil
	}

	ConsumeInlineArgs := func() (args TabList, err error) {
		for !IsPeek(TabEolToken) && !IsPeek(TabCloseToken) && !IsPeek(TabEofToken) && !IsPeek(TabDedentToken) && !IsPeek(TabIndentToken) {
			var _arg Tab
			_arg, err = ConsumeExpression()
			if err != nil {
				return
			}
			args = append(args, _arg)
		}
		return
	}

	ConsumeIndentArgs := func() (args TabList, err error) {
		if IsPeek(TabIndentToken) {
			_, err = Consume(TabIndentToken)
			if err != nil {
				return
			}
			for {
				if IsPeek(TabEolToken) {
					_, err = Consume(TabEolToken)
					if err != nil {
						return
					}
					continue
				}
				if IsPeek(TabDedentToken) {
					_, err = Consume(TabDedentToken)
					if err != nil {
						return
					}
					break
				}
				if IsPeek(TabEofToken) {
					break
				}
				var _arg Tab
				_arg, err = ConsumeList(false)
				if err != nil {
					return
				}
				args = append(args, _arg)
			}
		}
		return args, nil
	}

	ConsumeList = func(explicit bool) (Tab, error) {
		// fmt.Printf("ConsumeList %s\n", Print(Tab{Dict: next.Position}))
		first, err := ConsumeExpression()
		if err != nil {
			return Tab{}, err
		}
		inlineArgs, err := ConsumeInlineArgs()
		if err != nil {
			return Tab{}, err
		}
		indentArgs, err := ConsumeIndentArgs()
		if err != nil {
			return Tab{}, err
		}
		// fmt.Printf("Count inline and indent args: %d %d\n", len(inlineArgs), len(indentArgs))
		args := append(inlineArgs, indentArgs...)
		// fmt.Println("Count total args:", len(args))
		if len(args) == 0 && !explicit {
			// fmt.Println("Returning first only")
			return first, nil
		}
		list := append(TabList{first}, args...)
		// fmt.Println("List length:", len(list))
		var endLine Tab
		var endChar Tab
		if len(args) > 0 {
			endLine = (*args[len(args)-1].Position)["endLine"]
			endChar = (*args[len(args)-1].Position)["endChar"]
		} else {
			endLine = (*first.Position)["endLine"]
			endChar = (*first.Position)["endChar"]
		}
		result := FromList(list)
		result.Position = &TabDict{
			"filename":  (*first.Position)["filename"],
			"startLine": (*first.Position)["startLine"],
			"startChar": (*first.Position)["startChar"],
			"endLine":   endLine,
			"endChar":   endChar,
		}
		return result, nil
	}

	ConsumeLines := func() (Tab, error) {
		var expressions TabList
		for !IsPeek(TabEofToken) {
			if IsPeek(TabEolToken) {
				_, err := Consume(TabEolToken)
				if err != nil {
					return Tab{}, err
				}
				continue
			}
			if IsPeek(TabCommentToken) {
				_, err := Consume(TabCommentToken)
				if err != nil {
					return Tab{}, err
				}
				continue
			}
			exp, err := ConsumeList(false)
			if err != nil {
				return Tab{}, err
			}
			expressions = append(expressions, exp)
			for IsPeek(TabEolToken) {
				_, err := Consume(TabEolToken)
				if err != nil {
					return Tab{}, err
				}
			}
		}
		if len(expressions) == 0 {
			return Tab{}, nil
		}
		if len(expressions) == 1 {
			return expressions[0], nil
		}
		// wrap multiple lines in do
		expressions = append(TabList{FromSymbol("do")}, expressions...)
		result := FromList(expressions)
		result.Position = &TabDict{
			"filename":  (*expressions[1].Position)["filename"],
			"startLine": (*expressions[1].Position)["startLine"],
			"startChar": (*expressions[1].Position)["startChar"],
			"endLine":   (*expressions[len(expressions)-1].Position)["endLine"],
			"endChar":   (*expressions[len(expressions)-1].Position)["endChar"],
		}
		return result, nil
	}
	return ConsumeLines()
}
