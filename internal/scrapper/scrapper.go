package scrapper

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/gocolly/colly/v2"
)

func Scrap(letterboxdUsername string) []*string {
	movieCh := make(chan []string)
	var wg sync.WaitGroup

	pageCollector := colly.NewCollector(
		colly.AllowedDomains("letterboxd.com"),
	)

	var totalPages int

	pageCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("➡️ Visiting", r.URL.String())
	})

	pageCollector.OnHTML("div.paginate-pages ul", func(e *colly.HTMLElement) {
		e.ForEach("li.paginate-page a", func(_ int, el *colly.HTMLElement) {
			num := el.Text
			if n, err := strconv.Atoi(num); err == nil && n > totalPages {
				totalPages = n
			}
		})

		fmt.Println("Total pages:", totalPages)
	})

	err := pageCollector.Visit(fmt.Sprintf("https://letterboxd.com/%s/watchlist", letterboxdUsername))
	if err != nil {
		log.Fatal(err)
	}

	pageCollector.Wait()

	for i := 1; i <= totalPages; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			c := colly.NewCollector(colly.AllowedDomains("letterboxd.com"))
			var movies []string

			c.OnHTML("div.poster-grid li", func(e *colly.HTMLElement) {
				movieName := e.ChildAttr("div.react-component", "data-item-full-display-name")
				if movieName != "" {
					movies = append(movies, movieName)
				}
			})

			c.OnRequest(func(r *colly.Request) {
				fmt.Println("➡️ Visiting page", page, ":", r.URL.String())
			})

			err := c.Visit(fmt.Sprintf("https://letterboxd.com/%s/watchlist/page/%d", letterboxdUsername, page))
			if err != nil {
				log.Println("Error on the page ", page, ":", err)
			}

			c.Wait()
			movieCh <- movies
		}(i)
	}

	go func() {
		wg.Wait()
		close(movieCh)
	}()

	var movies []*string

	for ch := range movieCh {
		for _, movie := range ch {
			movies = append(movies, &movie)
		}
	}

	return movies
}
