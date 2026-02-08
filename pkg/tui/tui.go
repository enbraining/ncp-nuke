package tui

import (
	"fmt"
	"strings"

	"ncp-nuke/pkg/config"
	"ncp-nuke/pkg/excel"
	"ncp-nuke/pkg/ncp"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Application State
type sessionState int

const (
	stateSelectAccounts sessionState = iota
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
	confirmErr     bool
	accounts       []ncp.RootAccount
	selected       map[int]bool // Index of selected accounts in `accounts`
	cleanup        bool
	cfg            *config.Config
	logs           *strings.Builder
	logChan        chan string
	windowWidth    int
	windowHeight   int
}

func Start(filePath, configPath string) error {
	accounts, err := excel.ReadAccounts(filePath)
	if err != nil {
		return err
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

	ti := textinput.New()
	ti.Placeholder = "REAL DELETE"
	ti.CharLimit = 20
	ti.Width = 30

	return model{
		state:        stateSelectAccounts,
		table:        t,
		viewport:     vp,
		confirmInput: ti,
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

// Command to wait for next log
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
			if m.state != stateRunning {
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
					m.state = stateConfirm
				}
			}

		case stateConfirm:
			switch msg.String() {
			case "y", "Y":
				if m.cleanup {
					m.state = stateTypingConfirm
					m.confirmInput.Reset()
					m.confirmInput.Focus()
					m.confirmErr = false
					return m, m.confirmInput.Cursor.BlinkCmd()
				}
				m.state = stateRunning
				go func() {
					processSelectedAccounts(m.accounts, m.selected, m.cleanup, m.cfg, func(s string) {
						m.logChan <- s
					})
					close(m.logChan)
				}()
				return m, waitForLog(m.logChan)

			case "c", "C":
				m.cleanup = !m.cleanup
			case "b", "B", "esc":
				m.state = stateSelectAccounts
			}

		case stateTypingConfirm:
			switch msg.String() {
			case "enter":
				if m.confirmInput.Value() == "REAL DELETE" {
					m.state = stateRunning
					go func() {
						processSelectedAccounts(m.accounts, m.selected, m.cleanup, m.cfg, func(s string) {
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
		m.viewport.Height = msg.Height - 10 // Leave room for title

	case logMsg:
		m.logs.WriteString(string(msg) + "\n")
		m.viewport.SetContent(m.logs.String())
		m.viewport.GotoBottom()
		return m, waitForLog(m.logChan) // Wait for next

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

	case stateConfirm:
		cleanupStatus := "[ ] Cleanup (서비스 해지 및 리소스 삭제)"
		if m.cleanup {
			cleanupStatus = "[x] Cleanup (서비스 해지 및 리소스 삭제)"
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

		content := fmt.Sprintf(`
%s

선택된 계정:
%s

옵션:
%s (Toggle: 'c')

진행하시겠습니까? (y: 시작, b: 뒤로, q: 종료)
`, titleStyle.Render("작업 확인"), targets, cleanupStatus)
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
아래에 "REAL DELETE"를 정확히 입력하세요.

%s
%s

(Enter: 확인, Esc: 뒤로)
`, titleStyle.Render("안전 확인"), warningStyle.Render("⚠ 이 작업은 되돌릴 수 없습니다!"), m.confirmInput.View(), errMsg)
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

// Logic copied/adapted from cmd/deactivate.go for TUI usage
func processSelectedAccounts(accounts []ncp.RootAccount, selected map[int]bool, cleanup bool, cfg *config.Config, logFn func(string)) {
	logFn("작업 시작...")

	totalSuccess, totalFail := 0, 0

	for i, account := range accounts {
		if !selected[i] {
			continue
		}

		logFn(fmt.Sprintf("\n[루트 계정: %s]", account.AccountName))
		client := ncp.NewClient(account.AccessKey, account.SecretKey)

		// 1. Cleanup Phase
		if cleanup {
			logFn("  리소스 조회 중...")
			summary, errs := client.ListAllResources()
			for _, e := range errs {
				logFn(fmt.Sprintf("    [경고] 조회 오류: %v", e))
			}

			if cfg != nil {
				// Filter Logic Duplicated from applyConfigFilter (below) for safety
				applyFilter(summary, cfg)
			}

			if summary.TotalCount() > 0 {
				logFn(fmt.Sprintf("  총 %d개 서비스 해지 및 리소스 삭제 시작...", summary.TotalCount()))
				s, f := client.CleanupAllResources(summary, func(msg string) {
					logFn("    " + msg)
				})
				logFn(fmt.Sprintf("  서비스 해지 및 리소스 삭제 결과: 성공 %d, 실패 %d", s, f))
			} else {
				logFn("  삭제할 리소스 없음")
			}
		}

		// 2. Deactivate Phase
		logFn("  서브 계정 조회 중...")
		subAccounts, err := client.ListSubAccounts()
		if err != nil {
			logFn(fmt.Sprintf("    [실패] 서브 계정 조회: %v", err))
			continue
		}

		var targets []ncp.SubAccount
		if account.IamUsername != "" {
			for _, sa := range subAccounts {
				if sa.LoginId == account.IamUsername {
					targets = append(targets, sa)
					break
				}
			}
		} else {
			targets = subAccounts
		}

		if len(targets) == 0 {
			logFn("    대상 서브 계정 없음")
			continue
		}

		for _, sa := range targets {
			if !sa.Active {
				logFn(fmt.Sprintf("    [건너뜀] %s: 이미 비활성", sa.LoginId))
				continue
			}
			if err := client.DeactivateSubAccount(sa.SubAccountId); err != nil {
				logFn(fmt.Sprintf("    [실패] %s 비활성화: %v", sa.LoginId, err))
				totalFail++
			} else {
				logFn(fmt.Sprintf("    [성공] %s 비활성화 완료", sa.LoginId))
				totalSuccess++
			}
		}
	}

	logFn(fmt.Sprintf("\n최종 결과: 서브계정 비활성화 성공 %d, 실패 %d", totalSuccess, totalFail))
}

func applyFilter(summary *ncp.ResourceSummary, cfg *config.Config) {
	var servers []ncp.ServerInstance
	for _, s := range summary.Servers {
		if cfg.Servers.Match(s.ServerName, s.ServerInstanceNo) {
			servers = append(servers, s)
		}
	}
	summary.Servers = servers

	var storages []ncp.BlockStorageInstance
	for _, s := range summary.BlockStorages {
		if cfg.BlockStorages.Match(s.BlockStorageName, s.BlockStorageInstanceNo) {
			storages = append(storages, s)
		}
	}
	summary.BlockStorages = storages

	var bsSnaps []ncp.BlockStorageSnapshotInstance
	for _, s := range summary.BlockStorageSnapshots {
		if cfg.BlockStorageSnapshots.Match(s.BlockStorageSnapshotName, s.BlockStorageSnapshotInstanceNo) {
			bsSnaps = append(bsSnaps, s)
		}
	}
	summary.BlockStorageSnapshots = bsSnaps

	var ips []ncp.PublicIpInstance
	for _, s := range summary.PublicIps {
		if cfg.PublicIps.Match(s.PublicIp, s.PublicIpInstanceNo) {
			ips = append(ips, s)
		}
	}
	summary.PublicIps = ips

	var vols []ncp.NasVolumeInstance
	for _, s := range summary.NasVolumes {
		if cfg.NasVolumes.Match(s.VolumeName, s.NasVolumeInstanceNo) {
			vols = append(vols, s)
		}
	}
	summary.NasVolumes = vols

	var nasSnaps []ncp.NasVolumeSnapshot
	for _, s := range summary.NasVolumeSnapshots {
		if cfg.NasVolumeSnapshots.Match(s.NasVolumeSnapshotName, s.NasVolumeSnapshotInstanceNo) {
			nasSnaps = append(nasSnaps, s)
		}
	}
	summary.NasVolumeSnapshots = nasSnaps

	var lbs []ncp.LoadBalancerInstance
	for _, s := range summary.LoadBalancers {
		if cfg.LoadBalancers.Match(s.LoadBalancerName, s.LoadBalancerInstanceNo) {
			lbs = append(lbs, s)
		}
	}
	summary.LoadBalancers = lbs

	var tgs []ncp.TargetGroup
	for _, s := range summary.TargetGroups {
		if cfg.TargetGroups.Match(s.TargetGroupName, s.TargetGroupNo) {
			tgs = append(tgs, s)
		}
	}
	summary.TargetGroups = tgs

	var dbs []ncp.CloudDBInstance
	for _, s := range summary.CloudDBs {
		if cfg.CloudDBs.Match(s.CloudDBServiceName, s.CloudDBInstanceNo) {
			dbs = append(dbs, s)
		}
	}
	summary.CloudDBs = dbs

	var pgs []ncp.CloudPostgresqlInstance
	for _, s := range summary.CloudPostgresqls {
		if cfg.CloudPostgresqls.Match(s.CloudPostgresqlServiceName, s.CloudPostgresqlInstanceNo) {
			pgs = append(pgs, s)
		}
	}
	summary.CloudPostgresqls = pgs

	var mgs []ncp.CloudMongoDbInstance
	for _, s := range summary.CloudMongoDBs {
		if cfg.CloudMongoDBs.Match(s.CloudMongoDbServiceName, s.CloudMongoDbInstanceNo) {
			mgs = append(mgs, s)
		}
	}
	summary.CloudMongoDBs = mgs

	var mdbs []ncp.CloudMariaDbInstance
	for _, s := range summary.CloudMariaDBs {
		if cfg.CloudMariaDBs.Match(s.CloudMariaDbServiceName, s.CloudMariaDbInstanceNo) {
			mdbs = append(mdbs, s)
		}
	}
	summary.CloudMariaDBs = mdbs

	var mysqls []ncp.CloudMysqlInstance
	for _, s := range summary.CloudMySQLs {
		if cfg.CloudMySQLs.Match(s.CloudMysqlServiceName, s.CloudMysqlInstanceNo) {
			mysqls = append(mysqls, s)
		}
	}
	summary.CloudMySQLs = mysqls

	var redises []ncp.CloudRedisInstance
	for _, s := range summary.CloudRedises {
		if cfg.CloudRedises.Match(s.CloudRedisServiceName, s.CloudRedisInstanceNo) {
			redises = append(redises, s)
		}
	}
	summary.CloudRedises = redises

	var vpcs []ncp.Vpc
	for _, s := range summary.Vpcs {
		if cfg.Vpcs.Match(s.VpcName, s.VpcNo) {
			vpcs = append(vpcs, s)
		}
	}
	summary.Vpcs = vpcs

	var subnets []ncp.Subnet
	for _, s := range summary.Subnets {
		if cfg.Subnets.Match(s.SubnetName, s.SubnetNo) {
			subnets = append(subnets, s)
		}
	}
	summary.Subnets = subnets

	var nats []ncp.NatGatewayInstance
	for _, s := range summary.NatGateways {
		if cfg.NatGateways.Match(s.NatGatewayName, s.NatGatewayInstanceNo) {
			nats = append(nats, s)
		}
	}
	summary.NatGateways = nats

	var peerings []ncp.VpcPeeringInstance
	for _, s := range summary.VpcPeerings {
		if cfg.VpcPeerings.Match(s.VpcPeeringName, s.VpcPeeringInstanceNo) {
			peerings = append(peerings, s)
		}
	}
	summary.VpcPeerings = peerings

	var nacls []ncp.NetworkAcl
	for _, s := range summary.NetworkAcls {
		if cfg.NetworkAcls.Match(s.NetworkAclName, s.NetworkAclNo) {
			nacls = append(nacls, s)
		}
	}
	summary.NetworkAcls = nacls

	var acgs []ncp.AccessControlGroup
	for _, s := range summary.AccessControlGroups {
		if cfg.AccessControlGroups.Match(s.AccessControlGroupName, s.AccessControlGroupNo) {
			acgs = append(acgs, s)
		}
	}
	summary.AccessControlGroups = acgs

	var asgs []ncp.AutoScalingGroup
	for _, s := range summary.AutoScalingGroups {
		if cfg.AutoScalingGroups.Match(s.AutoScalingGroupName, s.AutoScalingGroupNo) {
			asgs = append(asgs, s)
		}
	}
	summary.AutoScalingGroups = asgs

	var lcs []ncp.LaunchConfiguration
	for _, s := range summary.LaunchConfigurations {
		if cfg.LaunchConfigurations.Match(s.LaunchConfigurationName, s.LaunchConfigurationNo) {
			lcs = append(lcs, s)
		}
	}
	summary.LaunchConfigurations = lcs

	var clusters []ncp.NksCluster
	for _, s := range summary.NksClusters {
		if cfg.NksClusters.Match(s.Name, s.Uuid) {
			clusters = append(clusters, s)
		}
	}
	summary.NksClusters = clusters

	var scripts []ncp.InitScript
	for _, s := range summary.InitScripts {
		if cfg.InitScripts.Match(s.InitScriptName, s.InitScriptNo) {
			scripts = append(scripts, s)
		}
	}
	summary.InitScripts = scripts

	var keys []ncp.LoginKey
	for _, s := range summary.LoginKeys {
		if cfg.LoginKeys.Match(s.KeyName, s.KeyName) {
			keys = append(keys, s)
		}
	}
	summary.LoginKeys = keys

	var placementGroups []ncp.PlacementGroup
	for _, s := range summary.PlacementGroups {
		if cfg.PlacementGroups.Match(s.PlacementGroupName, s.PlacementGroupNo) {
			placementGroups = append(placementGroups, s)
		}
	}
	summary.PlacementGroups = placementGroups
}
