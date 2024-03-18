package model

import (
	country "ad/model/country"
	gender "ad/model/gender"
	platform "ad/model/platform"
	"time"
)

type Ad struct {
	ID        int64
	Title     string    `gorm:"type:varchar(180);not null"`
	StartAt   time.Time `gorm:"type:DATETIME;not null;index:time_range"`
	EndAt     time.Time `gorm:"type:DATETIME;not null;index:time_range"`
	Age       Age
	Countries []Country
	Platforms []Platform
	Genders   []Gender
}

type Age struct {
	ID       int64
	AgeStart uint8 `gorm:"index:age"`
	AgeEnd   uint8 `gorm:"index:age"`
	AdID     int64
}

type Country struct {
	ID      int64
	Country country.Country `json:"country" gorm:"type:VARCHAR(255);index"`
	AdID    int64
}

type Platform struct {
	ID       int64
	Platform platform.Platform `json:"platform" gorm:"type:VARCHAR(255);index"`
	AdID     int64
}

type Gender struct {
	ID     int64
	Gender gender.Gender `json:"gender" gorm:"type:varchar(180);index"`
	AdID   int64
}
