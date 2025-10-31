package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/abroudoux/twinpick/internal/application"
	"github.com/abroudoux/twinpick/internal/domain"
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

type ProgramArgs struct {
	Usernames []string `json:"usernames"`
	Genres    []string `json:"genres,omitempty"`
	Platform  string   `json:"platform,omitempty"`
}

type Server struct {
	MatchService  *application.MatchService
	CommonService *application.CommonService
}

func NewServer(m *application.MatchService, c *application.CommonService) *Server {
	return &Server{MatchService: m, CommonService: c}
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
			Error: fmt.Sprintf("Invalid method: %s", req.Method),
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
	case "match_film":
		s.handleMatchFilm(req, call, encoder)
	case "common_films":
		s.handleCommonFilms(req, call, encoder)
	default:
		encoder.Encode(Response{ID: req.ID, Error: "Unknown tool"})
	}
}

func (s *Server) handleMatchFilm(req Request, call ToolCall, encoder *json.Encoder) {
	var args ProgramArgs
	if err := json.Unmarshal(call.Arguments, &args); err != nil {
		encoder.Encode(Response{ID: req.ID, Error: err.Error()})
		return
	}

	params := domain.NewScrapperParams(args.Genres, args.Platform)

	film, err := s.MatchService.MatchFilm(args.Usernames, params)
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
			"selected_film": film.Title,
		},
	})
}

func (s *Server) handleCommonFilms(req Request, call ToolCall, encoder *json.Encoder) {
	var args ProgramArgs
	if err := json.Unmarshal(call.Arguments, &args); err != nil {
		encoder.Encode(Response{ID: req.ID, Error: err.Error()})
		return
	}

	params := domain.NewScrapperParams(args.Genres, args.Platform)

	films, err := s.CommonService.GetCommonFilms(args.Usernames, params)
	if err != nil {
		encoder.Encode(Response{ID: req.ID, Error: err.Error()})
		return
	}

	encoder.Encode(Response{
		ID: req.ID,
		Result: map[string]interface{}{
			"usernames":    args.Usernames,
			"genres":       args.Genres,
			"platform":     args.Platform,
			"common_films": films,
		},
	})
}
