package commons

import (
	"fmt"
	"math/big"
)

func ToBig(value any) (*big.Int, error) {
	nonce := big.NewInt(0)
	nonceStr, ok := value.(string)
	if !ok {
		nonceFloat, ok := value.(float64)
		if !ok {
			return nil, fmt.Errorf("converting to big error")
		}
		nonce = nonce.SetUint64(uint64(nonceFloat))
	} else {
		nonce, ok = nonce.SetString(nonceStr, 10) // nolint
		if !ok {
			return nil, fmt.Errorf("converting to big error 2")
		}
	}
	return nonce, nil
}

func ToUint64(value interface{}) (uint64, error) {
	b, err := ToBig(value)
	if err != nil {
		return 0, err
	}
	return b.Uint64(), nil
}
