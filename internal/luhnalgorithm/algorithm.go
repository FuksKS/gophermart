package luhnalgorithm

import (
	"gophermart/internal/model"
	"strconv"
)

// LuhnCheck выполняет проверку по алгоритму Луна для произвольной длины строки цифр.
func LuhnCheck(number string) (bool, error) {
	var sum int
	var alternate bool

	// Идем по цифрам справа налево
	for i := len(number) - 1; i >= 0; i-- {
		// Преобразуем символ в число
		n, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false, model.ErrNotANumber
		}

		if alternate {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}

		sum += n
		alternate = !alternate
	}

	return sum%10 == 0, nil
}
