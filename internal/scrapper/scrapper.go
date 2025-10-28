package scrapper

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/gocolly/colly/v2"
)

func ScrapUsersWachtlists(usernames []string, genres []string) map[string][]string {
	watchlists := make(map[string][]string)

	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, username := range usernames {
		wg.Add(1)
		go func(user string) {
			defer wg.Done()
			watchlist := scrapWatchlist(user, genres)
			mu.Lock()
			watchlists[user] = watchlist
			mu.Unlock()
		}(username)
	}

	wg.Wait()

	return watchlists
}

func scrapWatchlist(letterboxdUsername string, genres []string) []string {
	filmCh := make(chan []string)

	var wg sync.WaitGroup

	pageCollector := colly.NewCollector(
		colly.AllowedDomains("letterboxd.com"),
	)

	var totalPages int

	pageCollector.OnRequest(func(r *colly.Request) {
		log.Infof("➡️ Visiting %s", r.URL.String())
	})

	pageCollector.OnHTML("div.paginate-pages ul", func(e *colly.HTMLElement) {
		e.ForEach("li.paginate-page a", func(_ int, el *colly.HTMLElement) {
			num := el.Text
			if n, err := strconv.Atoi(num); err == nil && n > totalPages {
				totalPages = n
			}
		})

		log.Infof("Total pages: %d", totalPages)
	})

	err := pageCollector.Visit(buildWatchlistURL(letterboxdUsername, genres))
	if err != nil {
		log.Fatal(err)
	}

	pageCollector.Wait()

	for i := 1; i <= totalPages; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			c := colly.NewCollector(colly.AllowedDomains("letterboxd.com"))
			var films []string

			c.OnHTML("div.poster-grid li", func(e *colly.HTMLElement) {
				filmName := e.ChildAttr("div.react-component", "data-item-full-display-name")
				if filmName != "" {
					films = append(films, filmName)
				}
			})

			c.OnRequest(func(r *colly.Request) {
				log.Infof("➡️ Visiting page %d : %s", page, r.URL.String())
			})

			err := c.Visit(fmt.Sprintf("https://letterboxd.com/%s/watchlist/page/%d", letterboxdUsername, page))
			if err != nil {
				log.Errorf("Error on the page %d : %v", page, err)
			}

			c.Wait()
			filmCh <- films
		}(i)
	}

	go func() {
		wg.Wait()
		close(filmCh)
	}()

	var watchlist []string

	for ch := range filmCh {
		for _, film := range ch {
			watchlist = append(watchlist, film)
		}
	}

	return watchlist
}

func buildWatchlistURL(username string, genres []string) string {
	url := fmt.Sprintf("https://letterboxd.com/%s/watchlist/", username)
	if len(genres) == 0 {
		return url
	}

	url += "genre/"
	for i, genre := range genres {
		url += genre
		if i < len(genres)-1 {
			url += "+"
		}
	}

	return url
}
