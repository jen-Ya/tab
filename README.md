# tab

**This is my personal toy language, used for exploring concepts of programming language design.**

It is an indentation-based programming language characterized by minimal syntax and maximum flexibility.


Implementations are currently available in Golang and JavaScript.

## Install

```sh
git clone https://github.com/jen-Ya/tab
cd tab
export TABHOME=`pwd`/tab

# Install Golang version (Golang must be already installed):
./bin/tabgo-install
# Run Golang tests
tabgo ./tab/tabgo/test-all.tab

# Run JavaScipt tests (NodeJS must be already installed):
node ./bin/tab ./tab/test-all.tab
```

You can optionally add the following line to your shell startup script (such as .bash_profile or .zshrc). Remember to substitute /path/to/tab with the actual path, paying attention to the double /tab/tab in TABHOME.

```sh
export TABHOME="/path/to/tab/tab"
export PATH="/path/to/tab/bin:${PATH}"
```

## Examples

### Comments

```
# Single line comment

#
	multiline
	comment
#

#
	# nested comments are okay #
	this # is also okay, multiline-comments are tokenized based on indentation
#

#
	this
is not allowed!
#
```

### Numbers and Strings

```
# Numbers always start with a digit or a minus and have arbitrary length and precision

2.2123123
2
0.1
-0.4

# Strings are either single or multi-line and are enclosed with ' or "

'string'
"string\nwith\nlinebreaks"

# multiline strings are stripped based on indentation

"
	multi
	line
		string
"

# is equal to "multiline\nline\n\tstring", with only one tab in the last line
```

### Indentation, parentheses

```
# A single identifier is treated as its value:

hello

# An identifier followed by something in the same line is treated as an lisp-style s-expression

print "hello"

# Indented arguments and parenthesis

print
	"the answer to everyhing is"
	* 21 2

print "the answer to everything is" (* 21 2)

```


### Fibonacci
```
fn fibo x
	if (< x 2) 1
		+
			fibo (- x 1)
			fibo (- x 2)
```

## Golang implementation

The source is in the `./go` directory.

The Implementation supports extensions written in go. You can find examples in `./tab/tabgo/plugins`.

## NodeJS implementation

The source is in the `./js` directory.

## Discussion

- Types
	- Null values
- Pattern Matching
- Host languages, backends, targets
- Currying in builtins
- Minial set of builtins
- Tooling
- Default function parameters
- Variadic functions, spreading
- Nested quasiquote
	- [How are nested quasiquotes implemented? - Reddit](https://www.reddit.com/r/lisp/comments/as0ch1/comment/egrcqdm/)
	- [Quasiquotation in Lisp](https://3e8.org/pub/scheme/doc/Quasiquotation%20in%20Lisp%20(Bawden).pdf)

### Inline comments?

might look like this, but it might not always be clear where the comment starts and ends

```
# inline comment # print "test"

#( you could indicate start and end with parens, but should the tokenizer be aware of it? Thinking of documentation generation etc. )#
```

## JS implementation

First implementation based on [mal](https://github.com/kanaka/mal).

### TABHOME env variable

For file imports to work, TABHOME environment variable must be set and point to this projects' tab **subfolder**.

Put that in your .bash_profile or whatever:

`export TABHOME='/path/to/tab/tab'`

To make tab command available everywhere in your shell add

`export PATH="path/to/tab/bin:${PATH}"`

try ```node bin/tab ./tab/hello.tab```


## Resources

### Lisps

- [Cirru](http://cirru.org/) - similar language
- [Slick](https://github.com/pcostanza/slick) - Lisp/Scheme-style s-expression surface syntax for the Go
- [awesome-lisp-languages](https://github.com/dundalek/awesome-lisp-languages)
- [Make a Lisp - mal](https://github.com/kanaka/mal)
- [Sweet-expressions (t-expressions)](https://srfi.schemers.org/srfi-110/srfi-110.html)
- [PicoLisp](https://picolisp.com)
- [Tinylisp - Lisp in 99 lines of C](https://github.com/Robert-van-Engelen/tinylisp)
- [Lisp in 1k lines of C, explained](https://github.com/Robert-van-Engelen/lisp)

### Language Design

- [Syntax Design](https://cs.lmu.edu/~ray/notes/syntaxdesign/)

### Logical Programming
- [How big an undertaking would it be to implement a Prolog interpreter?](https://news.ycombinator.com/item?id=2152964)
- [miniKanren](http://minikanren.org/) is a family of Domain Specific Languages for logic programming
	- [µKanren](http://webyrd.net/scheme-2013/papers/HemannMuKanren2013.pdf) : A Minimal Functional Core for Relational Programming
- [Paradigms of AI Programming, Chapter 11 - Logic Programming](https://github.com/norvig/paip-lisp/blob/main/docs/chapter11.md) - Simple implementaion in Common Lisp

### Compilers

- [Compiler Optimizations](https://predr.ag/blog/compiler-adventures-part1-no-op-instructions/)
- [Create your own compiler](https://citw.dev/tutorial/create-your-own-compiler)
- [A Compiler Writing Journey](https://github.com/DoctorWkt/acwj)

### LLVM

- [My First Language Frontend with LLVM Tutorial](https://llvm.org/docs/tutorial/MyFirstLanguageFrontend/index.html)
- [node-llvm](https://github.com/kevinmehall/node-llvm)
- [llvm-node](https://github.com/MichaReiser/llvm-node)
- [lvm-hello-world-example](https://github.com/zilder/llvm-hello-world-example/blob/master/main.cpp)

### VSC integration

- [Semantic Highlighting Guide](https://code.visualstudio.com/api/language-extensions/semantic-highlight-guide)
- [Language Configuration Guide](https://code.visualstudio.com/api/language-extensions/language-configuration-guide)
- [Writing Your Own Debugger And Language Extensions](https://www.codemag.com/article/1809051/Writing-Your-Own-Debugger-and-Language-Extensions-with-Visual-Studio-Code)
- [Debugger extension](https://code.visualstudio.com/api/extension-guides/debugger-extension)
- [Attach to Debug Server](https://github.com/microsoft/vscode/issues/23518)

### Webpack integration
- [Writing a loader](https://webpack.js.org/contribute/writing-a-loader/)

### WebAssemly
- [Node.js with WebAssembly](https://nodejs.dev/en/learn/nodejs-with-webassembly/)
- [WABT: The WebAssembly Binary Toolkit](https://github.com/webassembly/wabt)
- [Understanding WebAssembly text format](https://developer.mozilla.org/en-US/docs/WebAssembly/Understanding_the_text_format)
- [Schism](https://github.com/schism-lang/schism)
- [Low Level Lisp for WebAssembly](https://github.com/FemtoEmacs/wasCm)
- [WebAssembly/Instructions](https://wiki.freepascal.org/WebAssembly/Instructions)

### Misc

- [JS AST explorer](https://astexplorer.net)
- [Babel AST spec](https://github.com/babel/babel/blob/main/packages/babel-parser/ast/spec.md)