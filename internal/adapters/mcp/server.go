package mcp

import (
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
	Limit     int      `json:"limit,omitempty"`
}

type Server struct {
	PickService *application.PickService
}

func NewServer(ps *application.PickService) *Server {
	return &Server{PickService: ps}
}

func (s *Server) Run() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		return
	}

	var req Request
	if err := json.Unmarshal(input, &req); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid JSON: %v\n", err)
		return
	}

	var call ToolCall
	if err := json.Unmarshal(req.Params, &call); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid params: %v\n", err)
		return
	}

	switch call.Name {
	case "pick":
		var args ProgramArgs
		if err := json.Unmarshal(call.Arguments, &args); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid arguments: %v\n", err)
			return
		}

		params := domain.NewScrapperParams(args.Genres, args.Platform)
		films, err := s.PickService.Pick(args.Usernames, params, args.Limit)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Pick error: %v\n", err)
			return
		}

		resp := Response{
			ID: req.ID,
			Result: map[string]interface{}{
				"usernames": args.Usernames,
				"genres":    args.Genres,
				"platform":  args.Platform,
				"films":     films,
			},
		}
		json.NewEncoder(os.Stdout).Encode(resp)
	default:
		json.NewEncoder(os.Stdout).Encode(Response{ID: req.ID, Error: "Unknown tool"})
	}
}

func (s *Server) handleToolCall(req Request, encoder *json.Encoder) {
	var call ToolCall
	if err := json.Unmarshal(req.Params, &call); err != nil {
		encoder.Encode(Response{ID: req.ID, Error: err.Error()})
		return
	}

	switch call.Name {
	case "pick":
		s.handlePick(req, call, encoder)
	default:
		encoder.Encode(Response{ID: req.ID, Error: "Unknown tool"})
	}
}

func (s *Server) handlePick(req Request, call ToolCall, encoder *json.Encoder) {
	var args ProgramArgs
	if err := json.Unmarshal(call.Arguments, &args); err != nil {
		encoder.Encode(Response{ID: req.ID, Error: err.Error()})
		return
	}

	params := domain.NewScrapperParams(args.Genres, args.Platform)

	films, err := s.PickService.Pick(args.Usernames, params, args.Limit)
	if err != nil {
		encoder.Encode(Response{ID: req.ID, Error: err.Error()})
		return
	}

	encoder.Encode(Response{
		ID: req.ID,
		Result: map[string]interface{}{
			"usernames": args.Usernames,
			"genres":    args.Genres,
			"platform":  args.Platform,
			"films":     films,
		},
	})
}
