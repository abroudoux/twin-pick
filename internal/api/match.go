package api

import (
	"net/http"
	"strings"

	"github.com/abroudoux/twinpick/internal/match"
	"github.com/abroudoux/twinpick/internal/scrapper"
	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
)

func (s *server) handleMatch(c *gin.Context) {
	usernamesParam := c.Param("usernames")
	if usernamesParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Usernames parameter is required"})
		return
	}

	usernames := strings.Split(usernamesParam, ",")
	log.Infof("Matching watchlists for users: %v", usernames)

	var genres []string
	genresParam := c.Param("genres")
	if genresParam != "" {
		genres = strings.Split(genresParam, ",")
		log.Infof("Filtering watchlists by genres: %v", genres)
	}

	watchlists := scrapper.ScrapUsersWachtlists(usernames, genres)

	commonFilms, err := match.GetCommonFilms(watchlists)
	if err != nil {
		log.Error("Error while matching watchlists: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while matching watchlists"})
		return
	}
	if len(commonFilms) == 0 {
		log.Info("No common films found among the watchlists.")
		c.JSON(http.StatusOK, gin.H{"message": "No common films found among the watchlists."})
		return
	}

	selectedFilm, err := match.SelectRandomFilm(commonFilms)
	if err != nil {
		log.Error("Error while selecting a random film: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while selecting a random film"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"selected_film": selectedFilm})
}
