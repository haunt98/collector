package slack

type MessageRequest struct {
	Text string `json:"text,omitempty"`

	// Normally
	Channel  string `json:"channel,omitempty"`
	ThreadTS string `json:"thread_ts,omitempty"`

	// Webhook
	ResponseType string `json:"response_type,omitempty"`
}

type MessageResponse struct {
	TS      string  `json:"ts"`
	Message Message `json:"message"`
}

type MessagesResponse struct {
	Messages []Message `json:"messages"`
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
