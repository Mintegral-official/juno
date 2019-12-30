package query

import (
	"bytes"
	"github.com/Mintegral-official/juno/datastruct"
	"strconv"
)

var Char = []rune{
	'a', 'v', '=', 'X', 'c', 'R', 'B', '5', 'e', 'C', '[', 'h', '@', 'n', 'W', '<', '9', 'l', 'x', ']', 'I',
	'u', 'L', 'w', 'Q', 'b', 'y', '4', 'm', 'z', 'A', '6', '!', 's', 'r', '8', 'E', 'o', 'F', '7', 'J', 'O',
	'Y', 'P', 'd', 'H', 'S', 'f', 'T', 'U', 'q', 'G', '2', 'p', 'V', '0', 'k', 't', 'Z', '>', 'i', ',', 'g',
	'N', '#', '1', 'M', '3', 'D', 'j', 'K', '.',
}

type Expression struct {
	value string
}

func NewExpression(str string) *Expression {
	return &Expression{value: str}
}

func (e *Expression) GetValue() string {
	return e.value
}

func (e *Expression) string2Strings() []string {
	str, s, n := e.value, make([]string, 0), 0
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
					for s.Len() > 0 && s.Peek() != "(" && e.getSignValue(str)-e.getSignValue(s.Peek().(string)) <= 0 {
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

func (e *Expression) getSignValue(str string) int {
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
