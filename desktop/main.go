// Command ncp-nuke-desktop runs the ncp-nuke web UI inside a native desktop
// window. It starts the same local HTTP server (so fetch/SSE/upload all work
// unchanged) on a random loopback port and points a WebView at it.
package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"ncp-nuke/pkg/web"

	webview "github.com/webview/webview_go"
)

func main() {
	// Optional preload file: ncp-nuke-desktop [accounts.xlsx]
	filePath := ""
	if len(os.Args) > 1 {
		filePath = os.Args[1]
	}

	srv, err := web.NewServer(filePath, "")
	if err != nil {
		fmt.Fprintln(os.Stderr, "초기화 실패:", err)
		os.Exit(1)
	}
	srv.Desktop = true // use native OS file dialogs for upload/download

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Fprintln(os.Stderr, "로컬 서버 시작 실패:", err)
		os.Exit(1)
	}
	go http.Serve(ln, srv.Handler())
	url := "http://" + ln.Addr().String()

	w := webview.New(false)
	defer w.Destroy()
	w.SetTitle("NCP Nuke")
	w.SetSize(1120, 920, webview.HintNone)
	w.Navigate(url)
	w.Run()
}
