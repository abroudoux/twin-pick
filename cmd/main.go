package main

import (
	"fmt"

	scrapper "github.com/abroudoux/twinpick/internal/scrapper"
)

func main() {
	usernames := []string{"abroudoux", "66Sceptre", "mascim"}
	results := scrapper.ScrapUsers(usernames)

	for user, movies := range results {
		fmt.Println("=== Watchlist de", user, "===")
		for i, movie := range movies {
			fmt.Printf("%d: %s\n", i+1, *movie)
		}
	}
}
