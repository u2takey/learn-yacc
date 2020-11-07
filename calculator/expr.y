%{

package main
import (
	"fmt"
	"math/big"
)

%}

%union {
    num *big.Rat
}

%type <num> expr expr1 expr2 expr3 //  定义在后面的符号比定义在前面的符号具有更好的优先级

%token '+' '-' '*' '/' '(' ')'

%token	<num> NUM

%%

top:
expr
{
    if $1.IsInt() {
	fmt.Println($1.Num().String())
    } else {
        fmt.Println($1.String())
    }
}

expr:
expr1
{}
|
'+' expr
{
    $$ = $2
}
|
'-' expr
{
    $$ = $2.Neg($2)
}
;

expr1:
expr2
{}
|
expr1 '+' expr2
{
    $$ = $1.Add($1, $3)
}
|
expr1 '-' expr2
{
    $$ = $1.Sub($1, $3)
}
;

expr2:
expr3
{}
|
expr2 '*' expr3
{
    $$ = $1.Mul($1, $3)
}
|
expr2 '/' expr3
{
    $$ = $1.Quo($1, $3)
}
;

expr3:
NUM
|
'(' expr ')'
{
    $$ = $2
}

%%
