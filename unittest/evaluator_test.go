package unittest

import (
	"fmt"
	"testing"

	. "github.com/motto0808/go-calc/calc"
)

func evaluateContent(content string) int {
	evaluator := NewEvaluator()
	n, err := evaluator.Eval(content, Env{})
	assert(nil, err == nil, fmt.Sprintf("Evaluate %q failed: %s", content, err))
	return n
}

func TestEvaluate(t *testing.T) {
	{
		s, _ := Evaluate(&ExpressionStatement{&NumberExpression{123}}, Env{})
		assert(t, s == "123", fmt.Sprintf("Expect a statement \"123;\" to be evaluated as \"123\", but got %q", s))
	}
	{
		s, err := Evaluate(&ExpressionStatement{&IdentifierExpression{"a"}}, Env{})
		assert(t, err != nil, "Expect an undefined variable \"a\" not to be evaluated, but got %q", s)
	}
	{
		env := Env{"a": 16}
		s, _ := Evaluate(&VarDefStatement{"a", &NumberExpression{123}}, env)
		assert(t, s == "Assign 123 to a", "Expect a VarDef statement do an assignment, but it didn't")
		assert(t, env["a"] == 123, "Expect a VarDef statement do an assignment, but it didn't")
	}
	{
		s, err := Evaluate(&VarDefStatement{"a", &IdentifierExpression{"a"}}, Env{})
		assert(t, err != nil, fmt.Sprintf("Expect an undefined variable \"a\" not to be evaluated, but got %q", s))
	}
	{
		env := Env{"b": 16}
		s, err := Evaluate(&VarDefStatement{"a", &IdentifierExpression{"b"}}, env)
		assert(t, err == nil, "VarDef statement error")
		assert(t, s == "Assign 16 to a", "Expect a VarDef statement do an assignment, but it didn't")
		assert(t, env["a"] == 16, "Expect a VarDef statement do an assignment, but it didn't")
	}
}

func TestEvaluateExpr(t *testing.T) {
	if v, _ := EvaluateExpr(&NumberExpression{123}, Env{}); v != 123 {
		t.Error("Expect a number 123 to be evaluated as 123, but got", v)
	}
	if v, _ := EvaluateExpr(&IdentifierExpression{"a"}, Env{"a": 123}); v != 123 {
		t.Error("Expect a variable \"a\" to be evaluated as 123, but got", v)
	}
	if v, err := EvaluateExpr(&IdentifierExpression{"a"}, Env{}); err == nil {
		t.Error("Expect an undefined variable \"a\" not to be evaluated, but got", v)
	}
	if v, _ := EvaluateExpr(&UnaryMinusExpression{&NumberExpression{123}}, Env{}); v != -123 {
		t.Error("Expect a number -123 to be evaluated as -123, but got", v)
	}
	if v, err := EvaluateExpr(&UnaryMinusExpression{&IdentifierExpression{"a"}}, Env{}); err == nil {
		t.Error("Expect an undefined variable \"a\" not to be evaluated, but got", v)
	}
	if v, _ := EvaluateExpr(&ParenExpression{&NumberExpression{123}}, Env{}); v != 123 {
		t.Error("Expect an expression (123) to be evaluated as 123, but got", v)
	}
	if v, err := EvaluateExpr(&ParenExpression{&IdentifierExpression{"a"}}, Env{}); err == nil {
		t.Error("Expect an undefined variable \"a\" not to be evaluated, but got", v)
	}

	expr1 := &NumberExpression{42}
	expr2 := &NumberExpression{12}
	expr3 := &NumberExpression{26}
	exprE := &IdentifierExpression{"a"}
	if v, _ := EvaluateExpr(&BinOpExpression{expr1, '+', expr2}, Env{}); v != 54 {
		t.Error("Expect an expression 42+12 to be evaluated as 54, but got", v)
	}
	if v, _ := EvaluateExpr(&BinOpExpression{expr1, '-', expr2}, Env{}); v != 30 {
		t.Error("Expect an expression 42-12 to be evaluated as 30, but got", v)
	}
	if v, _ := EvaluateExpr(&BinOpExpression{expr1, '*', expr2}, Env{}); v != 504 {
		t.Error("Expect an expression 42*12 to be evaluated as 3, but got", v)
	}
	if v, _ := EvaluateExpr(&BinOpExpression{expr1, '/', expr2}, Env{}); v != 3 {
		t.Error("Expect an expression 42/12 to be evaluated as 3, but got", v)
	}
	if v, _ := EvaluateExpr(&BinOpExpression{expr1, '%', expr2}, Env{}); v != 6 {
		t.Error("Expect an expression 42%12 to be evaluated as 6, but got", v)
	}
	if v, _ := EvaluateExpr(&BinOpExpression{expr1, '+', expr3}, Env{}); v != 68 {
		t.Error("Expect an expression 42+0x1A to be evaluated as 68, but got", v)
	}
	if v, err := EvaluateExpr(&BinOpExpression{exprE, '+', expr2}, Env{}); err == nil {
		t.Error("Expect an undefined variable \"a\" not to be evaluated, but got", v)
	}
	if v, err := EvaluateExpr(&BinOpExpression{expr1, '+', exprE}, Env{}); err == nil {
		t.Error("Expect an undefined variable \"a\" not to be evaluated, but got", v)
	}
}

