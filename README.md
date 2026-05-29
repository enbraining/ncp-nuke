# NCP Nuke (ncp-nuke)

네이버 클라우드 플랫폼(NCP)의 루트 계정들을 대상으로 서브 계정(Sub Account)을 일괄 관리하거나 리소스를 파괴(Nuke)하기 위한 TUI 도구입니다.

엑셀 파일에 여러 루트 계정의 인증 키(Access Key, Secret Key)를 정의해두면, TUI에서 대상 계정을 선택하고 일괄 활성화/비활성화할 수 있습니다. 특히 교육용이나 다수의 계정을 한꺼번에 초기화/정리해야 할 때 유용합니다.

## 주요 기능

*   **계정 선택:** TUI에서 대상 루트 계정을 개별 선택하여 작업할 수 있습니다.
*   **일괄 활성화:** 선택한 계정의 서브 계정을 '활성(Active)' 상태로 변경하고, 비밀번호를 초기화합니다.
*   **일괄 비활성화:** 선택한 계정의 서브 계정을 '비활성(Inactive)' 상태로 변경하여 로그인을 차단합니다.
*   **리소스 전체 삭제 (Nuke):** 선택한 계정의 **모든 리소스(서버, 스토리지, IP, DB, VPC 등)를 영구 삭제**합니다. (서브 계정은 유지)
*   **리소스 목록 조회:** 삭제 없이, 계정별 리소스를 카테고리별 개수로 조회합니다. (읽기 전용)
*   **리소스 정리 (Cleanup):** 비활성화 작업 시 Cleanup 옵션을 켜면 비활성화와 동시에 모든 리소스를 삭제합니다.

## 삭제 지원 리소스 (Supported Deletion)

Nuke / Cleanup 작업은 아래 리소스를 **의존성 순서에 맞춰** 삭제합니다. (VPC 환경 기준)

**Compute**
- [x] Server (서버)
- [x] Block Storage (블록 스토리지)
- [x] Block Storage Snapshot (블록 스토리지 스냅샷)
- [x] Init Script (init 스크립트)
- [x] Login Key (로그인 키)
- [ ] Placement Group (물리 배치 그룹)

**Auto Scaling**
- [ ] Auto Scaling Group
- [ ] Launch Configuration

**Storage**
- [ ] NAS Volume (NAS 볼륨)
- [ ] NAS Volume Snapshot (NAS 스냅샷)
- [x] Object Storage Bucket (버킷의 모든 객체/버전/멀티파트 업로드를 비운 뒤 버킷 삭제)

**Database (Cloud DB)**
- [ ] Cloud DB
- [ ] Cloud PostgreSQL
- [ ] Cloud MongoDB
- [ ] Cloud MariaDB
- [ ] Cloud MySQL
- [ ] Cloud Redis

**Network**
- [x] Public IP (공인 IP)
- [ ] Load Balancer (로드밸런서)
- [ ] Target Group (타깃 그룹)
- [x] VPC
- [x] Subnet (서브넷)
- [x] NAT Gateway
- [ ] VPC Peering
- [ ] Network ACL
- [ ] Route Table (라우트 테이블)
- [x] Access Control Group (ACG)

**Kubernetes**
- [ ] NKS Cluster (Ncloud Kubernetes Service 클러스터)

> 특정 리소스를 삭제 대상에서 제외하려면 `--config` 옵션(JSON 필터)을 사용하세요. (`config_example.json` 참고)
> 삭제 전 `리소스 목록 조회` 작업으로 실제 대상 개수를 미리 확인할 수 있습니다.
>
> Object Storage는 기본적으로 KR 리전(`https://kr.object.ncloudstorage.com`)을 사용합니다.
> 다른 리전은 환경변수 `NCP_OBJECT_STORAGE_ENDPOINT` / `NCP_OBJECT_STORAGE_REGION`으로 변경할 수 있습니다.

## 설치 방법 (Installation)

### 1. 바이너리 다운로드

[Releases 페이지](https://github.com/enbraining/ncp-nuke/releases)에서 운영체제에 맞는 파일을 다운로드하세요.

### 2. 소스 코드 빌드

Go 언어(1.22 이상)가 설치되어 있어야 합니다.

```bash
git clone https://github.com/enbraining/ncp-nuke.git
cd ncp-nuke

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
| **IAM Username** | 대상 서브 계정 ID | **필수** (해당 LoginId만 제어) |
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
2. **작업 선택** - 다음 4가지 중 선택
   - 활성화 + 비밀번호 초기화
   - 비활성화 (+ Cleanup 옵션)
   - 리소스 전체 삭제 (Nuke)
   - 리소스 목록 조회 (읽기 전용)
3. **비밀번호 입력** - 활성화 시 엑셀에 비밀번호가 없는 계정이 있으면 공통 비밀번호 입력 (빈 값이면 자동 생성)
4. **확인** - 비활성화 시 Cleanup 옵션을 토글(c)할 수 있으며, y로 작업 시작
5. **안전 확인** - 리소스 삭제(Nuke / Cleanup)는 `CONFIRM DELETE`를 정확히 입력해야 진행됩니다
6. **실행** - 작업 진행 상황이 실시간으로 표시됩니다

**옵션:**

| 플래그 | 설명 |
| :--- | :--- |
| `-f, --file` | 루트 계정 목록 엑셀 파일 경로 (필수) |
| `-a, --account` | 특정 루트 계정만 대상 (AccountName 기준) |
| `--config` | 리소스 필터 설정 파일 경로 (JSON) |

### 3. 웹 애플리케이션 실행

TUI 대신 브라우저 UI로도 사용할 수 있습니다.

```bash
ncp-nuke serve -f ./accounts.xlsx          # 기본 포트 8080
ncp-nuke serve -f ./accounts.xlsx -p 9000  # 포트 지정
ncp-nuke serve -f ./accounts.xlsx --config ./config.json
```

실행 후 브라우저에서 `http://127.0.0.1:8080` 에 접속합니다.

1. **계정 선택** - 체크박스로 대상 계정 선택
2. **작업 선택** - Sub Account 활성화 / 비활성화 / 리소스 전체 삭제 / 리소스 전체 조회
3. **안전 확인** - 파괴적 작업은 `CONFIRM DELETE` 입력 후 실행
4. **진행 로그** - 작업 진행 상황이 실시간(SSE)으로 표시됩니다

> 로컬 전용(127.0.0.1) 서버이며 인증 키가 그대로 사용되므로 신뢰된 환경에서만 실행하세요.

## 주의사항

*   Nuke / Cleanup은 매우 강력한 파괴적 동작을 수행하므로 실제 운영 중인 계정에 사용할 때 각별히 주의하세요. 안전을 위해 "CONFIRM DELETE" 입력 확인이 필요합니다.
*   API 호출 횟수 제한이나 네트워크 오류 등으로 일부 작업이 실패할 수 있습니다. 실패 시 오류 메시지가 출력되므로 확인 후 재시도하십시오.

## License

MIT License
