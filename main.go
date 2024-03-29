package main

import (
	"ad/api"
	"ad/database"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/go-redis/redis/v8"

	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
)

func main() {
	version, isEnvSet := os.LookupEnv("APP_VERSION")
	if !isEnvSet {
		err := godotenv.Load(".env")
		if err != nil {
			// TODO: handle error
			log.Fatalf("Error loading .env file")
		}
	} else {
		log.Printf("version: %s", version)
	}

	db, err := database.New()
	if err != nil {
		// TODO: handle error
		log.Fatalf("Error database.New()")
	}

	redisStore := persist.NewRedisStore(redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "redis:6379",
	}))

	env := &api.AdminEnv{DB: &database.MongoDB{
		DB:                    db,
		AdCollections:         db.Database("dcard_ads").Collection("ads"),
		CurrentAdsCollections: db.Database("dcard_ads").Collection("current_ads"),
	}}

	r := gin.Default()
	r.RedirectFixedPath = true

	r.POST("/api/v1/ad", env.PostAdminAPIController)
	r.GET("/api/v1/ad",
		cache.CacheByRequestURI(redisStore, 20*time.Second),
		env.GetAdController)

	r.Run(":80")
}
