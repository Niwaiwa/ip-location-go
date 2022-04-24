package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/oschwald/maxminddb-golang"
)

func get_ip_location(search_ip string) (country_code, country_name_en, country_name_jp, country_name_cn string) {
	db, err := maxminddb.Open("./test-data/test-data/GeoIP2-Country-Test.mmdb")
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
	// log.Println(record)
	country_code = record.Country.ISOCode
	country_name_en = record.Country.Name.EN
	country_name_jp = record.Country.Name.JA
	country_name_cn = record.Country.Name.CN
	return
}

func setRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		ip := c.GetHeader("X-Real-Ip")
		log.Println(ip)
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
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/search", func(c *gin.Context) {
		ip := c.DefaultQuery("ip", "")
		log.Println(ip)
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
	return r
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	gin_mode := os.Getenv("GIN_MODE")
	gin.SetMode(gin_mode)

	r := setRouter()
	// r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	// need size 1 to get a signal from buffer.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
