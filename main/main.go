package main

import (
	"fmt"
	"log"
	"os"

	"github.com/motto0808/go-calc/calc"
)

func main() {
	env := calc.Env{}
	for _, arg := range os.Args[1:] {
		body, err := os.ReadFile(arg)
		if err != nil {
			log.Fatal(err)
		}
		p := calc.NewParser()
		evaluator := calc.NewEvaluator()
		stmts := p.Parse(string(body))
		for _, stmt := range stmts {
			fmt.Println(evaluator.EvaluateStmt(stmt, env))
		}
	}
}
