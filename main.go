package main

import (
	"bufio"
	"fmt"
	"labasqlredis/src"
	"log"
	"os"
	"strings"
)

const (
	dbConnectionString = "host=localhost port=5432 user=postgres password=yourpassword dbname=demo sslmode=disable"
	redisAddr          = "localhost:6379"
	redisPassword      = ""
)

func main() {
	app := src.NewApp(dbConnectionString, redisAddr, redisPassword)
	defer app.Close()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter your SQL query (or type 'exit' to quit):")

	for {
		fmt.Print("> ")
		scanner.Scan()
		query := scanner.Text()

		if strings.ToLower(query) == "exit" {
			log.Println("Exiting. . . ")
			break
		}

		result, err := app.QueryWithCache(query)
		if err != nil {
			log.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("Result: %s\n", result)
	}
}
