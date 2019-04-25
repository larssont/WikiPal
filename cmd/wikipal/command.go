package main

import (
	"fmt"
	"strings"
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

	s := strings.Split(queryParam, " ")

	var langCode string
	lastWord := s[len(s)-1]

	if string(lastWord[0]) == "-" {
		if len(lastWord) > 1 {
			langCode = lastWord[1:]
			queryParam = queryParam[:len(queryParam)-len(langCode)-1]
		}
	}

	if queryParam == "" {
		return "Search for what?"
	}

	wikiSearch := searchWiki(queryParam, langCode)

	if wikiSearch.URL != "" {
		return wikiSearch
	}

	return fmt.Sprintf(`Ehh, I couldn't find anything for "%s"`, queryParam)
}

func help() string {

	wHelp := "`!w help` Prints this information :information_source:"
	wSearch := "`!w search [title]` Searches wikipedia for a given title :mag:"
	myPrefix := "You can talk with me by using the prefix `!w`"

	response := fmt.Sprintf(`
		Hi my name is %s :robot:
		I can help you search Wikipedia.

		%s
		These are my commands:

		%s
		%s

		Have fun!`, discordBot.Name, myPrefix, wHelp, wSearch)
	return response
}
