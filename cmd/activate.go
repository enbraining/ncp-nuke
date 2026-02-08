package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"ncp-nuke/pkg/excel"
	"ncp-nuke/pkg/ncp"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var password string

var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "서브 계정 일괄 활성화 + 비밀번호 초기화",
	Long: `엑셀 파일에 등록된 루트 계정들의 하위 서브 계정들을
일괄로 활성화하고 비밀번호를 초기화합니다.

--password 플래그를 지정하지 않으면 대화형으로 비밀번호를 입력받습니다.`,
	RunE: runActivate,
}

func init() {
	activateCmd.Flags().StringVarP(&password, "password", "p", "", "설정할 비밀번호 (미지정 시 대화형 입력)")
	rootCmd.AddCommand(activateCmd)
}

func runActivate(cmd *cobra.Command, args []string) error {
	accounts, err := excel.ReadAccounts(filePath)
	if err != nil {
		return fmt.Errorf("엑셀 파일 읽기 실패: %w", err)
	}

	accounts = filterAccounts(accounts, accountFilter)
	if len(accounts) == 0 {
		return fmt.Errorf("대상 계정이 없습니다")
	}

	// Check if we need a global fallback password
	needsGlobalPassword := false
	for _, acc := range accounts {
		if acc.Password == "" {
			needsGlobalPassword = true
			break
		}
	}

	if needsGlobalPassword && password == "" {
		password, err = promptPassword("설정할 비밀번호를 입력하세요 (엑셀에 비번이 없는 계정에 적용됨): ")
		if err != nil {
			return fmt.Errorf("비밀번호 입력 실패: %w", err)
		}
		if password == "" {
			return fmt.Errorf("비밀번호가 비어있습니다 (엑셀에 비밀번호가 없는 계정이 있어 필수입니다)")
		}
	}

	// Confirm
	fmt.Printf("\n%d개 루트 계정의 서브 계정을 활성화하고 비밀번호를 초기화합니다.\n", len(accounts))
	fmt.Print("계속하시겠습니까? (y/N): ")

	reader := bufio.NewReader(os.Stdin)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm != "y" && confirm != "yes" {
		fmt.Println("취소되었습니다.")
		return nil
	}

	totalSuccess, totalFail := 0, 0

	for _, account := range accounts {
		fmt.Printf("\n[루트 계정: %s]\n", account.AccountName)

		effectivePassword := account.Password
		if effectivePassword == "" {
			effectivePassword = password
		}

		client := ncp.NewClient(account.AccessKey, account.SecretKey)
		subAccounts, err := client.ListSubAccounts()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  서브 계정 조회 오류: %v\n", err)
			continue
		}

		if len(subAccounts) == 0 {
			fmt.Println("  서브 계정이 없습니다.")
			continue
		}

		var targets []ncp.SubAccount
		if account.IamUsername != "" {
			found := false
			for _, sa := range subAccounts {
				if sa.LoginId == account.IamUsername {
					targets = append(targets, sa)
					found = true
					break // found the specific user
				}
			}
			if !found {
				fmt.Printf("  [경고] 지정된 IAM 사용자(%s)를 찾을 수 없습니다.\n", account.IamUsername)
				continue
			}
		} else {
			targets = subAccounts
		}

		for _, sa := range targets {
			err := client.ActivateSubAccount(sa.SubAccountId, effectivePassword)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  [실패] %s (%s): %v\n", sa.LoginId, sa.Name, err)
				totalFail++
			} else {
				fmt.Printf("  [성공] %s (%s): 활성화 + 비밀번호 초기화 완료\n", sa.LoginId, sa.Name)
				totalSuccess++
			}
		}
	}

	fmt.Printf("\n완료: 성공 %d, 실패 %d\n", totalSuccess, totalFail)
	return nil
}

func promptPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		pw, err := term.ReadPassword(fd)
		fmt.Println()
		return string(pw), err
	}
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	return strings.TrimSpace(line), err
}
