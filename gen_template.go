package main

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

func main() {
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
	if err := f.SaveAs("accounts_template.xlsx"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("accounts_template.xlsx created successfully.")
}
