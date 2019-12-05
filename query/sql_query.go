package query

import (
	"bytes"
	"github.com/Mintegral-official/juno/datastruct"
)

const (
	AND = "&"
	OR  = "|"
	IN  = "@"
	NOT = "^"
)

func String2Strings(str string) []string {
	s := make([]string, 0)
	var t bytes.Buffer
	n := 0
	for _, r := range str {
		if r == ' ' {
			continue
		}
		if isDigitOrChar(r) {
			t.WriteRune(r)
		} else {
			rs := string(r)
			if !isSign(rs) {
				panic("unknown sign: " + rs)
			}
			if t.Len() > 0 {
				s = append(s, t.String())
				t.Reset()
			}
			s = append(s, rs)
			if r == '(' {
				n++
			} else if r == ')' {
				n--
			}
		}
	}
	if t.Len() > 0 {
		s = append(s, t.String())
	}
	if n != 0 {
		panic("the number of '(' is not equal to the number of ')' ")
	}
	return s
}

func ToPostfix(exp []string) []string {
	result := make([]string, 0)
	s := datastruct.NewStack()
	for _, str := range exp {
		if isSign(str) {
			if str == "(" || s.Len() == 0 {
				s.Push(str)
			} else {
				if str == ")" {
					for s.Len() > 0 {
						if s.Peek().(string) == "(" {
							s.Pop()
							break
						}
						result = appendStr(result, s.Pop().(string))
					}
				} else {
					for s.Len() > 0 && s.Peek().(string) != "(" && signCompare(str, s.Peek().(string)) <= 0 {
						result = appendStr(result, s.Pop().(string))
					}
					s.Push(str)
				}
			}
		} else {
			result = appendStr(result, str)
		}
	}
	for s.Len() > 0 {
		result = appendStr(result, s.Pop().(string))
	}
	return result
}

func appendStr(slice []string, str string) []string {
	if str == "(" || str == ")" {
		return slice
	}
	return append(slice, str)
}

func signCompare(a, b string) int {
	return getSignValue(a) - getSignValue(b)
}

func getSignValue(a string) int {
	switch a {
	case "(", ")":
		return 2
	case "&", "|":
		return 1
	default:
		return 0
	}
}

func isDigitOrChar(r rune) bool {
	if r >= '0' && r <= '9' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
		return true
	}
	if r == '=' || r == '>' || r == '<' || r == '!' || r == '[' || r == ']' || r == ',' || r == '@' || r == '^' {
		return true
	}
	return false
}

func isSign(s string) bool {
	switch s {
	case "(", ")", "&", "|":
		return true
	default:
		return false
	}
}
