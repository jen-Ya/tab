package tabgo

import (
	"fmt"
	"strconv"
	"strings"
)

type TabToken int

const (
	TabOpenToken    TabToken = iota
	TabCloseToken   TabToken = iota
	TabEolToken     TabToken = iota
	TabIndentToken  TabToken = iota
	TabDedentToken  TabToken = iota
	TabCommentToken TabToken = iota
	TabStringToken  TabToken = iota
	TabNumberToken  TabToken = iota
	TabEofToken     TabToken = iota
	TabSymbolToken  TabToken = iota
	TabNilToken     TabToken = iota
	TabBooleanToken TabToken = iota
)

func (T TabToken) String() string {
	switch T {
	case TabOpenToken:
		return "TabOpenToken"
	case TabCloseToken:
		return "TabCloseToken"
	case TabEolToken:
		return "TabEolToken"
	case TabIndentToken:
		return "TabIndentToken"
	case TabDedentToken:
		return "TabDedentToken"
	case TabCommentToken:
		return "TabCommentToken"
	case TabStringToken:
		return "TabStringToken"
	case TabNumberToken:
		return "TabNumberToken"
	case TabEofToken:
		return "TabEofToken"
	case TabSymbolToken:
		return "TabSymbolToken"
	case TabNilToken:
		return "TabNilToken"
	case TabBooleanToken:
		return "TabBooleanToken"
	}
	return "Unknown"
}

func IsDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

func CanBeSymbol(char rune) bool {
	// check if char is not in '\n#\'"() \t'
	return char != '\n' && char != '#' && char != '\'' && char != '"' && char != '(' && char != ')' && char != ' ' && char != '\t'
}

func Unescape(str string) string {
	// unescape newlines and tabs
	str = strings.Replace(str, "\\n", "\n", -1)
	str = strings.Replace(str, "\\t", "\t", -1)
	return str
}

func MakeToken(kind TabToken, value Tab, filename string, startLine int, startChar int, endLine int, endChar int) Tab {
	result := FromDict(TabDict{
		"kind":  FromNumber(float64(kind)),
		"value": value,
	})
	result.Position = &TabDict{
		"filename":  FromString(filename),
		"startLine": FromNumber(float64(startLine)),
		"startChar": FromNumber(float64(startChar)),
		"endLine":   FromNumber(float64(endLine)),
		"endChar":   FromNumber(float64(endChar)),
	}
	return result
}

