package api

import (
	"ad/model"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// type Conditons struct {
// 	AgeStart uint8    `json:"ageStart"`
// 	AgeEnd   uint8    `json:"ageEnd"`
// 	Gender   []string `json:"gender"`
// 	Country  []string `json:"country"`
// 	Platform []string `json:"platform"`
// }

type PostData struct {
	Title     string         `json:"title" binding:"required"`
	StratAt   string         `json:"startAt" binding:"required"`
	EndAt     string         `json:"endAt" binding:"required"`
	Conditons model.Conditon `json:"conditions"`
}

func (e *Env) CreateAd(c *gin.Context) {
	var data PostData

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startAt, err := time.Parse(time.RFC3339, data.StratAt)
	if err != nil {
		log.Println(err)
	}

	endAt, err := time.Parse(time.RFC3339, data.EndAt)
	if err != nil {
		log.Println(err)
	}

	if startAt.After(endAt) {
		// TODO: error handling
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "startAt is after endAt",
		})
		return
	}

	if data.Conditons.AgeEnd != 0 && data.Conditons.AgeStart != 0 {
		if data.Conditons.AgeStart > data.Conditons.AgeEnd {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "ageStart is after ageEnd",
			})
			return
		}
	}

	condition := model.Conditon{
		Gender:   data.Conditons.Gender,
		AgeStart: data.Conditons.AgeStart,
		AgeEnd:   data.Conditons.AgeEnd,
		Country:  data.Conditons.Country,
		Platform: data.Conditons.Platform,
	}

	adData := model.Ad{
		Title:     data.Title,
		StartAt:   startAt,
		EndAt:     endAt,
		Condition: condition,
	}

	if data.Conditons.AgeEnd == 0 || data.Conditons.AgeStart == 0 {
		adData.Condition.Age = []int{}
	} else {
		ageArray := makeRange(int(data.Conditons.AgeStart), int(data.Conditons.AgeEnd))
		adData.Condition.Age = ageArray
	}

	result, err := e.DB.CreateAd(&adData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = e.Redis.AddAdToCache(fmt.Sprint(result.InsertedID), &adData)
	log.Println("insertedID", fmt.Sprint(result.InsertedID))
	if err != nil {
		log.Println(err)
	}

	// ad, err := e.Redis.GetAdFromCache(fmt.Sprint(result.InsertedID))
	// if err != nil {
	// 	log.Println(err)
	// }
	// log.Println(ad)

	for _, country := range adData.Condition.Country {
		err = e.Redis.AddAdToZSet("country", fmt.Sprint(country), fmt.Sprint(result.InsertedID), float64(adData.EndAt.Unix()))
		if err != nil {
			log.Println(err)
		}
	}

	for _, platform := range adData.Condition.Platform {
		err = e.Redis.AddAdToZSet("platform", fmt.Sprint(platform), fmt.Sprint(result.InsertedID), float64(adData.EndAt.Unix()))
		if err != nil {
			log.Println(err)
		}
	}

	for _, gender := range adData.Condition.Gender {
		err = e.Redis.AddAdToZSet("gender", fmt.Sprint(gender), fmt.Sprint(result.InsertedID), float64(adData.EndAt.Unix()))
		if err != nil {
			log.Println(err)
		}
	}

	for _, age := range adData.Condition.Age {
		err = e.Redis.AddAdToZSet("age", fmt.Sprint(age), fmt.Sprint(result.InsertedID), float64(adData.EndAt.Unix()))
		if err != nil {
			log.Println(err)
		}
	}

	if len(adData.Condition.Age) == 0 {
		err = e.Redis.AddAdToZSet("age", "All", fmt.Sprint(result.InsertedID), float64(adData.EndAt.Unix()))
		if err != nil {
			log.Println(err)
		}
	}

	if adData.Condition.Country == nil {
		err = e.Redis.AddAdToZSet("country", "All", fmt.Sprint(result.InsertedID), float64(adData.EndAt.Unix()))
		if err != nil {
			log.Println(err)
		}
	}

	if adData.Condition.Platform == nil {
		err = e.Redis.AddAdToZSet("platform", "All", fmt.Sprint(result.InsertedID), float64(adData.EndAt.Unix()))
		if err != nil {
			log.Println(err)
		}
	}

	if adData.Condition.Gender == nil {
		err = e.Redis.AddAdToZSet("gender", "All", fmt.Sprint(result.InsertedID), float64(adData.EndAt.Unix()))
		if err != nil {
			log.Println(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id": result.InsertedID,
	})
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
