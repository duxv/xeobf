package main

type TokenType int

const (
	Char TokenType = iota
	Ident
	Number
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
	Number:     "Number",
	Keyword:    "Keyword",
	Operator:   "Operator",
	String:     "String",
	Whitespace: "Whitespace",
}

func tokenToString(t TokenType) string {
	return tokNames[t]
}
