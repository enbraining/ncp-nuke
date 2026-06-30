//go:build darwin

package web

import (
	"os/exec"
	"strings"
)

func chooseFileDialog() (path string, cancelled bool, err error) {
	out, e := exec.Command("osascript", "-e",
		`POSIX path of (choose file with prompt "계정 엑셀 파일 선택" of type {"xlsx","org.openxmlformats.spreadsheetml.sheet"})`).Output()
	if e != nil { // osascript exits non-zero on cancel
		return "", true, nil
	}
	return strings.TrimSpace(string(out)), false, nil
}

func chooseFolderDialog() (dir string, cancelled bool, err error) {
	out, e := exec.Command("osascript", "-e",
		`POSIX path of (choose folder with prompt "템플릿을 저장할 폴더 선택")`).Output()
	if e != nil {
		return "", true, nil
	}
	return strings.TrimSpace(string(out)), false, nil
}

// openURL opens a URL in the default browser.
func openURL(url string) error { return exec.Command("open", url).Start() }
