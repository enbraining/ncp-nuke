package ncp

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// ResourceSummary holds a summary of resources found for a root account.
type ResourceSummary struct {
	Servers                 []ServerInstance
	BlockStorages           []BlockStorageInstance
	BlockStorageSnapshots   []BlockStorageSnapshotInstance
	PublicIps               []PublicIpInstance
	NasVolumes              []NasVolumeInstance
	NasVolumeSnapshots      []NasVolumeSnapshot
	LoadBalancers           []LoadBalancerInstance
	TargetGroups            []TargetGroup
	CloudDBs                []CloudDBInstance
	CloudPostgresqls        []CloudPostgresqlInstance
	CloudMongoDBs           []CloudMongoDbInstance
	CloudMariaDBs           []CloudMariaDbInstance
	CloudMySQLs             []CloudMysqlInstance
	CloudRedises            []CloudRedisInstance
	Vpcs                    []Vpc
	Subnets                 []Subnet
	NatGateways             []NatGatewayInstance
	VpcPeerings             []VpcPeeringInstance
	NetworkAcls             []NetworkAcl
	RouteTables             []RouteTable
	AccessControlGroups     []AccessControlGroup
	AutoScalingGroups       []AutoScalingGroup
	LaunchConfigurations    []LaunchConfiguration
	NksClusters             []NksCluster
	InitScripts             []InitScript
	LoginKeys               []LoginKey
	PlacementGroups         []PlacementGroup
}

// TotalCount returns total number of resources.
func (r *ResourceSummary) TotalCount() int {
	return len(r.Servers) + len(r.BlockStorages) + len(r.BlockStorageSnapshots) +
		len(r.PublicIps) + len(r.NasVolumes) + len(r.NasVolumeSnapshots) +
		len(r.LoadBalancers) + len(r.TargetGroups) +
		len(r.CloudDBs) + len(r.CloudPostgresqls) + len(r.CloudMongoDBs) +
		len(r.CloudMariaDBs) + len(r.CloudMySQLs) + len(r.CloudRedises) +
		len(r.Vpcs) + len(r.Subnets) + len(r.NatGateways) +
		len(r.VpcPeerings) + len(r.NetworkAcls) + len(r.RouteTables) +
		len(r.AccessControlGroups) + len(r.AutoScalingGroups) +
		len(r.LaunchConfigurations) + len(r.NksClusters) +
		len(r.InitScripts) + len(r.LoginKeys) + len(r.PlacementGroups)
}

// --- List APIs ---

