package confluence

import "fmt"

// https://confluence.atlassian.com/doc/confluence-wiki-markup-251003035.html#ConfluenceWikiMarkup-Links
func ComposeLinkFormat(url, description string) string {
	if description == "" {
		return fmt.Sprintf("[%s]", url)
	}

	return fmt.Sprintf("[%s|%s]", description, url)
}
