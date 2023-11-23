package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dark-steveneq/mobiapi"
)

/*
"Token" in this context is a value from a uniquely named cookie.
Each domain uses a different name, so make sure not to include it
if you don't want to dox yourself.
*/
func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <domain> <token>\n", os.Args[0])
		os.Exit(1)
	}

	api, err := mobiapi.New(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	signedin, err := api.TokenAuth(os.Args[2])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else if !signedin {
		fmt.Println("Couldn't sign in!")
		os.Exit(1)
	}

	defer api.Logout()

	for {
		time.Sleep(5 * time.Minute)
		if err := api.ExtendSession(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
