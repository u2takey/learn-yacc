%{
package main

func appendCode(yylex yyLexer, c Code) {
  yylex.(*lex).codes = append(yylex.(*lex).codes, c)
}
%}

%union {
    codes []Code
    code Code
}

%type <codes> expr
%token <code> Token

%%

expr:
{}
| Token
{
    appendCode(yylex, $1)
}
| expr Token
{
    appendCode(yylex, $2)
}
