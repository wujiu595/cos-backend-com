package flake

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

var (
	_ sql.Scanner   = (*IDs)(nil)
	_ driver.Valuer = (IDs)(nil)
)

// IDs used in db
type IDs []ID

// Value implement driver.Valuer
func (i IDs) Value() (driver.Value, error) {
	return json.Marshal(i)
}

// Scan implement sql.Sacnner
func (i *IDs) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	switch val := src.(type) {
	case string:
		if len(val) > 0 {
			return json.Unmarshal([]byte(val), i)
		}

	case []byte:
		if len(val) > 0 {
			return json.Unmarshal(val, i)
		}

	default:
		return fmt.Errorf("unexpected src type %T for flake.IDs", src)
	}

	return nil
}
