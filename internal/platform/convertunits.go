package platform

import (
	err "errors"
	"math/big"

	"github.com/alexrondon89/coinscan-transactions/internal/platform/errors"
)

func ConvertToUnitDesired(unit interface{}, exponent float64) (*big.Float, errors.Error) {
	switch unit.(type) {
	case uint64:
		return new(big.Float).Quo(new(big.Float).SetUint64(unit.(uint64)), big.NewFloat(exponent)), nil
	case *big.Int:
		return new(big.Float).Quo(new(big.Float).SetInt(unit.(*big.Int)), big.NewFloat(exponent)), nil
	case float64:
		return new(big.Float).Quo(new(big.Float).SetFloat64(unit.(float64)), big.NewFloat(exponent)), nil
	case string:
		floatSetted, err := new(big.Float).SetString(unit.(string))
		if !err {
			break
		}
		return new(big.Float).Quo(floatSetted, big.NewFloat(exponent)), nil
	}

	return nil, errors.NewError(errors.ConversionUnitError, err.New("error converting unit"))
}
