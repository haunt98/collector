package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dghubble/sling"
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
	type Params struct {
		Token   string `json:"token,omitempty"`
		Channel string `json:"channel,omitempty"`
		Cursor  string `json:"cursor,omitempty"`
	}

	log.Println(channel)
	log.Println(cursor)

	resultPointer := new(MessagesResponse)
	_, err = sling.New().Get(baseURL + "/conversations.history").
		QueryStruct(Params{
			Token:   token,
			Channel: channel,
			Cursor:  cursor,
		}).
		ReceiveSuccess(resultPointer)
	if err != nil {
		return
	}

	log.Printf("XXX %+v\n", resultPointer)

	result = *resultPointer
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

	return nil
}
