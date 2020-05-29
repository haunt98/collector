package slack

import (
	"bytes"
	"collector/pkg/httpwrap"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		})

	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, urlWithParams, nil)
	if err != nil {
		return
	}

	err = httpwrap.DoRequest(s.client, req, &result)
	return
}

// https://api.slack.com/methods/conversations.replies
func (s *Service) GetConversationsReplies(token, channel, threadTS string) (result MessagesResponse, err error) {
	url := fmt.Sprintf("%s%s?token=%s&channel=%s&ts=%s",
		baseURL, conversationsRepliesPath, token, channel, threadTS)

	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	err = httpwrap.DoRequest(s.client, req, &result)
	return
}

// https://api.slack.com/methods/users.list
func (s *Service) GetUsersList(token string) (result UsersResponse, err error) {
	url := fmt.Sprintf("%s%s?token=%s",
		baseURL, usersListPath, token)

	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	var rsp *http.Response
	rsp, err = s.client.Do(req)
	if err != nil {
		return
	}

	var body []byte
	body, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &result)
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

	rsp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	body, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	return nil
}
