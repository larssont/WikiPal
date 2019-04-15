package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var discordBot Bot

//Bot struct
type Bot struct {
	Name   string
	Prefix string
	Token  string
}

func getBotConf() {
	jsonFile, err := os.Open("../../configs/bot.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &discordBot)
}

func main() {

	getBotConf()

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + discordBot.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.

	content := strings.ToLower(m.Content)

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Author.Bot {
		return
	}

	if len(content) <= len(discordBot.Prefix) {
		return
	}

	if content[:2] != discordBot.Prefix {
		return
	}

	content = content[len(discordBot.Prefix)+1:]

	message := strings.SplitN(content, " ", 2)
	query := message[0]

	var queryParam string
	if len(message) > 1 {
		queryParam = message[1]
	}

	response := processCommand(query, queryParam)

	if str, ok := response.(string); ok {
		s.ChannelMessageSend(m.ChannelID, str)
	} else if wikiPage, ok := response.(WikiPage); ok {

		s.ChannelFileSend(m.ChannelID, wikiPage.ThumbnailFileName, wikiPage.ThumbnailFile)
		s.ChannelMessageSend(m.ChannelID, wikiPage.Snippet)

		defer wikiPage.ThumbnailFile.Close()
	}
}
