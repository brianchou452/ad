package api

import (
	"ad/model"
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

func (e *AdminEnv) PostAdminAPIController(c *gin.Context) {
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

	err = e.DB.CreateAd(&adData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
