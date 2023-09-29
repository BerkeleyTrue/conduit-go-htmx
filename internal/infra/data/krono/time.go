package krono

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Krono struct {
	Time time.Time
}

// define how to format the Krono to text for sqlite
func (ts Krono) Value() (driver.Value, error) {
	return ts.Time.Format(time.RFC3339), nil
}

// define a scan method for sql to use
func (ts *Krono) Scan(v interface{}) error {
	if v == nil {
		ts.Time = time.Time{}
		return nil
	}

	// first cast to string type
	s, ok := v.(string)

	if !ok {
		return fmt.Errorf("krono.Scan: unable to cast %v to string", v)
	}

	// now attempt to parse the string to time.Time
	t, err := time.Parse(time.RFC3339, s)

	if err != nil {
		return fmt.Errorf("krono.Scan: unable to parse %s to time.Time", s)
	}

	ts.Time = t

	return nil
}
