package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/larssont/WikiPal/internal/embed"
	"github.com/larssont/WikiPal/internal/parser"

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

	path := "configs/bot.json"

	jsonFile, err := os.Open(path)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &discordBot)

}

func main() {

	getBotConf()
	parser.GenerateCommands()

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

	if content[:len(discordBot.Prefix)] != discordBot.Prefix {
		return
	}

	content = content[len(discordBot.Prefix)+1:]

	message := strings.SplitN(content, " ", 2)
	query := message[0]

	var queryParam string
	if len(message) > 1 {
		queryParam = message[1]
	}

	response := parser.ProcessCommand(query, queryParam)
	switch response := response.(type) {
	case []string:
		sendText(m, s, response)
	case string:
		s.ChannelMessageSend(m.ChannelID, response)
	case embed.Message:
		sendEmbed(m, s, response)
	}

}

func sendEmbed(m *discordgo.MessageCreate, s *discordgo.Session, response embed.Message) {
	e := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       response.Color, // Green
		Description: response.Description,
		Title:       response.Title,
	}

	for k, v := range response.Fields {
		e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
			Name:   k,
			Value:  v,
			Inline: false,
		})
	}

	s.ChannelMessageSendEmbed(m.ChannelID, e)
}

func sendText(m *discordgo.MessageCreate, s *discordgo.Session, text []string) {
	for _, message := range text {
		s.ChannelMessageSend(m.ChannelID, message)
	}
}
