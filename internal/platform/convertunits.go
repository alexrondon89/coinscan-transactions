package platform

import "math/big"

func ConvertToUnitDesired(unit *big.Int, exponent float64) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(unit), big.NewFloat(exponent))
}
