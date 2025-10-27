package main

import (
	"fmt"
	"strconv"

	scrapper "github.com/abroudoux/twinpick/internal/scrapper"
)

func main() {
	movies := scrapper.Scrap("potatoze")
	for i, movie := range movies {
		fmt.Println(strconv.Itoa(i) + " " + *movie)
	}
}
