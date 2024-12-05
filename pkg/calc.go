package pkg

import (
	"errors"
	"strconv"
	"strings"
)

var precedence = map[rune]int{
	'+': 1,
	'-': 1,
	'*': 2,
	'/': 2,
}

func Calc(expression string) (float64, error) {
	tokens := tokenize(expression)
	if len(tokens) == 0 {
		return 0, ErrorEmptyExpr
	}

	output := []string{}
	operatorStack := []rune{}

	for _, token := range tokens {
		if isNumber(token) {
			output = append(output, token)
		} else if isOperator(rune(token[0])) {
			for len(operatorStack) > 0 && precedence[rune(operatorStack[len(operatorStack)-1])] >= precedence[rune(token[0])] {
				output = append(output, string(operatorStack[len(operatorStack)-1]))
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			operatorStack = append(operatorStack, rune(token[0]))
		} else if token == "(" {
			operatorStack = append(operatorStack, '(')
		} else if token == ")" {
			for len(operatorStack) > 0 && operatorStack[len(operatorStack)-1] != '(' {
				output = append(output, string(operatorStack[len(operatorStack)-1]))
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			if len(operatorStack) == 0 {
				return 0, ErrorMismatchedParh
			}
			operatorStack = operatorStack[:len(operatorStack)-1] // Удаляем '('
		} else {
			return 0, errors.New("invalid token: " + token)
		}
	}

	for len(operatorStack) > 0 {
		output = append(output, string(operatorStack[len(operatorStack)-1]))
		operatorStack = operatorStack[:len(operatorStack)-1]
	}

	return evaluateRPN(output)
}

func tokenize(expression string) []string {
	var tokens []string
	var currentToken strings.Builder

	for _, ch := range expression {
		if ch == ' ' {
			continue
		}
		if isOperator(ch) || ch == '(' || ch == ')' {
			if currentToken.Len() > 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			tokens = append(tokens, string(ch))
		} else if isDigit(ch) || ch == '.' {
			currentToken.WriteRune(ch)
		} else {
			return nil // Неверный токен
		}
	}
	if currentToken.Len() > 0 {
		tokens = append(tokens, currentToken.String())
	}
	return tokens
}

func isOperator(ch rune) bool {
	return ch == '+' || ch == '-' || ch == '*' || ch == '/'
}

func isNumber(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func evaluateRPN(tokens []string) (float64, error) {
	stack := []float64{}

	for _, token := range tokens {
		if isNumber(token) {
			num, _ := strconv.ParseFloat(token, 64)
			stack = append(stack, num)
		} else if isOperator(rune(token[0])) {
			if len(stack) < 2 {
				return 0, ErrInvalidExpr
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch token[0] {
			case '+':
				stack = append(stack, a+b)
			case '-':
				stack = append(stack, a-b)
			case '*':
				stack = append(stack, a*b)
			case '/':
				if b == 0 {
					return 0, ErrDivisonByZero
				}
				stack = append(stack, a/b)
			default:
				return 0, ErrIvalidOperat
			}
		} else {
			return 0, errors.New("invalid token: " + token)
		}
	}

	if len(stack) != 1 {
		return 0, ErrInvalidExpr
	}
	return stack[0], nil
}
