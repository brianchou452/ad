package platform

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type Platforms []Platform
type Platform string

const (
	WEB     Platform = "web"
	ANDROID Platform = "android"
	IOS     Platform = "ios"
)

func (o *Platforms) Scan(src any) error {
	bytes, ok := src.([]byte)
	if !ok {
		// TODO: return error
		return errors.New("src value cannot cast to []byte")
	}
	platforms := strings.Split(string(bytes), ",")
	*o = make(Platforms, len(platforms))
	for i, platform := range platforms {
		(*o)[i] = Platform(platform)
	}
	return nil
}

func (o Platforms) Value() (driver.Value, error) {
	if len(o) == 0 {
		return nil, nil
	}
	countries := make([]string, len(o))
	for i, country := range o {
		countries[i] = string(country)
	}
	return strings.Join(countries, ","), nil
}

func ToPlatforms(src []string) Platforms {
	if len(src) == 0 {
		return nil
	}
	r := make(Platforms, len(src))
	for i, platform := range src {
		(r)[i] = Platform(platform)
	}
	return r
}
