package main

import (
	"collector/pkg/report"
	"collector/pkg/slack"
	"log"
	"os"
	"time"
)

func main() {
	token := os.Getenv("TOKEN")
	if len(token) == 0 {
		log.Fatal("empty token")
	}

	channel := os.Getenv("CHANNEL_TEST")
	if len(channel) == 0 {
		log.Fatal("empty channel")
	}

	c := slack.NewSlack(token)

	// <@channel>
	channelMessageRsp, err := c.PostChannelMessage("Mọi người ơi cập nhập công việc đi ạ :licklick:", channel)
	if err != nil {
		log.Fatal(err)
	}

	// waiting
	time.Sleep(time.Duration(120) * time.Second)

	// get report
	threadMessagesRsp, err := c.GetThreadMessages(channel, channelMessageRsp.TS)
	if err != nil {
		log.Fatal(err)
	}

	usersRsp, err := c.GetUsers()
	if err != nil {
		log.Fatal(err)
	}

	// post report
	reportMessage := report.MakeMessage(threadMessagesRsp.Messages, usersRsp.Users)
	reportMessage = "Em xin phép tổng hợp :licklick:\n" + reportMessage
	if _, err = c.PostChannelMessage(reportMessage, channel); err != nil {
		log.Fatal(err)
	}
}
