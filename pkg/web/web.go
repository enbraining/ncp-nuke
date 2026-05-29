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

	// Each runner log line becomes a structured event. Lines are grouped by the
	// resource they refer to (Server, NAT Gateway, ...) rather than by step. The
	// resource is detected from keywords; lines without a keyword stick to the
	// most recently seen resource. "currentResource" carries that state.
	currentResource := ""
	send := func(line string) {
		ev := classifyLine(line, &currentResource)
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
	Type     string `json:"type"`     // account | global | resource
	Text     string `json:"text"`
	Resource string `json:"resource"` // section title for type=resource
	Depth    int    `json:"depth"`
	Status   string `json:"status"` // ok | fail | skip | info
}

func classifyLine(line string, currentResource *string) progressEvent {
	trimmedLeft := strings.TrimLeft(line, "\n")
	leading := len(trimmedLeft) - len(strings.TrimLeft(trimmedLeft, " "))
	depth := leading / 2
	text := strings.TrimSpace(line)

	ev := progressEvent{Text: text, Depth: depth, Status: classifyStatus(text)}

	if strings.HasPrefix(text, "[루트 계정") {
		ev.Type = "account"
		*currentResource = ""
		return ev
	}

	if res := detectResource(text); res != "" {
		*currentResource = res
		ev.Type = "resource"
		ev.Resource = res
		return ev
	}

	// No resource keyword: a deeper-indented line belongs to the current resource;
	// a top-level phase line (e.g. "리소스 조회 중", "총 N개 ... 시작") is global.
	if depth >= 2 && *currentResource != "" {
		ev.Type = "resource"
		ev.Resource = *currentResource
		return ev
	}

	ev.Type = "global"
	*currentResource = ""
	return ev
}

// detectResource maps a log line to a resource category by keyword. Order
// matters: more specific keywords are checked before broader ones.
func detectResource(t string) string {
	switch {
	case strings.Contains(t, "Object Storage") || strings.Contains(t, "버킷"):
		return "Object Storage"
	case strings.Contains(t, "NKS"):
		return "NKS Cluster"
	case strings.Contains(t, "Auto Scaling"):
		return "Auto Scaling Group"
	case strings.Contains(t, "Launch Configuration"):
		return "Launch Configuration"
	case strings.Contains(t, "PostgreSQL") || strings.Contains(t, "(Pg)"):
		return "Cloud PostgreSQL"
	case strings.Contains(t, "MongoDB") || strings.Contains(t, "(Mongo)"):
		return "Cloud MongoDB"
	case strings.Contains(t, "MariaDB"):
		return "Cloud MariaDB"
	case strings.Contains(t, "MySQL"):
		return "Cloud MySQL"
	case strings.Contains(t, "Redis"):
		return "Cloud Redis"
	case strings.Contains(t, "Cloud DB"):
		return "Cloud DB"
	case strings.Contains(t, "Load Balancer") || strings.Contains(t, "로드밸런서"):
		return "Load Balancer"
	case strings.Contains(t, "Target Group"):
		return "Target Group"
	case strings.Contains(t, "NAS 스냅샷") || strings.Contains(t, "NAS Volume Snapshot"):
		return "NAS Snapshot"
	case strings.Contains(t, "NAS"):
		return "NAS Volume"
	case strings.Contains(t, "Block Storage") || strings.Contains(t, "블록 스토리지") || strings.Contains(t, "스토리지"):
		return "Block Storage"
	case strings.Contains(t, "서버"):
		return "Server"
	case strings.Contains(t, "Route Table") || strings.Contains(t, "경로"):
		return "Route Table"
	case strings.Contains(t, "VPC Peering"):
		return "VPC Peering"
	case strings.Contains(t, "NAT Gateway"):
		return "NAT Gateway"
	case strings.Contains(t, "공인 IP") || strings.Contains(t, "Public IP"):
		return "Public IP"
	case strings.Contains(t, "ACG") || strings.Contains(t, "Access Control"):
		return "Access Control Group"
	case strings.Contains(t, "Network ACL"):
		return "Network ACL"
	case strings.Contains(t, "Subnet"):
		return "Subnet"
	case strings.Contains(t, "VPC"):
		return "VPC"
	case strings.Contains(t, "Init Script"):
		return "Init Script"
	case strings.Contains(t, "Login Key"):
		return "Login Key"
	case strings.Contains(t, "Placement Group"):
		return "Placement Group"
	case strings.Contains(t, "서브 계정") || strings.Contains(t, "서브계정"):
		return "Sub Account"
	}
	return ""
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
