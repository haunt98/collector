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

	messageType                 = "message"
	collectMessage              = "Update công việc mấy anh ơi :licklick: <!channel>"
	summaryForHumanMessage      = "Em xin tổng hợp công việc :licklick: <!channel>"
	summaryForConfluenceMessage = "Anh nhớ update vào Confluence nhé :licklick:"

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
		ctx.String(http.StatusOK, wrongCommand)
	}
}

func (s *Service) HandleGet(ctx *gin.Context) {
	channel := ctx.Query("channel")
	ts := ctx.Query("ts")

	// If empty thread, get latest thread in channel
	if len(ts) == 0 {
		botMsg := s.loopGetHistoryUntil(channel, maxLoop)
		ts = botMsg.TS
	}

	summary := s.composeThreadSummaryForConfluence(channel, ts)

	ctx.String(http.StatusOK, summary)
}

func (s *Service) handleCollect(ctx *gin.Context, payload slack.CommandPayload) {
	// slack need response as soon as possible
	ctx.String(http.StatusOK, "")

	if err := s.slackService.PostMessageByResponseURL(payload.ResponseURL, slack.MessageRequestByResponseURL{
		MessagePayload: slack.MessagePayload{
			Text:   collectMessage,
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

	// Human
	summaryForHuman := s.composeThreadSummaryForHuman(payload.ChannelID, botMsg.TS)

	if err := s.slackService.PostMessageByResponseURL(payload.ResponseURL, slack.MessageRequestByResponseURL{
		MessagePayload: slack.MessagePayload{
			Text:   "",
			Blocks: summaryForHuman,
		},
		ResponseType: slack.ResponseTypeInChannel,
	}); err != nil {
		log.Fatal("failed to post message by response url", err)
	}

	// Confluence
	summaryForConfluence := s.composeThreadSummaryForConfluence(payload.ChannelID, botMsg.TS)
	summaryForConfluenceExtra := summaryForConfluenceMessage + "\n```\n" + summaryForConfluence + "```"

	if err := s.slackService.PostMessageByResponseURL(payload.ResponseURL, slack.MessageRequestByResponseURL{
		MessagePayload: slack.MessagePayload{
			Text:   summaryForConfluenceExtra,
			Blocks: nil,
		},
		ResponseType: slack.ResponseTypeEphemeral,
	}); err != nil {
		log.Fatal("failed to post message by response url", err)
	}
}

func (s *Service) composeThreadSummaryForHuman(channel, thread string) []interface{} {
	conversationReplies, err := s.slackService.GetConversationsReplies(s.token, channel, thread)
	if err != nil {
		log.Fatal("failed to get conversations replies", err)
	}

	usersList, err := s.slackService.GetUsersList(s.token)
	if err != nil {
		log.Fatal("failed to get users list", err)
	}

	return composeSummaryForHuman(conversationReplies.Messages, usersList.Users)
}

func (s *Service) composeThreadSummaryForConfluence(channel, thread string) string {
	conversationReplies, err := s.slackService.GetConversationsReplies(s.token, channel, thread)
	if err != nil {
		log.Fatal("failed to get conversations replies", err)
	}

	usersList, err := s.slackService.GetUsersList(s.token)
	if err != nil {
		log.Fatal("failed to get users list", err)
	}

	return composeSummaryForConfluence(conversationReplies.Messages, usersList.Users)
}

func (s *Service) loopGetHistoryUntil(channel string, max int) (result slack.Message) {
	var cursor string
	for i := 0; i < max; i += 1 {
		conversationHistory, err := s.slackService.GetConversationsHistory(s.token, channel, cursor)
		if err != nil {
			log.Fatal("failed to get conversation history", err)
		}

		for _, msg := range conversationHistory.Messages {
			if msg.Type == messageType && msg.Text == collectMessage && msg.BotID == s.botID {
				result = msg
				return
			}
		}

		cursor = conversationHistory.ResponseMetadata.NextCursor
	}
	return
}
