package calc

import (
	"strconv"
	"strings"
)

const (
	EOF     = -1
	UNKNOWN = 0
)

var keywords = map[string]int{
	"var": VAR,
	"in":  IN,
}

type Position struct {
	Line   int
	Column int
}

type Scanner struct {
	src      []rune
	offset   int
	lineHead int
	line     int
}

func (s *Scanner) Init(src string) {
	l := len(src)
	if l == 0 {
		s.src = []rune{}
		return
	}
	s.src = []rune(src)
	//单行脚本自动加分号
	if !strings.Contains(src, "\n") && s.src[l-1] != ';' {
		s.src = append(s.src, ';')
	}
}

func (s *Scanner) Scan() (tok int, lit string, pos Position) {
	s.skipWhiteSpace()
	pos = s.position()
	switch ch := s.peek(); {
	case isLetter(ch):
		lit = s.scanIdentifier()
		if keyword, ok := keywords[lit]; ok {
			tok = keyword
		} else {
			tok = IDENT
		}
	case isDigit(ch):
		tok, lit = NUMBER, s.scanNumber()
	default:
		switch ch {
		case -1:
			tok = EOF
		case '(', ')', ';', '+', '-', '*', '/', '%', '[', ']', ',', '?', ':':
			tok = int(ch)
			lit = string(ch)
			s.next()
		case '&':
			tok = LAND
			lit = "&&"
			s.next()
			s.next()
		case '|':
			tok = LOR
			lit = "||"
			s.next()
			s.next()
		case '=':
			if s.peekNext() == '=' {
				tok = EQ
				lit = "=="
				s.next()
			} else {
				tok = int('=')
				lit = string(ch)
			}
			s.next()
		case '!':
			if s.peekNext() == '=' {
				tok = NE
				lit = "!="
				s.next()
			} else {
				tok = int('!')
				lit = string(ch)
			}
			s.next()
		case '>':
			if s.peekNext() == '=' {
				tok = GE
				lit = ">="
				s.next()
			} else {
				tok = GT
				lit = ">"
			}
			s.next()
		case '<':
			if s.peekNext() == '=' {
				tok = LE
				lit = "<="
				s.next()
			} else {
				tok = LT
				lit = "<"
			}
			s.next()
		}
	}
	return
}

// ========================================

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isHex(ch rune) bool {
	return '0' <= ch && ch <= '9' || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
}

func toUpper(ch rune) rune {
	if 'a' <= ch && ch <= 'z' {
		ch -= 'a' - 'A'
	}
	return ch
}

func isWhiteSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func (s *Scanner) peek() rune {
	if !s.reachEOF() {
		return s.src[s.offset]
	} else {
		return -1
	}
}

func (s *Scanner) peekNext() rune {
	if s.offset+1 < len(s.src) {
		return s.src[s.offset+1]
	} else {
		return -1
	}
}

func (s *Scanner) next() {
	if !s.reachEOF() {
		if s.peek() == '\n' {
			s.lineHead = s.offset + 1
			s.line++
		}
		s.offset++
	}
}

func (s *Scanner) reachEOF() bool {
	return len(s.src) <= s.offset
}

func (s *Scanner) position() Position {
	return Position{Line: s.line + 1, Column: s.offset - s.lineHead + 1}
}

func (s *Scanner) skipWhiteSpace() {
	for isWhiteSpace(s.peek()) {
		s.next()
	}
}

func (s *Scanner) scanIdentifier() string {
	var ret []rune
	for isLetter(s.peek()) || isDigit(s.peek()) {
		ret = append(ret, s.peek())
		s.next()
	}
	return string(ret)
}

/**
 * @description: 解析一个十进制数字或者十六进制数字.如果是十六进制，会自动转成大写
 * @param {*}
 * @return {*}
 */
func (s *Scanner) scanNumber() string {
	var ret []rune
	if s.peek() == '0' && s.peekNext() == 'x' || s.peekNext() == 'X' {
		ret = append(ret, s.peek())
		s.next()
		ret = append(ret, 'X')
		s.next()
		for isHex(s.peek()) {
			ret = append(ret, toUpper(s.peek()))
			s.next()
		}
	} else {
		for isDigit(s.peek()) {
			ret = append(ret, s.peek())
			s.next()
		}
	}
	return string(ret)
}

func toNumber(lit string) (int, error) {
	if len(lit) >= 2 && lit[1] == 'X' {
		src := []rune(lit)
		val := 0
		for i := 2; i < len(src); i++ {
			var ch rune = src[i]
			if ch >= 'A' {
				val = val*16 + int(ch)&15 + 9
			} else {
				val = val*16 + int(ch)&15
			}
		}
		return val, nil
	} else {
		return strconv.Atoi(lit)
	}
}
