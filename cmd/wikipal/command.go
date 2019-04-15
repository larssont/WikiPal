package main

import (
	"fmt"
)

func processCommand(query string, queryParam string) interface{} {

	var response interface{}

	switch query {
	case "find":
		response = find(queryParam)
	case "help":
		response = help()
	default:
		response = "Ehh, I don't know that command."
	}
	return response
}

func find(queryParam string) interface{} {
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
	response := fmt.Sprintf(`
		Hi my name is %s.
		I can help you find pictures on wikipedia.
		
		E.g. 
		You can find an image of Bruce Lee by typing "!w find bruce lee"

		Remember that you always need to use the prefix "!w" when you chat with me.
		Have fun!`, discordBot.Name)
	return response
}
