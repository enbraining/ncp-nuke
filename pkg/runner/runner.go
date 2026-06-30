// Package runner contains the core account/resource processing logic shared by
// the TUI and the web application. It is UI-agnostic: progress is reported via a
// logFn callback.
package runner

import (
	"context"
	"fmt"
	"strings"

	"ncp-nuke/pkg/config"
	"ncp-nuke/pkg/ncp"
)

// formatResourceErr classifies a resource-listing error. A 403 / access-key
// rejection means the API key lacks permission for that product (or the product
// is unused / not enabled), which is expected for many accounts, so it is shown
// as a benign skip rather than a warning.
func formatResourceErr(e error) string {
	msg := e.Error()
	for _, sig := range []string{"HTTP 403", "StatusCode: 403", "AccessDenied", "InvalidAccessKeyId"} {
		if strings.Contains(msg, sig) {
			return fmt.Sprintf("    [건너뜀] 권한 없음/미사용: %v", e)
		}
	}
	return fmt.Sprintf("    [경고] 조회 오류: %v", e)
}

// Process runs the selected action against the selected accounts.
// action is one of: "activate", "deactivate", "nuke", "list".
// ctx cancellation stops launching further work (in-flight API calls finish).
func Process(ctx context.Context, accounts []ncp.RootAccount, selected map[int]bool, action, globalPassword string, cleanup bool, cfg *config.Config, logFn func(string)) {
	logFn("작업 시작...")

	totalSuccess, totalFail := 0, 0
	totalCleanupSuccess, totalCleanupFail := 0, 0

	for i, account := range accounts {
		if !selected[i] {
			continue
		}
		if ctx.Err() != nil {
			logFn("\n[취소됨] 작업이 취소되었습니다.")
			return
		}

		logFn(fmt.Sprintf("\n[루트 계정: %s]", account.AccountName))
		client := ncp.NewClient(account.AccessKey, account.SecretKey)

		// Read-only resource listing.
		if action == "list" {
			logFn("  리소스 조회 중...")
			summary, errs := client.ListAllResources()
			for _, e := range errs {
				logFn(formatResourceErr(e))
			}
			if cfg != nil {
				applyFilter(summary, cfg)
			}
			if summary.TotalCount() == 0 {
				logFn("  리소스 없음")
			} else {
				logFn(fmt.Sprintf("  총 %d개 리소스:", summary.TotalCount()))
				for _, bc := range summary.Breakdown() {
					logFn(fmt.Sprintf("    - %s: %d개", bc.Name, bc.Count))
				}
			}
			continue
		}

		// Cleanup phase (deactivate + cleanup, or standalone nuke)
		if (action == "deactivate" && cleanup) || action == "nuke" {
			logFn("  리소스 조회 중...")
			summary, errs := client.ListAllResources()
			for _, e := range errs {
				logFn(formatResourceErr(e))
			}

			if cfg != nil {
				applyFilter(summary, cfg)
			}

			if summary.TotalCount() > 0 {
				logFn(fmt.Sprintf("  총 %d개 서비스 해지 및 리소스 삭제 시작...", summary.TotalCount()))
				s, f := client.CleanupAllResources(ctx, summary, logFn)
				totalCleanupSuccess += s
				totalCleanupFail += f
				logFn(fmt.Sprintf("  서비스 해지 및 리소스 삭제 결과: 성공 %d, 실패 %d", s, f))
			} else {
				logFn("  삭제할 리소스 없음")
			}
		}

		// Nuke only targets resources; sub accounts are left untouched.
		if action == "nuke" {
			continue
		}

		// Sub account operation
		logFn("  서브 계정 조회 중...")
		subAccounts, err := client.ListSubAccounts()
		if err != nil {
			logFn(fmt.Sprintf("    [실패] 서브 계정 조회: %v", err))
			continue
		}

		var targets []ncp.SubAccount
		if account.IamUsername != "" {
			for _, sa := range subAccounts {
				if strings.EqualFold(sa.LoginId, account.IamUsername) {
					targets = append(targets, sa)
					break
				}
			}
		} else {
			targets = subAccounts
		}

		if len(targets) == 0 {
			if account.IamUsername != "" {
				var ids []string
				for _, sa := range subAccounts {
					ids = append(ids, sa.LoginId)
				}
				logFn(fmt.Sprintf("    [오류] 지정된 IAM 사용자(%s)를 찾을 수 없습니다. 존재하는 LoginId: %v", account.IamUsername, ids))
			} else {
				logFn("    대상 서브 계정 없음")
			}
			continue
		}

		switch action {
		case "activate":
			for _, sa := range targets {
				effectivePassword := account.Password
				if effectivePassword == "" {
					effectivePassword = globalPassword
				}
				generatedPw, err := client.ActivateSubAccount(sa, effectivePassword)
				if err != nil {
					logFn(fmt.Sprintf("    [실패] %s (%s): %v", sa.LoginId, sa.Name, err))
					totalFail++
				} else {
					if generatedPw != "" {
						logFn(fmt.Sprintf("    [성공] %s (%s): 활성화 + 비밀번호 초기화 완료 (생성된 비밀번호: %s)", sa.LoginId, sa.Name, generatedPw))
					} else {
						logFn(fmt.Sprintf("    [성공] %s (%s): 활성화 + 비밀번호 초기화 완료", sa.LoginId, sa.Name))
					}
					totalSuccess++
				}
			}

		case "deactivate":
			for _, sa := range targets {
				if !sa.Active {
					logFn(fmt.Sprintf("    [건너뜀] %s: 이미 비활성", sa.LoginId))
					totalSuccess++
					continue
				}
				if err := client.DeactivateSubAccount(sa); err != nil {
					logFn(fmt.Sprintf("    [실패] %s 비활성화: %v", sa.LoginId, err))
					totalFail++
				} else {
					logFn(fmt.Sprintf("    [성공] %s 비활성화 완료", sa.LoginId))
					totalSuccess++
				}
			}
		}
	}

	if action == "list" {
		logFn("\n리소스 조회 완료")
		return
	}

	if action == "nuke" {
		logFn(fmt.Sprintf("\n최종 결과: 리소스 삭제 성공 %d, 실패 %d", totalCleanupSuccess, totalCleanupFail))
		return
	}

	actionLabel := "활성화"
	if action == "deactivate" {
		actionLabel = "비활성화"
	}
	logFn(fmt.Sprintf("\n최종 결과: 서브계정 %s 성공 %d, 실패 %d", actionLabel, totalSuccess, totalFail))
	if cleanup {
		logFn(fmt.Sprintf("리소스 삭제 성공 %d, 실패 %d", totalCleanupSuccess, totalCleanupFail))
	}
}

func applyFilter(summary *ncp.ResourceSummary, cfg *config.Config) {
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

	var buckets []ncp.Bucket
	for _, s := range summary.Buckets {
		if cfg.Buckets.Match(s.Name, s.Name) {
			buckets = append(buckets, s)
		}
	}
	summary.Buckets = buckets

	var apigw []ncp.ApiGatewayProduct
	for _, s := range summary.ApiGatewayProducts {
		if cfg.ApiGatewayProducts.Match(s.ProductName, s.ProductId) {
			apigw = append(apigw, s)
		}
	}
	summary.ApiGatewayProducts = apigw
}
