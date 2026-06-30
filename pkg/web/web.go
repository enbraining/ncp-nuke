// Package web serves a browser-based UI for ncp-nuke, reusing the shared
// runner logic. Live progress is streamed to the browser via Server-Sent Events.
package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"ncp-nuke/pkg/config"
	"ncp-nuke/pkg/excel"
	"ncp-nuke/pkg/ncp"
	"ncp-nuke/pkg/runner"
	"ncp-nuke/pkg/version"
)

//go:embed static/*
var staticFS embed.FS

const confirmPhrase = "CONFIRM DELETE"

// Server holds the loaded accounts and optional resource filter config.
type Server struct {
	accounts []ncp.RootAccount
	filePath string
	cfg      *config.Config
	Desktop  bool       // when true, file open/save use native OS dialogs
	mu       sync.Mutex // serialize destructive runs and account mutations
}

// NewServer optionally preloads accounts from an Excel file. filePath may be
// empty — accounts can then be uploaded from the browser via /api/upload.
func NewServer(filePath, configPath string) (*Server, error) {
	var accounts []ncp.RootAccount
	if filePath != "" {
		a, err := excel.ReadAccounts(filePath)
		if err != nil {
			return nil, err
		}
		accounts = a
	}
	var cfg *config.Config
	if configPath != "" {
		c, err := config.LoadConfig(configPath)
		if err != nil {
			return nil, err
		}
		cfg = c
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
	mux.HandleFunc("/api/upload", s.handleUpload)
	mux.HandleFunc("/api/template", s.handleTemplate)
	mux.HandleFunc("/api/desktop/pick-accounts", s.handlePickAccounts)
	mux.HandleFunc("/api/desktop/save-template", s.handleSaveTemplate)
	mux.HandleFunc("/api/env", s.handleEnv)
	mux.HandleFunc("/api/update/check", s.handleUpdateCheck)
	mux.HandleFunc("/api/open-url", s.handleOpenURL)
	mux.HandleFunc("/api/scan", s.handleScan)
	mux.HandleFunc("/api/execute", s.handleExecute)
	return mux
}

type accountDTO struct {
	Index       int    `json:"index"`
	AccountName string `json:"accountName"`
	IamUsername string `json:"iamUsername"`
	AccessKey   string `json:"accessKey"`
}

// handleUpload accepts an uploaded accounts .xlsx, parses it, and replaces the
// in-memory account list. The file is saved server-side so later "계정 추가"
// (AppendAccount) and downloads keep working.
func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseMultipartForm(20 << 20); err != nil {
		http.Error(w, "업로드 파싱 실패: "+err.Error(), http.StatusBadRequest)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "파일이 없습니다 (field: file)", http.StatusBadRequest)
		return
	}
	defer file.Close()

	tmp, err := os.CreateTemp("", "ncp-nuke-*.xlsx")
	if err != nil {
		http.Error(w, "임시 파일 생성 실패: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := io.Copy(tmp, file); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		http.Error(w, "파일 저장 실패: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmp.Close()

	accounts, err := excel.ReadAccounts(tmp.Name())
	if err != nil {
		os.Remove(tmp.Name())
		http.Error(w, "엑셀 파싱 실패: "+err.Error(), http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	if s.filePath != "" && strings.Contains(s.filePath, "ncp-nuke-") {
		os.Remove(s.filePath) // clean up a previously uploaded temp file
	}
	s.accounts = accounts
	s.filePath = tmp.Name()
	s.mu.Unlock()

	s.listAccounts(w)
}

// handleEnv tells the frontend whether it is running inside the desktop app
// (so it can use native file dialogs instead of browser upload/download).
func (s *Server) handleEnv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"desktop": s.Desktop, "version": version.Version})
}

