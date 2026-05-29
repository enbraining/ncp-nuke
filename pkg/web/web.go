// Package web serves a browser-based UI for ncp-nuke, reusing the shared
// runner logic. Live progress is streamed to the browser via Server-Sent Events.
package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	filePath string
	cfg      *config.Config
	mu       sync.Mutex // serialize destructive runs and account mutations
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
	return &Server{accounts: accounts, filePath: filePath, cfg: cfg}, nil
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
	switch r.Method {
	case http.MethodGet:
		s.listAccounts(w)
	case http.MethodPost:
		s.addAccount(w, r)
	default:
		http.Error(w, "GET 또는 POST만 지원", http.StatusMethodNotAllowed)
	}
}

func (s *Server) listAccounts(w http.ResponseWriter) {
	s.mu.Lock()
	defer s.mu.Unlock()
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

type newAccountRequest struct {
	AccountName string `json:"accountName"`
	AccessKey   string `json:"accessKey"`
	SecretKey   string `json:"secretKey"`
	IamUsername string `json:"iamUsername"`
	Password    string `json:"password"`
}

// addAccount appends a new account to both the in-memory list and the Excel file.
func (s *Server) addAccount(w http.ResponseWriter, r *http.Request) {
	var req newAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 요청: "+err.Error(), http.StatusBadRequest)
		return
	}
	acc := ncp.RootAccount{
		AccountName: req.AccountName,
		AccessKey:   req.AccessKey,
		SecretKey:   req.SecretKey,
		IamUsername: req.IamUsername,
		Password:    req.Password,
	}
	if acc.AccessKey == "" || acc.SecretKey == "" {
		http.Error(w, "AccessKey와 SecretKey는 필수입니다", http.StatusBadRequest)
		return
	}
	if acc.IamUsername == "" {
		http.Error(w, "IAM Username은 필수입니다", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if acc.AccountName == "" {
		acc.AccountName = fmt.Sprintf("Account-%d", len(s.accounts)+1)
	}

	// Persist to the Excel file first; only update memory if the write succeeds.
	if err := excel.AppendAccount(s.filePath, acc); err != nil {
		http.Error(w, "엑셀 저장 실패: "+err.Error(), http.StatusInternalServerError)
		return
	}
	s.accounts = append(s.accounts, acc)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accountDTO{
		Index:       len(s.accounts) - 1,
		AccountName: acc.AccountName,
		IamUsername: acc.IamUsername,
		AccessKey:   maskKey(acc.AccessKey),
	})
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

	// Each runner log line is turned into a structured event. The hierarchy is
	// derived from leading-space depth: depth 0 = global, 1 = section, 2+ = detail.
	send := func(line string) {
		ev := classifyLine(line)
		if ev.Text == "" {
			return
		}
		b, _ := json.Marshal(ev)
		fmt.Fprintf(w, "data: %s\n\n", b)
		flusher.Flush()
	}

	runner.Process(s.accounts, selected, req.Action, req.Password, req.Cleanup, s.cfg, send)
	fmt.Fprint(w, "event: done\ndata: end\n\n")
	flusher.Flush()
}

type progressEvent struct {
	Type   string `json:"type"`   // account | global | section | detail
	Text   string `json:"text"`
	Depth  int    `json:"depth"`
	Status string `json:"status"` // ok | fail | skip | info
}

func classifyLine(line string) progressEvent {
	trimmedLeft := strings.TrimLeft(line, "\n")
	leading := len(trimmedLeft) - len(strings.TrimLeft(trimmedLeft, " "))
	depth := leading / 2
	text := strings.TrimSpace(line)

	ev := progressEvent{Text: text, Depth: depth, Status: classifyStatus(text)}
	switch {
	case strings.HasPrefix(text, "[루트 계정"):
		ev.Type = "account"
	case depth <= 0:
		ev.Type = "global"
	case depth == 1:
		ev.Type = "section"
	default:
		ev.Type = "detail"
	}
	return ev
}

func classifyStatus(t string) string {
	switch {
	case strings.Contains(t, "[실패]") || strings.Contains(t, "[오류]"):
		return "fail"
	case strings.Contains(t, "[성공]"):
		return "ok"
	case strings.Contains(t, "[건너뜀]") || strings.Contains(t, "[경고]") || strings.Contains(t, "[대기]"):
		return "skip"
	default:
		return "info"
	}
}

func maskKey(k string) string {
	if len(k) <= 6 {
		return "******"
	}
	return k[:3] + "****" + k[len(k)-3:]
}
