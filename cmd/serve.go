package cmd

import (
	"fmt"
	"net/http"

	"ncp-nuke/pkg/web"

	"github.com/spf13/cobra"
)

var servePort int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "웹 애플리케이션 실행 (브라우저 UI)",
	Long: `브라우저에서 사용할 수 있는 웹 UI를 로컬에서 실행합니다.
엑셀 파일의 루트 계정 목록을 불러와 계정 선택/활성화/비활성화/리소스 삭제/조회를
브라우저에서 수행할 수 있습니다.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if filePath == "" {
			return fmt.Errorf("엑셀 파일 경로가 지정되지 않았습니다. -f 또는 --file 플래그를 사용하세요")
		}
		srv, err := web.NewServer(filePath, configPath)
		if err != nil {
			return err
		}
		addr := fmt.Sprintf("127.0.0.1:%d", servePort)
		fmt.Printf("🧨 NCP Nuke 웹 콘솔 실행 중: http://%s\n", addr)
		fmt.Println("종료하려면 Ctrl+C 를 누르세요.")
		return http.ListenAndServe(addr, srv.Handler())
	},
}

func init() {
	serveCmd.Flags().StringVar(&configPath, "config", "", "리소스 필터 설정 파일 경로 (JSON)")
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8080, "웹 서버 포트")
	rootCmd.AddCommand(serveCmd)
}
