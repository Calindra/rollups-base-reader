package commons

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"
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

type ToHex interface {
	Hex() string
}

func PrepareKeyArguments(args ...any) ([]string, []any) {
	placeholders := make([]string, len(args))
	output := make([]any, len(args))
	for i, arg := range args {
		// Use reflect to check if it's a null pointer
		metadata := reflect.ValueOf(arg)
		isNilPtr := metadata.Kind() == reflect.Pointer && metadata.IsNil()

		// Check if the argument is nil or a nil pointer
		if arg == nil || isNilPtr {
			output[i] = nil
			placeholders[i] = fmt.Sprintf("$%d", i+1)
			continue
		}

		switch v := arg.(type) {
		case ToHex:
			hex := strings.TrimPrefix(v.Hex(), "0x")
			output[i] = hex
			placeholders[i] = fmt.Sprintf("decode($%d, 'hex')", i+1)
		default:
			output[i] = arg
			placeholders[i] = fmt.Sprintf("$%d", i+1)
		}
	}
	return placeholders, output
}
