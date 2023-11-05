package unittest

import (
	"testing"

	. "github.com/motto0808/go-calc/calc"
)

func testScanner(t *testing.T, src string, expectTok int) {
	s := new(Scanner)
	s.Init(src)
	tok, lit, _ := s.Scan()
	if tok != expectTok {
		t.Errorf("Expect Scanner{%q}.Scan() = %#v, _ want %#v", src, tok, expectTok)
	}
	if lit != src {
		t.Logf("Scanner{%q}.Scan() = _, %#v want %#v", src, lit, src)
	}
}

func TestScanner(t *testing.T) {
	testScanner(t, "var", VAR)
	testScanner(t, "in", IN)
	testScanner(t, "abc", IDENT)
	testScanner(t, "123", NUMBER)
	testScanner(t, "0xff", NUMBER)
	testScanner(t, "0x123ABC", NUMBER)
	testScanner(t, "(", '(')
	testScanner(t, ")", ')')
	testScanner(t, ";", ';')
	testScanner(t, "+", '+')
	testScanner(t, "-", '-')
	testScanner(t, "*", '*')
	testScanner(t, "/", '/')
	testScanner(t, "%", '%')
	testScanner(t, "=", '=')
	testScanner(t, "!", '!')
	testScanner(t, "==", EQ)
	testScanner(t, "!=", NE)
	testScanner(t, "&&", LAND)
	testScanner(t, "||", LOR)
	testScanner(t, "[", '[')
	testScanner(t, "]", ']')
	testScanner(t, ",", ',')
}
