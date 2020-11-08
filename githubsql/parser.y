%{

package main
import (
)

%}

%union {
    sql *Sql
    num int
    str string
    strs []string
}

%type <sql> expr expr1
%type <strs> columns
%token TokSelect TokFrom TokLimit TokPage

%token	<str> TokColumn TokTable
%token	<num> TokNum

%%

top:
expr
{
    yylex.(*lex).sql = $1
}

expr:
expr1
{
	$$=$1
}
|expr1 TokLimit TokNum
{
	$1.count = $3
	$$=$1
}
| expr1 TokPage TokNum
{
	$1.page = $3
	$$=$1
}
| expr1 TokLimit TokNum TokPage TokNum
{
	$1.count = $3
	$1.page = $5
	$$=$1
}

expr1:
{
	$$ = &Sql{}
}
| TokSelect columns TokFrom TokTable
{
	$$ = &Sql{columns: $2, table: $4}
}

columns:
TokColumn
{
	$$ = []string{$1}
}
| columns TokColumn
{
	$$ = append($$, $2)
}

%%
