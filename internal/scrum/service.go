package scrum

import (
	"collector/pkg/clock"
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

	wrongMessage   = "Sai câu lệnh rồi anh ơi"
	collectMessage = "Update công việc mấy anh ơi " + slack.MentionChannel

	maxLoop = 2
)

func (s *Service) HandlePost(ctx *gin.Context) {
	var payload slack.CommandPayload
	if err := ctx.Bind(&payload); err != nil {
		log.Fatal("failed to bind", err)
	}

	switch payload.Text {
	case collectCommand:
		s.handleCollect(ctx, payload)
	case summaryCommand:
		s.handleSummary(ctx, payload)
	default:
		ctx.String(http.StatusOK, wrongMessage)
	}
}

func (s *Service) handleCollect(ctx *gin.Context, payload slack.CommandPayload) {
	// slack need response as soon as possible
	ctx.String(http.StatusOK, "")
	date, err := clock.NowDateInSaiGon()
	if err != nil {
		log.Fatal("failed to get date", err)
	}
	collectMessageWithDate := collectMessage + " " + date

	if err := s.slackService.PostMessageByResponseURL(payload.ResponseURL, slack.MessageRequestByResponseURL{
		MessagePayload: slack.MessagePayload{
			Text:   collectMessageWithDate,
			Blocks: nil,
		},
		ResponseType: slack.ResponseTypeInChannel,
	}); err != nil {
		log.Fatal("failed to post message by response url", err)
	}
}

func (s *Service) handleSummary(ctx *gin.Context, payload slack.CommandPayload) {
	// slack need response as soon as possible
	ctx.String(http.StatusOK, "")

	botMsg := s.loopGetHistoryUntil(payload.ChannelID, maxLoop)

	humanSummary, confluenceSummary, err := s.composeThreadSummary(payload.ChannelID, botMsg.TS)
	if err != nil {
		log.Fatal("failed to compose thread summary", err)
	}

	// human
	if err := s.slackService.PostMessageByResponseURL(payload.ResponseURL, slack.MessageRequestByResponseURL{
		MessagePayload: slack.MessagePayload{
			Text:   humanMessageIntro,
			Blocks: humanSummary,
		},
		ResponseType: slack.ResponseTypeInChannel,
	}); err != nil {
		log.Fatal("failed to post message by response url", err)
	}

	// confluence
	if err := s.slackService.PostMessageByResponseURL(payload.ResponseURL, slack.MessageRequestByResponseURL{
		MessagePayload: slack.MessagePayload{
			Text:   confluenceSummary,
			Blocks: nil,
		},
		ResponseType: slack.ResponseTypeInChannel,
	}); err != nil {
		log.Fatal("failed to post message by response url", err)
	}
}

func (s *Service) composeThreadSummary(channel, thread string) (
	humanSummary []interface{}, confluenceSummary string, err error) {
	var conversationReplies slack.MessagesResponse
	conversationReplies, err = s.slackService.GetConversationsReplies(s.token, channel, thread)
	if err != nil {
		// log.Fatal("failed to get conversations replies", err)
		return
	}

	var usersList slack.UsersResponse
	usersList, err = s.slackService.GetUsersList(s.token)
	if err != nil {
		// log.Fatal("failed to get users list", err)
		return
	}

	humanSummary, confluenceSummary = composeSummary(conversationReplies.Messages, usersList.Users)
	return
}

func (s *Service) loopGetHistoryUntil(channel string, max int) (result slack.Message) {
	var cursor string
	for i := 0; i < max; i += 1 {
		conversationHistory, err := s.slackService.GetConversationsHistory(s.token, channel, cursor)
		if err != nil {
			log.Fatal("failed to get conversation history", err)
		}

		for _, msg := range conversationHistory.Messages {
			if msg.Type == slack.TypeMessage && msg.Text == collectMessage && msg.BotID == s.botID {
				result = msg
				return
			}
		}

		cursor = conversationHistory.ResponseMetadata.NextCursor
	}
	return
}
