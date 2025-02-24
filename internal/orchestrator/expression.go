package orchestrator

import (
	"fmt"
	"strconv"
	"time"

	"github.com/TimofeySar/ya_go_calculate.go/internal/calculation"
)

type Expression struct {
	ID      string  `json:"id"`
	Status  string  `json:"status"`
	Result  float64 `json:"result,omitempty"`
	tasks   map[string]float64
	expr    string
	postfix []string
}

func NewExpression(id, expr string) *Expression {
	return &Expression{
		ID:     id,
		Status: "pending",
		tasks:  make(map[string]float64),
		expr:   expr,
	}
}
func (s *Expression) Start(tasksChan chan<- calculation.Task) {
	postfix, err := calculation.InfixToPostfix(s.expr)
	if err != nil {
		s.Status = "error"
		return
	}
	s.postfix = postfix

	tasks, err := calculation.GenerateTasks(postfix)
	if err != nil {
		s.Status = "error"
		return
	}

	for _, task := range tasks {
		s.tasks[task.ID] = 0
		tasksChan <- task
	}

	go s.checkCompletion()
}

func (e *Expression) UpdateTaskResult(taskID string, result float64) bool {
	if _, ok := e.tasks[taskID]; ok {
		e.tasks[taskID] = result
		return true
	}
	return false
}

func (e *Expression) checkCompletion() {
	for {
		time.Sleep(500 * time.Millisecond)
		allDone := true
		for _, result := range e.tasks {
			if result == 0 { // Если есть невыполненные задачи
				allDone = false
				break
			}
		}
		if allDone {
			// Пересчитываем результат на основе постфиксной записи
			e.Result = e.calculateResult()
			e.Status = "completed"
			break
		}
	}
}

func (e *Expression) calculateResult() float64 {
	var stack []float64

	for _, token := range e.postfix {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0
			}
			arg2 := stack[len(stack)-1]
			arg1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result float64
			switch token {
			case "+":
				result = arg1 + arg2
			case "-":
				result = arg1 - arg2
			case "*":
				result = arg1 * arg2
			case "/":
				if arg2 == 0 {
					return 0 // Ошибка: деление на ноль
				}
				result = arg1 / arg2
			default:
				return 0 // Ошибка: неизвестная операция
			}
			stack = append(stack, result)
		}
	}

	if len(stack) != 1 {
		return 0 // Ошибка: лишние операнды
	}

	return stack[0]
}

func generateID() string {
	return fmt.Sprintf("expr-%d", time.Now().UnixNano())
}
