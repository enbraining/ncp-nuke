package cmd

import (
	"fmt"
	"ncp-nuke/pkg/tui"

	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "TUI 모드로 실행",
	Long:  `텍스트 사용자 인터페이스(TUI)를 통해 대화형으로 계정을 관리합니다.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if filePath == "" {
			return fmt.Errorf("엑셀 파일 경로가 지정되지 않았습니다. -f 또는 --file 플래그를 사용하세요")
		}

		// TUI 실행
		return tui.Start(filePath, configPath)
	},
}

func init() {
	tuiCmd.Flags().StringVar(&configPath, "config", "", "리소스 필터 설정 파일 경로 (JSON)")
	rootCmd.AddCommand(tuiCmd)
}
