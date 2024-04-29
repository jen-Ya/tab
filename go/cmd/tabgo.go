package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	t "jen-ya.de/tabgo"
)

const initScript = `eval
	read-string
		file-read (str-join (list .tabhome '/tabgo/goinit.tab') '')
		dict
			'filename'
			str-join (list .tabhome '/tabgo/goinit.tab') ''
`

func EnvDefault() t.Tab {
	env := t.Env(t.Tab{})
	tru := t.SymbolToTab("tru")
	t.EnvSet(env, tru, tru)
	t.EnvSet(env, t.SymbolToTab(".env-root"), env)
	t.EnvSet(env, t.SymbolToTab(".tabhome"), t.StringToTab(os.Getenv("TABHOME")))
	t.AddCore(env)
	t.Eval(t.Read(initScript, false, "init"), env)
	return env
}

func ClearRepl() {
	fmt.Print("\033[H\033[2J")
}

func Rep(env t.Tab, reader *bufio.Reader) error {
	defer func() {
		if r := recover(); r != nil {
			if tab, ok := r.(t.Tab); ok {
				fmt.Printf("error: %s at %s\n", t.Print(tab, false), t.AstPositionToString(tab))
				return
			}
			fmt.Println("error:", r)
		}
	}()
	fmt.Print("tabgo: ")
	input, err := reader.ReadString('\n')
	switch input {
	case ".clear\n":
		ClearRepl()
		return nil
	}
	if err != nil {
		return err
	}
	ast := t.Read(input, false, "repl")
	if err != nil {
		return err
	}
	result := t.Eval(ast, env)
	fmt.Println(t.Print(result, false))
	return nil
}

func RunFile(filepath string) {
	defer func() {
		if r := recover(); r != nil {
			if tab, ok := r.(t.Tab); ok {
				fmt.Printf("error: %s at %s\n", t.Print(tab, false), t.AstPositionToString(tab))
			}
			panic(r)
		}
	}()
	bytes, error := os.ReadFile(filepath)
	if error != nil {
		panic(error)
	}
	t.Eval(t.Read(string(bytes), false, filepath), EnvDefault())
}

func RunRepl() {
	defer func() {
		if r := recover(); r != nil {
			if tab, ok := r.(t.Tab); ok {
				fmt.Printf("error: %s at %s\n", t.Print(tab, false), t.AstPositionToString(tab))
				return
			}
			fmt.Println("error:", r)
		}
	}()
	reader := bufio.NewReader(os.Stdin)
	env := EnvDefault()
	for {
		if err := Rep(env, reader); err != nil {
			if err == io.EOF {
				fmt.Println("")
				break
			}
			panic(err)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		RunRepl()
		return
	}
	script := os.Args[1]
	RunFile(script)
}
