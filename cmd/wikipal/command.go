package main

import (
	"fmt"
)

func processCommand(query string, queryParam string) interface{} {

	var response interface{}

	switch query {
	case "search":
		response = search(queryParam)
	case "help":
		response = help()
	default:
		response = "Ehh, I don't know that command."
	}
	return response
}

func search(queryParam string) interface{} {
	if queryParam == "" {
		return "Find what?"
	}

	wikiPage := SearchWiki(queryParam)

	if wikiPage.ThumbnailFile != nil || wikiPage.Snippet != "" {
		return wikiPage
	}

	return fmt.Sprintf(`Ehh, I couldn't find anything for "%s"`, queryParam)
}

func help() string {

	wHelp := "`!w help` Prints this information :information_source:"
	wSearch := "`!w search [title]` Searches wikipedia for a given title :mag:"

	response := fmt.Sprintf(`
		Hi my name is %s :robot:
		I can help you search Wikipedia.

		You can talk with me by using the prefix !w
		These are my commands:

		%s
		%s

		Have fun!`, discordBot.Name, wHelp, wSearch)
	return response
}