func Tokenize(text string, keepKomments bool, filename string) (TabTokens Tab, err error) {
	var tokens TabList
	text = text + "\n"
	runes := []rune(text)
	cursor := 0
	linecount := 0
	charcount := 0
	indentation := 0
	startLine := 0
	startChar := 0
	textLength := len(runes)
	TokenizeError := func(message string) error {
		return fmt.Errorf(
			"TokenizeError at %s:%d:%d\nToken Started at %s:%d:%d\n%s",
			filename, linecount+1, charcount+1,
			filename, startLine+1, startChar+1,
			message,
		)
	}
	IncCursor := func() error {
		if cursor >= textLength {
			return TokenizeError("Unexpected end of file")
		}
		if runes[cursor] == '\n' {
			linecount++
			charcount = 0
		} else {
			charcount++
		}
		cursor++
		return nil
	}
	AddToken := func(kind TabToken, value Tab) {
		tokens = append(tokens, MakeToken(kind, value, filename, startLine, startChar, linecount, charcount))
		startLine = linecount
		startChar = charcount
	}
	Consume := func(chars string) error {
		for _, char := range chars {
			if cursor >= textLength {
				return TokenizeError("Unexpected end of file consuming " + chars)
			}
			if runes[cursor] != char {
				return TokenizeError(fmt.Sprintf("Unexpected character %s, expected %s", runes[cursor], char))
			}
			if err := IncCursor(); err != nil {
				return err
			}
		}
		return nil
	}
	IsPeek := func(chars string) bool {
		for offset, char := range chars {
			if (cursor+offset >= textLength) || runes[cursor+offset] != char {
				return false
			}
		}
		return true
	}
	IsPeekWord := func(chars string) bool {
		return IsPeek(chars) && (cursor+len(chars) >= textLength ||
			!CanBeSymbol(runes[cursor+len(chars)]))
	}
	ConsumeIndent := func(indent int) error {
		// fmt.Println("ConsumeIndent")
		for i := 0; i < indent; i++ {
			if err := Consume("\t"); err != nil {
				return err
			}
		}
		return nil
	}
	GetIndent := func() int {
		indent := 0
		for ; indent < textLength-cursor; indent++ {
			if runes[cursor+indent] != '\t' {
				return indent
			}
		}
		return indent
	}
	ConsumeEol := func() error {
		// fmt.Println("ConsumeEol")
		for cursor < textLength && runes[cursor] == '\n' {
			if err := Consume("\n"); err != nil {
				return err
			}
		}
		indent := GetIndent()
		if indent == indentation {
			if err := ConsumeIndent(indent); err != nil {
				return err
			}
			AddToken(TabEolToken, Tab{})
		} else if indent < indentation {
			if err := ConsumeIndent(indent); err != nil {
				return err
			}
			for i := indent; i < indentation; i++ {
				AddToken(TabDedentToken, Tab{})
			}
			indentation = indent
		} else if indent == indentation+1 {
			if err := ConsumeIndent(indent); err != nil {
				return err
			}
			AddToken(TabIndentToken, Tab{})
			indentation = indent
		} else {
			return TokenizeError(fmt.Sprintf("IndentError: unexpected indentation of %d from %d", indent, indentation))
		}
		return nil
	}

	ConsumeMultiline := func(separator string) (string, error) {
		// fmt.Println("ConsumeMultiline")
		if err := Consume(separator); err != nil {
			return "", err
		}
		if err := Consume("\n"); err != nil {
			return "", err
		}
		value := ""
		if err := ConsumeIndent(indentation + 1); err != nil {
			return "", err
		}
		for {
			for runes[cursor] != '\n' {
				value += string(runes[cursor])
				if err := IncCursor(); err != nil {
					return "", err
				}
			}
			if err := Consume("\n"); err != nil {
				return "", err
			}
			for runes[cursor] == '\n' {
				if err := Consume("\n"); err != nil {
					return "", err
				}
				value += "\n"
			}
			indent := GetIndent()
			if indent == indentation {
				if err := ConsumeIndent(indentation); err != nil {
					return "", err
				}
				if err := Consume(separator); err != nil {
					return "", err
				}
				return Unescape(value), nil
			}
			value += "\n"
			if err := ConsumeIndent(indentation + 1); err != nil {
				return "", err
			}
		}
	}

	ConsumeToEol := func() (string, error) {
		// fmt.Println("ConsumeToEol")
		value := ""
		for runes[cursor] != '\n' {
			value += string(runes[cursor])
			if err := IncCursor(); err != nil {
				return "", err
			}
		}
		return value, nil
	}

	ConsumeComment := func() error {
		// fmt.Println("ConsumeComments")
		if runes[cursor+1] == '\n' {
			value, err := ConsumeMultiline("#")
			if err != nil {
				return err
			}
			if keepKomments {
				AddToken(TabCommentToken, FromString(value))
			}
		} else {
			if err := Consume("#"); err != nil {
				return err
			}
			if runes[cursor] == ' ' {
				if err := Consume(" "); err != nil {
					return err
				}
			}
			value, err := ConsumeToEol()
			if err != nil {
				return err
			}
			if keepKomments {
				AddToken(TabCommentToken, FromString(value))
			}
		}
		return nil
	}

	ConsumeString := func(quote string) error {
		// fmt.Println("ConsumeString")
		if runes[cursor+1] == '\n' {
			value, err := ConsumeMultiline(quote)
			if err != nil {
				return err
			}
			AddToken(TabStringToken, FromString(value))
		} else {
			if err := Consume(quote); err != nil {
				return err
			}
			value := ""
			for !IsPeek(quote) {
				if IsPeek("\n") {
					return TokenizeError(fmt.Sprintf("Unterminated single line string %s", "\""))
				}
				value += string(runes[cursor])
				if err := IncCursor(); err != nil {
					return err
				}
			}
			if err := Consume(quote); err != nil {
				return err
			}
			// Should we unescape escaped quotes in single line strings?
			// e.g. "foo\"bar" -> foo"bar?
			value = Unescape(value)
			AddToken(TabStringToken, FromString(value))
		}
		return nil
	}

	ConsumeNumber := func() error {
		// fmt.Println("ConsumeNumber")
		value := ""
		if IsPeek("-") {
			value += string(runes[cursor])
			if err := Consume("-"); err != nil {
				return err
			}
		}
		for IsDigit(runes[cursor]) {
			value += string(runes[cursor])
			if err := IncCursor(); err != nil {
				return err
			}
		}
		if IsPeek(".") {
			value += string(runes[cursor])
			if err := Consume("."); err != nil {
				return err
			}
		}
		for IsDigit(runes[cursor]) {
			value += string(runes[cursor])
			if err := IncCursor(); err != nil {
				return err
			}
		}
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return TokenizeError(fmt.Sprintf("Could not parse number %s", value))
		}
		AddToken(TabNumberToken, FromNumber(num))
		return nil
	}

	ConsumeEof := func() {
		// fmt.Println("ConsumeEof")
		AddToken(TabEofToken, Tab{})
	}

	ConsumeSymbol := func() error {
		// fmt.Println("ConsumeSymbol")
		value := ""
		for CanBeSymbol(runes[cursor]) {
			value += string(runes[cursor])
			if err := IncCursor(); err != nil {
				return err
			}
		}
		if value == "" {
			return TokenizeError("Symbol has no value")
		}
		AddToken(TabSymbolToken, FromSymbol(value))
		return nil
	}

	ConsumeOpen := func() error {
		// fmt.Println("ConsumeOpen")
		if err := Consume("("); err != nil {
			return err
		}
		AddToken(TabOpenToken, Tab{})
		return nil
	}

	ConsumeClose := func() error {
		// fmt.Println("ConsumeClose")
		if err := Consume(")"); err != nil {
			return err
		}
		AddToken(TabCloseToken, Tab{})
		return nil
	}

	for {
		if cursor >= textLength-1 {
			ConsumeEof()
			return FromList(tokens), nil
		} else if IsPeek(" ") {
			if err = Consume(" "); err != nil {
				return
			}
		} else if IsPeek("#") {
			if err = ConsumeComment(); err != nil {
				return
			}
		} else if IsPeek("'") || IsPeek("\"") {
			if err = ConsumeString(string(runes[cursor])); err != nil {
				return
			}
		} else if IsDigit(runes[cursor]) ||
			(IsPeek("-") && (cursor < textLength-1) && IsDigit(runes[cursor+1])) {
			if err = ConsumeNumber(); err != nil {
				return
			}
		} else if IsPeek("(") {
			if err = ConsumeOpen(); err != nil {
				return
			}
		} else if IsPeek(")") {
			if err = ConsumeClose(); err != nil {
				return
			}
		} else if IsPeek("\n") {
			if err = ConsumeEol(); err != nil {
				return
			}
		} else if IsPeekWord("nil") {
			if err = Consume("nil"); err != nil {
				return
			}
			AddToken(TabNilToken, Tab{})
		} else if IsPeekWord("_") {
			if err = Consume("_"); err != nil {
				return
			}
			AddToken(TabNilToken, Tab{})
		} else if IsPeekWord("true") {
			if err = Consume("true"); err != nil {
				return
			}
			AddToken(TabBooleanToken, FromBool(true))
		} else if IsPeekWord("false") {
			if err = Consume("false"); err != nil {
				return
			}
			AddToken(TabBooleanToken, FromBool(false))
		} else {
			if err = ConsumeSymbol(); err != nil {
				return
			}
		}
	}
}
