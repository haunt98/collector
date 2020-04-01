package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// PORT
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("PORT is empty")
		return
	}
	log.Printf("PORT is %s\n", port)

	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ping")
	})
	r.POST("/collect", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "collect")
	})
	r.POST("/summary", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "summary")
	})

	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal(err)
	}
}
