package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/oschwald/maxminddb-golang"
)

func get_ip_location(search_ip string) (country_code, country_name_en, country_name_jp, country_name_cn string) {
	db, err := maxminddb.Open("./test_data/test_data/GeoIP2-Country-Test.mmdb")
	// db, err := maxminddb.Open("./static/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ip := net.ParseIP(search_ip)

	var record struct {
		Country struct {
			ISOCode string `maxminddb:"iso_code"`
			Name    struct {
				EN string `maxminddb:"en"`
				JA string `maxminddb:"ja"`
				CN string `maxminddb:"zh-CN"`
			} `maxminddb:"names"`
		} `maxminddb:"country"`
	} // Or any appropriate struct

	err = db.Lookup(ip, &record)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(record)
	country_code = record.Country.ISOCode
	country_name_en = record.Country.Name.EN
	country_name_jp = record.Country.Name.JA
	country_name_cn = record.Country.Name.CN
	return
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	gin_mode := os.Getenv("GIN_MODE")
	gin.SetMode(gin_mode)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/search", func(c *gin.Context) {
		ip := c.DefaultQuery("ip", "")
		if ip == "" {
			c.JSON(200, gin.H{})
		} else {
			country_code, country_name_en, country_name_jp, country_name_cn := get_ip_location(ip)
			c.JSON(200, gin.H{
				"ip":              ip,
				"country_code":    country_code,
				"country_name_en": country_name_en,
				"country_name_jp": country_name_jp,
				"country_name_cn": country_name_cn,
			})
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
