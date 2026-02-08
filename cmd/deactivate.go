package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"ncp-nuke/pkg/config"
	"ncp-nuke/pkg/excel"
	"ncp-nuke/pkg/ncp"

	"github.com/spf13/cobra"
)

var cleanup bool
var configPath string

var deactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "서브 계정 일괄 비활성화",
	Long: `엑셀 파일에 등록된 루트 계정들의 하위 서브 계정들을 일괄로 비활성화(정지)합니다.

--cleanup 옵션을 사용하면 서브 계정 비활성화 전에
서버, 블록 스토리지, 공인 IP, NAS, 로드밸런서 등 모든 리소스를 삭제합니다.

--config 옵션으로 JSON 설정 파일을 지정하면 삭제할 리소스를 필터링할 수 있습니다.
설정 파일 예시 (JSON):
{
  "servers": {
    "exclude": ["my-server-1", "S-12345"]
  }
}
`,
	RunE: runDeactivate,
}

func init() {
	deactivateCmd.Flags().BoolVar(&cleanup, "cleanup", false, "리소스 전체 삭제 (서버, 스토리지, 공인IP, NAS, 로드밸런서)")
	deactivateCmd.Flags().StringVar(&configPath, "config", "", "리소스 필터 설정 파일 경로 (JSON)")
	rootCmd.AddCommand(deactivateCmd)
}

func runDeactivate(cmd *cobra.Command, args []string) error {
	if filePath == "" {
		return fmt.Errorf("엑셀 파일 경로가 지정되지 않았습니다. -f 또는 --file 플래그를 사용하세요")
	}

	accounts, err := excel.ReadAccounts(filePath)
	if err != nil {
		return fmt.Errorf("엑셀 파일 읽기 실패: %w", err)
	}

	accounts = filterAccounts(accounts, accountFilter)
	if len(accounts) == 0 {
		return fmt.Errorf("대상 계정이 없습니다")
	}

	if cleanup {
		return runDeactivateWithCleanup(accounts)
	}
	return runDeactivateOnly(accounts)
}

func runDeactivateOnly(accounts []ncp.RootAccount) error {
	fmt.Printf("\n%d개 루트 계정의 서브 계정을 비활성화합니다.\n", len(accounts))
	if !confirmPrompt() {
		return nil
	}

	totalSuccess, totalFail := 0, 0

	for _, account := range accounts {
		fmt.Printf("\n[루트 계정: %s]\n", account.AccountName)

		client := ncp.NewClient(account.AccessKey, account.SecretKey)
		subAccounts, err := client.ListSubAccounts()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  서브 계정 조회 오류: %v\n", err)
			continue
		}

		if len(subAccounts) == 0 {
			fmt.Println("  서브 계정이 없습니다.")
			continue
		}

		var targets []ncp.SubAccount
		if account.IamUsername != "" {
			found := false
			for _, sa := range subAccounts {
				if sa.LoginId == account.IamUsername {
					targets = append(targets, sa)
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("  [경고] 지정된 IAM 사용자(%s)를 찾을 수 없습니다.\n", account.IamUsername)
				continue
			}
		} else {
			targets = subAccounts
		}

		for _, sa := range targets {
			if !sa.Active {
				fmt.Printf("  [건너뜀] %s (%s): 이미 비활성 상태\n", sa.LoginId, sa.Name)
				continue
			}

			err := client.DeactivateSubAccount(sa.SubAccountId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  [실패] %s (%s): %v\n", sa.LoginId, sa.Name, err)
				totalFail++
			} else {
				fmt.Printf("  [성공] %s (%s): 비활성화 완료\n", sa.LoginId, sa.Name)
				totalSuccess++
			}
		}
	}

	fmt.Printf("\n완료: 성공 %d, 실패 %d\n", totalSuccess, totalFail)
	return nil
}

