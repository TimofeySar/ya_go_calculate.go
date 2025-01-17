package calculation

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Calc(expression string) (float64, error) {
	expression = strings.ReplaceAll(expression, " ", "")
	if expression == "" {
		return 0, errors.New("пустое выражение")
	}

	postfix, err := infixToPostfix(expression)
	if err != nil {
		return 0, err
	}

	return evaluatePostfix(postfix)
}

func infixToPostfix(expression string) ([]string, error) {
	var postfix []string
	var stack []rune
	precedence := map[rune]int{
		'+': 1,
		'-': 1,
		'*': 2,
		'/': 2,
	}

	isOperator := func(ch rune) bool {
		_, exists := precedence[ch]
		return exists
	}

	openParentheses := 0
	lastChar := ' '

	for i, ch := range expression {
		switch {
		case ch >= '0' && ch <= '9' || ch == '.':
			if lastChar >= '0' && lastChar <= '9' {
				postfix[len(postfix)-1] += string(ch)
			} else {
				postfix = append(postfix, string(ch))
			}
			lastChar = ch
		case isOperator(ch):
			if isOperator(lastChar) || lastChar == '(' || i == 0 {
				return nil, fmt.Errorf("некорректный оператор: %c", ch)
			}
			for len(stack) > 0 && isOperator(stack[len(stack)-1]) && precedence[stack[len(stack)-1]] >= precedence[ch] {
				postfix = append(postfix, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, ch)
			lastChar = ch
		case ch == '(':
			stack = append(stack, ch)
			openParentheses++
			lastChar = ch
		case ch == ')':
			if openParentheses == 0 {
				return nil, errors.New("некорректное выражение: несогласованные скобки")
			}
			for len(stack) > 0 && stack[len(stack)-1] != '(' {
				postfix = append(postfix, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			stack = stack[:len(stack)-1]
			openParentheses--
			lastChar = ch
		default:
			return nil, fmt.Errorf("некорректный символ: %c", ch)
		}
	}

	if openParentheses > 0 {
		return nil, errors.New("некорректное выражение: несогласованные скобки")
	}

	for len(stack) > 0 {
		if stack[len(stack)-1] == '(' {
			return nil, errors.New("некорректное выражение: несогласованные скобки")
		}
		postfix = append(postfix, string(stack[len(stack)-1]))
		stack = stack[:len(stack)-1]
	}

	return postfix, nil
}

func evaluatePostfix(postfix []string) (float64, error) {
	var stack []float64

	for _, token := range postfix {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, errors.New("некорректное выражение: недостаточно операндов")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch token {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				if b == 0 {
					return 0, errors.New("деление на ноль")
				}
				stack = append(stack, a/b)
			default:
				return 0, fmt.Errorf("некорректный оператор: %s", token)
			}
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("некорректное выражение: лишние операнды")
	}

	return stack[0], nil
}
