package slack

import (
	"bytes"
	"collector/pkg/httpwrap"
	"encoding/json"
	"net/http"
)

const (
	baseURL                  = "https://slack.com/api"
	conversationsHistoryPath = "/conversations.history"
	conversationsRepliesPath = "/conversations.replies"
	usersListPath            = "/users.list"
)

type Service struct {
	client *http.Client
}

func NewService() *Service {
	return &Service{
		client: &http.Client{},
	}
}

// https://api.slack.com/methods/conversations.history
func (s *Service) GetConversationsHistory(token, channel, cursor string) (result MessagesResponse, err error) {
	var urlWithParams string
	urlWithParams, err = httpwrap.AddParams(baseURL+conversationsHistoryPath,
		httpwrap.Param{
			Name:  "token",
			Value: token,
		},
		httpwrap.Param{
			Name:  "channel",
			Value: channel,
		},
		httpwrap.Param{
			Name:  "cursor",
			Value: cursor,
		},
	)
	if err != nil {
		return
	}

	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, urlWithParams, nil)
	if err != nil {
		return
	}

	err = httpwrap.DoRequestWithResult(s.client, req, &result)
	return
}

// https://api.slack.com/methods/conversations.replies
func (s *Service) GetConversationsReplies(token, channel, threadTS string) (result MessagesResponse, err error) {
	var urlWithParams string
	urlWithParams, err = httpwrap.AddParams(baseURL+conversationsRepliesPath,
		httpwrap.Param{
			Name:  "token",
			Value: token,
		},
		httpwrap.Param{
			Name:  "channel",
			Value: channel,
		},
		httpwrap.Param{
			Name:  "ts",
			Value: threadTS,
		},
	)
	if err != nil {
		return
	}

	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, urlWithParams, nil)
	if err != nil {
		return
	}

	err = httpwrap.DoRequestWithResult(s.client, req, &result)
	return
}

// https://api.slack.com/methods/users.list
func (s *Service) GetUsersList(token string) (result UsersResponse, err error) {
	var urlWithParams string
	urlWithParams, err = httpwrap.AddParams(baseURL+usersListPath,
		httpwrap.Param{
			Name:  "token",
			Value: token,
		},
	)
	if err != nil {
		return
	}

	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, urlWithParams, nil)
	if err != nil {
		return
	}

	err = httpwrap.DoRequestWithResult(s.client, req, &result)
	return
}

// https://api.slack.com/interactivity/handling#message_responses
func (s *Service) PostMessageByResponseURL(responseURL string, msgReq MessageRequestByResponseURL) error {
	body, err := json.Marshal(msgReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, responseURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	_, err = s.client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