func runDeactivateWithCleanup(accounts []ncp.RootAccount) error {
	// Load config if provided
	var cfg *config.Config
	if configPath != "" {
		var err error
		cfg, err = config.LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("설정 파일 로드 실패: %w", err)
		}
		fmt.Printf("설정 파일 로드됨: %s\n", configPath)
	}

	fmt.Printf("\n%d개 루트 계정의 모든 서비스를 해지하고 리소스를 삭제한 뒤 서브 계정을 비활성화합니다.\n", len(accounts))
	fmt.Println("해지/삭제 대상: 서버, 블록 스토리지, 공인 IP, NAS 볼륨, 로드밸런서, Cloud DB, VPC 등")
	fmt.Println()

	// Phase 1: 리소스 조회
	fmt.Println("=== 리소스 조회 중 ===")
	type accountResources struct {
		account ncp.RootAccount
		client  *ncp.Client
		summary *ncp.ResourceSummary
	}
	var targets []accountResources

	for _, account := range accounts {
		fmt.Printf("\n[루트 계정: %s]\n", account.AccountName)
		client := ncp.NewClient(account.AccessKey, account.SecretKey)
		summary, errs := client.ListAllResources()

		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "  경고: %v\n", e)
		}

		// Apply filter if config exists
		if cfg != nil {
			applyConfigFilter(summary, cfg)
		}

		dbCount := len(summary.CloudDBs) + len(summary.CloudPostgresqls) + len(summary.CloudMongoDBs) +
			len(summary.CloudMariaDBs) + len(summary.CloudMySQLs) + len(summary.CloudRedises)
		fmt.Printf("  서버: %d대, 블록스토리지: %d개(스냅샷 %d), 공인IP: %d개, NAS: %d개(스냅샷 %d), LB: %d개, TG: %d개\n"+
			"  CloudDB(All): %d개, VPC: %d개, Subnet: %d개, NAT: %d개, Peering: %d개, NACL: %d개\n"+
			"  RouteTable: %d개, ACG: %d개, ASG: %d개, LC: %d개, NKS: %d개\n"+
			"  InitScript: %d개, LoginKey: %d개, PlacementGroup: %d개\n",
			len(summary.Servers),
			len(summary.BlockStorages), len(summary.BlockStorageSnapshots),
			len(summary.PublicIps),
			len(summary.NasVolumes), len(summary.NasVolumeSnapshots),
			len(summary.LoadBalancers), len(summary.TargetGroups),
			dbCount,
			len(summary.Vpcs),
			len(summary.Subnets),
			len(summary.NatGateways),
			len(summary.VpcPeerings),
			len(summary.NetworkAcls),
			len(summary.RouteTables),
			len(summary.AccessControlGroups),
			len(summary.AutoScalingGroups),
			len(summary.LaunchConfigurations),
			len(summary.NksClusters),
			len(summary.InitScripts),
			len(summary.LoginKeys),
			len(summary.PlacementGroups),
		)

		targets = append(targets, accountResources{
			account: account,
			client:  client,
			summary: summary,
		})
	}

	// Confirm with full resource count
	totalResources := 0
	for _, t := range targets {
		totalResources += t.summary.TotalCount()
	}

	if totalResources == 0 {
		fmt.Println("\n삭제할 리소스가 없습니다.")
	} else {
		fmt.Printf("\n총 %d개 서비스 해지 및 리소스를 삭제합니다.\n", totalResources)
		fmt.Println("이 작업은 되돌릴 수 없습니다!")
		if !confirmPrompt() {
			return nil
		}

		// Phase 2: 리소스 삭제
		fmt.Println("\n=== 서비스 해지 및 리소스 삭제 중 ===")
		for _, t := range targets {
			if t.summary.TotalCount() == 0 {
				continue
			}
			fmt.Printf("\n[루트 계정: %s]\n", t.account.AccountName)
			logFn := func(msg string) { fmt.Println(msg) }
			s, f := t.client.CleanupAllResources(t.summary, logFn)
			fmt.Printf("  서비스 해지 및 리소스 삭제 결과: 성공 %d, 실패 %d\n", s, f)
		}
	}

	// Phase 3: 서브 계정 비활성화
	fmt.Println("\n=== 서브 계정 비활성화 ===")
	totalSuccess, totalFail := 0, 0

	for _, t := range targets {
		fmt.Printf("\n[루트 계정: %s]\n", t.account.AccountName)
		subAccounts, err := t.client.ListSubAccounts()
		if err != nil {
			fmt.Fprintf(os.Stderr, "  서브 계정 조회 오류: %v\n", err)
			continue
		}

		if len(subAccounts) == 0 {
			fmt.Println("  서브 계정이 없습니다.")
			continue
		}

		var targetSubAccounts []ncp.SubAccount
		if t.account.IamUsername != "" {
			found := false
			for _, sa := range subAccounts {
				if sa.LoginId == t.account.IamUsername {
					targetSubAccounts = append(targetSubAccounts, sa)
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("  [경고] 지정된 IAM 사용자(%s)를 찾을 수 없습니다.\n", t.account.IamUsername)
				continue
			}
		} else {
			targetSubAccounts = subAccounts
		}

		for _, sa := range targetSubAccounts {
			if !sa.Active {
				fmt.Printf("  [건너뜀] %s (%s): 이미 비활성 상태\n", sa.LoginId, sa.Name)
				continue
			}
			err := t.client.DeactivateSubAccount(sa.SubAccountId)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  [실패] %s (%s): %v\n", sa.LoginId, sa.Name, err)
				totalFail++
			} else {
				fmt.Printf("  [성공] %s (%s): 비활성화 완료\n", sa.LoginId, sa.Name)
				totalSuccess++
			}
		}
	}

	fmt.Printf("\n완료: 서브 계정 비활성화 성공 %d, 실패 %d\n", totalSuccess, totalFail)
	return nil
}

