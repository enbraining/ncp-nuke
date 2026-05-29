// Package web serves a browser-based UI for ncp-nuke, reusing the shared
// runner logic. Live progress is streamed to the browser via Server-Sent Events.
package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"ncp-nuke/pkg/config"
	"ncp-nuke/pkg/excel"
	"ncp-nuke/pkg/ncp"
	"ncp-nuke/pkg/runner"
)

//go:embed static/*
var staticFS embed.FS

const confirmPhrase = "CONFIRM DELETE"

// Server holds the loaded accounts and optional resource filter config.
type Server struct {
	accounts []ncp.RootAccount
	cfg      *config.Config
	mu       sync.Mutex // serialize destructive runs
}

// NewServer loads accounts from the given Excel file (and optional config JSON).
func NewServer(filePath, configPath string) (*Server, error) {
	accounts, err := excel.ReadAccounts(filePath)
	if err != nil {
		return nil, err
	}
	var cfg *config.Config
	if configPath != "" {
		cfg, err = config.LoadConfig(configPath)
		if err != nil {
			return nil, err
		}
	}
	return &Server{accounts: accounts, cfg: cfg}, nil
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	staticContent, _ := staticFS.ReadFile("static/index.html")
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(staticContent)
	})

	mux.HandleFunc("/api/accounts", s.handleAccounts)
	mux.HandleFunc("/api/run", s.handleRun)
	return mux
}

type accountDTO struct {
	Index       int    `json:"index"`
	AccountName string `json:"accountName"`
	IamUsername string `json:"iamUsername"`
	AccessKey   string `json:"accessKey"`
}

func (s *Server) handleAccounts(w http.ResponseWriter, r *http.Request) {
	out := make([]accountDTO, 0, len(s.accounts))
	for i, a := range s.accounts {
		out = append(out, accountDTO{
			Index:       i,
			AccountName: a.AccountName,
			IamUsername: a.IamUsername,
			AccessKey:   maskKey(a.AccessKey),
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

type runRequest struct {
	Selected []int  `json:"selected"`
	Action   string `json:"action"` // activate | deactivate | nuke | list
	Password string `json:"password"`
	Cleanup  bool   `json:"cleanup"`
	Confirm  string `json:"confirm"`
}

// handleRun streams progress as Server-Sent Events.
func (s *Server) handleRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	var req runRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 요청: "+err.Error(), http.StatusBadRequest)
		return
	}

	switch req.Action {
	case "activate", "deactivate", "nuke", "list":
	default:
		http.Error(w, "알 수 없는 작업: "+req.Action, http.StatusBadRequest)
		return
	}
	if len(req.Selected) == 0 {
		http.Error(w, "선택된 계정이 없습니다", http.StatusBadRequest)
		return
	}

	// Destructive actions require the exact confirmation phrase.
	destructive := req.Action == "nuke" || (req.Action == "deactivate" && req.Cleanup)
	if destructive && req.Confirm != confirmPhrase {
		http.Error(w, fmt.Sprintf("삭제 확인 문구가 일치하지 않습니다. \"%s\" 를 정확히 입력하세요.", confirmPhrase), http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "스트리밍 미지원", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	selected := make(map[int]bool)
	for _, idx := range req.Selected {
		if idx >= 0 && idx < len(s.accounts) {
			selected[idx] = true
		}
	}

	// Serialize runs so concurrent destructive operations don't interleave.
	s.mu.Lock()
	defer s.mu.Unlock()

	send := func(line string) {
		fmt.Fprintf(w, "data: %s\n\n", escapeSSE(line))
		flusher.Flush()
	}

	runner.Process(s.accounts, selected, req.Action, req.Password, req.Cleanup, s.cfg, send)
	fmt.Fprint(w, "event: done\ndata: end\n\n")
	flusher.Flush()
}

func maskKey(k string) string {
	if len(k) <= 6 {
		return "******"
	}
	return k[:3] + "****" + k[len(k)-3:]
}

// escapeSSE keeps multi-line log messages valid in the SSE "data:" framing by
// turning newlines into spaces (the client renders one event per line).
func escapeSSE(s string) string {
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if r == '\n' || r == '\r' {
			out = append(out, ' ')
			continue
		}
		out = append(out, r)
	}
	return string(out)
}
