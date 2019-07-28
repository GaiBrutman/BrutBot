package bot

import (
	"brutBot/config"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"os"
	"strings"
	"time"
)

var BotID string
var goBot *discordgo.Session

func Start() {
	var err error
	goBot, err = discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	goBot.AddHandler(messageHandler)

	if err = goBot.Open(); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running!")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	if strings.HasPrefix(m.Content, config.BotPrefix) {
		fmt.Printf("%s: %q\n", m.Author.Username, m.Content)

		content := strings.Split(strings.ToLower(m.Content[1:]), " ")

		if len(content) == 0 {
			return
		}

		command, args := content[0], content[1:]

		var response string

		switch command {
		case "help", "h":
			response = "Hello I am BOT."
		case "ping":
			response = "Pong"
		case "time", "T":
			response = time.Now().Format(time.UnixDate)
		case "gopher", "go":
			f, err := os.Open("./gopher.jpg")
			if err == nil {
				_, _ = s.ChannelFileSend(m.ChannelID, "gopher.jpg", f)
			}
		case "dog", "cat":
			sendRedditImage(command+"pictures", s, m)
		case "pewds", "pewdiepie":
			sendRedditImage("PewdiepieSubmissions", s, m)
		case "meme", "mm":
			subreddit := "memes"
			if rand.Seed(time.Now().UnixNano()); rand.Intn(2) == 1 {
				subreddit = "dank" + subreddit
			}
			sendRedditImage(subreddit, s, m)
		case "reddit", "rd":
			if len(args) > 0 {
				sendRedditImage(args[0], s, m)
			}
		}

		if response != "" {
			_, _ = s.ChannelMessageSend(m.ChannelID, response)
			fmt.Printf("Sending %q\n", response)
		}
	}
}

func sendRedditImage(subreddit string, s *discordgo.Session, m *discordgo.MessageCreate) {
	a := time.Now()

	item, err := GetRandImage(subreddit)
	fmt.Printf("Get: %.2f secs\n", time.Since(a).Seconds())

	a = time.Now()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("__%s__", item.Title))
		if item.Body != "" {
			_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s", item.Body))
		}
		if item.Image != nil {
			_, _ = s.ChannelFileSend(m.ChannelID, subreddit+".jpg", item.Image)
		}
	}

	fmt.Printf("Send: %.2f secs\n", time.Since(a).Seconds())
}
