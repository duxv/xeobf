package main

import (
	"fmt"
	"io"
	"unicode"
)

/*
TODO: add inFunction property to Obfuscator
	for knowing whether we are in a function
	or not

TODO: add way to delete elements of the
	obfuscated map when went out of the
	specific scope
*/

func (o *Obfuscator) read() rune {
	ch, _, err := o.reader.ReadRune()

	switch err {
	case nil:
	case io.EOF:
		o.eof = true
	default:
		o.eof = true
		o.errorMessage = err.Error()
	}
	return ch
}

func (o *Obfuscator) unread() {
	_ = o.reader.UnreadRune()
	// if err != nil {
	// 	// if any error occurred, means that there's been
	// 	// something wrong with the written code
	// 	//panic(err)
	// }
}

func (o *Obfuscator) peek() rune {
	ch := o.read()
	if !o.eof {
		o.unread()
	}
	return ch
}

func (o *Obfuscator) readKeywordOrIdent() string {
	var ident string
	for ch := o.read(); unicode.IsLetter(ch) || unicode.IsDigit(ch); ch = o.read() {
		ident += string(ch)
	}
	o.unread()
	return ident
}

func isNumberPart(c rune) bool {
	switch {
	case unicode.IsDigit(c):
	case c == '.':
	case c == 'x' || c == 'a' || c == 'b' ||
		c == 'c' || c == 'd' || c == 'e' ||
		c == 'f':
	default:
		return false
	}

	return true
}

// Only working with int for now.
func (o *Obfuscator) readNumber() string {
	var num string

	for ch := o.read(); isNumberPart(ch); ch = o.read() {
		num += string(ch)
	}
	return num
}

// reads a character from the reader
// not to be mixed up with read method
func (o *Obfuscator) readChar() string {
	var lit string
	for ch := o.read(); ch != '\''; ch = o.read() {
		lit += string(ch)
	}
	o.unread()
	return lit
}

func (o *Obfuscator) readString() string {
	var lit string
	for ch := o.read(); ch != '"' && !o.eof; ch = o.read() {
		lit += string(ch)
	}
	return lit
}

var lastch rune

func (o *Obfuscator) nextToken() (tok Token) {
	var ch0, ch1 rune
	if lastch == 0 {
		ch0, ch1 = o.read(), o.peek()
	} else {
		ch0 = lastch
		ch1 = o.peek()
		lastch = 0
	}
	conc := string(ch0) + string(ch1)

	switch {
	case unicode.IsSpace(ch0):
		lit := string(ch0)
		var ch rune
		for ch = o.read(); unicode.IsSpace(ch); ch = o.read() {
			lit += string(ch)
		}
		lastch = ch // do not throw the ch away, keep it for the next iteration
		tok.Literal = lit
		tok.Type = Whitespace
	case isOperator(conc) || conc == "..":
		tok.Type = Operator
		_ = o.read()
		x := conc + string(o.peek())
		if isOperator(x) {
			tok.Literal = x
			_ = o.read()
		} else {
			tok.Literal = conc
		}
	case isOperator(string(ch0)):
		tok.Type = Operator
		tok.Literal = string(ch0)
	case ch0 == '\'':
		tok.Type = Char
		tok.Literal = o.readChar()
	case ch0 == '"':
		tok.Type = String
		tok.Literal = o.readString()
	default:
		if unicode.IsLetter(ch0) || ch0 == '_' {
			tok.Literal = string(ch0) + o.readKeywordOrIdent()
			if isKeyword(tok.Literal) {
				tok.Type = Keyword
			} else {
				tok.Type = Ident
			}
			return tok
		} else if unicode.IsDigit(ch0) {
			tok.Literal = string(ch0) + o.readNumber()
			tok.Type = Number
			return tok
		} else {
			fmt.Println(string(ch0))
			panic("Something went wrong in the next token func... Please contact a developer with this error message.")
		}
	}
	return tok
}

// List of every delimiter/operator
var operators = []string{
	"+", "-", "*", "/", "%",
	"&", "|", "^", "<<", ">>",
	"&^", "+=", "-=", "*=",
	"/=", "%=", "&=", "|=",
	"^=", "<<=", ">>=", "&^=",
	"&&", "||", "<-", "++", "--",
	"==", "<", ">", "=", "!",
	"!=", "<=", ">=", ":=",
	"...", "(", "[", "{", ",",
	".", ")", "]", "}", ";", ":",
}

// Check if a string is whether an operator
// or a delimiter
func isOperator(s string) bool {
	for _, op := range operators {
		if op == s {
			return true
		}
	}
	return false
}

// List of every keyword (Some may be missing though.)
var keywords = []string{
	"break", "case", "chan", "const",
	"continue", "default", "defer",
	"else", "fallthrough", "for",
	"func", "go", "goto", "if", "import",
	"interface", "map", "package", "range",
	"return", "select", "struct", "switch",
	"type", "var",

	// Data types related keywords
	"int8", "int16", "int32",
	"int64", "uint8", "uint16", "uint32", "uint64",
	"int", "uint", "rune", "byte", "uintptr", "float32",
	"float64", "complex64", "complex128", "false", "true",
	"bool", "string",

	// Built-in functions are counted as keywords
	"append", "cap", "close", "delete", "copy", "imag",
	"len", "make", "new", "panic", "print", "println",
	"real",
}

func isKeyword(s string) bool {
	for _, key := range keywords {
		if key == s {
			return true
		}
	}
	return false
}
