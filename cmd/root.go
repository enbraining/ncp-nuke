package cmd

import (
	"fmt"
	"os"

	"ncp-nuke/pkg/tui"

	"github.com/spf13/cobra"
)

var filePath string
var accountFilter string
var configPath string

var rootCmd = &cobra.Command{
	Use:   "ncp-nuke",
	Short: "NCP Sub Account 일괄 관리 TUI 도구",
	Long: `네이버 클라우드 플랫폼(NCP) 루트 계정들의 Sub Account를
일괄로 조회, 활성화/비밀번호 초기화, 비활성화할 수 있는 TUI 도구입니다.

엑셀 파일에 루트 계정 정보(AccountName, AccessKey, SecretKey)를 입력하고,
각 루트 계정의 서브 계정들을 일괄 관리합니다.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if filePath == "" {
			return fmt.Errorf("엑셀 파일 경로가 지정되지 않았습니다. -f 또는 --file 플래그를 사용하세요")
		}
		return tui.Start(filePath, configPath, accountFilter)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "루트 계정 목록 엑셀 파일 경로 (필수)")
	rootCmd.PersistentFlags().StringVarP(&accountFilter, "account", "a", "", "특정 루트 계정만 대상 (AccountName 기준)")
	rootCmd.Flags().StringVar(&configPath, "config", "", "리소스 필터 설정 파일 경로 (JSON)")
}

