package excel

import (
	"fmt"
	"strings"

	"ncp-nuke/pkg/ncp"

	"github.com/xuri/excelize/v2"
)

// ReadAccounts reads root account information from an Excel file.
// Expected columns: AccountName, AccessKey, SecretKey (first row is header).
func ReadAccounts(filePath string) ([]ncp.RootAccount, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening excel file: %w", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	if sheetName == "" {
		return nil, fmt.Errorf("no sheets found in excel file")
	}

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("reading rows: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("excel file must have a header row and at least one data row")
	}

	// Find column indices from header row
	header := rows[0]
	colIdx := map[string]int{
		"accountname": -1,
		"accesskey":   -1,
		"secretkey":   -1,
		"iamusername": -1,
		"password":    -1,
	}

	for i, cell := range header {
		normalized := strings.ToLower(strings.TrimSpace(cell))
		// Support various header formats
		switch {
		case strings.Contains(normalized, "account") && strings.Contains(normalized, "name"):
			colIdx["accountname"] = i
		case normalized == "name" || normalized == "계정명" || normalized == "계정이름" || normalized == "계정 이름":
			if colIdx["accountname"] == -1 {
				colIdx["accountname"] = i
			}
		case strings.Contains(normalized, "access") && strings.Contains(normalized, "key"):
			colIdx["accesskey"] = i
		case strings.Contains(normalized, "secret") && strings.Contains(normalized, "key"):
			colIdx["secretkey"] = i
		case strings.Contains(normalized, "iam") || normalized == "id" || normalized == "아이디" || normalized == "loginid":
			colIdx["iamusername"] = i
		case normalized == "password" || normalized == "pw" || normalized == "비밀번호" || normalized == "비번":
			colIdx["password"] = i
		}
	}

	if colIdx["accesskey"] == -1 {
		return nil, fmt.Errorf("column 'AccessKey' not found in header row")
	}
	if colIdx["secretkey"] == -1 {
		return nil, fmt.Errorf("column 'SecretKey' not found in header row")
	}

	var accounts []ncp.RootAccount
	for i, row := range rows[1:] {
		lineNum := i + 2 // 1-indexed, skip header

		accessKey := getCell(row, colIdx["accesskey"])
		secretKey := getCell(row, colIdx["secretkey"])

		if accessKey == "" || secretKey == "" {
			fmt.Printf("[WARN] Row %d: AccessKey or SecretKey is empty, skipping\n", lineNum)
			continue
		}

		name := ""
		if colIdx["accountname"] != -1 {
			name = getCell(row, colIdx["accountname"])
		}
		if name == "" {
			name = fmt.Sprintf("Account-%d", lineNum-1)
		}

		iamUsername := ""
		if colIdx["iamusername"] != -1 {
			iamUsername = getCell(row, colIdx["iamusername"])
		}

		password := ""
		if colIdx["password"] != -1 {
			password = getCell(row, colIdx["password"])
		}

		accounts = append(accounts, ncp.RootAccount{
			AccountName: name,
			AccessKey:   accessKey,
			SecretKey:   secretKey,
			IamUsername: iamUsername,
			Password:    password,
		})
	}

	if len(accounts) == 0 {
		return nil, fmt.Errorf("no valid accounts found in excel file")
	}

	return accounts, nil
}

func getCell(row []string, idx int) string {
	if idx < 0 || idx >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[idx])
}
