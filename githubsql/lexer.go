package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func Parse(input string) (*Sql, error) {
	l := newLex(input)
	_ = yyParse(l)
	return l.sql, l.err
}

type lex struct {
	sql         *Sql
	inputTokens []string
	index       int
	err         error
	state       string
}

type Sql struct {
	columns []string
	table   string
	count   int
	page    int
}

func (s *Sql) String() string {
	return fmt.Sprintf("select %s from %s with page_count %d, page %d", s.columns, s.table, s.count, s.page)
}

func newLex(input string) *lex {
	re := regexp.MustCompile("[\\s,;]+")
	l := re.Split(input, -1)
	var a []string
	for _, b := range l {
		if b != "" {
			a = append(a, b)
		}
	}
	return &lex{inputTokens: a, index: 0}
}

// Error satisfies yyLexer.
func (l *lex) Error(s string) {
	l.err = errors.New(s)
}

var literal = map[string]int{
	"select": TokSelect,
	"from":   TokFrom,
	"count":  TokLimit,
	"page":   TokPage,
}

// Lex satisfies yyLexer.
func (l *lex) Lex(lval *yySymType) int {
	w := strings.ToLower(l.next())
	if w == "" {
		return 0
	}
	if a, ok := literal[w]; ok {
		l.state = w
		return a
	}
	//log.Println("word:", w, ", state:", l.state)
	switch l.state {
	case "select":
		lval.str = w
		return TokColumn
	case "from":
		lval.str = w
		return TokTable
	case "count":
		a, err := strconv.Atoi(w)
		if err != nil {
			log.Println("limit not valid", err)
			return 0
		}
		lval.num = a
		return TokNum
	case "page":
		a, err := strconv.Atoi(w)
		if err != nil {
			log.Println("offset not valid", err)
			return 0
		}
		lval.num = a
		return TokNum
	default:
		return 0
	}
}

func (l *lex) next() string {
	if l.index >= len(l.inputTokens) {
		return ""
	}
	defer func() { l.index++ }()
	return l.inputTokens[l.index]
}