func (c *Client) ListServers() ([]ServerInstance, error) {
	path := "/getServerInstanceList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VServerBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}

	var resp struct {
		GetServerInstanceListResponse struct {
			TotalRows          int              `json:"totalRows"`
			ServerInstanceList []ServerInstance  `json:"serverInstanceList"`
		} `json:"getServerInstanceListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.GetServerInstanceListResponse.ServerInstanceList, nil
}

func (c *Client) ListBlockStorages() ([]BlockStorageInstance, error) {
	path := "/getBlockStorageInstanceList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VServerBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}

	var resp struct {
		GetBlockStorageInstanceListResponse struct {
			TotalRows                int                    `json:"totalRows"`
			BlockStorageInstanceList []BlockStorageInstance  `json:"blockStorageInstanceList"`
		} `json:"getBlockStorageInstanceListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.GetBlockStorageInstanceListResponse.BlockStorageInstanceList, nil
}

func (c *Client) ListPublicIps() ([]PublicIpInstance, error) {
	path := "/getPublicIpInstanceList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VServerBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}

	var resp struct {
		GetPublicIpInstanceListResponse struct {
			TotalRows            int                `json:"totalRows"`
			PublicIpInstanceList []PublicIpInstance   `json:"publicIpInstanceList"`
		} `json:"getPublicIpInstanceListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.GetPublicIpInstanceListResponse.PublicIpInstanceList, nil
}

func (c *Client) ListNasVolumes() ([]NasVolumeInstance, error) {
	path := "/getNasVolumeInstanceList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VNASBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}

	var resp struct {
		GetNasVolumeInstanceListResponse struct {
			TotalRows             int                 `json:"totalRows"`
			NasVolumeInstanceList []NasVolumeInstance  `json:"nasVolumeInstanceList"`
		} `json:"getNasVolumeInstanceListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.GetNasVolumeInstanceListResponse.NasVolumeInstanceList, nil
}

func (c *Client) ListLoadBalancers() ([]LoadBalancerInstance, error) {
	path := "/getLoadBalancerInstanceList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VLBBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}

	var resp struct {
		GetLoadBalancerInstanceListResponse struct {
			TotalRows                  int                      `json:"totalRows"`
			LoadBalancerInstanceList   []LoadBalancerInstance    `json:"loadBalancerInstanceList"`
		} `json:"getLoadBalancerInstanceListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.GetLoadBalancerInstanceListResponse.LoadBalancerInstanceList, nil
}

// ListAllResources collects all resources for a root account.
func (c *Client) ListAllResources() (*ResourceSummary, []error) {
	summary := &ResourceSummary{}
	var errs []error

	// 1. Existing Resources
	if servers, err := c.ListServers(); err != nil {
		errs = append(errs, fmt.Errorf("서버 조회: %w", err))
	} else {
		summary.Servers = servers
	}
	if storages, err := c.ListBlockStorages(); err != nil {
		errs = append(errs, fmt.Errorf("블록 스토리지 조회: %w", err))
	} else {
		summary.BlockStorages = storages
	}
	if ips, err := c.ListPublicIps(); err != nil {
		errs = append(errs, fmt.Errorf("공인 IP 조회: %w", err))
	} else {
		summary.PublicIps = ips
	}
	if vols, err := c.ListNasVolumes(); err != nil {
		errs = append(errs, fmt.Errorf("NAS 볼륨 조회: %w", err))
	} else {
		summary.NasVolumes = vols
	}
	if lbs, err := c.ListLoadBalancers(); err != nil {
		errs = append(errs, fmt.Errorf("로드밸런서 조회: %w", err))
	} else {
		summary.LoadBalancers = lbs
	}

	// 2. Snapshots
	if snaps, err := c.ListBlockStorageSnapshotInstances(); err != nil {
		errs = append(errs, fmt.Errorf("블록 스토리지 스냅샷 조회: %w", err))
	} else {
		summary.BlockStorageSnapshots = snaps
	}
	if snaps, err := c.ListNasVolumeSnapshots(); err != nil {
		errs = append(errs, fmt.Errorf("NAS 스냅샷 조회: %w", err))
	} else {
		summary.NasVolumeSnapshots = snaps
	}

	// 3. Cloud DB (all types)
	if dbs, err := c.ListCloudDBInstances(); err != nil {
		errs = append(errs, fmt.Errorf("Cloud DB 조회: %w", err))
	} else {
		summary.CloudDBs = dbs
	}
	if pgs, err := c.ListCloudPostgresqlInstances(); err != nil {
		errs = append(errs, fmt.Errorf("Cloud DB(Pg) 조회: %w", err))
	} else {
		summary.CloudPostgresqls = pgs
	}
	if mgs, err := c.ListCloudMongoDBInstances(); err != nil {
		errs = append(errs, fmt.Errorf("Cloud DB(Mongo) 조회: %w", err))
	} else {
		summary.CloudMongoDBs = mgs
	}
	if mdbs, err := c.ListCloudMariaDbInstances(); err != nil {
		errs = append(errs, fmt.Errorf("Cloud DB(MariaDB) 조회: %w", err))
	} else {
		summary.CloudMariaDBs = mdbs
	}
	if mysqls, err := c.ListCloudMysqlInstances(); err != nil {
		errs = append(errs, fmt.Errorf("Cloud DB(MySQL) 조회: %w", err))
	} else {
		summary.CloudMySQLs = mysqls
	}
	if redises, err := c.ListCloudRedisInstances(); err != nil {
		errs = append(errs, fmt.Errorf("Cloud DB(Redis) 조회: %w", err))
	} else {
		summary.CloudRedises = redises
	}

	// 4. Load Balancer / Target Group
	if tgs, err := c.ListTargetGroups(); err != nil {
		errs = append(errs, fmt.Errorf("Target Group 조회: %w", err))
	} else {
		summary.TargetGroups = tgs
	}

	// 5. VPC 관련
	if vpcs, err := c.ListVpcs(); err != nil {
		errs = append(errs, fmt.Errorf("VPC 조회: %w", err))
	} else {
		summary.Vpcs = vpcs
	}
	if subnets, err := c.ListSubnets(); err != nil {
		errs = append(errs, fmt.Errorf("Subnet 조회: %w", err))
	} else {
		summary.Subnets = subnets
	}
	if nats, err := c.ListNatGateways(); err != nil {
		errs = append(errs, fmt.Errorf("NAT Gateway 조회: %w", err))
	} else {
		summary.NatGateways = nats
	}
	if peerings, err := c.ListVpcPeeringInstances(); err != nil {
		errs = append(errs, fmt.Errorf("VPC Peering 조회: %w", err))
	} else {
		summary.VpcPeerings = peerings
	}
	if nacls, err := c.ListNetworkAcls(); err != nil {
		errs = append(errs, fmt.Errorf("Network ACL 조회: %w", err))
	} else {
		summary.NetworkAcls = nacls
	}
	if rts, err := c.ListRouteTables(); err != nil {
		errs = append(errs, fmt.Errorf("Route Table 조회: %w", err))
	} else {
		summary.RouteTables = rts
	}
	if acgs, err := c.ListAccessControlGroups(); err != nil {
		errs = append(errs, fmt.Errorf("ACG 조회: %w", err))
	} else {
		summary.AccessControlGroups = acgs
	}

	// 6. Auto Scaling
	if asgs, err := c.ListAutoScalingGroups(); err != nil {
		errs = append(errs, fmt.Errorf("Auto Scaling 조회: %w", err))
	} else {
		summary.AutoScalingGroups = asgs
	}
	if lcs, err := c.ListLaunchConfigurations(); err != nil {
		errs = append(errs, fmt.Errorf("Launch Configuration 조회: %w", err))
	} else {
		summary.LaunchConfigurations = lcs
	}

	// 7. NKS
	if clusters, err := c.ListNksClusters(); err != nil {
		errs = append(errs, fmt.Errorf("NKS Cluster 조회: %w", err))
	} else {
		summary.NksClusters = clusters
	}

	// 8. Server 부가 리소스
	if scripts, err := c.ListInitScripts(); err != nil {
		errs = append(errs, fmt.Errorf("Init Script 조회: %w", err))
	} else {
		summary.InitScripts = scripts
	}
	if keys, err := c.ListLoginKeys(); err != nil {
		errs = append(errs, fmt.Errorf("Login Key 조회: %w", err))
	} else {
		summary.LoginKeys = keys
	}
	if pgs, err := c.ListPlacementGroups(); err != nil {
		errs = append(errs, fmt.Errorf("Placement Group 조회: %w", err))
	} else {
		summary.PlacementGroups = pgs
	}

	return summary, errs
}

// CleanupAllResources deletes all resources for a root account.
// Order: NKS -> ASG -> LaunchConfig -> All DBs -> LB -> TargetGroup -> Server -> BS Snapshot -> BS -> NAS Snapshot -> NAS
//        -> NAT -> VPC Peering -> Routes -> IP -> ACG -> Network ACL -> Subnet -> VPC -> InitScript -> LoginKey -> PlacementGroup
func (c *Client) CleanupAllResources(summary *ResourceSummary, logFn func(string)) (int, int) {
	success, fail := 0, 0

	// 1. NKS Clusters
	for _, k := range summary.NksClusters {
		logFn(fmt.Sprintf("  NKS 클러스터 서비스 해지: %s (%s)", k.Name, k.Uuid))
		if err := c.DeleteNksCluster(k.Uuid); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공] 삭제 요청 완료")
			success++
		}
	}
	if len(summary.NksClusters) > 0 {
		logFn("  NKS 클러스터 삭제 대기 (60초)...")
		time.Sleep(60 * time.Second)
	}

	// 2. Auto Scaling Groups
	for _, asg := range summary.AutoScalingGroups {
		logFn(fmt.Sprintf("  Auto Scaling Group 삭제: %s (%s)", asg.AutoScalingGroupName, asg.AutoScalingGroupNo))
		if err := c.DeleteAutoScalingGroup(asg.AutoScalingGroupNo); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공]")
			success++
		}
	}
	if len(summary.AutoScalingGroups) > 0 {
		logFn("  ASG 삭제 및 서버 종료 대기 (60초)...")
		time.Sleep(60 * time.Second)
	}

	// 3. Launch Configurations (must delete after ASG)
	for _, lc := range summary.LaunchConfigurations {
		logFn(fmt.Sprintf("  Launch Configuration 삭제: %s (%s)", lc.LaunchConfigurationName, lc.LaunchConfigurationNo))
		if err := c.DeleteLaunchConfiguration(lc.LaunchConfigurationNo); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공]")
			success++
		}
	}

	// 4. Cloud DBs (All types)
	for _, db := range summary.CloudDBs {
		logFn(fmt.Sprintf("  Cloud DB 서비스 해지: %s (%s)", db.CloudDBServiceName, db.CloudDBInstanceNo))
		if err := c.DeleteCloudDBInstance(db.CloudDBInstanceNo); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공] 삭제 요청 완료")
			success++
		}
	}
	for _, pg := range summary.CloudPostgresqls {
		logFn(fmt.Sprintf("  Cloud DB(Pg) 서비스 해지: %s (%s)", pg.CloudPostgresqlServiceName, pg.CloudPostgresqlInstanceNo))
		if err := c.DeleteCloudPostgresqlInstance(pg.CloudPostgresqlInstanceNo); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공] 삭제 요청 완료")
			success++
		}
	}
	for _, mg := range summary.CloudMongoDBs {
		logFn(fmt.Sprintf("  Cloud DB(Mongo) 서비스 해지: %s (%s)", mg.CloudMongoDbServiceName, mg.CloudMongoDbInstanceNo))
		if err := c.DeleteCloudMongoDBInstance(mg.CloudMongoDbInstanceNo); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공] 삭제 요청 완료")
			success++
		}
	}
	for _, mdb := range summary.CloudMariaDBs {
		logFn(fmt.Sprintf("  Cloud DB(MariaDB) 서비스 해지: %s (%s)", mdb.CloudMariaDbServiceName, mdb.CloudMariaDbInstanceNo))
		if err := c.DeleteCloudMariaDbInstance(mdb.CloudMariaDbInstanceNo); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공] 삭제 요청 완료")
			success++
		}
	}
	for _, mysql := range summary.CloudMySQLs {
		logFn(fmt.Sprintf("  Cloud DB(MySQL) 서비스 해지: %s (%s)", mysql.CloudMysqlServiceName, mysql.CloudMysqlInstanceNo))
		if err := c.DeleteCloudMysqlInstance(mysql.CloudMysqlInstanceNo); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공] 삭제 요청 완료")
			success++
		}
	}
	for _, redis := range summary.CloudRedises {
		logFn(fmt.Sprintf("  Cloud DB(Redis) 서비스 해지: %s (%s)", redis.CloudRedisServiceName, redis.CloudRedisInstanceNo))
		if err := c.DeleteCloudRedisInstance(redis.CloudRedisInstanceNo); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공] 삭제 요청 완료")
			success++
		}
	}

	hasDBs := len(summary.CloudDBs) > 0 || len(summary.CloudPostgresqls) > 0 || len(summary.CloudMongoDBs) > 0 ||
		len(summary.CloudMariaDBs) > 0 || len(summary.CloudMySQLs) > 0 || len(summary.CloudRedises) > 0
	if hasDBs {
		logFn("  Cloud DB 삭제 대기 (30초)...")
		time.Sleep(30 * time.Second)
	}

	// 5. Load Balancers
	for _, lb := range summary.LoadBalancers {
		logFn(fmt.Sprintf("  로드밸런서 삭제: %s (%s)", lb.LoadBalancerName, lb.LoadBalancerInstanceNo))
		if err := c.DeleteLoadBalancer(lb.LoadBalancerInstanceNo); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공]")
			success++
		}
	}

	// 6. Target Groups (must delete after LB)
	for _, tg := range summary.TargetGroups {
		logFn(fmt.Sprintf("  Target Group 삭제: %s (%s)", tg.TargetGroupName, tg.TargetGroupNo))
		if err := c.DeleteTargetGroup(tg.TargetGroupNo); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공]")
			success++
		}
	}

	// 7. Stop & Terminate Servers
	if len(summary.Servers) > 0 {
		var runningNos []string
		for _, s := range summary.Servers {
			if s.ServerInstanceStatus.Code == "RUN" {
				runningNos = append(runningNos, s.ServerInstanceNo)
			}
		}
		if len(runningNos) > 0 {
			logFn(fmt.Sprintf("  서버 %d대 정지 중...", len(runningNos)))
			if err := c.StopServers(runningNos); err != nil {
				logFn(fmt.Sprintf("    [실패] 서버 정지: %v", err))
			} else {
				logFn("    서버 정지 요청 완료, 30초 대기...")
				time.Sleep(30 * time.Second)
			}
		}

		var allNos []string
		for _, s := range summary.Servers {
			allNos = append(allNos, s.ServerInstanceNo)
		}
		logFn(fmt.Sprintf("  서버 %d대 반납(삭제) 중...", len(allNos)))
		if err := c.TerminateServers(allNos); err != nil {
			logFn(fmt.Sprintf("    [실패] 서버 반납: %v", err))
			fail += len(allNos)
		} else {
			logFn("    [성공] 서버 반납 요청 완료")
			success += len(allNos)
			logFn("    서버 반납 대기 중 (30초)...")
			time.Sleep(30 * time.Second)
		}
	}

	// 8. Block Storage Snapshots (must delete before block storages)
	if len(summary.BlockStorageSnapshots) > 0 {
		var snapNos []string
		for _, snap := range summary.BlockStorageSnapshots {
			snapNos = append(snapNos, snap.BlockStorageSnapshotInstanceNo)
		}
		logFn(fmt.Sprintf("  블록 스토리지 스냅샷 %d개 삭제 중...", len(snapNos)))
		if err := c.DeleteBlockStorageSnapshotInstances(snapNos); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail += len(snapNos)
		} else {
			logFn("    [성공]")
			success += len(snapNos)
		}
	}

	// 9. Delete Block Storages (non-basic)
	var storagesToDelete []string
	for _, bs := range summary.BlockStorages {
		if bs.BlockStorageDiskDetailType.Code == "BASIC" {
			continue
		}
		storagesToDelete = append(storagesToDelete, bs.BlockStorageInstanceNo)
	}
	if len(storagesToDelete) > 0 {
		logFn(fmt.Sprintf("  블록 스토리지 %d개 삭제 중...", len(storagesToDelete)))
		if err := c.DeleteBlockStorages(storagesToDelete); err != nil {
			logFn(fmt.Sprintf("    [실패] 블록 스토리지 삭제: %v", err))
			fail += len(storagesToDelete)
		} else {
			logFn("    [성공]")
			success += len(storagesToDelete)
		}
	}

	// 10. NAS Volume Snapshots (must delete before NAS volumes)
	for _, snap := range summary.NasVolumeSnapshots {
		logFn(fmt.Sprintf("  NAS 스냅샷 삭제: %s (%s)", snap.NasVolumeSnapshotName, snap.NasVolumeSnapshotInstanceNo))
		if err := c.DeleteNasVolumeSnapshot(snap.NasVolumeSnapshotInstanceNo); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공]")
			success++
		}
	}

	// 11. NAS Volumes
	for _, vol := range summary.NasVolumes {
		logFn(fmt.Sprintf("  NAS 볼륨 삭제: %s (%s)", vol.VolumeName, vol.NasVolumeInstanceNo))
		if err := c.DeleteNasVolume(vol.NasVolumeInstanceNo); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공]")
			success++
		}
	}

	// 12. NAT Gateways
	for _, nat := range summary.NatGateways {
		logFn(fmt.Sprintf("  NAT Gateway 삭제: %s (%s)", nat.NatGatewayName, nat.NatGatewayInstanceNo))
		if err := c.DeleteNatGateway(nat.NatGatewayInstanceNo); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공]")
			success++
		}
	}
	if len(summary.NatGateways) > 0 {
		logFn("  NAT Gateway 삭제 대기 (20초)...")
		time.Sleep(20 * time.Second)
	}

	// 13. VPC Peering
	for _, p := range summary.VpcPeerings {
		logFn(fmt.Sprintf("  VPC Peering 삭제: %s (%s)", p.VpcPeeringName, p.VpcPeeringInstanceNo))
		if err := c.DeleteVpcPeeringInstance(p.VpcPeeringInstanceNo); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공]")
			success++
		}
	}

	// 14. Route Tables (Clean up NAT/Peering routes)
	logFn("  Route Table 정리 (NAT/Peering 경로 삭제)...")
	for _, rt := range summary.RouteTables {
		for _, route := range rt.RouteList {
			if route.TargetTypeCode.Code == "NATGW" || route.TargetTypeCode.Code == "VPCPEERING" {
				logFn(fmt.Sprintf("    경로 삭제: Table %s -> %s", rt.RouteTableName, route.DestinationCidrBlock))
				if err := c.RemoveRoute(rt.RouteTableNo, route); err != nil {
					logFn(fmt.Sprintf("      [실패] %v", err))
				} else {
					logFn("      [성공]")
				}
			}
		}
	}

	// 15. Release Public IPs
	for _, ip := range summary.PublicIps {
		logFn(fmt.Sprintf("  공인 IP 해제: %s (%s)", ip.PublicIp, ip.PublicIpInstanceNo))
		if ip.ServerInstanceNo != "" {
			if err := c.DisassociatePublicIp(ip.PublicIpInstanceNo); err != nil {
				logFn(fmt.Sprintf("    [실패] 연결 해제: %v", err))
			} else {
				time.Sleep(3 * time.Second)
			}
		}
		if err := c.DeletePublicIp(ip.PublicIpInstanceNo); err != nil {
			logFn(fmt.Sprintf("    [실패] IP 삭제: %v", err))
			fail++
		} else {
			logFn("    [성공]")
			success++
		}
	}

	// 16. Access Control Groups (skip default - deleted with VPC)
	var acgsToDelete []string
	for _, acg := range summary.AccessControlGroups {
		if !acg.IsDefault {
			acgsToDelete = append(acgsToDelete, acg.AccessControlGroupNo)
		}
	}
	if len(acgsToDelete) > 0 {
		logFn(fmt.Sprintf("  ACG %d개 삭제 중 (Default 제외)...", len(acgsToDelete)))
		for _, acgNo := range acgsToDelete {
			if err := c.DeleteAccessControlGroup(acgNo); err != nil {
				logFn(fmt.Sprintf("    [실패] ACG(%s) 삭제: %v", acgNo, err))
				fail++
			} else {
				logFn(fmt.Sprintf("    [성공] ACG(%s) 삭제", acgNo))
				success++
			}
		}
	}

	// 17. Network ACLs (skip default - deleted with VPC)
	for _, nacl := range summary.NetworkAcls {
		if nacl.IsDefault {
			continue
		}
		logFn(fmt.Sprintf("  Network ACL 삭제: %s (%s)", nacl.NetworkAclName, nacl.NetworkAclNo))
		if err := c.DeleteNetworkAcl(nacl.NetworkAclNo); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공]")
			success++
		}
	}

	// 18. Subnets
	for _, subnet := range summary.Subnets {
		logFn(fmt.Sprintf("  Subnet 삭제: %s (%s)", subnet.SubnetName, subnet.SubnetNo))
		if err := c.DeleteSubnet(subnet.SubnetNo); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공]")
			success++
		}
	}

	// 19. VPCs
	for _, vpc := range summary.Vpcs {
		logFn(fmt.Sprintf("  VPC 삭제: %s (%s)", vpc.VpcName, vpc.VpcNo))
		if err := c.DeleteVpc(vpc.VpcNo); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공]")
			success++
		}
	}

	// 20. Init Scripts
	if len(summary.InitScripts) > 0 {
		var scriptNos []string
		for _, s := range summary.InitScripts {
			scriptNos = append(scriptNos, s.InitScriptNo)
		}
		logFn(fmt.Sprintf("  Init Script %d개 삭제 중...", len(scriptNos)))
		if err := c.DeleteInitScripts(scriptNos); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail += len(scriptNos)
		} else {
			logFn("    [성공]")
			success += len(scriptNos)
		}
	}

	// 21. Login Keys
	for _, key := range summary.LoginKeys {
		logFn(fmt.Sprintf("  Login Key 삭제: %s", key.KeyName))
		if err := c.DeleteLoginKey(key.KeyName); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공]")
			success++
		}
	}

	// 22. Placement Groups
	for _, pg := range summary.PlacementGroups {
		logFn(fmt.Sprintf("  Placement Group 삭제: %s (%s)", pg.PlacementGroupName, pg.PlacementGroupNo))
		if err := c.DeletePlacementGroup(pg.PlacementGroupNo); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공]")
			success++
		}
	}

	return success, fail
}

// --- Delete/Terminate APIs ---

func (c *Client) StopServers(instanceNos []string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	for i, no := range instanceNos {
		params.Set(fmt.Sprintf("serverInstanceNoList.%d", i+1), no)
	}
	path := "/stopServerInstances?" + params.Encode()
	body, status, err := c.doRequestWithBase(VServerBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

func (c *Client) TerminateServers(instanceNos []string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	for i, no := range instanceNos {
		params.Set(fmt.Sprintf("serverInstanceNoList.%d", i+1), no)
	}
	path := "/terminateServerInstances?" + params.Encode()
	body, status, err := c.doRequestWithBase(VServerBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

func (c *Client) DeleteBlockStorages(instanceNos []string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	for i, no := range instanceNos {
		params.Set(fmt.Sprintf("blockStorageInstanceNoList.%d", i+1), no)
	}
	path := "/deleteBlockStorageInstances?" + params.Encode()
	body, status, err := c.doRequestWithBase(VServerBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

func (c *Client) DisassociatePublicIp(publicIpInstanceNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("publicIpInstanceNo", publicIpInstanceNo)
	path := "/disassociatePublicIpFromServerInstance?" + params.Encode()
	body, status, err := c.doRequestWithBase(VServerBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

func (c *Client) DeletePublicIp(publicIpInstanceNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("publicIpInstanceNo", publicIpInstanceNo)
	path := "/deletePublicIpInstance?" + params.Encode()
	body, status, err := c.doRequestWithBase(VServerBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

func (c *Client) DeleteNasVolume(nasVolumeInstanceNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("nasVolumeInstanceNo", nasVolumeInstanceNo)
	path := "/deleteNasVolumeInstance?" + params.Encode()
	body, status, err := c.doRequestWithBase(VNASBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

func (c *Client) DeleteLoadBalancer(loadBalancerInstanceNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("loadBalancerInstanceNo", loadBalancerInstanceNo)
	path := "/deleteLoadBalancerInstances?" + params.Encode()
	body, status, err := c.doRequestWithBase(VLBBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}
