package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"ad/model"
	"os"
)

type GormDatabase struct {
	DB *gorm.DB
}

func New() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	charset := os.Getenv("DB_CHARSET")

	dsn := user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=" + charset + "&parseTime=True&loc=Local"
	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.Ad{})

	return db, nil
}

func NewTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.Ad{})

	return db, nil
}
