package main

import (
	"collector/pkg/report"
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

	s := slack.NewSlack(token)

	r := gin.Default()
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ping")
	})
	r.POST("/", func(ctx *gin.Context) {
		var payload CommandPayload
		if err := ctx.Bind(&payload); err != nil {
			log.Fatal(err)
		}
		log.Printf("CommandPayload: %+v\n", payload)

		switch payload.Text {
		case "collect":
			if err := s.PostThreadMessageByWebhook(payload.ResponseURL, "update please", "in_channel"); err != nil {
				log.Fatal(err)
			}
		case "summary":
			ctx.String(http.StatusOK, "got it wait it")

			history, err := s.GetChannelHistory(payload.ChannelID)
			if err != nil {
				log.Fatal(err)
			}

			var botMsg slack.Message
			for _, msg := range history.Messages {
				if msg.Type == "message" && msg.Text == "update please" && msg.BotID == botID {
					botMsg = msg
					break
				}
			}

			threadMessages, err := s.GetThreadMessages(payload.ChannelID, botMsg.TS)
			if err != nil {
				log.Fatal(err)
			}

			users, err := s.GetUsers()
			if err != nil {
				log.Fatal(err)
			}

			reportMsg := report.MakeMessage(threadMessages.Messages, users.Users)
			if err := s.PostThreadMessageByWebhook(payload.ResponseURL, reportMsg, "in_channel"); err != nil {
				log.Fatal(err)
			}
		default:
			ctx.String(http.StatusOK, "wrong hole")
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
