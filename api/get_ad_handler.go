package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Query struct {
	Offset   int64  `form:"offset,default=0" binding:"gte=0"`
	Limit    int64  `form:"limit,default=5" binding:"gte=1,lte=100"`
	Age      int    `form:"age,default=0" binding:"omitempty,gte=0,lte=100"`
	Gender   string `form:"gender" binding:"omitempty,oneof=M F"`
	Country  string `form:"country" binding:"omitempty,len=2,iso3166_1_alpha2"`
	Platform string `form:"platform" binding:"omitempty,oneof=android ios web"`
}

func (e *Env) GetAds(c *gin.Context) {
	var query Query

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if query.Limit == 0 {
		query.Limit = 5
	}

	result, err := e.Redis.GetAdsFromCondition(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": result,
	})
}
