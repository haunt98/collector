package scrum

import (
	"collector/pkg/slack"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service struct {
	slackService *slack.Service

	token, botID string
}

func NewService(slackService *slack.Service, token, botID string) *Service {
	return &Service{
		slackService: slackService,
		token:        token,
		botID:        botID,
	}
}

const (
	collectCommand = "collect"
	summaryCommand = "summary"
	wrongCommand   = "Sai câu lệnh rồi anh ơi"

	message        = "message"
	collectMessage = "Update công việc tại đây nha mấy anh ơi :licklick: <!channel>"
	summaryMessage = "Em xin tổng hợp nhẹ :licklick:, copy rồi quăng qua wiki confluence nha mấy anh :licklick: <!channel>\n"

	responseInChannel = "in_channel"

	maxLoop = 2
)

func (s *Service) HandlerPost(ctx *gin.Context) {
	var payload slack.CommandPayload
	if err := ctx.Bind(&payload); err != nil {
		log.Fatal(err)
	}

	switch payload.Text {
	case collectCommand:
		s.collect(ctx, payload)
	case summaryCommand:
		s.summary(ctx, payload)
	default:
		ctx.String(http.StatusOK, wrongCommand)
	}
}

func (s *Service) HandlerGet(ctx *gin.Context) {
	channel := ctx.Query("channel")
	ts := ctx.Query("ts")

	if len(ts) == 0 {
		botMsg := s.loopGetHistoryUntil(channel, maxLoop)
		ts = botMsg.TS
	}

	summary := s.makeSummary(channel, ts)

	ctx.String(http.StatusOK, summary)
}

func (s *Service) collect(ctx *gin.Context, payload slack.CommandPayload) {
	ctx.String(http.StatusOK, "")

	if err := s.slackService.PostThreadMessageByWebhook(payload.ResponseURL, collectMessage, "in_channel"); err != nil {
		log.Fatal(err)
	}
}

func (s *Service) summary(ctx *gin.Context, payload slack.CommandPayload) {
	ctx.String(http.StatusOK, "")

	botMsg := s.loopGetHistoryUntil(payload.ChannelID, maxLoop)
	summary := s.makeSummary(payload.ChannelID, botMsg.TS)
	summaryExtra := summaryMessage + "```\n" + summary + "```"

	if err := s.slackService.PostThreadMessageByWebhook(payload.ResponseURL, summaryExtra, responseInChannel); err != nil {
		log.Fatal(err)
	}
}

func (s *Service) makeSummary(channel, thread string) string {
	conversationReplies, err := s.slackService.GetConversationReplies(s.token, channel, thread)
	if err != nil {
		log.Fatal(err)
	}

	usersList, err := s.slackService.GetUsersList(s.token)
	if err != nil {
		log.Fatal(err)
	}

	return makeConfluenceSummary(conversationReplies.Messages, usersList.Users)
}

func (s *Service) loopGetHistoryUntil(channel string, max int) (result slack.Message) {
	var cursor string
	for i := 0; i < max; i += 1 {
		conversationHistory, err := s.slackService.GetConversationHistory(s.token, channel, cursor)
		if err != nil {
			log.Fatal(err)
		}

		for _, msg := range conversationHistory.Messages {
			if msg.Type == message && msg.Text == collectMessage && msg.BotID == s.botID {
				result = msg
				return
			}
		}

		cursor = conversationHistory.ResponseMetadata.NextCursor
	}
	return
}
