package main

import (
	"collector/internal/scrum"
	"collector/pkg/slack"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		log.Fatal("PORT is empty")
	}
	log.Printf("PORT is %s\n", port)

	token := os.Getenv("TOKEN")
	if len(token) == 0 {
		log.Fatal("TOKEN is empty")
	}

	botID := os.Getenv("BOT_ID")
	if len(botID) == 0 {
		log.Fatal("BOT_ID is empty")
	}

	slackService := slack.NewService()
	scrumService := scrum.NewService(slackService, token, botID)
	r := gin.Default()

	r.POST("/scrum", scrumService.Handle)
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal(err)
	}
}
