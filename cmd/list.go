package cmd

import (
	"fmt"
	"os"
	"strconv"

	"ncp-nuke/pkg/excel"
	"ncp-nuke/pkg/ncp"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "모든 루트 계정의 서브 계정 목록 조회",
	Long:  "엑셀 파일에 등록된 루트 계정들의 하위 서브 계정들을 조회하여 테이블로 출력합니다.",
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	accounts, err := excel.ReadAccounts(filePath)
	if err != nil {
		return fmt.Errorf("엑셀 파일 읽기 실패: %w", err)
	}

	accounts = filterAccounts(accounts, accountFilter)
	if len(accounts) == 0 {
		return fmt.Errorf("대상 계정이 없습니다")
	}

	for _, account := range accounts {
		fmt.Printf("\n[루트 계정: %s]\n", account.AccountName)

		client := ncp.NewClient(account.AccessKey, account.SecretKey)
		subAccounts, err := client.ListSubAccounts()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  오류: %v\n", err)
			continue
		}

		if len(subAccounts) == 0 {
			fmt.Println("  서브 계정이 없습니다.")
			continue
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "로그인 ID", "이름", "이메일", "상태", "콘솔 접근", "API 접근"})
		table.SetBorder(true)
		table.SetRowLine(false)

		for _, sa := range subAccounts {
			status := "비활성"
			if sa.Active {
				status = "활성"
			}
			console := "X"
			if sa.CanConsoleAccess {
				console = "O"
			}
			api := "X"
			if sa.CanAPIGatewayAccess {
				api = "O"
			}
			table.Append([]string{
				strconv.FormatInt(sa.SubAccountId, 10),
				sa.LoginId,
				sa.Name,
				sa.Email,
				status,
				console,
				api,
			})
		}
		table.Render()
		fmt.Printf("  총 %d개 서브 계정\n", len(subAccounts))
	}

	return nil
}

func filterAccounts(accounts []ncp.RootAccount, filter string) []ncp.RootAccount {
	if filter == "" {
		return accounts
	}
	var filtered []ncp.RootAccount
	for _, a := range accounts {
		if a.AccountName == filter {
			filtered = append(filtered, a)
		}
	}
	return filtered
}
