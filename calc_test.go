package main

import "testing"

func TestCalc(t *testing.T) {
	tests := []struct {
		expression string
		want       float64
		wantErr    bool
	}{
		{"1+1", 2, false},
		{"2*2", 4, false},
		{"2+2*2", 6, false},
		{"(2+2)*2", 8, false},
		{"8/2", 5, false},
		{"1/(5-5)", 0, true}, // Деление на 0
		{"", 0, true},        // Пустое выражение
	}

	for _, tt := range tests {
		got, err := Calc(tt.expression)
		if (err != nil) != tt.wantErr {
			t.Errorf("Calc(%q) error = %v, wantErr %v", tt.expression, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("Calc(%q) = %v, want %v", tt.expression, got, tt.want)
		}
	}
}
