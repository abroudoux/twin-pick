package scrapper

import (
	"github.com/gocolly/colly/v2"

	"github.com/abroudoux/twinpick/internal/domain"
)

type CollectorFactory func() *colly.Collector

type LetterboxdScrapper struct {
	NewCollector    CollectorFactory
	GetTotalPages   func(url string) (int, error)
	GetFilmsOnPage  func(url string, page int) ([]domain.Film, error)
	GetPopularFilms func(url string) ([]domain.Film, error)
}

func NewLetterboxdScrapper() *LetterboxdScrapper {
	scrapper := &LetterboxdScrapper{}
	scrapper.NewCollector = func() *colly.Collector { return colly.NewCollector() }

	scrapper.GetTotalPages = func(url string) (int, error) {
		return scrapper.getTotalPagesImpl(url)
	}
	scrapper.GetFilmsOnPage = func(url string, page int) ([]domain.Film, error) {
		return scrapper.getFilmsOnPageImpl(url, page)
	}
	scrapper.GetPopularFilms = func(url string) ([]domain.Film, error) {
		return scrapper.getPopularFilmsImpl(url)
	}
	return scrapper
}
