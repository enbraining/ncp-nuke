package tui

import (
	"context"
	"fmt"
	"strings"

	"ncp-nuke/pkg/config"
	"ncp-nuke/pkg/excel"
	"ncp-nuke/pkg/ncp"
	"ncp-nuke/pkg/runner"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// confirmPhrase is the exact text a user must type to authorize a destructive action.
const confirmPhrase = "CONFIRM DELETE"

// Application State
type sessionState int

const (
	stateSelectAccounts sessionState = iota
	stateSelectAction
	statePasswordInput
	stateConfirm
	stateTypingConfirm
	stateRunning
	stateDone
)

var (
	baseStyle  = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")).Bold(true)
)

type model struct {
	state          sessionState
	table          table.Model
	viewport       viewport.Model
	confirmInput   textinput.Model
	passwordInput  textinput.Model
	confirmErr     bool
	accounts       []ncp.RootAccount
	selected       map[int]bool
	action         string // "activate" or "deactivate"
	actionCursor   int
	globalPassword string
	cleanup        bool
	cfg            *config.Config
	logs           *strings.Builder
	logChan        chan string
	windowWidth    int
	windowHeight   int
}

func Start(filePath, configPath, accountFilter string) error {
	accounts, err := excel.ReadAccounts(filePath)
	if err != nil {
		return err
	}

	// Apply account filter
	if accountFilter != "" {
		var filtered []ncp.RootAccount
		for _, a := range accounts {
			if a.AccountName == accountFilter {
				filtered = append(filtered, a)
			}
		}
		accounts = filtered
	}
	if len(accounts) == 0 {
		return fmt.Errorf("대상 계정이 없습니다")
	}

	var cfg *config.Config
	if configPath != "" {
		cfg, err = config.LoadConfig(configPath)
		if err != nil {
			return err
		}
	}

	m := initialModel(accounts, cfg)
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		return err
	}
	return nil
}

