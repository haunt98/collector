package slack

type MessageRequest struct {
	Text string `json:"text,omitempty"`
}

type WebhookMessageRequest struct {
	MessageRequest

	ResponseType string `json:"response_type,omitempty"`
}

type MessagesResponse struct {
	Messages         []Message        `json:"messages"`
	ResponseMetadata ResponseMetadata `json:"response_metadata"`
}

type ResponseMetadata struct {
	NextCursor string `json:"next_cursor"`
}

type Message struct {
	Type string `json:"type"`
	User string `json:"user"`
	Text string `json:"text"`

	// https://api.slack.com/messaging/retrieving#finding_threads
	TS       string `json:"ts"`
	ThreadTS string `json:"thread_ts"`

	BotID string `json:"bot_id"`
}

type UsersResponse struct {
	Users []User `json:"members"`
}

type User struct {
	ID        string  `json:"id"`
	Profile   Profile `json:"profile"`
	IsBot     bool    `json:"is_bot"`
	IsAppUser bool    `json:"is_app_user"`
}

type Profile struct {
	DisplayName string `json:"display_name"`
}

type CommandPayload struct {
	Text        string `form:"text"`
	ResponseURL string `form:"response_url"`
	ChannelID   string `form:"channel_id"`
}
