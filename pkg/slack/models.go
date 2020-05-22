package slack

// https://api.slack.com/reference/messaging/payload
type MessagePayload struct {
	Text   string        `json:"text,omitempty"`
	Blocks []interface{} `json:"blocks,omitempty"`
}

type BlockType struct {
	Type string `json:"type"`
}

const (
	TypeMarkdown = "mrkdwn"
	TypeSection  = "section"
	TypeDivider  = "divider"
	TypeImage    = "image"
)

// https://api.slack.com/reference/block-kit/blocks#section
type SectionBlock struct {
	BlockType

	Text interface{} `json:"text"`
	// https://api.slack.com/reference/block-kit/block-elements
	Accessory interface{} `json:"accessory,omitempty"`
}

type TextBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// https://api.slack.com/reference/block-kit/block-elements#image
type ImageElement struct {
	BlockType

	ImageURL string `json:"image_url"`
	AltText  string `json:"alt_text"`
}

// https://api.slack.com/reference/block-kit/blocks#divider
type DividerBlock struct {
	BlockType
}

const (
	ResponseTypeInChannel = "in_channel" // everyone can see response
	ResponseTypeEphemeral = "ephemeral"  // only the one who send message can see response
)

type MessageRequestByResponseURL struct {
	MessagePayload

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

	TS       string `json:"ts"`
	ThreadTS string `json:"thread_ts"`

	BotID string `json:"bot_id"`
}

// https://api.slack.com/types/user
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
	Image48     string `json:"image_48"`
}

// https://api.slack.com/interactivity/slash-commands#command_payload_descriptions
type CommandPayload struct {
	Text        string `form:"text"`
	ResponseURL string `form:"response_url"`
	ChannelID   string `form:"channel_id"`
}
