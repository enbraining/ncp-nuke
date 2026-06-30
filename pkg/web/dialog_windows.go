//go:build windows

package web

import (
	"os/exec"
	"strings"
	"syscall"
)

// psDialog runs a PowerShell WinForms dialog with the console window hidden
// (CREATE_NO_WINDOW) so only the GUI dialog appears.
func psDialog(script string) (result string, cancelled bool, err error) {
	// Force UTF-8 stdout so Korean (non-ASCII) file paths are not mangled by the
	// console's default OEM code page (e.g. CP949).
	script = `[Console]::OutputEncoding = [System.Text.Encoding]::UTF8; ` + script
	cmd := exec.Command("powershell", "-NoProfile", "-STA", "-Command", script)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000} // CREATE_NO_WINDOW
	out, e := cmd.Output()
	if e != nil {
		return "", false, e
	}
	p := strings.TrimSpace(string(out))
	return p, p == "", nil
}

func chooseFileDialog() (path string, cancelled bool, err error) {
	return psDialog(`Add-Type -AssemblyName System.Windows.Forms; ` +
		`$f = New-Object System.Windows.Forms.OpenFileDialog; ` +
		`$f.Filter = 'Excel (*.xlsx)|*.xlsx|All files (*.*)|*.*'; ` +
		`if ($f.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) { [Console]::Out.Write($f.FileName) }`)
}

func chooseFolderDialog() (dir string, cancelled bool, err error) {
	return psDialog(`Add-Type -AssemblyName System.Windows.Forms; ` +
		`$f = New-Object System.Windows.Forms.FolderBrowserDialog; ` +
		`if ($f.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) { [Console]::Out.Write($f.SelectedPath) }`)
}

// openURL opens a URL in the default browser (no console window).
func openURL(url string) error {
	cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: 0x08000000}
	return cmd.Start()
}
