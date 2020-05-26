package confluence

import "fmt"

// https://confluence.atlassian.com/doc/confluence-wiki-markup-251003035.html#ConfluenceWikiMarkup-Links
func ComposeLinkFormat(URL, description string) string {
	if description == "" {
		return fmt.Sprintf("[%s]", URL)
	}

	return fmt.Sprintf("[%s|%s]", description, URL)
}
