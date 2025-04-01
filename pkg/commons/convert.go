package commons

import (
	"database/sql/driver"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// Custom type to handle conversion
type Address common.Address

// Scan implements the sql.Scanner interface for reading from DB
func (a *Address) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Address: expected []byte but got %T", value)
	}
	*a = Address(common.BytesToAddress(bytes))
	return nil
}

// Value implements the driver.Valuer interface for writing to DB
func (a Address) Value() (driver.Value, error) {
	return common.Address(a).Bytes(), nil
}

// A struct using the custom Address type
type A struct {
	Address Address `db:"address"`
}
