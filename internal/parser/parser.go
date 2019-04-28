package parser

import (
	"WikiPal/internal/embed"
	"WikiPal/pkg/wiki"
	"fmt"
	"math/rand"
	"strings"
)

var input string
var commands []command

type addCommand func(keyword string, description string, param string, option string, fn func() interface{})

type command struct {
	keyword, description, param, option string
	fn                                  func() interface{}
}

func search() interface{} {

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

func help() interface{} {

	var e embed.Message

	var cmdExpression = func(cmd, param, option string) (out string) {
		out = cmd
		if param != "" {
			out += fmt.Sprintf(" <%s>", param)
		}
		if option != "" {
			out += fmt.Sprintf(" -%s", option)
		}
		out = "`" + out + "`"
		return
	}

	e.Description = fmt.Sprintf(`
	Hi!
	My name is wikipal. :robot:
	I can help you find stuff on wikipedia. :mag:


	You can see all my commands down below.
	Obligatory paramaters are marked with %s.
	Optional paramaters are preceded with a hyphen %s.




	`, "`<>`", "`-`")

	e.Fields = make(map[string]string)
	e.Color = 0x900000
	for _, cmd := range commands {
		e.Fields[cmdExpression(cmd.keyword, cmd.param, cmd.option)] = cmd.description
	}

	return e
}

func langs() interface{} {

	text := "_ _\n"
	text += fmt.Sprintf(`Default language: %s`, wiki.Languages[wiki.DefaultLanguage])
	text += "\nThe following languages are available:"

	for k, v := range wiki.Languages {
		text += fmt.Sprintf("\n`%s = %s`", k, v)
	}

	return text
}

func setlang() interface{} {
	if input == "" {
		return "No language paramater given."
	}
	for k, v := range wiki.Languages {
		if k == input || v == input {
			wiki.DefaultLanguage = k
			return fmt.Sprintf("Default language changed to %s.", v)
		}
	}

	return "Sorry i don't know that language :("
}

//GenerateCommands generates all commands for use in ProcessCommand
func GenerateCommands() {
	var cmd addCommand = func(a string, b string, c string, d string, fn func() interface{}) {
		commands = append(commands, command{a, b, c, d, fn})
	}

	cmd("search", "Search wikipedia for a given query. Use language codes to specify which site you want to use (such as -de for germany)", "query", "langcode", search)
	cmd("help", "Shows you this :)", "", "", help)
	cmd("langs", "Lists all available languages to search with.", "", "", langs)
	cmd("setlang", "Set the default language for wiki searches.", "language", "", setlang)
}

//ProcessCommand processes text queries and returns a response
func ProcessCommand(query string, queryParam string) interface{} {

	input = queryParam

	for _, cmd := range commands {
		if cmd.keyword == query {
			return cmd.fn()
		}
	}

	responses := []string{
		"I don't know that command.",
		"I'm not quite sure what you meant there.",
		"I have no idea what you said.",
	}

	return []string{responses[rand.Intn(len(responses))]}
}
