package query

import (
	"bytes"
)

var Char = []rune{
	'a', 'v', '=', 'X', 'c', 'R', 'B', '5', 'e', 'C', '[', 'h', '@', 'n', 'W', '<', '9', 'l', 'x', ']', 'I',
	'u', 'L', 'w', 'Q', 'b', 'y', '4', 'm', 'z', 'A', '6', '!', 's', 'r', '8', 'E', 'o', 'F', '7', 'J', 'O',
	'Y', 'P', 'd', 'H', 'S', 'f', 'T', 'U', 'q', 'G', '2', 'p', 'V', '0', 'k', 't', 'Z', '>', 'i', ',', 'g',
	'N', '#', '1', 'M', '3', 'D', 'j', 'K',
}

func String2Strings(str string) []string {
	s := make([]string, 0)
	var t bytes.Buffer
	n := 0
	for _, v := range str {
		if v == ' ' {
			continue
		}
		if isDigitOrChar(v) {
			t.WriteRune(v)
		} else {
			vs := string(v)
			if !isSign(vs) {
				panic("Illegal Character: " + vs)
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
		panic("the expression '(' is not equal to ')' ")
	}
	return s
}

func ToPostfix(exp []string) []string {
	result, s := make([]string, 0), NewStack()
	for _, str := range exp {
		if isSign(str) {
			if str == "(" || s.Len() == 0 {
				s.Push(str)
			} else {
				if str == ")" {
					for s.Len() > 0 {
						if s.Peek() == "(" {
							s.Pop()
							break
						}
						result = appendStr(result, s.Pop())
					}
				} else {
					for s.Len() > 0 && s.Peek() != "(" && getSignValue(str)-getSignValue(s.Peek()) <= 0 {
						result = appendStr(result, s.Pop())
					}
					s.Push(str)
				}
			}
		} else {
			result = appendStr(result, str)
		}
	}
	for s.Len() > 0 {
		result = appendStr(result, s.Pop())
	}
	return result
}

func appendStr(slice []string, str string) []string {
	if str == "(" || str == ")" {
		return slice
	}
	return append(slice, str)
}

func getSignValue(str string) int {
	if str == "(" || str == ")" {
		return 2
	} else if str == "&" || str == "|" {
		return 1
	}
	return 0
}

func isDigitOrChar(r rune) bool {
	for _, v := range Char {
		if v == r {
			return true
		}
	}
	return false
}

func isSign(str string) bool {
	if str == "(" || str == ")" || str == "&" || str == "|" {
		return true
	}
	return false
}
