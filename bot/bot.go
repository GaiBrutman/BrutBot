// bot.go is responsible for handling chat messages and sending the bot's responses
package bot

import (
	"brutBot/config"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var BotID string
var goBot *discordgo.Session

func Start() {
	// Initialize new bot
	var err error
	goBot, err = discordgo.New("Bot " + config.Token)

	if err != nil {
		log.Fatal(err)
	}

	// Gets the bot's user details
	u, err := goBot.User("@me")

	if err != nil {
		log.Fatal(err)
	}

	BotID = u.ID

	// Adds messageHandler to bot's handlers
	goBot.AddHandler(messageHandler)

	if err = goBot.Open(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bot is running!")
}

// messageHandler handles chat messages and sends back the bot's responses
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Return if handling the bot's messages
	if m.Author.ID == BotID {
		return
	}

	// Check if the message is a bot command (starts with BotPrefix)
	if strings.HasPrefix(m.Content, config.BotPrefix) {
		fmt.Printf("%s: %q\n", m.Author.Username, m.Content)

		// Split message by spaces
		content := strings.Split(m.Content[1:], " ")

		if len(content) == 0 {
			return
		}

		// Divide to command and command arguments
		command, args := strings.ToLower(content[0]), content[1:]

		switch command {
		case "help", "h":
			sendText("Hello I am BOT.", s, m)
		case "ping":
			sendText("Pong", s, m)
		case "time", "t":
			sendText(time.Now().Format(time.UnixDate), s, m)
		case "gopher", "go": // Send the image in ./gopher.jpg
			if f, err := os.Open("./gopher.jpg"); err == nil {
				_, _ = s.ChannelFileSend(m.ChannelID, "gopher.jpg", f)
			}
		case "dog", "cat": // sends a dogs/cats related Reddit post
			sendRedditImage(command+"pictures", s, m)
		case "pewdiepie", "pewds": // sends a r/PewdiepieSubmissions Reddit post
			sendRedditImage("PewdiepieSubmissions", s, m)
		case "meme", "mm": // sends a meme Reddit post
			subreddit := "memes"
			if rand.Seed(time.Now().UnixNano()); rand.Intn(2) == 1 {
				subreddit = "dank" + subreddit
			}
			sendRedditImage(subreddit, s, m)
		case "reddit", "rd": // sends a Reddit post of a given subreddit
			if len(args) > 0 {
				sendRedditImage(args[0], s, m)
			}
		}
	}
}

// sendText send msg to a given channel
func sendText(msg string, s *discordgo.Session, m *discordgo.MessageCreate) {
	_, _ = s.ChannelMessageSend(m.ChannelID, msg)
	fmt.Printf("Sending %q\n", msg)
}

// sendText send a random post from a given subreddit to a given channel
func sendRedditImage(subreddit string, s *discordgo.Session, m *discordgo.MessageCreate) {
	a := time.Now()

	// Get random post from subreddit
	post, err := GetRandPost(subreddit)
	fmt.Printf("Get: %.2f secs\n", time.Since(a).Seconds())

	a = time.Now()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("__%s__", post.Title)) // Send post title
		if post.Body != "" {
			_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s", post.Body)) // Send post body
		}
		if post.Image != nil {
			_, _ = s.ChannelFileSend(m.ChannelID, subreddit+".jpg", post.Image) // Send post image
		}
	}

	fmt.Printf("Send: %.2f secs\n", time.Since(a).Seconds())
}
