package utils

import (
	"fmt"
	"strconv"
)

func Decimal(val float64, dot int) float64 {
	format := fmt.Sprintf("%%.%df", dot)
	ret, _ := strconv.ParseFloat(fmt.Sprintf(format, val), 64)
	return ret
}