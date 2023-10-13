package krono

import (
	"database/sql"
	"fmt"
	"time"
)

type Krono struct {
	Time time.Time
}

// create a new Krono of the current time
func Now() Krono {
	return Krono{Time: time.Now()}
}

func (ts *Krono) IsZero() bool {
	return ts.Time.IsZero()
}

func (ts Krono) ToString() string {
	return ts.Time.Format(time.RFC3339)
}

func (ts Krono) ToNullString() sql.NullString {
	return sql.NullString{
		String: ts.ToString(),
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

func FromNullString(s sql.NullString) (Krono, error) {
	if !s.Valid {
		return Krono{}, nil
	}

	return FromString(s.String)
}

// Krono implements fmt.Stringer
func (ts Krono) String() string {
	if ts.IsZero() {
		return ""
	}

	return ts.Time.Format(time.RFC3339)
}
