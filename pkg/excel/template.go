package excel

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

// templateHeaders are the columns of the accounts template, in order.
var templateHeaders = []string{"AccountName", "AccessKey", "SecretKey", "IAM Username", "Password"}

var templateSamples = [][]string{
	{"Student-01", "YOUR_ACCESS_KEY_HERE_1", "YOUR_SECRET_KEY_HERE_1", "student-id-01", "InitialPassword123!"},
	{"Student-02", "YOUR_ACCESS_KEY_HERE_2", "YOUR_SECRET_KEY_HERE_2", "student-id-02", "InitialPassword123!"},
}

// buildTemplateFile returns a new accounts template workbook.
func buildTemplateFile() (*excelize.File, error) {
	f := excelize.NewFile()
	sheet := "Accounts"
	f.SetSheetName("Sheet1", sheet)

	for i, h := range templateHeaders {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}
	for r, row := range templateSamples {
		for c, v := range row {
			cell, _ := excelize.CoordinatesToCellName(c+1, r+2)
			f.SetCellValue(sheet, cell, v)
		}
	}
	f.SetColWidth(sheet, "A", "E", 30)
	return f, nil
}

// WriteTemplate creates the accounts template at the given path.
func WriteTemplate(path string) error {
	f, err := buildTemplateFile()
	if err != nil {
		return err
	}
	defer f.Close()
	if err := f.SaveAs(path); err != nil {
		return fmt.Errorf("템플릿 생성 실패: %w", err)
	}
	return nil
}

// TemplateBytes returns the accounts template as an .xlsx byte slice
// (for serving over HTTP).
func TemplateBytes() ([]byte, error) {
	f, err := buildTemplateFile()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("템플릿 생성 실패: %w", err)
	}
	return buf.Bytes(), nil
}
