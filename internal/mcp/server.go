package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/abroudoux/twinpick/internal/core"
	"github.com/abroudoux/twinpick/internal/scrapper"
)

type Request struct {
	ID     int             `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type Response struct {
	ID     int         `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

type ToolCall struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

type FindCommonFilmArgs struct {
	Usernames []string `json:"usernames"`
	Genres    []string `json:"genres,omitempty"`
	Platform  string   `json:"platform,omitempty"`
}

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run() {
	reader := bufio.NewReader(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	input, err := io.ReadAll(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		return
	}

	if len(input) == 0 {
		encoder.Encode(Response{Error: "Empty input"})
		return
	}

	var req Request
	if err := json.Unmarshal(input, &req); err != nil {
		encoder.Encode(Response{Error: fmt.Sprintf("Invalid JSON: %v", err)})
		return
	}

	switch req.Method {
	case "tools/call":
		s.handleToolCall(req, encoder)
	default:
		encoder.Encode(Response{
			ID:    req.ID,
			Error: fmt.Sprintf("Invalid method : %s", req.Method),
		})
	}
}

func (s *Server) handleToolCall(req Request, encoder *json.Encoder) {
	var call ToolCall
	if err := json.Unmarshal(req.Params, &call); err != nil {
		encoder.Encode(Response{ID: req.ID, Error: err.Error()})
		return
	}

	switch call.Name {
	case "find_common_film":
		s.handleFindCommonFilm(req, call, encoder)
	default:
		encoder.Encode(Response{ID: req.ID, Error: "Tool inconnu"})
	}
}

func (s *Server) handleFindCommonFilm(req Request, call ToolCall, encoder *json.Encoder) {
	var args FindCommonFilmArgs
	if err := json.Unmarshal(call.Arguments, &args); err != nil {
		encoder.Encode(Response{ID: req.ID, Error: err.Error()})
		return
	}

	scrapperParams := scrapper.NewScrapperParams(args.Usernames, args.Genres, args.Platform)
	watchlists := scrapper.ScrapUsersWachtlists(scrapperParams)
	commonFilms, err := core.GetCommonFilms(watchlists)
	if err != nil {
		encoder.Encode(Response{ID: req.ID, Error: err.Error()})
		return
	}

	selectedFilm, err := core.SelectRandomFilm(commonFilms)
	if err != nil {
		encoder.Encode(Response{ID: req.ID, Error: err.Error()})
		return
	}

	encoder.Encode(Response{
		ID: req.ID,
		Result: map[string]interface{}{
			"usernames":     args.Usernames,
			"genres":        args.Genres,
			"platform":      args.Platform,
			"common_films":  commonFilms,
			"selected_film": selectedFilm,
		},
	})
}
