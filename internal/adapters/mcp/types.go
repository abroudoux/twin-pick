package mcp

import (
	"encoding/json"

	"github.com/abroudoux/twinpick/internal/application"
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
	Duration  int      `json:"duration,omitempty"`
}

type Server struct {
	PickService *application.PickService
	SpotService *application.SpotService
}
