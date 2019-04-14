package main

import (
	"fmt"
	"io"
	"os"
)

// ImageFile struct
type ImageFile struct {
	Name string
	File io.Reader
}

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

	imgFound, imgURL := findImage(queryParam)

	if !imgFound {
		return "I'm sorry, I couldn't find anything for " + queryParam + "."
	}

	fileName := downloadImage(imgURL)
	file, _ := os.Open(fileName)

	img := ImageFile{fileName, file}

	os.Remove(fileName)
	return img
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
