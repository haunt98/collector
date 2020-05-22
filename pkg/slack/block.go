package slack

// https://api.slack.com/block-kit

func BuildSectionBlock(text string) SectionBlock {
	return SectionBlock{
		BlockType: BlockType{
			Type: TypeSection,
		},
		Text: TextBlock{
			Type: TypeMarkdown,
			Text: text,
		},
		Accessory: nil,
	}
}

func BuildSectionBlockWithImage(text, imageURL, imageAltText string) SectionBlock {
	return SectionBlock{
		BlockType: BlockType{
			Type: TypeSection,
		},
		Text: text,
		Accessory: ImageElement{
			BlockType: BlockType{
				Type: TypeImage,
			},
			ImageURL: imageURL,
			AltText:  imageAltText,
		},
	}
}

func BuildDividerBlock() DividerBlock {
	return DividerBlock{
		BlockType: BlockType{
			Type: TypeDivider,
		},
	}
}
