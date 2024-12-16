package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

	for _, ch := range expression {
		switch {
		case ch >= '0' && ch <= '9' || ch == '.':
			postfix = append(postfix, string(ch))
		case isOperator(ch):
			for len(stack) > 0 && isOperator(stack[len(stack)-1]) && precedence[stack[len(stack)-1]] >= precedence[ch] {
				postfix = append(postfix, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, ch)
		case ch == '(':
			stack = append(stack, ch)
		case ch == ')':
			for len(stack) > 0 && stack[len(stack)-1] != '(' {
				postfix = append(postfix, string(stack[len(stack)-1]))
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, errors.New("некорректное выражение: несогласованные скобки")
			}

			stack = stack[:len(stack)-1]
		default:
			return nil, fmt.Errorf("некорректный символ: %c", ch)
		}
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

func calcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error": "Only POST method is supported"}`))
		return
	}

	type request struct {
		Expression string `json:"expression"`
	}

	type response struct {
		Result float64 `json:"result,omitempty"`
		Error  string  `json:"error,omitempty"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Expression == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response{Error: "Invalid request or empty expression"})
		return
	}

	result, err := Calc(req.Expression)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(response{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response{Result: result})
}

func main() {
	http.HandleFunc("/api/v1/calculat", calcHandler)
	fmt.Println("Server is running on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
