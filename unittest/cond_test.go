package unittest

import (
	"testing"

	"github.com/motto0808/go-calc/calc"
)

// 这里模仿游戏业务做了一些简单的条件
type ICondition interface {
	Calc(args interface{}) int
}

type typeCondMap map[string]ICondition

var condMap = make(typeCondMap, 32)

type CondHelper struct {
}

var condHelper = CondHelper{}

func (*CondHelper) Eval(name string, args interface{}) int {
	cond := condMap[name]
	if cond == nil {
		return 0
	}
	return int(cond.Calc(args))
}

type ChargeCond struct {
}
type AgeCond struct {
}

func init() {
	condMap["charge"] = &ChargeCond{}
	condMap["age"] = &AgeCond{}

}
func (c ChargeCond) Calc(args interface{}) int {
	return 500
}

func (c AgeCond) Calc(args interface{}) int {
	return 20
}

func TestCondition(t *testing.T) {
	eva := calc.NewEvaluator()
	eva.SetCondHelper(&condHelper, nil)
	// 求值时在局部环境中找不到变量，会尝试通过条件求值
	v, err := eva.Eval("charge>=200 && age<=30", calc.Env{})
	assert(t, err == nil)
	assert(t, v == 1)
}
