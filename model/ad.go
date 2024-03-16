package model

import (
	country "ad/model/country"
	gender "ad/model/gender"
	platform "ad/model/platform"
	"time"
)

type Ad struct {
	UUID      string    `gorm:"primaryKey;default:(UUID_TO_BIN(UUID(),true));type:BINARY(16)"`
	Title     string    `gorm:"type:varchar(180);not null"`
	StartAt   time.Time `gorm:"type:DATETIME;not null;index:tima_range"`
	EndAt     time.Time `gorm:"type:DATETIME;not null;index:tima_range"`
	Condition Conditon  `gorm:"embedded"`
}

type Conditon struct {
	Gender   gender.Genders     `json:"gender" gorm:"type:varchar(180);index"`
	AgeStart uint8              `json:"ageStart" gorm:"index:age"`
	AgeEnd   uint8              `json:"ageEnd" gorm:"index:age"`
	Country  country.Countrys   `json:"country" gorm:"type:VARCHAR(255);index"`
	Platform platform.Platforms `json:"platform" gorm:"type:VARCHAR(255);index"`
}