func confirmPrompt() bool {
	fmt.Print("계속하시겠습니까? (y/N): ")
	reader := bufio.NewReader(os.Stdin)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))
	if confirm != "y" && confirm != "yes" {
		fmt.Println("취소되었습니다.")
		return false
	}
	return true
}

func applyConfigFilter(summary *ncp.ResourceSummary, cfg *config.Config) {
	var servers []ncp.ServerInstance
	for _, s := range summary.Servers {
		if cfg.Servers.Match(s.ServerName, s.ServerInstanceNo) {
			servers = append(servers, s)
		}
	}
	summary.Servers = servers

	var storages []ncp.BlockStorageInstance
	for _, s := range summary.BlockStorages {
		if cfg.BlockStorages.Match(s.BlockStorageName, s.BlockStorageInstanceNo) {
			storages = append(storages, s)
		}
	}
	summary.BlockStorages = storages

	var bsSnaps []ncp.BlockStorageSnapshotInstance
	for _, s := range summary.BlockStorageSnapshots {
		if cfg.BlockStorageSnapshots.Match(s.BlockStorageSnapshotName, s.BlockStorageSnapshotInstanceNo) {
			bsSnaps = append(bsSnaps, s)
		}
	}
	summary.BlockStorageSnapshots = bsSnaps

	var ips []ncp.PublicIpInstance
	for _, s := range summary.PublicIps {
		if cfg.PublicIps.Match(s.PublicIp, s.PublicIpInstanceNo) {
			ips = append(ips, s)
		}
	}
	summary.PublicIps = ips

	var vols []ncp.NasVolumeInstance
	for _, s := range summary.NasVolumes {
		if cfg.NasVolumes.Match(s.VolumeName, s.NasVolumeInstanceNo) {
			vols = append(vols, s)
		}
	}
	summary.NasVolumes = vols

	var nasSnaps []ncp.NasVolumeSnapshot
	for _, s := range summary.NasVolumeSnapshots {
		if cfg.NasVolumeSnapshots.Match(s.NasVolumeSnapshotName, s.NasVolumeSnapshotInstanceNo) {
			nasSnaps = append(nasSnaps, s)
		}
	}
	summary.NasVolumeSnapshots = nasSnaps

	var lbs []ncp.LoadBalancerInstance
	for _, s := range summary.LoadBalancers {
		if cfg.LoadBalancers.Match(s.LoadBalancerName, s.LoadBalancerInstanceNo) {
			lbs = append(lbs, s)
		}
	}
	summary.LoadBalancers = lbs

	var tgs []ncp.TargetGroup
	for _, s := range summary.TargetGroups {
		if cfg.TargetGroups.Match(s.TargetGroupName, s.TargetGroupNo) {
			tgs = append(tgs, s)
		}
	}
	summary.TargetGroups = tgs

	var dbs []ncp.CloudDBInstance
	for _, s := range summary.CloudDBs {
		if cfg.CloudDBs.Match(s.CloudDBServiceName, s.CloudDBInstanceNo) {
			dbs = append(dbs, s)
		}
	}
	summary.CloudDBs = dbs

	var pgs []ncp.CloudPostgresqlInstance
	for _, s := range summary.CloudPostgresqls {
		if cfg.CloudPostgresqls.Match(s.CloudPostgresqlServiceName, s.CloudPostgresqlInstanceNo) {
			pgs = append(pgs, s)
		}
	}
	summary.CloudPostgresqls = pgs

	var mgs []ncp.CloudMongoDbInstance
	for _, s := range summary.CloudMongoDBs {
		if cfg.CloudMongoDBs.Match(s.CloudMongoDbServiceName, s.CloudMongoDbInstanceNo) {
			mgs = append(mgs, s)
		}
	}
	summary.CloudMongoDBs = mgs

	var mdbs []ncp.CloudMariaDbInstance
	for _, s := range summary.CloudMariaDBs {
		if cfg.CloudMariaDBs.Match(s.CloudMariaDbServiceName, s.CloudMariaDbInstanceNo) {
			mdbs = append(mdbs, s)
		}
	}
	summary.CloudMariaDBs = mdbs

	var mysqls []ncp.CloudMysqlInstance
	for _, s := range summary.CloudMySQLs {
		if cfg.CloudMySQLs.Match(s.CloudMysqlServiceName, s.CloudMysqlInstanceNo) {
			mysqls = append(mysqls, s)
		}
	}
	summary.CloudMySQLs = mysqls

	var redises []ncp.CloudRedisInstance
	for _, s := range summary.CloudRedises {
		if cfg.CloudRedises.Match(s.CloudRedisServiceName, s.CloudRedisInstanceNo) {
			redises = append(redises, s)
		}
	}
	summary.CloudRedises = redises

	var vpcs []ncp.Vpc
	for _, s := range summary.Vpcs {
		if cfg.Vpcs.Match(s.VpcName, s.VpcNo) {
			vpcs = append(vpcs, s)
		}
	}
	summary.Vpcs = vpcs

	var subnets []ncp.Subnet
	for _, s := range summary.Subnets {
		if cfg.Subnets.Match(s.SubnetName, s.SubnetNo) {
			subnets = append(subnets, s)
		}
	}
	summary.Subnets = subnets

	var nats []ncp.NatGatewayInstance
	for _, s := range summary.NatGateways {
		if cfg.NatGateways.Match(s.NatGatewayName, s.NatGatewayInstanceNo) {
			nats = append(nats, s)
		}
	}
	summary.NatGateways = nats

	var peerings []ncp.VpcPeeringInstance
	for _, s := range summary.VpcPeerings {
		if cfg.VpcPeerings.Match(s.VpcPeeringName, s.VpcPeeringInstanceNo) {
			peerings = append(peerings, s)
		}
	}
	summary.VpcPeerings = peerings

	var nacls []ncp.NetworkAcl
	for _, s := range summary.NetworkAcls {
		if cfg.NetworkAcls.Match(s.NetworkAclName, s.NetworkAclNo) {
			nacls = append(nacls, s)
		}
	}
	summary.NetworkAcls = nacls

	var acgs []ncp.AccessControlGroup
	for _, s := range summary.AccessControlGroups {
		if cfg.AccessControlGroups.Match(s.AccessControlGroupName, s.AccessControlGroupNo) {
			acgs = append(acgs, s)
		}
	}
	summary.AccessControlGroups = acgs

	var asgs []ncp.AutoScalingGroup
	for _, s := range summary.AutoScalingGroups {
		if cfg.AutoScalingGroups.Match(s.AutoScalingGroupName, s.AutoScalingGroupNo) {
			asgs = append(asgs, s)
		}
	}
	summary.AutoScalingGroups = asgs

	var lcs []ncp.LaunchConfiguration
	for _, s := range summary.LaunchConfigurations {
		if cfg.LaunchConfigurations.Match(s.LaunchConfigurationName, s.LaunchConfigurationNo) {
			lcs = append(lcs, s)
		}
	}
	summary.LaunchConfigurations = lcs

	var clusters []ncp.NksCluster
	for _, s := range summary.NksClusters {
		if cfg.NksClusters.Match(s.Name, s.Uuid) {
			clusters = append(clusters, s)
		}
	}
	summary.NksClusters = clusters

	var scripts []ncp.InitScript
	for _, s := range summary.InitScripts {
		if cfg.InitScripts.Match(s.InitScriptName, s.InitScriptNo) {
			scripts = append(scripts, s)
		}
	}
	summary.InitScripts = scripts

	var keys []ncp.LoginKey
	for _, s := range summary.LoginKeys {
		if cfg.LoginKeys.Match(s.KeyName, s.KeyName) {
			keys = append(keys, s)
		}
	}
	summary.LoginKeys = keys

	var placementGroups []ncp.PlacementGroup
	for _, s := range summary.PlacementGroups {
		if cfg.PlacementGroups.Match(s.PlacementGroupName, s.PlacementGroupNo) {
			placementGroups = append(placementGroups, s)
		}
	}
	summary.PlacementGroups = placementGroups
}