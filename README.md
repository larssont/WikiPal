# WikiPal
A discord bot for fetching information from wikipedia using the DiscordGo package.

### Prerequisites

In order to run this bot on your own machine, you need to install Go.

Follow the instructions here for your system: [golang.org/doc/install](https://golang.org/doc/install)

### Installing/Running

#### Normal

```shell
git clone https://github.com/larssont/WikiPal.git

cd WikiPal
```

Edit `WikiPal/configs/bot.json`. Make sure to set your own bot token.

```shell
go run cmd/wikipal/main.go
```

#### Docker


## Built With

* [MediaWiki action API](https://www.mediawiki.org/wiki/API:Main_page) - Used to get data from wikipedia
* [DiscordGo](https://github.com/bwmarrin/discordgo) - Go package for discord chat client API

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
