# NCP Nuke (ncp-nuke)

네이버 클라우드 플랫폼(NCP)의 루트 계정들을 대상으로 서브 계정(Sub Account)을 일괄 관리하거나 리소스를 파괴(Nuke)하기 위한 TUI 도구입니다.

엑셀 파일에 여러 루트 계정의 인증 키(Access Key, Secret Key)를 정의해두면, TUI에서 대상 계정을 선택하고 일괄 활성화/비활성화할 수 있습니다. 특히 교육용이나 다수의 계정을 한꺼번에 초기화/정리해야 할 때 유용합니다.

## 주요 기능

*   **계정 선택:** TUI에서 대상 루트 계정을 개별 선택하여 작업할 수 있습니다.
*   **일괄 활성화:** 선택한 계정의 서브 계정을 '활성(Active)' 상태로 변경하고, 비밀번호를 초기화합니다.
*   **일괄 비활성화:** 선택한 계정의 서브 계정을 '비활성(Inactive)' 상태로 변경하여 로그인을 차단합니다.
*   **리소스 정리 (Cleanup):** 비활성화 시 Cleanup 옵션을 활성화하면, **모든 리소스(서버, 스토리지, IP, DB, VPC 등)를 영구 삭제**합니다. (Nuke 기능)

## 설치 방법 (Installation)

### 1. Winget (Windows)

```bash
winget install Enbraining.NCPNuke
```

### 2. 바이너리 다운로드

[Releases 페이지](https://github.com/enbraining/ncp-subaccount-cli/releases)에서 운영체제에 맞는 파일을 다운로드하세요.

### 3. 소스 코드 빌드

Go 언어(1.22 이상)가 설치되어 있어야 합니다.

```bash
git clone https://github.com/enbraining/ncp-subaccount-cli.git
cd ncp-subaccount-cli

go mod tidy
go build -o ncp-nuke main.go
```

## 설정 파일 (Excel)

관리 대상 루트 계정들의 정보를 담은 엑셀 파일(`.xlsx`)이 필요합니다.

| 헤더 명 | 설명 | 필수 여부 |
| :--- | :--- | :--- |
| **AccountName** | 계정을 식별하기 위한 이름 | 선택 (없으면 자동 생성) |
| **AccessKey** | NCP API Access Key | **필수** |
| **SecretKey** | NCP API Secret Key | **필수** |
| **IAM Username** | 특정 서브 계정 ID | 선택 (지정 시 해당 유저만 제어) |
| **Password** | 설정할 비밀번호 | 선택 (활성화 시 사용) |

## 사용 방법 (Usage)

### 1. 엑셀 템플릿 생성

```bash
ncp-nuke template
```

현재 디렉토리에 `accounts_template.xlsx` 파일이 생성됩니다. 이 파일을 열어 계정 정보를 입력하세요.

### 2. TUI 실행

```bash
ncp-nuke -f ./accounts.xlsx
```

TUI가 실행되면 다음 흐름으로 진행됩니다:

1. **계정 선택** - Space로 대상 계정을 선택/해제하고 Enter로 다음 단계
2. **작업 선택** - 활성화 + 비밀번호 초기화 또는 비활성화 선택
3. **비밀번호 입력** - 활성화 시 엑셀에 비밀번호가 없는 계정이 있으면 공통 비밀번호 입력 (빈 값이면 자동 생성)
4. **확인** - 비활성화 시 Cleanup 옵션을 토글(c)할 수 있으며, y로 작업 시작
5. **실행** - 작업 진행 상황이 실시간으로 표시됩니다

**옵션:**

| 플래그 | 설명 |
| :--- | :--- |
| `-f, --file` | 루트 계정 목록 엑셀 파일 경로 (필수) |
| `-a, --account` | 특정 루트 계정만 대상 (AccountName 기준) |
| `--config` | 리소스 필터 설정 파일 경로 (JSON) |

## 주의사항

*   Cleanup 옵션은 매우 강력한 파괴적 동작을 수행하므로 실제 운영 중인 계정에 사용할 때 각별히 주의하세요. 안전을 위해 "REAL DELETE" 입력 확인이 필요합니다.
*   API 호출 횟수 제한이나 네트워크 오류 등으로 일부 작업이 실패할 수 있습니다. 실패 시 오류 메시지가 출력되므로 확인 후 재시도하십시오.

## License

MIT License
