package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	baseURL                  = "https://slack.com/api"
	conversationsHistoryPath = "/conversations.history"
	conversationsRepliesPath = "/conversations.replies"
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
	url := fmt.Sprintf("%s%s?token=%s&channel=%s",
		baseURL, conversationsHistoryPath, token, channel)
	if cursor != "" {
		url += fmt.Sprintf("&cursor=%s", cursor)
	}

	var req *http.Request
	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	err = s.Do(req, &result)
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

	err = s.Do(req, &result)
	return
}

// https://api.slack.com/methods/users.list
func (s *Service) GetUsersList(token string) (result UsersResponse, err error) {
	url := fmt.Sprintf("%s/users.list?token=%s",
		baseURL, token)

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

func (s *Service) Do(req *http.Request, result interface{}) error {
	rsp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	log.Printf("XXX %s", string(body))

	if err = json.Unmarshal(body, result); err != nil {
		return err
	}

	return nil
}
