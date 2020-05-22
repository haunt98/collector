package slack

// https://api.slack.com/reference/messaging/payload
type MessagePayload struct {
	Text string `json:"text,omitempty"`
}

const (
	ResponseTypeInChannel = "in_channel"
	ResponseTypeEphemeral = "ephemeral"
)

type MessageRequestByResponseURL struct {
	MessagePayload

	ResponseType string `json:"response_type,omitempty"`
}

type Block struct {
	Type string `json:"type"`
}

// https://api.slack.com/reference/block-kit/blocks#section
type SectionBlock struct {
	Block

	Text interface{} `json:"text"`
}

// https://api.slack.com/reference/block-kit/blocks#divider
type DividerBlock struct {
	Block
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

// https://api.slack.com/interactivity/slash-commands#command_payload_descriptions
type CommandPayload struct {
	Text        string `form:"text"`
	ResponseURL string `form:"response_url"`
	ChannelID   string `form:"channel_id"`
}
