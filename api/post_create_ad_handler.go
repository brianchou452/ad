package api

import (
	"ad/model"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostData struct {
	Title     string           `json:"title" binding:"required"`
	StratAt   string           `json:"startAt" binding:"required"`
	EndAt     string           `json:"endAt" binding:"required"`
	Conditons model.Conditions `json:"conditions"`
}

func (e *Env) CreateAd(c *gin.Context) {
	var data PostData

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startAt, err := time.Parse(time.RFC3339, data.StratAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	endAt, err := time.Parse(time.RFC3339, data.EndAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if startAt.After(endAt) {
		// TODO: 改用統一的方法回傳錯誤、並提供錯誤代碼
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

	condition := model.Conditions{
		Gender:    data.Conditons.Gender,
		AgeStart:  data.Conditons.AgeStart,
		AgeEnd:    data.Conditons.AgeEnd,
		Countries: data.Conditons.Countries,
		Platform:  data.Conditons.Platform,
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		go e.Redis.AddAd(adData, id.Hex())
	} else {
		log.Println("inserted id is not primitive.ObjectID")
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
