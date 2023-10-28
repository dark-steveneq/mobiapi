package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dark-steveneq/mobiapi"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s <domain> <login> <package>\n", os.Args[0])
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)

	api, err := mobiapi.New(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	signedin, err := api.PasswordAuth(os.Args[2], os.Args[3])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else if !signedin {
		fmt.Println("Couldn't log in!")
		os.Exit(1)
	}
	defer api.Logout()

	messages, err := api.GetReceivedMessages(false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	unread := map[int]mobiapi.MessageInfo{}
	for _, message := range messages {
		if message.Read == true {
			unread[len(unread)] = message
		}
	}

	fmt.Printf("Hello, %s!\nYou Have %d unread messages. Do you want to read them one by one? [Y/n]\n", api.GetName(), len(unread))
	text, _ := reader.ReadString('\n')
	text = strings.ToLower(text)
	if text == "n" || text == "no" {
		fmt.Println("Okeh.")
		os.Exit(0)
	}
	for _, message := range unread {
		messagecontent, err := api.GetMessageContent(message)
		if err != nil {
			fmt.Println("Couldn't read message!")
		} else {
			fmt.Printf("Title: %s\nFrom: %s\nRead: %t\n'''\n%s\n'''\nDownloads: %d", message.Title, message.Author, message.Read, messagecontent.Content, len(messagecontent.Downloads))
			fmt.Printf("\n")
			for name := range messagecontent.Downloads {
				fmt.Printf("- %s\n", name)
			}
		}
		reader.ReadString('\n')
	}
}
