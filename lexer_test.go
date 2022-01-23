package main

import (
	"strings"
	"testing"
)

func TestReadChar(t *testing.T) {
	chars := "abcdefghijklmnopqrstuvwxyz12345678908**~!"

	for _, char := range chars {
		schar := "'" + string(char) + "'"
		obf := NewObfuscator(schar)
		obf.read()
		x := obf.readChar()
		if x != string(schar[1]) {
			t.Fatalf("Got %q, expected %q", x, schar[1])
		}
	}
}

func TestReadString(t *testing.T) {
	totest := `testing this onehundered333 timesaday`
	split := strings.Split(totest, " ")

	for _, str := range split {
		this := "\"" + str + "\""
		obf := NewObfuscator(this)
		obf.read()
		x := obf.readString()
		expected := this[1 : len(this)-1]
		if x != expected {
			t.Fatalf("Got %q, expected %q", x, expected)
		}
	}
}

func TestReadNumber(t *testing.T) {
	totest := []string{
		"0", "1000", "393", "842",
		"432898423",
	}
	for _, str := range totest {
		obf := NewObfuscator(str)
		num := obf.readNumber()
		if num != str {
			t.Fatalf("Got %q, expected %q", num, str)
		}
	}
}

func TestReadKeywordOrIdent(t *testing.T) {
	totest := []string{
		"break", "package", "main", "etc",
	}
	for _, str := range totest {
		obf := NewObfuscator(str)
		ident := obf.readKeywordOrIdent()
		if ident != str {
			t.Fatalf("Got %q, expected %q", ident, str)
		}
	}
}

func TestNextToken(t *testing.T) {
	type testCase struct {
		input   string
		typ     TokenType
		correct string
	}
	totest := []testCase{
		{"break", Keyword, "break"},
		{"43289", Int, "43289"},
		{"\"STRINGY\"", String, "STRINGY"},
		{"'c'", Char, "c"},
		{"+", Operator, "+"},
		{"			   ", Whitespace, "			   "},
	}
	for _, key := range keywords {
		totest = append(totest, testCase{
			input:   key,
			typ:     Keyword,
			correct: key,
		})
	}
	for _, op := range operators {
		totest = append(totest, testCase{
			input:   op,
			typ:     Operator,
			correct: op,
		})
	}
	for idx, test := range totest {
		obf := NewObfuscator(test.input)
		tok := obf.nextToken()
		if tok.Literal != test.correct {
			t.Fatalf("[%d] Expected literal %q, got %q", idx, test.correct, tok.Literal)
		}
		if tok.Type != test.typ {
			t.Fatalf("[%d] Expected type %s, got %s", idx, tokenToString(test.typ), tokenToString(tok.Type))
		}

	}
}
