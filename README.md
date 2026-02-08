# NCP Nuke (ncp-nuke)

네이버 클라우드 플랫폼(NCP)의 루트 계정들을 대상으로 서브 계정(Sub Account)을 일괄 관리하거나 리소스를 파괴(Nuke)하기 위한 CLI 도구입니다.

엑셀 파일에 여러 루트 계정의 인증 키(Access Key, Secret Key)를 정의해두면, 각 계정의 서브 계정 목록을 조회하거나 일괄 활성화/비활성화할 수 있습니다. 특히 교육용이나 다수의 계정을 한꺼번에 초기화/정리해야 할 때 유용합니다.

## 주요 기능

*   **일괄 조회 (`list`):** 등록된 모든 루트 계정 하위의 서브 계정 정보를 테이블 형태로 출력합니다.
*   **일괄 활성화 (`activate`):** 모든 서브 계정을 '활성(Active)' 상태로 변경하고, 지정한 비밀번호로 초기화합니다.
*   **일괄 비활성화 (`deactivate`):** 모든 서브 계정을 '비활성(Inactive)' 상태로 변경하여 로그인을 차단합니다.
    *   `--cleanup` 옵션 사용 시, 비활성화 전에 **모든 리소스(서버, 스토리지, IP 등)를 영구 삭제**합니다. (Nuke 기능)

## 설치 방법 (Installation)

Go 언어(1.22 이상)가 설치되어 있어야 합니다.

```bash
# 레포지토리 클론
git clone https://github.com/your-repo/ncp-nuke.git
cd ncp-nuke

# 의존성 설치 및 빌드
go mod tidy
go build -o ncp-nuke main.go
```

## 설정 파일 (Excel)

관리 대상 루트 계정들의 정보를 담은 엑셀 파일(`.xlsx`)이 필요합니다.  
프로젝트 루트에 생성된 `accounts_template.xlsx` 파일을 참고하여 작성하세요.

첫 번째 행은 헤더여야 하며, 다음 컬럼들이 포함되어야 합니다 (순서 무관, 대소문자 구분 없음).

| 헤더 명 (예시) | 설명 | 필수 여부 |
| :--- | :--- | :--- |
| **AccountName** (또는 Name, 계정명) | 계정을 식별하기 위한 이름 | 선택 (없으면 자동 생성) |
| **AccessKey** | NCP API Access Key | **필수** |
| **SecretKey** | NCP API Secret Key | **필수** |
| **IAM User** (또는 ID, LoginId) | 특정 서브 계정 ID | 선택 (지정 시 해당 유저만 제어) |
| **Password** (또는 비밀번호) | 설정할 비밀번호 | 선택 (Activate 시 사용) |

**예시 (`accounts.xlsx`):**

| AccountName | AccessKey | SecretKey | IAM User | Password |
| :--- | :--- | :--- | :--- | :--- |
| Student-01 | ABCDEFGHIJKLMNOPQR | 1234567890abcdef | student-01 | P@ssword123! |
| Student-02 | STUVWXYZABCDEFGHIJ | abcdef1234567890 | student-02 | P@ssword123! |

## 사용 방법 (Usage)

모든 명령은 `-f` (또는 `--file`) 플래그로 엑셀 파일 경로를 지정해야 합니다.

### 1. 서브 계정 목록 조회

```bash
./ncp-nuke list -f ./accounts.xlsx
```

*   **옵션:**
    *   `-a, --account <이름>`: 특정 루트 계정(AccountName)만 필터링하여 조회합니다.

### 2. 서브 계정 활성화 및 비밀번호 초기화

서브 계정을 활성화하고 비밀번호를 강제로 재설정합니다.

```bash
# 비밀번호 대화형 입력 (엑셀에 없는 경우)
./ncp-nuke activate -f ./accounts.xlsx

# 전역 비밀번호 지정 실행 (엑셀에 없는 경우 적용)
./ncp-nuke activate -f ./accounts.xlsx -p "FallbackPassword123!"
```

### 3. 서브 계정 비활성화 (및 리소스 삭제)

서브 계정의 접근을 차단합니다.

```bash
# 단순 비활성화 (로그인 차단)
./ncp-nuke deactivate -f ./accounts.xlsx
```

**[주의] 리소스 전체 삭제 (Nuke)**  
`--cleanup` 플래그를 사용하면 해당 계정의 **서버, 블록 스토리지, 공인 IP, NAS, 로드밸런서**를 모두 삭제한 후 서브 계정을 비활성화합니다. 이 작업은 되돌릴 수 없습니다.

```bash
# 리소스 삭제 후 비활성화
./ncp-nuke deactivate -f ./accounts.xlsx --cleanup
```

## 주의사항

*   `--cleanup` 옵션은 매우 강력한 파괴적 동작을 수행하므로 실제 운영 중인 계정에 사용할 때 각별히 주의하세요.
*   API 호출 횟수 제한이나 네트워크 오류 등으로 일부 작업이 실패할 수 있습니다. 실패 시 오류 메시지가 출력되므로 확인 후 재시도하십시오.

## License

MIT License