func initialModel(accounts []ncp.RootAccount, cfg *config.Config) model {
	columns := []table.Column{
		{Title: "선택", Width: 6},
		{Title: "Account Name", Width: 20},
		{Title: "IAM Username", Width: 15},
		{Title: "Access Key", Width: 25},
	}

	rows := []table.Row{}
	for _, acc := range accounts {
		rows = append(rows, table.Row{"[ ]", acc.AccountName, acc.IamUsername, acc.AccessKey})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	ci := textinput.New()
	ci.Placeholder = confirmPhrase
	ci.CharLimit = 20
	ci.Width = 30

	pi := textinput.New()
	pi.Placeholder = "비밀번호 입력 (빈 값이면 자동 생성)"
	pi.CharLimit = 100
	pi.Width = 50
	pi.EchoMode = textinput.EchoPassword
	pi.EchoCharacter = '*'

	return model{
		state:        stateSelectAccounts,
		table:        t,
		viewport:     vp,
		confirmInput: ci,
		passwordInput: pi,
		accounts:     accounts,
		selected:     make(map[int]bool),
		cfg:          cfg,
		logs:         &strings.Builder{},
		logChan:      make(chan string, 100),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

type logMsg string
type doneMsg struct{}

func waitForLog(sub <-chan string) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-sub
		if !ok {
			return doneMsg{}
		}
		return logMsg(msg)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			if m.state != stateRunning && m.state != statePasswordInput && m.state != stateTypingConfirm {
				return m, tea.Quit
			}
		}

		switch m.state {
		case stateSelectAccounts:
			if msg.String() == " " {
				idx := m.table.Cursor()
				if m.selected[idx] {
					delete(m.selected, idx)
				} else {
					m.selected[idx] = true
				}
				m.table.SetRows(updateRows(m.accounts, m.selected))
			} else if msg.String() == "enter" {
				if len(m.selected) > 0 {
					m.state = stateSelectAction
					m.actionCursor = 0
				}
			}

		case stateSelectAction:
			switch msg.String() {
			case "up", "k":
				if m.actionCursor > 0 {
					m.actionCursor--
				}
			case "down", "j":
				if m.actionCursor < 3 {
					m.actionCursor++
				}
			case "enter":
				switch m.actionCursor {
				case 0:
					m.action = "activate"
					// Check if any selected account has no password in excel
					needsPassword := false
					for i := range m.selected {
						if m.accounts[i].Password == "" {
							needsPassword = true
							break
						}
					}
					if needsPassword {
						m.state = statePasswordInput
						m.passwordInput.Focus()
						return m, m.passwordInput.Cursor.BlinkCmd()
					}
					m.state = stateConfirm
				case 1:
					m.action = "deactivate"
					m.state = stateConfirm
				case 2:
					m.action = "nuke"
					m.state = stateConfirm
				case 3:
					m.action = "list"
					m.state = stateConfirm
				}
			case "b", "B", "esc":
				m.state = stateSelectAccounts
			}

		case statePasswordInput:
			switch msg.String() {
			case "enter":
				m.globalPassword = m.passwordInput.Value()
				m.state = stateConfirm
			case "esc":
				m.state = stateSelectAction
				m.passwordInput.Blur()
				m.passwordInput.Reset()
			default:
				m.passwordInput, cmd = m.passwordInput.Update(msg)
				return m, cmd
			}

		case stateConfirm:
			switch msg.String() {
			case "y", "Y":
				if (m.action == "deactivate" && m.cleanup) || m.action == "nuke" {
					m.state = stateTypingConfirm
					m.confirmInput.Reset()
					m.confirmInput.Focus()
					m.confirmErr = false
					return m, m.confirmInput.Cursor.BlinkCmd()
				}
				m.state = stateRunning
				go func() {
					runner.Process(context.Background(), m.accounts, m.selected, m.action, m.globalPassword, m.cleanup, m.cfg, func(s string) {
						m.logChan <- s
					})
					close(m.logChan)
				}()
				return m, waitForLog(m.logChan)

			case "c", "C":
				if m.action == "deactivate" {
					m.cleanup = !m.cleanup
				}
			case "b", "B", "esc":
				m.state = stateSelectAction
			}

		case stateTypingConfirm:
			switch msg.String() {
			case "enter":
				if m.confirmInput.Value() == confirmPhrase {
					m.state = stateRunning
					go func() {
						runner.Process(context.Background(), m.accounts, m.selected, m.action, m.globalPassword, m.cleanup, m.cfg, func(s string) {
							m.logChan <- s
						})
						close(m.logChan)
					}()
					return m, waitForLog(m.logChan)
				}
				m.confirmErr = true
			case "esc":
				m.state = stateConfirm
				m.confirmInput.Blur()
			}

		case stateDone:
			if msg.String() == "enter" {
				return m, tea.Quit
			}
		}

	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
		m.table.SetWidth(msg.Width - 10)
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - 10

	case logMsg:
		m.logs.WriteString(string(msg) + "\n")
		m.viewport.SetContent(m.logs.String())
		m.viewport.GotoBottom()
		return m, waitForLog(m.logChan)

	case doneMsg:
		m.state = stateDone
		m.logs.WriteString("\n=== 모든 작업 완료 ===\n[Enter]를 눌러 종료하거나 [q]로 나가세요.")
		m.viewport.SetContent(m.logs.String())
		m.viewport.GotoBottom()
	}

	switch m.state {
	case stateSelectAccounts:
		m.table, cmd = m.table.Update(msg)
	case stateTypingConfirm:
		m.confirmInput, cmd = m.confirmInput.Update(msg)
	case statePasswordInput:
		m.passwordInput, cmd = m.passwordInput.Update(msg)
	case stateRunning, stateDone:
		m.viewport, cmd = m.viewport.Update(msg)
	}

	return m, cmd
}

func updateRows(accounts []ncp.RootAccount, selected map[int]bool) []table.Row {
	rows := []table.Row{}
	for i, acc := range accounts {
		mark := "[ ]"
		if selected[i] {
			mark = "[x]"
		}

		rows = append(rows, table.Row{mark, acc.AccountName, acc.IamUsername, acc.AccessKey})
	}
	return rows
}

