package luhnalgorithm

import (
	"errors"
	"gophermart/internal/model"
	"testing"
)

func TestLuhnCheck(t *testing.T) {
	type want struct {
		isCorrect bool
		err       error
	}

	tests := []struct {
		name    string
		orderID string
		want    want
	}{
		{
			name:    "correct number",
			orderID: "49927398716",
			want: want{
				isCorrect: true,
				err:       nil,
			},
		},
		{
			name:    "incorrect number",
			orderID: "49927398715",
			want: want{
				isCorrect: false,
				err:       nil,
			},
		},
		{
			name:    "not a number",
			orderID: "49927398716a",
			want: want{
				isCorrect: false,
				err:       model.ErrNotANumber,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isCorrect, err := LuhnCheck(tt.orderID)
			if !errors.Is(err, tt.want.err) {
				t.Errorf("LuhnCheck() error = %v, wantErr %v", err, tt.want.err)
				return
			}
			if isCorrect != tt.want.isCorrect {
				t.Errorf("LuhnCheck() isCorrect = %v, want %v", isCorrect, tt.want.isCorrect)
			}
		})
	}
}
