package excel

import (
	"fmt"

	"ncp-nuke/pkg/ncp"

	"github.com/xuri/excelize/v2"
)

// AppendAccount appends a new account as a row to an existing accounts Excel
// file, writing each value into the column matching the header. AccessKey,
// SecretKey and IAM Username are required (mirroring ReadAccounts validation).
func AppendAccount(filePath string, acc ncp.RootAccount) error {
	if acc.AccessKey == "" || acc.SecretKey == "" {
		return fmt.Errorf("AccessKey와 SecretKey는 필수입니다")
	}
	if acc.IamUsername == "" {
		return fmt.Errorf("IAM Username은 필수입니다")
	}

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("opening excel file: %w", err)
	}
	defer f.Close()

	sheet := f.GetSheetName(0)
	if sheet == "" {
		return fmt.Errorf("no sheets found in excel file")
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		return fmt.Errorf("reading rows: %w", err)
	}
	if len(rows) == 0 {
		return fmt.Errorf("excel file has no header row")
	}

	colIdx := detectColumns(rows[0])
	if colIdx["accesskey"] == -1 || colIdx["secretkey"] == -1 {
		return fmt.Errorf("header must contain AccessKey and SecretKey columns")
	}

	newRow := len(rows) + 1 // 1-based; append after the last existing row

	set := func(field, value string) error {
		idx := colIdx[field]
		if idx < 0 || value == "" {
			return nil
		}
		cell, err := excelize.CoordinatesToCellName(idx+1, newRow)
		if err != nil {
			return err
		}
		return f.SetCellValue(sheet, cell, value)
	}

	for field, value := range map[string]string{
		"accountname": acc.AccountName,
		"accesskey":   acc.AccessKey,
		"secretkey":   acc.SecretKey,
		"iamusername": acc.IamUsername,
		"password":    acc.Password,
	} {
		if err := set(field, value); err != nil {
			return fmt.Errorf("writing %s: %w", field, err)
		}
	}

	if err := f.Save(); err != nil {
		return fmt.Errorf("saving excel file: %w", err)
	}
	return nil
}
