package main

import (
	"bytes"
	"log"
	"math/big"
)

const eof = 0

type exprLex struct {
	input *bytes.Buffer
}

func newExprLex(in []byte) *exprLex {
	return &exprLex{
		input: bytes.NewBuffer(in),
	}
}

func (x *exprLex) Lex(yylval *exprSymType) int {
	for {
		c, _, err := x.input.ReadRune()
		if err != nil {
			return eof
		}
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return x.num(c, yylval)
		case '+', '-', '*', '/', '(', ')':
			return int(c)
		case '×':
			return '*'
		case '÷':
			return '/'

		case ' ', '\t', '\n', '\r':
		default:
			log.Printf("unrecognized character %q", c)
		}
	}
}

// Lex a number.
func (x *exprLex) num(c rune, yylval *exprSymType) int {
	var b bytes.Buffer
	b.WriteRune(c)
L:
	for {
		c, _, err := x.input.ReadRune()
		if err != nil {
			return eof
		}
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', 'e', 'E':
			b.WriteRune(c)
		default:
			_ = x.input.UnreadRune()
			break L
		}
	}
	yylval.num = &big.Rat{}
	_, ok := yylval.num.SetString(b.String())
	if !ok {
		log.Printf("bad number %q", b.String())
		return eof
	}
	return NUM
}

func (x *exprLex) Error(s string) {
	log.Println("parse error: ", s)
}
