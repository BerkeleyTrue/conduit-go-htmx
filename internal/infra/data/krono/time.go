package krono

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type Krono struct {
	Time time.Time
}

func Now() Krono {
	return Krono{Time: time.Now()}
}

func (ts *Krono) IsZero() bool {
	return ts.Time.IsZero()
}

func (ts Krono) String() string {
	return ts.Time.Format(time.RFC3339)
}

func (ts Krono) ToNullString() sql.NullString {
	return sql.NullString{
		String: ts.String(),
		Valid:  !ts.IsZero(),
	}
}

func FromString(s string) (Krono, error) {
	t, err := time.Parse(time.RFC3339, s)

	if err != nil {
		return Krono{}, fmt.Errorf("krono.FromString: unable to parse %s to time.Time", s)
	}

	return Krono{Time: t}, nil
}

// define how to format the Krono to text for sqlite
func (ts Krono) Value() (driver.Value, error) {
	return ts.String(), nil
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

	t, err := FromString(s)

	if err != nil {
		return err
	}

	ts.Time = t.Time

	return nil
}
