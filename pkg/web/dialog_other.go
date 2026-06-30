//go:build !windows && !darwin

package web

import (
	"os"
	"os/exec"
	"strings"
)

// relaunchApp restarts the app after a self-update.
func relaunchApp() {
	if exe, err := os.Executable(); err == nil {
		exec.Command(exe).Start()
	}
}

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

// openURL opens a URL in the default browser.
func openURL(url string) error { return exec.Command("xdg-open", url).Start() }

// elevatedReplaceAndRelaunch tries pkexec to replace the running binary as root.
func elevatedReplaceAndRelaunch(src string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	if err := exec.Command("pkexec", "cp", "-f", src, exe).Run(); err != nil {
		return err
	}
	exec.Command(exe).Start()
	return nil
}
