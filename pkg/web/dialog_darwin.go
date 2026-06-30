//go:build darwin

package web

import (
	"os"
	"os/exec"
	"strings"
)

// relaunchApp restarts the app after a self-update (re-opens the .app bundle).
func relaunchApp() {
	exe, err := os.Executable()
	if err != nil {
		return
	}
	if i := strings.Index(exe, ".app/"); i >= 0 {
		exec.Command("open", "-n", exe[:i+4]).Start()
		return
	}
	exec.Command(exe).Start()
}

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

// elevatedReplaceAndRelaunch copies the new binary (src) over the running
// executable with administrator privileges (osascript prompts for password),
// then relaunches the app. Used when the install location isn't user-writable.
func elevatedReplaceAndRelaunch(src string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	shell := `cp -f "` + src + `" "` + exe + `" && chmod +x "` + exe + `"`
	as := `do shell script "` + strings.ReplaceAll(shell, `"`, `\"`) + `" with administrator privileges`
	if err := exec.Command("osascript", "-e", as).Run(); err != nil {
		return err
	}
	relaunchApp()
	return nil
}
