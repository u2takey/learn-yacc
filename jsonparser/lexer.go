package jsonparser

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// Parse parses the input and returs the result.
func Parse(input []byte) (map[string]interface{}, error) {
	l := newLex(input)
	_ = yyParse(l)
	return l.result, l.err
}

type lex struct {
	input  *bytes.Buffer
	result map[string]interface{}
	err    error
}

func newLex(input []byte) *lex {
	return &lex{
		input: bytes.NewBuffer(input),
	}
}

const eof = 0

// Lex satisfies yyLexer.
func (l *lex) Lex(lval *yySymType) int {
	for {
		r, _, err := l.input.ReadRune()
		if err != nil {
			return eof
		}
		switch {
		case unicode.IsSpace(r):
			continue
		case r == '"':
			return l.scanString(lval)
		case unicode.IsDigit(r) || r == '+' || r == '-':
			_ = l.input.UnreadRune()
			return l.scanNum(lval)
		case unicode.IsLetter(r):
			_ = l.input.UnreadRune()
			return l.scanLiteral(lval)
		default:
			return int(r)
		}
	}
}

var escape = map[byte]byte{
	'"':  '"',
	'\\': '\\',
	'/':  '/',
	'b':  '\b',
	'f':  '\f',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
}

func (l *lex) scanString(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	for {
		r, _, err := l.input.ReadRune()
		if err != nil {
			break
		}
		switch r {
		case '\\':
			// TODO(sougou): handle \uxxxx construct.
			b2 := escape[byte(r)]
			if b2 == 0 {
				return LexError
			}
			buf.WriteByte(b2)
		case '"':
			lval.val = buf.String()
			return String
		default:
			buf.WriteRune(r)
		}
	}
	return LexError
}

func (l *lex) scanNum(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	for {
		r, _, err := l.input.ReadRune()
		if err != nil {
			break
		}
		switch {
		case unicode.IsDigit(r):
			buf.WriteRune(r)
		case strings.IndexRune(".+-eE", r) != -1:
			buf.WriteRune(r)
		default:
			_ = l.input.UnreadRune()
			val, err := strconv.ParseFloat(buf.String(), 64)
			if err != nil {
				return LexError
			}
			lval.val = val
			return Number
		}
	}
	return LexError
}

var literal = map[string]interface{}{
	"true":  true,
	"false": false,
	"null":  nil,
}

func (l *lex) scanLiteral(lval *yySymType) int {
	buf := bytes.NewBuffer(nil)
	for {
		r, _, err := l.input.ReadRune()
		if err != nil {
			break
		}
		switch {
		case unicode.IsLetter(r):
			buf.WriteRune(r)
		default:
			_ = l.input.UnreadRune()
			val, ok := literal[buf.String()]
			if !ok {
				return LexError
			}
			lval.val = val
			return Literal
		}
	}
	return LexError
}

// Error satisfies yyLexer.
func (l *lex) Error(s string) {
	l.err = errors.New(s)
}
