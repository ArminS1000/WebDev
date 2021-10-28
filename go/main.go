package main

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type MyData struct {
	Hash string `gorm:"primary_key"`
	Text string
}

var db *gorm.DB
var err error
var ctx = context.Background()
var rdb *redis.Client

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	dbURL := "host=localhost user=docker dbname=docker sslmode=disable password=docker port=5432"

	db, err = gorm.Open("postgres", dbURL)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("success")
	}

	defer db.Close()

	db.AutoMigrate(&MyData{})

	r := gin.Default()
	r.GET("/sha256", hashFind)
	r.POST("/sha256", hashInsert)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func hashFind(c *gin.Context) {
	hash := c.Query("message")
	if len(hash) != 44 {
		c.JSON(400, gin.H{
			"error": "hash is invalid",
		})
		return
	}
	val, err := rdb.Get(ctx, hash).Result()
	if err == nil {
		c.JSON(200, gin.H{
			"message": val,
		})
		return
	}

	var data MyData
	result := db.Where("Hash = ?", hash).First(&data)
	if result.Error == nil {
		c.JSON(200, gin.H{
			"message": data.Text,
		})
	} else {
		c.JSON(2400, gin.H{
			"error": "hash not found",
		})
	}

}

func hashInsert(c *gin.Context) {
	text := c.Query("message")
	if len(text) < 8 {
		c.JSON(400, gin.H{
			"error": "text has less than 8 chars",
		})
		return
	}
	sha256 := sha256.New()
	sha256.Write([]byte(text))

	hash := base64.URLEncoding.EncodeToString(sha256.Sum(nil))
	var data = &MyData{Text: text, Hash: hash}

	db.Create(&data)
	err := rdb.Set(ctx, hash, text, 0).Err()
	if err != nil {
		panic(err)
	}
	c.JSON(200, gin.H{
		"message": hash,
	})
}