func (m model) View() string {
	if m.windowWidth == 0 {
		return "Loading..."
	}

	switch m.state {
	case stateSelectAccounts:
		return baseStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				titleStyle.Render("대상 계정 선택 (Space: 선택/해제, Enter: 다음)"),
				m.table.View(),
				fmt.Sprintf("\n선택된 계정: %d개", len(m.selected)),
			),
		)

	case stateSelectAction:
		actions := []string{"Sub Account 활성화", "Sub Account 비활성화", "리소스 전체 삭제", "리소스 전체 조회"}
		var items string
		for i, a := range actions {
			cursor := "  "
			if i == m.actionCursor {
				cursor = "> "
			}
			items += fmt.Sprintf("%s%s\n", cursor, a)
		}

		content := fmt.Sprintf(`
%s

선택된 계정: %d개

%s
(위/아래: 선택, Enter: 다음, b: 뒤로, q: 종료)
`, titleStyle.Render("수행할 작업 선택"), len(m.selected), items)
		return baseStyle.Render(content)

	case statePasswordInput:
		content := fmt.Sprintf(`
%s

엑셀에 비밀번호가 없는 계정이 있습니다.
공통으로 적용할 비밀번호를 입력하세요.
(빈 값으로 Enter 시 자동 생성)

%s

(Enter: 다음, Esc: 뒤로)
`, titleStyle.Render("비밀번호 입력"), m.passwordInput.View())
		return baseStyle.Render(content)

	case stateConfirm:
		actionLabel := "활성화 + 비밀번호 초기화"
		switch m.action {
		case "deactivate":
			actionLabel = "비활성화"
		case "nuke":
			actionLabel = "리소스 전체 삭제 (Nuke)"
		case "list":
			actionLabel = "리소스 목록 조회"
		}

		targets := ""
		count := 0
		for i := range m.selected {
			if count < 5 {
				targets += fmt.Sprintf("- %s\n", m.accounts[i].AccountName)
			}
			count++
		}
		if count > 5 {
			targets += fmt.Sprintf("... 외 %d개\n", count-5)
		}

		options := ""
		if m.action == "deactivate" {
			cleanupStatus := "[ ] Cleanup (서비스 해지 및 리소스 삭제)"
			if m.cleanup {
				cleanupStatus = "[x] Cleanup (서비스 해지 및 리소스 삭제)"
			}
			options = fmt.Sprintf("\n옵션:\n%s (Toggle: 'c')\n", cleanupStatus)
		} else if m.action == "nuke" {
			warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true)
			options = "\n" + warningStyle.Render("⚠ 선택한 계정의 모든 리소스(서버/스토리지/IP/DB/VPC 등)를 영구 삭제합니다. (서브 계정은 유지)") + "\n"
		}

		content := fmt.Sprintf(`
%s

작업: %s

선택된 계정:
%s
%s
진행하시겠습니까? (y: 시작, b: 뒤로, q: 종료)
`, titleStyle.Render("작업 확인"), actionLabel, targets, options)
		return baseStyle.Render(content)

	case stateTypingConfirm:
		warningStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true)
		errMsg := ""
		if m.confirmErr {
			errMsg = warningStyle.Render("\n입력이 일치하지 않습니다. 다시 입력해주세요.")
		}
		content := fmt.Sprintf(`
%s

%s
정말로 모든 리소스를 삭제하시겠습니까?
아래에 "%s"를 정확히 입력하세요.

%s
%s

(Enter: 확인, Esc: 뒤로)
`, titleStyle.Render("안전 확인"), warningStyle.Render("⚠ 이 작업은 되돌릴 수 없습니다!"), confirmPhrase, m.confirmInput.View(), errMsg)
		return baseStyle.Render(content)

	case stateRunning, stateDone:
		return baseStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				titleStyle.Render("작업 진행 중..."),
				m.viewport.View(),
			),
		)
	}

	return ""
}
