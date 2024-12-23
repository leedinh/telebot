package main

import (
	"fmt"

	config "github.com/leedinh/telebot/bitly/internal/config/url-shortener"
)

func main() {
	// Load the configuration
	config := config.LoadConfig()
	fmt.Println(config)

}
