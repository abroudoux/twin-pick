package mcp

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/abroudoux/twinpick/internal/domain"
)

func (s *Server) spotTool(req Request, call ToolCall) {
	var args ProgramArgs
	if err := json.Unmarshal(call.Arguments, &args); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid arguments: %v\n", err)
		return
	}

	spotParams := domain.NewSpotParams(domain.NewScrapperParams(args.Genres, args.Platform), args.Limit)

	films, err := s.SpotService.Spot(spotParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Spot error: %v\n", err)
		return
	}

	resp := Response{
		ID: req.ID,
		Result: map[string]interface{}{
			"genres":   args.Genres,
			"platform": args.Platform,
			"films":    films,
		},
	}
	json.NewEncoder(os.Stdout).Encode(resp)
}
