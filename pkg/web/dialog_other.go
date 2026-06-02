//go:build !windows && !darwin

package web

import (
	"os/exec"
	"strings"
)

func chooseFileDialog() (path string, cancelled bool, err error) {
	out, e := exec.Command("zenity", "--file-selection", "--title=계정 엑셀 파일 선택", "--file-filter=*.xlsx").Output()
	if e != nil {
		return "", true, nil
	}
	return strings.TrimSpace(string(out)), false, nil
}

func chooseFolderDialog() (dir string, cancelled bool, err error) {
	out, e := exec.Command("zenity", "--file-selection", "--directory", "--title=템플릿 저장 폴더 선택").Output()
	if e != nil {
		return "", true, nil
	}
	return strings.TrimSpace(string(out)), false, nil
}
