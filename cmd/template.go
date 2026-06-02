package cmd

import (
	"fmt"
	"os"

	"ncp-nuke/pkg/excel"

	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "엑셀 템플릿 파일 생성",
	Long:  `계정 정보를 입력할 수 있는 엑셀 템플릿 파일(accounts_template.xlsx)을 현재 디렉토리에 생성합니다.`,
	RunE:  runTemplate,
}

func init() {
	rootCmd.AddCommand(templateCmd)
}

func runTemplate(cmd *cobra.Command, args []string) error {
	filename := "accounts_template.xlsx"

	// 파일 존재 여부 확인
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("'%s' 파일이 이미 존재합니다. 덮어쓰지 않습니다", filename)
	}

	if err := excel.WriteTemplate(filename); err != nil {
		return err
	}

	fmt.Printf("✅ '%s' 파일이 생성되었습니다.\n", filename)
	return nil
}
