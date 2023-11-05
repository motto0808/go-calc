package unittest

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/motto0808/go-calc/calc"
)

func parseStmt(t *testing.T, src string, expect interface{}) {
	s := new(Scanner)
	s.Init(src)
	statements := Parse(s)
	if len(statements) != 1 {
		t.Errorf("Expect %q to be parsed", src)
		return
	}
	if !reflect.DeepEqual(statements[0], expect) {
		t.Errorf("Expect %+#v to be %+#v", statements[0], expect)
		return
	}
}

func TestParseStatement(t *testing.T) {
	parseStmt(t, "123;", &ExpressionStatement{&NumberExpression{123}})
	parseStmt(t, "var a = 123;", &VarDefStatement{"a", &NumberExpression{123}})
}

func parseExpr(t *testing.T, src string, expect interface{}) {
	s := new(Scanner)
	s.Init(src + ";")
	statements := Parse(s)
	if len(statements) != 1 {
		t.Errorf("Expect %q to be parsed", src)
		return
	}
	if exprStmt, isExprStmt := statements[0].(*ExpressionStatement); isExprStmt {
		if !reflect.DeepEqual(exprStmt.Expr, expect) {
			t.Errorf("Expect %+#v to be %+#v", exprStmt.Expr, expect)
			return
		}
	}
}

/**
 * @description: 传入的语句必定发生错误，如果不发生错误，则测试用例失败
 * @param {*testing.T} t
 * @param {string} src
 * @return {*}
 */
func expectError(t *testing.T, src string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered while compile error", r)
		}
	}()
	s := new(Scanner)
	s.Init(src + ";")
	statements := Parse(s)
	if len(statements) > 0 {
		t.Fail()
	}
}

func TestParseExpr(t *testing.T) {
	parseExpr(t, "123", &NumberExpression{123})
	parseExpr(t, "0XFF", &NumberExpression{255})
	parseExpr(t, "0xaa", &NumberExpression{170})
	parseExpr(t, "abc", &IdentifierExpression{"abc"})
	parseExpr(t, "-abc", &UnaryMinusExpression{&IdentifierExpression{"abc"}})
	parseExpr(t, "(abc)", &ParenExpression{&IdentifierExpression{"abc"}})

	aExp := &IdentifierExpression{"a"}
	bExp := &IdentifierExpression{"b"}
	parseExpr(t, "a+b", &BinOpExpression{aExp, '+', bExp})
	parseExpr(t, "a-b", &BinOpExpression{aExp, '-', bExp})
	parseExpr(t, "a*b", &BinOpExpression{aExp, '*', bExp})
	parseExpr(t, "a/b", &BinOpExpression{aExp, '/', bExp})
	parseExpr(t, "a%b", &BinOpExpression{aExp, '%', bExp})

	parseExpr(t, "!a", &UnaryNotExpression{aExp})
	parseExpr(t, "a==b", &BinOpExpression{aExp, EQ, bExp})
	parseExpr(t, "a!=b", &BinOpExpression{aExp, NE, bExp})
	parseExpr(t, "a>=b", &BinOpExpression{aExp, GE, bExp})
	parseExpr(t, "a>b", &BinOpExpression{aExp, GT, bExp})
	parseExpr(t, "a<=b", &BinOpExpression{aExp, LE, bExp})
	parseExpr(t, "a<b", &BinOpExpression{aExp, LT, bExp})

	// condition expr
	parseExpr(t, "a?1:3", &TernaryExpression{aExp, &NumberExpression{1}, &NumberExpression{3}})
}

func TestParseError(t *testing.T) {
	expectError(t, "a<b<c")
	expectError(t, "a<b>c")
}
