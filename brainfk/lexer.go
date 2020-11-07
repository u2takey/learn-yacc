package main

import (
	"bytes"
	"errors"
)

// Parse parses the input and returs the result.
func Parse(input []byte) ([]Code, error) {
	l := newLex(input)
	_ = yyParse(l)
	return l.codes, l.err
}

type OpType int

type Code struct {
	opType   OpType
	argument int
}

const (
	INVALID_OP OpType = iota
	INC_PTR
	DEC_PTR
	INC_DATA
	DEC_DATA
	WRITE_STDOUT
	READ_STDIN
	JUMP_IF_DATA_ZERO
	JUMP_IF_DATA_NOT_ZERO
	LOOP_SET_TO_ZERO
	LOOP_MOVE_PTR
	LOOP_MOVE_DATA
)

var OpMap = map[string]OpType{
	">": INC_PTR,
	"<": DEC_PTR,
	"+": INC_DATA,
	"-": DEC_DATA,
	".": WRITE_STDOUT,
	",": READ_STDIN,
	"[": JUMP_IF_DATA_ZERO,
	"]": JUMP_IF_DATA_NOT_ZERO,
}

type lex struct {
	input *bytes.Buffer
	codes []Code
	err   error
}

func newLex(input []byte) *lex {
	return &lex{
		input: bytes.NewBuffer(input),
	}
}

// Error satisfies yyLexer.
func (l *lex) Error(s string) {
	l.err = errors.New(s)
}

// Lex satisfies yyLexer.
func (l *lex) Lex(lval *yySymType) int {
	for {
		b, err := l.input.ReadByte()
		if err == nil {
			if OpMap[string(b)] != INVALID_OP {
				lval.code = Code{
					opType:   OpMap[string(b)],
					argument: 1,
				}
				return Token
			}
			continue
		}
		break
	}

	return 0
}
