package httpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// MessageRequest represents the incoming JSON payload.
type MessageRequest struct {
	Type string `json:"type"`
	Msg  string `json:"msg,omitempty"`
}

// MessageResponse for repeat
type RepeatResponse struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

// MessageResponse for time
type TimeResponse struct {
	Type string `json:"type"`
	Time string `json:"time"`
}

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	log := getLogger(r).Named("hello")

	var req MessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		log.Error("invalid JSON", zap.Error(err))
		return
	}

	switch req.Type {
	case "repeat":
		resp := RepeatResponse{Type: "repeat", Msg: req.Msg}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			log.Error("failed to encode response", zap.Error(err))
			return
		}
	case "time":
		resp := TimeResponse{Type: "time", Time: time.Now().UTC().Format(time.RFC3339)}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			log.Error("failed to write response", zap.Error(err))
			return
		}
	default:
		http.Error(w, "unknown type", http.StatusBadRequest)
		log.Error("failed to write response")
	}
	log.Debug("handled message request")
}
