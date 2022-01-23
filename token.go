package main

type TokenType int

const (
	Char TokenType = iota
	Ident
	Int
	Keyword
	Operator
	String
	Whitespace
)

type Token struct {
	Type    TokenType
	Literal string
}

var tokNames = map[TokenType]string{
	Char:       "Char",
	Ident:      "Ident",
	Int:        "Int",
	Keyword:    "Keyword",
	Operator:   "Operator",
	String:     "String",
	Whitespace: "Whitespace",
}

func tokenToString(t TokenType) string {
	return tokNames[t]
}
