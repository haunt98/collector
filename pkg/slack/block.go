package slack

// https://api.slack.com/block-kit

func BuildSectionBlock(text string) SectionBlock {
	return SectionBlock{
		Block: Block{
			Type: TypeSection,
		},
		Text: TextBlock{
			Type: TypeMarkdown,
			Text: text,
		},
	}
}

func BuildDividerBlock() DividerBlock {
	return DividerBlock{
		Block: Block{
			Type: TypeDivider,
		},
	}
}
