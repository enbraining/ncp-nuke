package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"ncp-nuke/pkg/excel"
	"ncp-nuke/pkg/ncp"

	"github.com/spf13/cobra"
)

var cleanup bool

var deactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "서브 계정 일괄 비활성화",
	Long: `엑셀 파일에 등록된 루트 계정들의 하위 서브 계정들을 일괄로 비활성화(정지)합니다.

--cleanup 옵션을 사용하면 서브 계정 비활성화 전에
서버, 블록 스토리지, 공인 IP, NAS, 로드밸런서 등 모든 리소스를 삭제합니다.`,
	RunE: runDeactivate,
}

func init() {
	deactivateCmd.Flags().BoolVar(&cleanup, "cleanup", false, "리소스 전체 삭제 (서버, 스토리지, 공인IP, NAS, 로드밸런서)")
	rootCmd.AddCommand(deactivateCmd)
}

func runDeactivate(cmd *cobra.Command, args []string) error {
	accounts, err := excel.ReadAccounts(filePath)
	if err != nil {
		return fmt.Errorf("엑셀 파일 읽기 실패: %w", err)
	}

	accounts = filterAccounts(accounts, accountFilter)
	if len(accounts) == 0 {
		return fmt.Errorf("대상 계정이 없습니다")
	}

	if cleanup {
		return runDeactivateWithCleanup(accounts)
	}
	return runDeactivateOnly(accounts)
}

func runDeactivateOnly(accounts []ncp.RootAccount) error {
	fmt.Printf("\n%d개 루트 계정의 서브 계정을 비활성화합니다.\n", len(accounts))
	if !confirmPrompt() {
		return nil
	}

	totalSuccess, totalFail := 0, 0

	for _, account := range accounts {
		fmt.Printf("\n[루트 계정: %s]\n", account.AccountName)

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
					break
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
			if !sa.Active {
				fmt.Printf("  [건너뜀] %s (%s): 이미 비활성 상태\n", sa.LoginId, sa.Name)
				continue
			}

			err := client.DeactivateSubAccount(sa.SubAccountId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  [실패] %s (%s): %v\n", sa.LoginId, sa.Name, err)
				totalFail++
			} else {
				fmt.Printf("  [성공] %s (%s): 비활성화 완료\n", sa.LoginId, sa.Name)
				totalSuccess++
			}
		}
	}

	fmt.Printf("\n완료: 성공 %d, 실패 %d\n", totalSuccess, totalFail)
	return nil
}

func runDeactivateWithCleanup(accounts []ncp.RootAccount) error {
	fmt.Printf("\n%d개 루트 계정의 모든 리소스를 삭제하고 서브 계정을 비활성화합니다.\n", len(accounts))
	fmt.Println("삭제 대상: 서버, 블록 스토리지, 공인 IP, NAS 볼륨, 로드밸런서")
	fmt.Println()

	// Phase 1: 리소스 조회
	fmt.Println("=== 리소스 조회 중 ===")
	type accountResources struct {
		account ncp.RootAccount
		client  *ncp.Client
		summary *ncp.ResourceSummary
	}
	var targets []accountResources

	for _, account := range accounts {
		fmt.Printf("\n[루트 계정: %s]\n", account.AccountName)
		client := ncp.NewClient(account.AccessKey, account.SecretKey)
		summary, errs := client.ListAllResources()

		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "  경고: %v\n", e)
		}

		fmt.Printf("  서버: %d대, 블록스토리지: %d개, 공인IP: %d개, NAS: %d개, 로드밸런서: %d개\n",
			len(summary.Servers),
			len(summary.BlockStorages),
			len(summary.PublicIps),
			len(summary.NasVolumes),
			len(summary.LoadBalancers),
		)

		targets = append(targets, accountResources{
			account: account,
			client:  client,
			summary: summary,
		})
	}

	// Confirm with full resource count
	totalResources := 0
	for _, t := range targets {
		totalResources += t.summary.TotalCount()
	}

	if totalResources == 0 {
		fmt.Println("\n삭제할 리소스가 없습니다.")
	} else {
		fmt.Printf("\n총 %d개 리소스를 삭제합니다.\n", totalResources)
		fmt.Println("이 작업은 되돌릴 수 없습니다!")
		if !confirmPrompt() {
			return nil
		}

		// Phase 2: 리소스 삭제
		fmt.Println("\n=== 리소스 삭제 중 ===")
		for _, t := range targets {
			if t.summary.TotalCount() == 0 {
				continue
			}
			fmt.Printf("\n[루트 계정: %s]\n", t.account.AccountName)
			logFn := func(msg string) { fmt.Println(msg) }
			s, f := t.client.CleanupAllResources(t.summary, logFn)
			fmt.Printf("  리소스 삭제 결과: 성공 %d, 실패 %d\n", s, f)
		}
	}

	// Phase 3: 서브 계정 비활성화
	fmt.Println("\n=== 서브 계정 비활성화 ===")
	totalSuccess, totalFail := 0, 0

	for _, t := range targets {
		fmt.Printf("\n[루트 계정: %s]\n", t.account.AccountName)
		subAccounts, err := t.client.ListSubAccounts()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  서브 계정 조회 오류: %v\n", err)
			continue
		}

		if len(subAccounts) == 0 {
			fmt.Println("  서브 계정이 없습니다.")
			continue
		}

		var targetSubAccounts []ncp.SubAccount
		if t.account.IamUsername != "" {
			found := false
			for _, sa := range subAccounts {
				if sa.LoginId == t.account.IamUsername {
					targetSubAccounts = append(targetSubAccounts, sa)
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("  [경고] 지정된 IAM 사용자(%s)를 찾을 수 없습니다.\n", t.account.IamUsername)
				continue
			}
		} else {
			targetSubAccounts = subAccounts
		}

		for _, sa := range targetSubAccounts {
			if !sa.Active {
				fmt.Printf("  [건너뜀] %s (%s): 이미 비활성 상태\n", sa.LoginId, sa.Name)
				continue
			}
			err := t.client.DeactivateSubAccount(sa.SubAccountId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  [실패] %s (%s): %v\n", sa.LoginId, sa.Name, err)
				totalFail++
			} else {
				fmt.Printf("  [성공] %s (%s): 비활성화 완료\n", sa.LoginId, sa.Name)
				totalSuccess++
			}
		}
	}

	fmt.Printf("\n완료: 서브 계정 비활성화 성공 %d, 실패 %d\n", totalSuccess, totalFail)
	return nil
}

func confirmPrompt() bool {
	fmt.Print("계속하시겠습니까? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm != "y" && confirm != "yes" {
		fmt.Println("취소되었습니다.")
		return false
	}
	return true
}
