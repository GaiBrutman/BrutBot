// The main file that responsible for starting and running the bot
package main

import (
	"brutBot/bot"
	"brutBot/config"
	"fmt"
)

func main() {
	// Read the bot's configurations
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Starting the bot
	bot.Start()

	// Blocks the main execution (to prevent it from returning)
	<-make(chan struct{})
	return
}
