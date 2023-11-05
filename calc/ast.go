package calc

type (
	Statement interface {
		statement()
	}

	Expression interface {
		expression()
	}
)

type (
	ExpressionStatement struct {
		Expr Expression
	}

	VarDefStatement struct {
		VarName string
		Expr    Expression
	}
)

func (x *ExpressionStatement) statement() {}
func (x *VarDefStatement) statement()     {}

type (
	NumberExpression struct {
		Val int
	}

	ArrayExpression struct {
		Arr []NumberExpression
	}

	InExpression struct {
		LHS Expression
		Arr []NumberExpression
	}

	TernaryExpression struct {
		Cond      Expression
		TrueExpr  Expression
		FalseExpr Expression
	}

	IdentifierExpression struct {
		Lit string
	}

	UnaryMinusExpression struct {
		SubExpr Expression
	}

	UnaryNotExpression struct {
		SubExpr Expression
	}

	ParenExpression struct {
		SubExpr Expression
	}

	BinOpExpression struct {
		LHS      Expression
		Operator int
		RHS      Expression
	}

	BinOpLogicExpression struct {
		LHS      Expression
		Operator int
		RHS      Expression
	}
)

func (x *NumberExpression) expression()     {}
func (x *ArrayExpression) expression()      {}
func (x *IdentifierExpression) expression() {}
func (x *UnaryMinusExpression) expression() {}
func (x *UnaryNotExpression) expression()   {}
func (x *ParenExpression) expression()      {}
func (x *BinOpExpression) expression()      {}
func (x *BinOpLogicExpression) expression() {}
func (x *InExpression) expression()         {}
func (x *TernaryExpression) expression()    {}
