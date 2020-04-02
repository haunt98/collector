package main

import (
	"collector/pkg/slack"
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
	r.POST("/", func(ctx *gin.Context) {
		var commandPayload CommandPayload
		if err := ctx.Bind(&commandPayload); err != nil {
			log.Fatal(err)
		}
		log.Println(commandPayload)

		ctx.String(http.StatusOK, "<3")

		s := slack.NewSlack("")
		if err := s.PostThreadMessageByWebhook(commandPayload.ResponseURL, "another <3", "in_channel"); err != nil {
			log.Fatal(err)
		}
	})

	if err := r.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal(err)
	}
}

type CommandPayload struct {
	Command     string `form:"command"`
	Text        string `form:"text"`
	ResponseURL string `form:"response_url"`
	UserID      string `form:"user_id"`
	ChannelID   string `form:"channel_id"`
}
