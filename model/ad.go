package model

import (
	country "ad/model/country"
	gender "ad/model/gender"
	platform "ad/model/platform"
	"time"
)

type Ad struct {
	Title     string    `bson:"title" json:"title"`
	StartAt   time.Time `bson:"startAt" json:"startAt"`
	EndAt     time.Time `bson:"endAt" json:"endAt"`
	Condition Conditon  `bson:"conditions" json:"conditions"`
}

type Conditon struct {
	Gender   gender.Genders     `json:"gender" gorm:"type:varchar(180);index"`
	AgeStart uint8              `json:"ageStart" gorm:"index:age"`
	AgeEnd   uint8              `json:"ageEnd" gorm:"index:age"`
	Country  country.Countrys   `json:"country" gorm:"type:VARCHAR(255);index"`
	Platform platform.Platforms `json:"platform" gorm:"type:VARCHAR(255);index"`
	Age      []int              `json:"age" gorm:"-"`
}
