package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/larssont/WikiPal/internal/embed"
	"github.com/larssont/WikiPal/internal/parser"

	"github.com/bwmarrin/discordgo"
)

var conf BotConfig

//BotConfig struct
type BotConfig struct {
	prefix string
	token  string
}

func getBotConfig() BotConfig {

	var token, prefix string

	flag.StringVar(&token, "token", "", "bot token")
	flag.StringVar(&prefix, "prefix", "!w", "bot prefix")

	flag.Parse()

	c := BotConfig{
		prefix: prefix,
		token:  token}

	return c

}

func main() {

	conf = getBotConfig()
	if len(conf.token) == 0 {
		fmt.Println("No bot token supplied. Use flag -token=MYTOKEN.")
		return
	}

	parser.GenerateCommands()

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + conf.token)
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

	// Update status
	var status = fmt.Sprintf("%s help", conf.prefix)
	dg.UpdateListeningStatus(status)

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

	if len(content) <= len(conf.prefix) {
		return
	}

	if content[:len(conf.prefix)] != conf.prefix {
		return
	}

	content = content[len(conf.prefix)+1:]

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
