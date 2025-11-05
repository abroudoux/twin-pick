package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/abroudoux/twinpick/internal/application"
)

func NewServer(ps *application.PickService, ss *application.SpotService) *Server {
	return &Server{PickService: ps, SpotService: ss}
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
		s.pickTool(req, call)
	case "spot":
		s.spotTool(req, call)
	default:
		json.NewEncoder(os.Stdout).Encode(Response{ID: req.ID, Error: "Unknown tool"})
	}
}
