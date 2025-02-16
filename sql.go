package timex

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Scan implements the sql.Scanner interface.
func (d *Date) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case []byte:
		*d, err = ParseDate(RFC3339Date, string(v))
	case string:
		*d, err = ParseDate(RFC3339Date, v)
	case time.Time:
		*d = DateFromTime(v)
	default:
		err = fmt.Errorf("unsupported type %T", value)
	}
	return err
}

// Value implements the driver.Valuer interface.
func (d Date) Value() (driver.Value, error) {
	return d.Format(RFC3339Date), nil
}

// NullDate represents a specific day in Gregorian calendar that may be null.
// NullDate implements the sql.Scanner interface, so it can be used as a scan destination, similar to sql.NullString.
type NullDate struct {
	Date  Date
	Valid bool // Valid is true if Date is not NULL
}

// Scan implements the sql.Scanner interface.
func (d *NullDate) Scan(value interface{}) error {
	if value == nil {
		d.Date, d.Valid = Date{}, false
		return nil
	}
	d.Valid = true
	return d.Date.Scan(value)
}

// Value implements the driver.Valuer interface.
func (d NullDate) Value() (driver.Value, error) {
	if !d.Valid {
		return nil, nil
	}
	return d.Date.Value()
}
