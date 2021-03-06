package postgres

import (
	"fmt"
	"database/sql/driver"
	"errors"
	"github.com/hhh0pE/decimal"
)

const (
	MaxIntegralDigits   = 131072 // max digits before the decimal point
	MaxFractionalDigits = 16383  // max digits after the decimal point
)

// LengthError is returned from Decimal.Value when either its integral (digits
// before the decimal point) or fractional (digits after the decimal point)
// parts are too long for PostgresSQL.
type LengthError struct {
	Part string // "integral" or "fractional"
	N    int    // length of invalid part
	max  int
}

func (e LengthError) Error() string {
	return fmt.Sprintf("%s (%d digits) is too long (%d max)", e.Part, e.N, e.max)
}

type Decimal struct {
	decimal.Decimal
}


// Value implements driver.Valuer.
func (d *Decimal) Value() (driver.Value, error) {
	if d == nil {
		return nil, nil
	}
	if d.IsNaN(0) {
		return "NaN", nil
	}
	if d.IsInf(0) {
		return nil, errors.New("Decimal.Value: DECIMAL does not accept Infinities")
	}

	dl := d.Precision()  // length of d
	sl := int(d.Scale()) // length of fractional part

	if il := dl - sl; il > MaxIntegralDigits {
		return nil, &LengthError{Part: "integral", N: il, max: MaxIntegralDigits}
	}
	if sl > MaxFractionalDigits {
		return nil, &LengthError{Part: "fractional", N: sl, max: MaxFractionalDigits}
	}
	return d.String(), nil
}

// Scan implements sql.Scanner.
func (d *Decimal) Scan(val interface{}) error {
	if d == nil {
		d.Decimal = decimal.NewDecimal()
	}
	switch t := val.(type) {
	case string:
		if _, ok := d.SetString(t); !ok {
			if err := d.Context.Err(); err != nil {
				return err
			}
			return fmt.Errorf("Decimal.Scan: invalid syntax: %q", t)
		}
		return nil
	case []byte:
		return d.UnmarshalText(t)
	default:
		return fmt.Errorf("Decimal.Scan: unknown value: %#v", val)
	}
}