func TestEvaluateCompareExpr(t *testing.T) {
	expr1 := &NumberExpression{42}
	expr2 := &NumberExpression{12}
	expr3 := &UnaryMinusExpression{&NumberExpression{123}}
	env := Env{"a": 42, "b": -1}
	var v int = 0
	v, _ = EvaluateExpr(&BinOpExpression{expr1, NE, expr2}, Env{})
	assert(t, v == 1, "test NE")

	v, _ = EvaluateExpr(&BinOpExpression{expr1, EQ, &IdentifierExpression{"a"}}, env)
	assert(t, v == 1, "test EQ")
	v, _ = EvaluateExpr(&BinOpExpression{expr1, GT, expr2}, Env{})
	assert(t, v == 1, "test GT")
	v, _ = EvaluateExpr(&BinOpExpression{expr1, GT, expr1}, Env{})
	assert(t, v == 0, "test GT")
	v, _ = EvaluateExpr(&BinOpExpression{expr1, GE, expr1}, Env{})
	assert(t, v == 1, "test GE")

	v, _ = EvaluateExpr(&BinOpExpression{expr1, LT, &IdentifierExpression{"a"}}, env)
	assert(t, v == 0, "test expr1 < a")
	v, _ = EvaluateExpr(&BinOpExpression{expr1, LE, &IdentifierExpression{"a"}}, env)
	assert(t, v == 1, "test LE")
	v, _ = EvaluateExpr(&BinOpExpression{expr3, LE, &IdentifierExpression{"b"}}, env)
	assert(t, v == 1, "test LE")

}
func TestEvaluateLogicExpr(t *testing.T) {
	env := Env{"a": 42, "b": -1}
	expr1 := &NumberExpression{42}
	expr2 := &NumberExpression{0}
	var v int = 0
	v, _ = EvaluateExpr(&BinOpLogicExpression{expr1, LAND, expr2}, env)
	assert(t, v == 0, "test expr1 && expr2")
	v, _ = EvaluateExpr(&BinOpLogicExpression{expr2, LAND, expr1}, env)
	assert(t, v == 0, "test expr2 && expr1")
	v, _ = EvaluateExpr(&BinOpLogicExpression{expr1, LOR, expr2}, env)
	assert(t, v == 1, "test expr1 || expr2")
	v, _ = EvaluateExpr(&BinOpLogicExpression{expr2, LOR, expr1}, env)
	assert(t, v == 1, "test expr2 || expr1")
}
func TestLargeNum(t *testing.T) {
	expr1 := &NumberExpression{999999999}
	expr2 := &NumberExpression{999999999}
	v, _ := EvaluateExpr(&BinOpExpression{expr1, '+', expr2}, Env{})
	assert(t, v == 1999999998)
}

func TestArray(t *testing.T) {
	expr1 := &NumberExpression{1}
	expr2 := &NumberExpression{2}
	expr3 := &InExpression{LHS: expr1}
	expr3.Arr = append(expr3.Arr, *expr1)
	expr3.Arr = append(expr3.Arr, *expr2)

	v, _ := EvaluateExpr(expr3, Env{})
	assert(t, v == 1)
}

func TestConditionExpr(t *testing.T) {
	{
		env := Env{"a": 16}
		te := TernaryExpression{
			Cond:      &IdentifierExpression{Lit: "a"},
			TrueExpr:  &NumberExpression{Val: 2},
			FalseExpr: &NumberExpression{Val: 3},
		}
		s, err := Evaluate(&ExpressionStatement{Expr: &te}, env)
		assert(t, err == nil, "TernaryExpression error")
		assert(t, s == "2", "Expect expr return 2, but it didn't")
	}

	{
		s := new(Scanner)
		s.Init("a>20?b:20;")
		statements := Parse(s)
		assert(t, len(statements) == 1, "Expect 1 statement, but it didn't")
		r1, _ := Evaluate(statements[0], Env{"a": 16, "b": 10})
		assert(t, r1 == "20", "Expect expr return 20, but it didn't")
		r2, _ := Evaluate(statements[0], Env{"a": 30, "b": 10})
		assert(t, r2 == "10", "Expect expr return 10, but it didn't")
	}

	// 优先级测试
	{
		s := new(Scanner)
		s.Init("a>20?b+c:20*5;")
		statements := Parse(s)
		assert(t, len(statements) == 1, "Expect 1 statement, but it didn't")
		r1, _ := Evaluate(statements[0], Env{"a": 16, "b": 10, "c": 4})
		assert(t, r1 == "100", "Expect expr return 100, but it didn't")
		r2, _ := Evaluate(statements[0], Env{"a": 30, "b": 10, "c": 4})
		assert(t, r2 == "14", "Expect expr return 14, but it didn't")
	}
}

func TestMultiStatement(t *testing.T) {
	n := evaluateContent("var a=1;var b=a>10?a:10;a+b")
	assert(t, n == 11, "Expect 11, but it didn't")
}
