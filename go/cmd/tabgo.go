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
		file-read (str-join (li .tabhome '/tabgo/goinit.tab') '')
		dict
			'filename'
			str-join (li .tabhome '/tabgo/goinit.tab') ''
`

func EnvDefault() t.Tab {
	env := t.Env(t.TabNil)
	// Tru is a symbol that evaluates to itself
	// Actually this is not needed, because tabgo has booleans
	// But I saw this in other lisps and wanted to check it out
	tru := t.FromSymbol("tru")
	t.EnvSet(env, tru, tru)
	// Save the environment as .env-root, so that it can be accessed later
	t.EnvSet(env, t.FromSymbol(".env-root"), env)
	// Set the TABHOME environment variable for require
	t.EnvSet(env, t.FromSymbol(".tabhome"), t.FromString(os.Getenv("TABHOME")))
	// Set command line arguments
	var args t.TabList
	for _, arg := range os.Args[1:] {
		args = append(args, t.FromString(arg))
	}
	t.EnvSet(env, t.FromSymbol(".argv"), t.FromList(args))
	// Add core functions
	t.AddCore(env)
	// Run init script
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
