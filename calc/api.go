package calc

import (
	"fmt"
	"log"
)

type Evaluator struct {
	condFac  ICondHelper
	condArgs interface{}
}

func NewEvaluator() *Evaluator {
	e := new(Evaluator)
	return e
}

/**
 * @description: 条件辅助类，负责对条件求值
 * @param {*}
 * @return {*}
 */
type ICondHelper interface {
	Eval(name string, args interface{}) int
}

func (e *Evaluator) SetCondHelper(condFac ICondHelper, condArgs interface{}) {
	e.condFac = condFac
	e.condArgs = condArgs
}

func (e Evaluator) Eval(content string, env Env) (n int, err error) {
	scanner := new(Scanner)
	scanner.Init(content)
	statements := Parse(scanner)
	for _, s := range statements {
		n, err = e.EvaluateStmt(s, env)
		if err != nil {
			err = fmt.Errorf("evaluator failed to eval: %s", err)
			break
		}
	}
	return
}

func (e Evaluator) EvaluateStmt(statement Statement, env Env) (int, error) {
	switch stmt := statement.(type) {
	case *ExpressionStatement:
		v, err := e.evaluateExpr(stmt.Expr, env)
		if err != nil {
			return 0, err
		}
		return v, nil
	case *VarDefStatement:
		v, err := e.evaluateExpr(stmt.Expr, env)
		if err != nil {
			return 0, err
		}
		env[stmt.VarName] = v
		return v, nil
	default:
		panic("Unknown Statement type")
	}
}

func (eva Evaluator) evalIdWithCond(id string) (ret int, ok bool) {
	if eva.condFac == nil {
		return 0, false
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("error while eval cond [%s]: %s\n", id, r)
		}
	}()
	ret = eva.condFac.Eval(id, eva.condArgs)
	ok = true
	return
}

func (eva Evaluator) evaluateExpr(expr Expression, env Env) (int, error) {
	switch e := expr.(type) {
	case *NumberExpression:
		return e.Val, nil
	case *IdentifierExpression:
		if v, ok := env[e.Lit]; ok {
			return v, nil
		}
		if v, ok := eva.evalIdWithCond(e.Lit); ok {
			env[e.Lit] = v
			return v, nil
		} else {
			return 0, fmt.Errorf("undefined variable: %s", e.Lit)
		}
	case *UnaryMinusExpression:
		v, err := eva.evaluateExpr(e.SubExpr, env)
		if err != nil {
			return 0, err
		}
		return -v, nil
	case *UnaryNotExpression:
		v, err := eva.evaluateExpr(e.SubExpr, env)
		if err != nil {
			return 0, err
		}
		if v == 0 {
			v = 1
		} else {
			v = 0
		}
		return v, nil
	case *ParenExpression:
		v, err := eva.evaluateExpr(e.SubExpr, env)
		if err != nil {
			return 0, err
		}
		return v, nil
	case *BinOpExpression:
		lhsV, err := eva.evaluateExpr(e.LHS, env)
		if err != nil {
			return 0, err
		}
		rhsV, err := eva.evaluateExpr(e.RHS, env)
		if err != nil {
			return 0, err
		}
		switch e.Operator {
		case EQ:
			return boolToInt(lhsV == rhsV), nil
		case NE:
			return boolToInt(lhsV != rhsV), nil
		case GE:
			return boolToInt(lhsV >= rhsV), nil
		case GT:
			return boolToInt(lhsV > rhsV), nil
		case LE:
			return boolToInt(lhsV <= rhsV), nil
		case LT:
			return boolToInt(lhsV < rhsV), nil
		case '+':
			return lhsV + rhsV, nil
		case '-':
			return lhsV - rhsV, nil
		case '*':
			return lhsV * rhsV, nil
		case '/':
			return lhsV / rhsV, nil
		case '%':
			return lhsV % rhsV, nil
		default:
			panic("Unknown operator")
		}
	case *BinOpLogicExpression:
		lhsV, err := eva.evaluateExpr(e.LHS, env)
		if err != nil {
			return 0, err
		}
		if e.Operator == LAND {
			if lhsV == 0 {
				return 0, nil
			}
		} else {
			if lhsV != 0 {
				return 1, nil
			}
		}

		rhsV, err := eva.evaluateExpr(e.RHS, env)
		if err != nil {
			return 0, err
		}
		return boolToInt(rhsV != 0), nil
	case *InExpression:
		lhsV, err := eva.evaluateExpr(e.LHS, env)
		if err != nil {
			return 0, err
		}
		var v int = 0
		for _, ele := range e.Arr {
			if lhsV == ele.Val {
				v = 1
				break
			}
		}
		return v, nil
	case *TernaryExpression:
		condV, err := eva.evaluateExpr(e.Cond, env)
		if err != nil {
			return 0, err
		}
		if condV != 0 {
			return eva.evaluateExpr(e.TrueExpr, env)
		}
		return eva.evaluateExpr(e.FalseExpr, env)

	default:
		panic("Unknown Expression type")
	}
}

type Parser struct {
}

func NewParser() *Parser {
	p := new(Parser)
	return p
}

func (p *Parser) Parse(content string) (stmts []Statement) {
	scanner := new(Scanner)
	scanner.Init(content)
	return Parse(scanner)
}
