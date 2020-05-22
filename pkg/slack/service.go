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

func NewService() *Service {
	return &Service{
		client: &http.Client{},
	}
}

const (
	baseURL = "https://slack.com/api"
)

// https://api.slack.com/methods/conversations.history
func (c *Service) GetConversationsHistory(token, channel, cursor string) (result MessagesResponse, err error) {
	url := fmt.Sprintf("%s/conversations.history?token=%s&channel=%s",
		baseURL, token, channel)
	if len(cursor) != 0 {
		url += fmt.Sprintf("&cursor=%s", cursor)
	}

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

// https://api.slack.com/methods/conversations.replies
func (c *Service) GetConversationsReplies(token, channel, threadTS string) (result MessagesResponse, err error) {
	url := fmt.Sprintf("%s/conversations.replies?token=%s&channel=%s&ts=%s",
		baseURL, token, channel, threadTS)

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

// https://api.slack.com/methods/users.list
func (c *Service) GetUsersList(token string) (result UsersResponse, err error) {
	url := fmt.Sprintf("%s/users.list?token=%s",
		baseURL, token)

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

// https://api.slack.com/interactivity/handling#message_responses
func (c *Service) PostMessageByResponseURL(responseURL string, msgReq MessageRequestByResponseURL) error {
	log.Print("XXX")
	log.Printf("%+v\n", msgReq)

	body, err := json.Marshal(msgReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, responseURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	rsp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	body, err = ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	log.Println(string(body))
	return nil
}
