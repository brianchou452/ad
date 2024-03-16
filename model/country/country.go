package conutry

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type Countrys []Country
type Country string

const (
	TW Country = "TW"
	JP Country = "JP"
	HK Country = "HK"
	US Country = "US"
)

func (o *Countrys) Scan(src any) error {
	bytes, ok := src.([]byte)
	if !ok {
		// TODO: return error
		return errors.New("src value cannot cast to []byte")
	}
	countries := strings.Split(string(bytes), ",")
	*o = make(Countrys, len(countries))
	for i, country := range countries {
		(*o)[i] = Country(country)
	}
	return nil
}

func (o Countrys) Value() (driver.Value, error) {
	if len(o) == 0 {
		return nil, nil
	}
	countries := make([]string, len(o))
	for i, country := range o {
		countries[i] = string(country)
	}
	return strings.Join(countries, ","), nil
}

func ToCountrys(src []string) Countrys {
	if len(src) == 0 {
		return nil
	}
	r := make(Countrys, len(src))
	for i, country := range src {
		(r)[i] = Country(country)
	}
	return r
}
