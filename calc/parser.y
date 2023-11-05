// vim: noet sw=8 sts=8
%{
package calc

import (
	"log"
)

type Token struct {
	tok int
	lit string
	val int
	pos Position
}

%}

%union{
	statements []Statement
	statement  Statement
	expr       Expression
	tok        Token
	arr        []NumberExpression
}

%type<statements> statements
%type<statement> statement
%type<expr> expr
%type<arr> array array_element 

%token<tok> IDENT NUMBER VAR 

/* conditional operator TernaryExpression */
%left '?' ':'
%left LOR
%left LAND
%nonassoc EQ NE LE LT GE GT
/* 判断元素是否存在于数组中 */
%left IN
%left ','
%left '+' '-'
%left '*' '/' '%'
%right UNARY

%%

statements
	:
	{
		$$ = nil
		if l, isLexerWrapper := yylex.(*LexerWrapper); isLexerWrapper {
			l.statements = $$
		}
	}
	| statement statements
	{
		$$ = append([]Statement{$1}, $2...)
		if l, isLexerWrapper := yylex.(*LexerWrapper); isLexerWrapper {
			l.statements = $$
		}
	}

statement
	: expr ';'
	{
		$$ = &ExpressionStatement{Expr: $1}
	}
	| VAR IDENT '=' expr ';'
	{
		$$ = &VarDefStatement{VarName: $2.lit, Expr: $4}
	}

expr	: NUMBER
	{
		$$ = &NumberExpression{Val: $1.val}
	}
	| IDENT
	{
		$$ = &IdentifierExpression{Lit: $1.lit}
	}
	| expr '?' expr ':' expr
	{
		$$ = &TernaryExpression{Cond: $1, TrueExpr: $3, FalseExpr: $5}
	}
	| expr IN array
	{
		$$ = &InExpression{LHS: $1, Arr: $3}
	}
	| '!' expr      %prec UNARY
	{
		$$ = &UnaryNotExpression{SubExpr: $2}
	}
	| '-' expr      %prec UNARY
	{
		$$ = &UnaryMinusExpression{SubExpr: $2}
	}
	| '(' expr ')'
	{
		$$ = &ParenExpression{SubExpr: $2}
	}
	| expr LAND expr
	{ $$ = &BinOpLogicExpression{LHS: $1, Operator: LAND, RHS: $3} }
	| expr LOR expr
	{ $$ = &BinOpLogicExpression{LHS: $1, Operator: LOR, RHS: $3} }
	| expr EQ expr
	{ $$ = &BinOpExpression{LHS: $1, Operator: EQ, RHS: $3} }
	| expr NE expr
	{ $$ = &BinOpExpression{LHS: $1, Operator: NE, RHS: $3} }
	| expr LE expr
	{ $$ = &BinOpExpression{LHS: $1, Operator: LE, RHS: $3} }
	| expr LT expr
	{ $$ = &BinOpExpression{LHS: $1, Operator: LT, RHS: $3} }
	| expr GE expr
	{ $$ = &BinOpExpression{LHS: $1, Operator: GE, RHS: $3} }
	| expr GT expr
	{ $$ = &BinOpExpression{LHS: $1, Operator: GT, RHS: $3} }
	| expr '+' expr
	{ $$ = &BinOpExpression{LHS: $1, Operator: int('+'), RHS: $3} }
	| expr '-' expr
	{ $$ = &BinOpExpression{LHS: $1, Operator: int('-'), RHS: $3} }
	| expr '*' expr
	{ $$ = &BinOpExpression{LHS: $1, Operator: int('*'), RHS: $3} }
	| expr '/' expr
	{ $$ = &BinOpExpression{LHS: $1, Operator: int('/'), RHS: $3} }
	| expr '%' expr
	{ $$ = &BinOpExpression{LHS: $1, Operator: int('%'), RHS: $3} }

array
	: '[' array_element ']'
	{
		$$ = $2
	}
	|
	'[' ']'
	{

	}

array_element
	: NUMBER
	{
		$$ = []NumberExpression{NumberExpression{Val: $1.val}}
	}
	| array_element ',' NUMBER
	{
		$$ = append($1, NumberExpression{Val: $3.val})
	}

%%

type LexerWrapper struct {
	s          *Scanner
	recentLit  string
	recentPos  Position
	statements []Statement
}

func (l *LexerWrapper) Lex(lval *yySymType) int {
	tok, lit, pos := l.s.Scan()
	if tok == EOF {
		return 0
	}
	lval.tok = Token{tok: tok, lit: lit, pos: pos}
	if tok == NUMBER {
		lval.tok.val, _ = toNumber(lit)
	}
	l.recentLit = lit
	l.recentPos = pos
	return tok
}

func (l *LexerWrapper) Error(e string) {
	err := __yyfmt__.Sprintf("Line %d, Column %d: %q %s",
		l.recentPos.Line, l.recentPos.Column, l.recentLit, e)
	log.Print(err)
	panic(err)
}

func Parse(s *Scanner) []Statement {
	l := LexerWrapper{s: s}
	if yyParse(&l) != 0 {
		panic("Parse error")
	}
	return l.statements
}
