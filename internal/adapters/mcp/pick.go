package mcp

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/abroudoux/twinpick/internal/domain"
)

func (s *Server) pickTool(req Request, call ToolCall) {
	var args ProgramArgs
	if err := json.Unmarshal(call.Arguments, &args); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid arguments: %v\n", err)
		return
	}

	pickParams := domain.NewPickParams(args.Usernames, domain.NewScrapperParams(args.Genres, args.Platform), args.Limit)

	films, err := s.PickService.Pick(pickParams)
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
}
