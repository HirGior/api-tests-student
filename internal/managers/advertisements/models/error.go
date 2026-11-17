package models

import "encoding/json"

// ErrorResponse содержит сообщение и подробности ошибки API.
type ErrorResponse struct {
	Error   string          `json:"error"`
	Message string          `json:"message"`
	Details json.RawMessage `json:"details"`
}
