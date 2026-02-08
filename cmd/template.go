package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
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

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheetName := "Accounts"
	f.SetSheetName("Sheet1", sheetName)

	// Set Headers
	headers := []string{"AccountName", "AccessKey", "SecretKey", "IAM Username", "Password"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// Set Sample Data
	samples := [][]string{
		{"Student-01", "YOUR_ACCESS_KEY_HERE_1", "YOUR_SECRET_KEY_HERE_1", "student-id-01", "InitialPassword123!"},
		{"Student-02", "YOUR_ACCESS_KEY_HERE_2", "YOUR_SECRET_KEY_HERE_2", "student-id-02", "InitialPassword123!"},
	}

	for i, row := range samples {
		for j, value := range row {
			cell, _ := excelize.CoordinatesToCellName(j+1, i+2)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// Set column widths
	f.SetColWidth(sheetName, "A", "E", 30)

	// Save
	if err := f.SaveAs(filename); err != nil {
		return fmt.Errorf("템플릿 생성 실패: %w", err)
	}

	fmt.Printf("✅ '%s' 파일이 생성되었습니다.\n", filename)
	return nil
}
