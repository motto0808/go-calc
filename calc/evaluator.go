package calc

import (
	"fmt"
	"strconv"
)

type Env map[string]int

/**
 * @description: 单句求值
 * @param {Statement} statement
 * @param {Env} env
 * @return {*}
 */
func Evaluate(statement Statement, env Env) (string, error) {
	eva := NewEvaluator()
	switch stmt := statement.(type) {
	case *ExpressionStatement:
		v, err := eva.evaluateExpr(stmt.Expr, env)
		if err != nil {
			return "", err
		}
		return strconv.Itoa(v), nil
	case *VarDefStatement:
		v, err := eva.evaluateExpr(stmt.Expr, env)
		if err != nil {
			return "", err
		}
		env[stmt.VarName] = v
		return fmt.Sprintf("Assign %v to %s", v, stmt.VarName), nil
	default:
		panic("Unknown Statement type")
	}
}

func EvaluateExpr(expr Expression, env Env) (int, error) {
	eva := NewEvaluator()
	return eva.evaluateExpr(expr, env)
}

func boolToInt(cond bool) int {
	if cond {
		return 1
	}
	return 0
}
