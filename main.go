package main

import (
	"ad/api"
	"ad/database"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	env := &api.AdminEnv{DB: &database.GormDatabase{DB: db}}
	r := gin.Default()
	r.RedirectFixedPath = true

	r.POST("/api/v1/ad", env.PostAdminAPIController)
	r.GET("/api/v1/ad", env.GetAdController)

	r.Run(":80")
}
