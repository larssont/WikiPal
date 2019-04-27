package parser

import (
	"WikiPal/pkg/wiki"
	"fmt"
	"strings"
)

var input string
var commands []command

type addCommand func(keyword string, description string, fn func() []string)

type test struct {
	fn func() []string
}

type command struct {
	keyword, description string
	fn                   func() []string
}

func search() []string {

	s := strings.Split(input, " ")

	var langCode string
	searchQuery := input

	lastWord := s[len(s)-1]

	if string(lastWord[0]) == "-" {
		if len(lastWord) > 1 {
			langCode = lastWord[1:]
			searchQuery = input[:len(input)-len(langCode)-1]
		}
	}

	if input == "" {
		return []string{"Search for what?"}
	}

	wikiSearch := wiki.Search(searchQuery, langCode)
	var res string

	if wikiSearch.URL != "" {

		var altPages string

		for i := range wikiSearch.AlternativeHits {
			altPages += fmt.Sprintf("<%s>\n", wikiSearch.AlternativeHits[i])
		}

		text := fmt.Sprintf(" ᠌᠌᠌᠌᠌᠌᠌᠌\n**Alternative pages:**\n%s", altPages)

		return []string{wikiSearch.URL, text}

	}

	res = fmt.Sprintf(`Ehh, I couldn't find anything for "%s"`, input)
	return []string{res}
}

func help() []string {
	text := `
	Hi!
	My name is wikipal. :robot:
	I can help you find stuff on wikipedia. :mag:
	These are my commands.
	
	`
	for _, cmd := range commands {
		text += fmt.Sprintf("\n`%s` %s", cmd.keyword, cmd.description)
	}
	return []string{text}
}

//GenerateCommands generates all commands for use in ProcessCommand
func GenerateCommands() {
	var a addCommand = func(a string, b string, fn func() []string) {
		commands = append(commands, command{a, b, fn})
	}

	a("search", "searches for stuff", search)
	a("help", "shows you this", help)
}

//ProcessCommand processes text queries and returns a response
func ProcessCommand(query string, queryParam string) []string {

	input = queryParam

	for _, cmd := range commands {
		if cmd.keyword == query {
			return cmd.fn()
		}
	}

	return []string{"Ehh, I don't know that command."}
}