// handleUpdateCheck queries the latest GitHub release and reports whether an
// update is available, with an OS-appropriate download URL.
func (s *Server) handleUpdateCheck(w http.ResponseWriter, r *http.Request) {
	cli := &http.Client{Timeout: 8 * time.Second}
	req, _ := http.NewRequest("GET", "https://api.github.com/repos/"+version.Repo+"/releases/latest", nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := cli.Do(req)
	if err != nil {
		http.Error(w, "업데이트 확인 실패: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		http.Error(w, fmt.Sprintf("GitHub 응답 %d", resp.StatusCode), http.StatusBadGateway)
		return
	}
	var rel struct {
		TagName string `json:"tag_name"`
		HTMLURL string `json:"html_url"`
		Assets  []struct {
			Name string `json:"name"`
			URL  string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		http.Error(w, "응답 파싱 실패: "+err.Error(), http.StatusBadGateway)
		return
	}

	latest := strings.TrimPrefix(rel.TagName, "v")
	download := rel.HTMLURL
	for _, a := range rel.Assets {
		if matchOSAsset(a.Name) {
			download = a.URL
			break
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"current":         version.Version,
		"latest":          latest,
		"updateAvailable": semverLess(version.Version, latest),
		"htmlUrl":         rel.HTMLURL,
		"downloadUrl":     download,
	})
}

// matchOSAsset picks the release asset matching the current OS.
func matchOSAsset(name string) bool {
	switch runtime.GOOS {
	case "darwin":
		return strings.HasSuffix(name, ".dmg")
	case "windows":
		return strings.Contains(name, "_setup.exe")
	}
	return false
}

// semverLess reports whether a < b for dotted numeric versions.
func semverLess(a, b string) bool {
	pa, pb := strings.Split(a, "."), strings.Split(b, ".")
	for i := 0; i < 3; i++ {
		var x, y int
		if i < len(pa) {
			fmt.Sscanf(pa[i], "%d", &x)
		}
		if i < len(pb) {
			fmt.Sscanf(pb[i], "%d", &y)
		}
		if x != y {
			return x < y
		}
	}
	return false
}

type openURLRequest struct {
	URL string `json:"url"`
}

// handleOpenURL opens a URL in the system browser (desktop only).
func (s *Server) handleOpenURL(w http.ResponseWriter, r *http.Request) {
	if !s.Desktop {
		http.Error(w, "데스크톱 모드에서만 사용 가능", http.StatusBadRequest)
		return
	}
	var req openURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.URL == "" {
		http.Error(w, "url이 필요합니다", http.StatusBadRequest)
		return
	}
	if !strings.HasPrefix(req.URL, "https://") {
		http.Error(w, "https URL만 허용", http.StatusBadRequest)
		return
	}
	if err := openURL(req.URL); err != nil {
		http.Error(w, "브라우저 열기 실패: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// handlePickAccounts (desktop) opens a native file chooser, loads the picked
// .xlsx, and replaces the in-memory accounts.
func (s *Server) handlePickAccounts(w http.ResponseWriter, r *http.Request) {
	if !s.Desktop {
		http.Error(w, "데스크톱 모드에서만 사용 가능", http.StatusBadRequest)
		return
	}
	path, cancelled, err := chooseFileDialog()
	if cancelled {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		http.Error(w, "파일 선택 실패: "+err.Error(), http.StatusInternalServerError)
		return
	}
	accounts, err := excel.ReadAccounts(path)
	if err != nil {
		http.Error(w, "엑셀 파싱 실패: "+err.Error(), http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	s.accounts = accounts
	s.filePath = path // their real file — "계정 추가" appends back here
	s.mu.Unlock()
	s.listAccounts(w)
}

// handleSaveTemplate (desktop) writes the template to a folder the user picks
// (falls back to ~/Downloads) and returns the saved path.
func (s *Server) handleSaveTemplate(w http.ResponseWriter, r *http.Request) {
	if !s.Desktop {
		http.Error(w, "데스크톱 모드에서만 사용 가능", http.StatusBadRequest)
		return
	}
	dir, cancelled, err := chooseFolderDialog()
	if cancelled {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil || dir == "" {
		home, _ := os.UserHomeDir()
		dir = filepath.Join(home, "Downloads")
	}
	dest := filepath.Join(dir, "accounts_template.xlsx")
	if err := excel.WriteTemplate(dest); err != nil {
		http.Error(w, "템플릿 저장 실패: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"path": dest})
}

// chooseFileDialog / chooseFolderDialog are implemented per-OS in
// dialog_{windows,darwin,other}.go.

// handleTemplate serves the accounts Excel template as a download.
func (s *Server) handleTemplate(w http.ResponseWriter, r *http.Request) {
	b, err := excel.TemplateBytes()
	if err != nil {
		http.Error(w, "템플릿 생성 실패: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="accounts_template.xlsx"`)
	w.Write(b)
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

func (s *Server) selectedMap(idxs []int) map[int]bool {
	m := make(map[int]bool)
	for _, idx := range idxs {
		if idx >= 0 && idx < len(s.accounts) {
			m[idx] = true
		}
	}
	return m
}

// --- Scan: list resources for the selected accounts ---

type scanRequest struct {
	Selected []int `json:"selected"`
}

type resourceCountDTO struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

type itemDTO struct {
	Account string `json:"account"`
	Name    string `json:"name"`
	ID      string `json:"id"`
}

type scanResponse struct {
	Accounts []string             `json:"accounts"`
	Types    []resourceCountDTO   `json:"types"`
	Warnings []string             `json:"warnings"`
	Details  map[string][]itemDTO `json:"details"` // resource type key -> items across accounts
}

// handleScan aggregates resource counts (by type) across the selected accounts.
func (s *Server) handleScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	var req scanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 요청: "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(req.Selected) == 0 {
		http.Error(w, "선택된 계정이 없습니다", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	selected := s.selectedMap(req.Selected)

	// Collect the accounts to scan (preserving order for deterministic output).
	type acctScan struct {
		name    string
		summary *ncp.ResourceSummary
		errs    []error
	}
	var jobs []*acctScan
	for i, acc := range s.accounts {
		if selected[i] {
			jobs = append(jobs, &acctScan{name: acc.AccountName})
		}
	}

	// Scan each account in parallel — each makes its own (sequential) set of
	// resource list calls with an independent client.
	var wg sync.WaitGroup
	idx := 0
	for i, acc := range s.accounts {
		if !selected[i] {
			continue
		}
		job := jobs[idx]
		idx++
		wg.Add(1)
		go func(a ncp.RootAccount, j *acctScan) {
			defer wg.Done()
			client := ncp.NewClient(a.AccessKey, a.SecretKey)
			j.summary, j.errs = client.ListAllResources()
		}(acc, job)
	}
	wg.Wait()

	// Aggregate in account order so type ordering is stable.
	counts := map[string]int{}
	var order []string
	var names []string
	var warnings []string
	details := map[string][]itemDTO{}
	for _, j := range jobs {
		names = append(names, j.name)
		for _, e := range j.errs {
			warnings = append(warnings, friendlyScanErr(j.name, e))
		}
		if j.summary == nil {
			continue
		}
		for _, bc := range j.summary.Breakdown() {
			if _, ok := counts[bc.Name]; !ok {
				order = append(order, bc.Name)
			}
			counts[bc.Name] += bc.Count
		}
		for typeName, items := range j.summary.Items() {
			for _, it := range items {
				details[typeName] = append(details[typeName], itemDTO{Account: j.name, Name: it.Name, ID: it.ID})
			}
		}
	}

	resp := scanResponse{Accounts: names, Warnings: warnings, Details: details}
	for _, k := range order {
		resp.Types = append(resp.Types, resourceCountDTO{Key: k, Count: counts[k]})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func friendlyScanErr(account string, e error) string {
	msg := e.Error()
	for _, sig := range []string{"HTTP 403", "StatusCode: 403", "AccessDenied", "InvalidAccessKeyId"} {
		if strings.Contains(msg, sig) {
			return fmt.Sprintf("[%s] 권한 없음/미사용: %v", account, e)
		}
	}
	return fmt.Sprintf("[%s] %v", account, e)
}

// --- Execute: optional sub-account action + delete of selected resource types ---

type executeRequest struct {
	Selected    []int    `json:"selected"`
	SubAction   string   `json:"subAction"` // none | activate | deactivate
	Password    string   `json:"password"`
	DeleteTypes []string `json:"deleteTypes"` // resource type keys (Breakdown names)
	Confirm     string   `json:"confirm"`
}

func (s *Server) handleExecute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	var req executeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "잘못된 요청: "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(req.Selected) == 0 {
		http.Error(w, "선택된 계정이 없습니다", http.StatusBadRequest)
		return
	}
	switch req.SubAction {
	case "", "none", "activate", "deactivate":
	default:
		http.Error(w, "알 수 없는 서브계정 작업: "+req.SubAction, http.StatusBadRequest)
		return
	}
	if req.SubAction == "" {
		req.SubAction = "none"
	}
	if len(req.DeleteTypes) == 0 && req.SubAction == "none" {
		http.Error(w, "수행할 작업이 없습니다 (삭제할 리소스 또는 서브계정 작업을 선택하세요)", http.StatusBadRequest)
		return
	}
	if len(req.DeleteTypes) > 0 && req.Confirm != confirmPhrase {
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

	s.mu.Lock()
	defer s.mu.Unlock()

	selected := s.selectedMap(req.Selected)
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

	// 1) Sub-account action (independent of resource deletion).
	if req.SubAction == "activate" || req.SubAction == "deactivate" {
		runner.Process(s.accounts, selected, req.SubAction, req.Password, false, nil, send)
		currentResource = ""
	}

	// 2) Delete only the selected resource types (nuke with a filter that disables
	//    every other type).
	if len(req.DeleteTypes) > 0 {
		cfg := buildDeleteConfig(req.DeleteTypes)
		runner.Process(s.accounts, selected, "nuke", "", false, cfg, send)
	}

	fmt.Fprint(w, "event: done\ndata: end\n\n")
	flusher.Flush()
}

// buildDeleteConfig returns a config that enables ONLY the given resource types
// (by Breakdown name) and disables all others, so a nuke deletes just those.
func buildDeleteConfig(types []string) *config.Config {
	want := make(map[string]bool, len(types))
	for _, t := range types {
		want[t] = true
	}
	disabled := func(name string) config.ResourceFilter {
		on := want[name]
		return config.ResourceFilter{Enabled: &on}
	}
	return &config.Config{
		Servers:               disabled("Server"),
		BlockStorages:         disabled("Block Storage"),
		BlockStorageSnapshots: disabled("Block Storage Snapshot"),
		PublicIps:             disabled("Public IP"),
		NasVolumes:            disabled("NAS Volume"),
		NasVolumeSnapshots:    disabled("NAS Volume Snapshot"),
		LoadBalancers:         disabled("Load Balancer"),
		TargetGroups:          disabled("Target Group"),
		CloudDBs:              disabled("Cloud DB"),
		CloudPostgresqls:      disabled("Cloud PostgreSQL"),
		CloudMongoDBs:         disabled("Cloud MongoDB"),
		CloudMariaDBs:         disabled("Cloud MariaDB"),
		CloudMySQLs:           disabled("Cloud MySQL"),
		CloudRedises:          disabled("Cloud Redis"),
		Vpcs:                  disabled("VPC"),
		Subnets:               disabled("Subnet"),
		NatGateways:           disabled("NAT Gateway"),
		VpcPeerings:           disabled("VPC Peering"),
		NetworkAcls:           disabled("Network ACL"),
		AccessControlGroups:   disabled("Access Control Group"),
		AutoScalingGroups:     disabled("Auto Scaling Group"),
		LaunchConfigurations:  disabled("Launch Configuration"),
		NksClusters:           disabled("NKS Cluster"),
		InitScripts:           disabled("Init Script"),
		LoginKeys:             disabled("Login Key"),
		PlacementGroups:       disabled("Placement Group"),
		Buckets:               disabled("Object Storage Bucket"),
	}
}

type progressEvent struct {
	Type     string `json:"type"` // account | global | resource
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
		return "Object Storage Bucket"
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
		return "NAS Volume Snapshot"
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
