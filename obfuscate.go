package main

import (
	"bufio"
	"bytes"
	"math/rand"
	"strings"
	"time"
)

/*
	TODO: make a way to obfuscate numbers
	maybe converting them to 16 base.
*/

type Obfuscator struct {
	reader *bufio.Reader
	// Contains each identifier and its obfuscation
	obfuscated map[string]string
	// True if met the end of the reader
	// or if any error occurs
	eof bool
	// Message of the error, if any
	errorMessage string
	// This is a list of all of the imports,
	// since we cannot change the name of imports.
	imported []string
}

func NewObfuscator(src string) *Obfuscator {
	obf := &Obfuscator{
		reader:     bufio.NewReader(strings.NewReader(src)),
		obfuscated: make(map[string]string),
	}
	return obf
}

/*
	Char:       "Char",
	Ident:      "Ident",
	Int:        "Int",
	Keyword:    "Keyword",
	Operator:   "Operator",
	String:     "String",
	Whitespace: "Whitespace",
*/

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

// Add a string with a matching obfuscation in the
// obfuscated map, and also return the generated string
func (o *Obfuscator) addObfuscated(s string) (a string) {
	a = randString(12)
	o.obfuscated[s] = a
	return a
}

func (o *Obfuscator) addImported(s string) {
	o.imported = append(o.imported, s)
}

func (o *Obfuscator) isImported(s string) bool {
	for _, pack := range o.imported {
		if pack == s {
			return true
		}
	}
	return false
}

// Generate a buffer with the obfuscated source code
func (o *Obfuscator) Obfuscate() *bytes.Buffer {
	out := &bytes.Buffer{}
	var lastKeyword string // last keyword or ident

	for !o.eof {
		tok := o.nextToken()
		switch tok.Type {
		case Ident:
			// do not obfuscate package names
			// since we are not ready for this yet
			if lastKeyword == "package" {
				out.WriteString(tok.Literal)
				continue
			} else if o.isImported(lastKeyword) {
				out.WriteString(tok.Literal)
				continue
			}
			if obf, ok := o.obfuscated[tok.Literal]; ok {
				out.WriteString(obf)
			} else if tok.Literal != "main" && !o.isImported(tok.Literal) {
				obf := o.addObfuscated(tok.Literal)
				out.WriteString(obf)
			} else {
				out.WriteString(tok.Literal)
			}

			/*
				TODO: make a way not to obfuscate imported functions
				XXX already implemented something for this, doesn't work for
				more than one func.
			*/
		case Keyword:
			if tok.Literal == "import" {
				out.WriteString("import")
				// skip until the first (
				for tok = o.nextToken(); tok.Literal != "("; tok = o.nextToken() {
					out.WriteString(tok.Literal)
				}
				out.WriteRune('(')
				// The token is now (
				for tok = o.nextToken(); tok.Literal != ")"; tok = o.nextToken() {
					if tok.Type == Ident {
						// means it is an alias
						o.addImported(tok.Literal)
						out.WriteString(tok.Literal)
					} else {
						// Assumes it is a string
						split := strings.Split(tok.Literal, "/")
						pckgName := split[len(split)-1]
						o.addImported(pckgName)
						if tok.Type == String {
							out.WriteString("\"" + tok.Literal + "\"")
						} else {
							out.WriteString(tok.Literal)
						}
					}
				}
				out.WriteRune(')')
			} else {
				out.WriteString(tok.Literal)
			}
		case String:
			out.WriteString("\"" + tok.Literal + "\"")
		case Char:
			out.WriteString("'" + tok.Literal + "'")
		default:
			out.WriteString(tok.Literal)
		}
		if tok.Type == Keyword || tok.Type == Ident {
			lastKeyword = tok.Literal
		}
	}
	return out
}
