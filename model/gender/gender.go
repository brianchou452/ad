package gender

import (
	"database/sql/driver"
	"errors"
	"sort"
	"strings"
)

type Gender string
type Genders []Gender

const (
	M Gender = "M"
	F Gender = "F"
)

func (o *Genders) Scan(src any) error {
	bytes, ok := src.([]byte)
	if !ok {
		// TODO: return error
		return errors.New("src value cannot cast to []byte")
	}
	genders := strings.Split(string(bytes), ",")
	*o = make(Genders, len(genders))
	for i, gender := range genders {
		(*o)[i] = Gender(gender)
	}
	return nil
}

func (o Genders) Value() (driver.Value, error) {
	if len(o) == 0 {
		return nil, nil
	}
	countries := make([]string, len(o))
	for i, country := range o {
		countries[i] = string(country)
	}
	return strings.Join(countries, ","), nil
}

func ToGenders(src []string) Genders {
	if len(src) == 0 {
		return nil
	}
	sort.Strings(src)
	r := make(Genders, len(src))
	for i, gender := range src {
		(r)[i] = Gender(gender)
	}
	return r
}
