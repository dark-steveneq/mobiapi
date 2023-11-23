package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dark-steveneq/mobiapi"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Printf("Usage: %s <domain> <login> <password> <semester>\n", os.Args[0])
		os.Exit(1)
	}
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
		fmt.Println("Couldn't sign in!")
		os.Exit(1)
	}

	defer api.Logout()

	semester, err := strconv.Atoi(os.Args[4])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	grades, err := api.GetGrades(semester)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for subject, grades := range grades {
		fmt.Println(subject)
		for _, grade := range grades {
			fmt.Printf("  [%s] %s (%s)\n", grade.Category, grade.Value, grade.Description)
		}
	}
}
