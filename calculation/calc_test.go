package calc_test

import (
	"testing"
)

func TestCalc(t *testing.T) {
	// Успешные кейсы
	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult float64
	}{
		{
			name:           "simple addition",
			expression:     "1+1",
			expectedResult: 2,
		},
		{
			name:           "parentheses priority",
			expression:     "(2+2)*2",
			expectedResult: 8,
		},
		{
			name:           "multiplication priority",
			expression:     "2+2*2",
			expectedResult: 6,
		},
		{
			name:           "division",
			expression:     "1/2",
			expectedResult: 0.5,
		},
		{
			name:           "complex expression",
			expression:     "(3+5)*(2-1)/4",
			expectedResult: 2,
		},
	}

	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := calc.Calc(testCase.expression)
			if err != nil {
				t.Fatalf("expression %s returned unexpected error: %v", testCase.expression, err)
			}
			if result != testCase.expectedResult {
				t.Fatalf("expected %f but got %f for expression %s", testCase.expectedResult, result, testCase.expression)
			}
		})
	}

	// Ошибочные кейсы
	testCasesFail := []struct {
		name        string
		expression  string
		expectedErr string
	}{
		{
			name:        "invalid operator at the end",
			expression:  "1+1*",
			expectedErr: "некорректное выражение",
		},
		{
			name:        "double operator",
			expression:  "2+2**2",
			expectedErr: "некорректный оператор",
		},
		{
			name:        "unmatched parentheses",
			expression:  "((2+2)-*(2",
			expectedErr: "некорректное выражение: несогласованные скобки",
		},
		{
			name:        "empty expression",
			expression:  "",
			expectedErr: "пустое выражение",
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := calc.Calc(testCase.expression)
			if err == nil {
				t.Fatalf("expression %s should have failed but did not", testCase.expression)
			}
			if err.Error() != testCase.expectedErr {
				t.Fatalf("expected error %s but got %s for expression %s", testCase.expectedErr, err.Error(), testCase.expression)
			}
		})
	}
}
