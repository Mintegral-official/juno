package query

import (
	"bytes"
	"github.com/Mintegral-official/juno/datastruct"
	"strconv"
	"strings"
)

var Char = []rune{
	'a', 'v', '=', 'X', 'c', 'R', 'B', '5', 'e', 'C', '[', 'h', '@', 'n', 'W', '<', '9', 'l',
	'u', 'L', 'w', 'Q', 'b', 'y', '4', 'm', 'z', 'A', '6', '!', 's', 'r', '8', 'E', 'o', 'F',
	'Y', 'P', 'd', 'H', 'S', 'f', 'T', 'U', 'q', 'G', '2', 'p', 'x', 'V', '0', 'k', 't', 'Z',
	'N', '#', '1', 'M', '3', 'D', 'j', 'K', '.', 'I', 'g', 'O', ',', 'J', ']', 'i', '7', '>',
}

type Expression struct {
	str string
}

func NewExpression(str string) *Expression {
	return &Expression{str: str}
}

func (e *Expression) GetStr() string {
	return e.str
}

func (e *Expression) string2Strings() []string {
	if strings.Contains(e.str, " and ") {
		e.str = strings.Replace(e.str, "and", " & ", -1)
	}
	if strings.Contains(e.str, " or ") {
		e.str = strings.Replace(e.str, " or ", " | ", -1)
	}
	if strings.Contains(e.str, " in ") {
		e.str = strings.Replace(e.str, " in ", " @ ", -1)
	}
	if strings.Contains(e.str, " !in ") {
		e.str = strings.Replace(e.str, " !in ", " # ", -1)
	}
	str, s, n := e.str, make([]string, 0), 0
	var t bytes.Buffer
	for i, v := range str {
		if v == ' ' {
			continue
		}
		if e.isDigitOrChar(v) {
			t.WriteRune(v)
		} else {
			vs := string(v)
			if !e.isSign(vs) {
				panic("index: " + strconv.Itoa(i) + ", Illegal Character: " + vs)
			}
			if t.Len() > 0 {
				s = append(s, t.String())
				t.Reset()
			}
			s = append(s, vs)
			if v == '(' {
				n++
			} else if v == ')' {
				n--
			}
		}
	}
	if t.Len() > 0 {
		s = append(s, t.String())
	}
	if n != 0 {
		panic("the expression is wrong ")
	}
	return s
}

func (e *Expression) ToPostfix(exp []string) []string {
	result, s := make([]string, 0), datastruct.NewStack()
	for _, str := range exp {
		if e.isSign(str) {
			if str == "(" || s.Len() == 0 {
				s.Push(str)
			} else {
				if str == ")" {
					for s.Len() > 0 {
						if s.Peek() == "(" {
							s.Pop()
							break
						}
						result = e.appendStr(result, s.Pop().(string))
					}
				} else {
					for s.Len() > 0 && s.Peek() != "(" && e.getSignstr(str)-e.getSignstr(s.Peek().(string)) <= 0 {
						result = e.appendStr(result, s.Pop().(string))
					}
					s.Push(str)
				}
			}
		} else {
			result = e.appendStr(result, str)
		}
	}
	for s.Len() > 0 {
		result = e.appendStr(result, s.Pop().(string))
	}
	return result
}

func (e *Expression) appendStr(slice []string, str string) []string {
	if str == "(" || str == ")" {
		return slice
	}
	return append(slice, str)
}

func (e *Expression) getSignstr(str string) int {
	if str == "(" || str == ")" {
		return 2
	} else if str == "&" || str == "|" {
		return 1
	}
	return 0
}

func (e *Expression) isDigitOrChar(r rune) bool {
	for _, v := range Char {
		if v == r {
			return true
		}
	}
	return false
}

func (e *Expression) isSign(str string) bool {
	if str == "(" || str == ")" || str == "&" || str == "|" {
		return true
	}
	return false
}
