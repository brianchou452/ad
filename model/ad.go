package model

import (
	gender "ad/model/gender"
	platform "ad/model/platform"
	"time"
)

type Ad struct {
	Title     string     `bson:"title" json:"title"`
	StartAt   time.Time  `bson:"startAt" json:"startAt"`
	EndAt     time.Time  `bson:"endAt" json:"endAt"`
	Condition Conditions `bson:"conditions" json:"conditions"`
}

type Conditions struct {
	Gender    gender.Genders     `json:"gender"`
	AgeStart  uint8              `json:"ageStart"`
	AgeEnd    uint8              `json:"ageEnd"`
	Countries []string           `json:"country"`
	Platform  platform.Platforms `json:"platform"`
	Age       []int              `json:"age"`
}

type AdResponse struct {
	Title string    `json:"title"`
	EndAt time.Time `json:"endAt"`
}

type AdSet struct {
	Id    string    `bson:"_id"`
	EndAt time.Time `bson:"endAt"`
}
