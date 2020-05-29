package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dghubble/sling"
)

const (
	baseURL = "https://slack.com/api"
)

type Service struct {
	client *http.Client
	sl     *sling.Sling
}

func NewService() *Service {
	return &Service{
		client: &http.Client{},
		sl:     sling.New().Base(baseURL),
	}
}

// https://api.slack.com/methods/conversations.history
func (s *Service) GetConversationsHistory(token, channel, cursor string) (result MessagesResponse, err error) {
	type Params struct {
		Token   string `json:"token,omitempty"`
		Channel string `json:"channel,omitempty"`
		Cursor  string `json:"cursor,omitempty"`
	}

	var req *http.Request
	req, err = s.sl.New().Get("/conversations.history").
		QueryStruct(Params{
			Token:   token,
			Channel: channel,
			Cursor:  cursor,
		}).
		Request()
	if err != nil {
		return
	}

	err = s.Do(req, &result)
	return
}

// https://api.slack.com/methods/conversations.replies
func (s *Service) GetConversationsReplies(token, channel, threadTS string) (result MessagesResponse, err error) {
	url := fmt.Sprintf("%s/conversations.replies?token=%s&channel=%s&ts=%s",
		baseURL, token, channel, threadTS)

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

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, result); err != nil {
		return err
	}

	return nil
}
