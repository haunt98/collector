package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Slack struct {
	token string

	client      *http.Client
	bearerToken string
}

const (
	bearer                  = "Bearer"
	baseURL                 = "https://slack.com/api/"
	chatPostMessageURL      = "chat.postMessage"      // https://api.slack.com/methods/chat.postMessage
	conversationsRepliesURL = "conversations.replies" // https://api.slack.com/methods/conversations.replies
	usersListURL            = "users.list"            // https://api.slack.com/methods/users.list
	conversationsHistory    = "conversations.history"
)

func NewSlack(token string) *Slack {
	c := &Slack{
		client:      &http.Client{},
		token:       token,
		bearerToken: bearer + " " + token,
	}

	return c
}

func (c *Slack) PostChannelMessage(text, channel string) (MessageResponse, error) {
	msgReq := MessageRequest{
		Channel: channel,
		Text:    text,
	}

	return c.postMessage(msgReq)
}

func (c *Slack) PostThreadMessage(text, channel, threadTS string) (MessageResponse, error) {
	msgReq := MessageRequest{
		Channel:  channel,
		Text:     text,
		ThreadTS: threadTS,
	}

	return c.postMessage(msgReq)
}

func (c *Slack) postMessage(msgReq MessageRequest) (result MessageResponse, err error) {
	var body []byte
	body, err = json.Marshal(msgReq)
	if err != nil {
		return
	}

	var req *http.Request
	req, err = http.NewRequest(http.MethodPost, baseURL+chatPostMessageURL, bytes.NewBuffer(body))
	if err != nil {
		return
	}
	req.Header.Set("Authorization", c.bearerToken)
	req.Header.Set("Content-Type", "application/json")

	var rsp *http.Response
	rsp, err = c.client.Do(req)
	if err != nil {
		return
	}

	body, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &result)
	return
}

func (c *Slack) GetThreadMessages(channel, threadTS string) (result MessagesResponse, err error) {
	url := fmt.Sprintf("%s%s?token=%s&channel=%s&ts=%s",
		baseURL, conversationsRepliesURL, c.token, channel, threadTS)

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

func (c *Slack) GetUsers() (result UsersResponse, err error) {
	url := fmt.Sprintf("%s%s?token=%s",
		baseURL, usersListURL, c.token)

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

func (c *Slack) PostThreadMessageByWebhook(webhookURL, text, responseType string) error {
	msgReq := MessageRequest{
		Text:         text,
		ResponseType: responseType,
	}

	return c.postMessageByWebhook(webhookURL, msgReq)
}

func (c *Slack) postMessageByWebhook(webhookURL string, msgReq MessageRequest) (err error) {
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

func (c *Slack) GetChannelHistory(channel string) (result MessagesResponse, err error) {
	url := fmt.Sprintf("%s%s?token=%s&channel=%s",
		baseURL, conversationsHistory, c.token, channel)

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
