package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Service struct {
	client *http.Client
}

const (
	baseURL                      = "https://slack.com/api/"
	conversationsHistoryEndpoint = "conversations.history"
	conversationsRepliesEndpoint = "conversations.replies"
	usersListEndpoint            = "users.list"
)

func NewService(token string) *Service {
	return &Service{
		client: &http.Client{},
	}
}

func (c *Service) GetConversationHistory(token, channel string) (result MessagesResponse, err error) {
	url := fmt.Sprintf("%s%s?token=%s&channel=%s",
		baseURL, conversationsHistoryEndpoint, token, channel)

	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	var rsp *http.Response
	rsp, err = c.client.Do(req)
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

func (c *Service) GetConversationReplies(token, channel, threadTS string) (result MessagesResponse, err error) {
	url := fmt.Sprintf("%s%s?token=%s&channel=%s&ts=%s",
		baseURL, conversationsRepliesEndpoint, channel, threadTS)

	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	var rsp *http.Response
	rsp, err = c.client.Do(req)
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

func (c *Service) GetUsersList(token string) (result UsersResponse, err error) {
	url := fmt.Sprintf("%s%s?token=%s",
		baseURL, usersListEndpoint, token)

	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	var rsp *http.Response
	rsp, err = c.client.Do(req)
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

func (c *Service) PostThreadMessageByWebhook(webhookURL, text, responseType string) (err error) {
	msgReq := WebhookMessageRequest{
		MessageRequest: MessageRequest{
			Text: text,
		},
		ResponseType: responseType,
	}

	var body []byte
	body, err = json.Marshal(msgReq)
	if err != nil {
		return
	}

	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	var rsp *http.Response
	rsp, err = c.client.Do(req)
	if err != nil {
		return
	}

	body, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return
	}

	log.Println(string(body))
	return
}
