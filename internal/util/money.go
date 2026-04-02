package util

import (
	"fmt"
	"math/big"
)

func ParseAmountToPaise(amountStr string) (int64, error) {
	r := new(big.Rat)
	if _, ok := r.SetString(amountStr); !ok {
		return 0, fmt.Errorf("invalid amount string")
	}
	hundred := big.NewRat(100, 1)
	val := new(big.Rat).Mul(r, hundred)

	if !val.IsInt() {
		return 0, fmt.Errorf("amount has smaller denominator than paise")
	}
	res := val.Num().Int64()
	return res, nil
}

func FormatPaiseToAmount(paise int64) string {
	val := big.NewRat(paise, 100)
	return val.FloatString(2)
}
