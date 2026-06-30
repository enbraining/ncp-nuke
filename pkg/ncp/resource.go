package ncp

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
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
	Buckets                 []Bucket
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
		len(r.InitScripts) + len(r.LoginKeys) + len(r.PlacementGroups) +
		len(r.Buckets)
}

// ResourceCount is a single category's resource count for display.
type ResourceCount struct {
	Name  string
	Count int
}

// Breakdown returns a per-category resource count in display order,
// including only categories that currently have at least one resource.
func (r *ResourceSummary) Breakdown() []ResourceCount {
	all := []ResourceCount{
		{"Server", len(r.Servers)},
		{"Block Storage", len(r.BlockStorages)},
		{"Block Storage Snapshot", len(r.BlockStorageSnapshots)},
		{"Public IP", len(r.PublicIps)},
		{"NAS Volume", len(r.NasVolumes)},
		{"NAS Volume Snapshot", len(r.NasVolumeSnapshots)},
		{"Load Balancer", len(r.LoadBalancers)},
		{"Target Group", len(r.TargetGroups)},
		{"Cloud DB", len(r.CloudDBs)},
		{"Cloud PostgreSQL", len(r.CloudPostgresqls)},
		{"Cloud MongoDB", len(r.CloudMongoDBs)},
		{"Cloud MariaDB", len(r.CloudMariaDBs)},
		{"Cloud MySQL", len(r.CloudMySQLs)},
		{"Cloud Redis", len(r.CloudRedises)},
		{"VPC", len(r.Vpcs)},
		{"Subnet", len(r.Subnets)},
		{"NAT Gateway", len(r.NatGateways)},
		{"VPC Peering", len(r.VpcPeerings)},
		{"Network ACL", len(r.NetworkAcls)},
		{"Route Table", len(r.RouteTables)},
		{"Access Control Group", len(r.AccessControlGroups)},
		{"Auto Scaling Group", len(r.AutoScalingGroups)},
		{"Launch Configuration", len(r.LaunchConfigurations)},
		{"NKS Cluster", len(r.NksClusters)},
		{"Init Script", len(r.InitScripts)},
		{"Login Key", len(r.LoginKeys)},
		{"Placement Group", len(r.PlacementGroups)},
		{"Object Storage Bucket", len(r.Buckets)},
	}
	var out []ResourceCount
	for _, c := range all {
		if c.Count > 0 {
			out = append(out, c)
		}
	}
	return out
}

// ResourceItem is a single resource's display identity (name + id).
type ResourceItem struct {
	Name string
	ID   string
}

// Items returns the individual resources per category (keyed by the same names
// as Breakdown), for showing a detailed per-resource list.
func (r *ResourceSummary) Items() map[string][]ResourceItem {
	m := map[string][]ResourceItem{}
	add := func(key, name, id string) { m[key] = append(m[key], ResourceItem{Name: name, ID: id}) }

	for _, x := range r.Servers {
		add("Server", x.ServerName, x.ServerInstanceNo)
	}
	for _, x := range r.BlockStorages {
		add("Block Storage", x.BlockStorageName, x.BlockStorageInstanceNo)
	}
	for _, x := range r.BlockStorageSnapshots {
		add("Block Storage Snapshot", x.BlockStorageSnapshotName, x.BlockStorageSnapshotInstanceNo)
	}
	for _, x := range r.PublicIps {
		add("Public IP", x.PublicIp, x.PublicIpInstanceNo)
	}
	for _, x := range r.NasVolumes {
		add("NAS Volume", x.VolumeName, x.NasVolumeInstanceNo)
	}
	for _, x := range r.NasVolumeSnapshots {
		add("NAS Volume Snapshot", x.NasVolumeSnapshotName, x.NasVolumeSnapshotInstanceNo)
	}
	for _, x := range r.LoadBalancers {
		add("Load Balancer", x.LoadBalancerName, x.LoadBalancerInstanceNo)
	}
	for _, x := range r.TargetGroups {
		add("Target Group", x.TargetGroupName, x.TargetGroupNo)
	}
	for _, x := range r.CloudDBs {
		add("Cloud DB", x.CloudDBServiceName, x.CloudDBInstanceNo)
	}
	for _, x := range r.CloudPostgresqls {
		add("Cloud PostgreSQL", x.CloudPostgresqlServiceName, x.CloudPostgresqlInstanceNo)
	}
	for _, x := range r.CloudMongoDBs {
		add("Cloud MongoDB", x.CloudMongoDbServiceName, x.CloudMongoDbInstanceNo)
	}
	for _, x := range r.CloudMariaDBs {
		add("Cloud MariaDB", x.CloudMariaDbServiceName, x.CloudMariaDbInstanceNo)
	}
	for _, x := range r.CloudMySQLs {
		add("Cloud MySQL", x.CloudMysqlServiceName, x.CloudMysqlInstanceNo)
	}
	for _, x := range r.CloudRedises {
		add("Cloud Redis", x.CloudRedisServiceName, x.CloudRedisInstanceNo)
	}
	for _, x := range r.Vpcs {
		add("VPC", x.VpcName, x.VpcNo)
	}
	for _, x := range r.Subnets {
		add("Subnet", x.SubnetName, x.SubnetNo)
	}
	for _, x := range r.NatGateways {
		add("NAT Gateway", x.NatGatewayName, x.NatGatewayInstanceNo)
	}
	for _, x := range r.VpcPeerings {
		add("VPC Peering", x.VpcPeeringName, x.VpcPeeringInstanceNo)
	}
	for _, x := range r.NetworkAcls {
		add("Network ACL", x.NetworkAclName, x.NetworkAclNo)
	}
	for _, x := range r.RouteTables {
		add("Route Table", x.RouteTableName, x.RouteTableNo)
	}
	for _, x := range r.AccessControlGroups {
		add("Access Control Group", x.AccessControlGroupName, x.AccessControlGroupNo)
	}
	for _, x := range r.AutoScalingGroups {
		add("Auto Scaling Group", x.AutoScalingGroupName, x.AutoScalingGroupNo)
	}
	for _, x := range r.LaunchConfigurations {
		add("Launch Configuration", x.LaunchConfigurationName, x.LaunchConfigurationNo)
	}
	for _, x := range r.NksClusters {
		add("NKS Cluster", x.Name, x.Uuid)
	}
	for _, x := range r.InitScripts {
		add("Init Script", x.InitScriptName, x.InitScriptNo)
	}
	for _, x := range r.LoginKeys {
		add("Login Key", x.KeyName, "")
	}
	for _, x := range r.PlacementGroups {
		add("Placement Group", x.PlacementGroupName, x.PlacementGroupNo)
	}
	for _, x := range r.Buckets {
		add("Object Storage Bucket", x.Name, "")
	}
	return m
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
	// NAS snapshots must be queried per volume (nasVolumeInstanceNo is required).
	for _, vol := range summary.NasVolumes {
		if snaps, err := c.ListNasVolumeSnapshots(vol.NasVolumeInstanceNo); err != nil {
			errs = append(errs, fmt.Errorf("NAS 스냅샷 조회(%s): %w", vol.VolumeName, err))
		} else {
			summary.NasVolumeSnapshots = append(summary.NasVolumeSnapshots, snaps...)
		}
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

	// 9. Object Storage (S3 호환)
	if buckets, err := c.ListBuckets(); err != nil {
		errs = append(errs, fmt.Errorf("Object Storage 조회: %w", err))
	} else {
		summary.Buckets = buckets
	}

	return summary, errs
}

// CleanupAllResources deletes all resources for a root account.
// 올바른 의존성 순서 (하위 리소스 -> 상위 리소스):
// 1. Application Layer: NKS -> ASG -> LaunchConfig
// 2. Data Layer: All DBs
// 3. Network Services: LB -> TargetGroup
// 4. Compute: Server -> BS Snapshot -> BS -> NAS Snapshot -> NAS
// 5. Network Infrastructure (순서 중요!):
//    - Routes (의존성 제거) -> VPC Peering -> NAT Gateway
//    - Public IP -> ACG -> Network ACL
//    - Subnet (완전 삭제 대기) -> VPC
// 6. Server Resources: InitScript -> LoginKey -> PlacementGroup
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
	// A non-empty ASG cannot be deleted (returnCode 1250600). Set its capacity to
	// 0 first so its servers terminate, wait, then delete (retrying while servers
	// drain). Termination protection on ASG-managed servers would block draining,
	// so disable it up front.
	if len(summary.AutoScalingGroups) > 0 {
		if len(summary.Servers) > 0 {
			logFn("  ASG 서버 반납 보호 해제 중...")
			c.disableServerProtection(summary.Servers, logFn)
		}
		for _, asg := range summary.AutoScalingGroups {
			logFn(fmt.Sprintf("  ASG 용량(최소/최대/기대) 0으로 설정: %s (%s)", asg.AutoScalingGroupName, asg.AutoScalingGroupNo))
			if err := c.SetAutoScalingGroupSizeZero(asg.AutoScalingGroupNo); err != nil {
				logFn(fmt.Sprintf("    [경고] 용량 0 설정 실패: %v", err))
			}
		}
		logFn("  ASG 서버 종료 대기 (60초)...")
		time.Sleep(60 * time.Second)

		for _, asg := range summary.AutoScalingGroups {
			logFn(fmt.Sprintf("  Auto Scaling Group 삭제: %s (%s)", asg.AutoScalingGroupName, asg.AutoScalingGroupNo))
			maxRetries := 5
			var lastErr error
			for retry := 0; retry < maxRetries; retry++ {
				if retry > 0 {
					logFn(fmt.Sprintf("    서버 종료 대기 후 재시도 %d/%d...", retry+1, maxRetries))
					time.Sleep(30 * time.Second)
				}
				if err := c.DeleteAutoScalingGroup(asg.AutoScalingGroupNo); err != nil {
					lastErr = err
					if retry < maxRetries-1 {
						logFn(fmt.Sprintf("    [대기] %v", err))
					}
				} else {
					lastErr = nil
					break
				}
			}
			if lastErr != nil {
				logFn(fmt.Sprintf("    [실패] %v", lastErr))
				fail++
			} else {
				logFn("    [성공]")
				success++
			}
		}
		logFn("  ASG 삭제 및 서버 종료 대기 (30초)...")
		time.Sleep(30 * time.Second)
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

	// 7. Stop -> disable termination protection -> Terminate Servers.
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
				logFn("    서버 정지 요청 완료, 정지 완료 대기 중...")
				c.waitForServersStopped(runningNos, logFn)
			}
		}

		// 반납 보호 해제 (보호된 서버는 반납이 막히므로 먼저 전부 해제)
		logFn("  서버 반납 보호 해제 중...")
		c.disableServerProtection(summary.Servers, logFn)

		var allNos []string
		for _, s := range summary.Servers {
			allNos = append(allNos, s.ServerInstanceNo)
		}
		logFn(fmt.Sprintf("  서버 %d대 반납(삭제) 중...", len(allNos)))
		var err error
		for retry := 0; retry < 3; retry++ {
			if retry > 0 {
				time.Sleep(5 * time.Second)
				logFn(fmt.Sprintf("    재시도 %d/3...", retry+1))
			}
			if err = c.TerminateServers(allNos); err == nil {
				break
			}
			logFn(fmt.Sprintf("    [대기] %v", err))
		}
		if err != nil {
			logFn(fmt.Sprintf("    [실패] 서버 반납: %v", err))
			fail += len(allNos)
		} else {
			logFn("    [성공] 서버 반납 요청 완료")
			success += len(allNos)
			logFn("    서버 반납 완료 대기 중...")
			c.waitForServersTerminated(summary.Servers, logFn)
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

	// 12. Route Tables: remove routes that target the NAT Gateways / VPC Peerings
	// we are about to delete. getRouteList does not reliably populate
	// targetTypeCode, so match each route by targetNo against the in-scope
	// resources and supply the type code explicitly to removeRoute.
	if len(summary.RouteTables) > 0 {
		natNos := make(map[string]bool)
		for _, n := range summary.NatGateways {
			natNos[n.NatGatewayInstanceNo] = true
		}
		peeringNos := make(map[string]bool)
		for _, p := range summary.VpcPeerings {
			peeringNos[p.VpcPeeringInstanceNo] = true
		}

		logFn(fmt.Sprintf("  Route Table 정리 (NAT/Peering 경로 삭제)... 테이블 %d개", len(summary.RouteTables)))
		routesToDelete := 0
		for _, rt := range summary.RouteTables {
			// getRouteTableList does not include routes; fetch them per table.
			routes := rt.RouteList
			if fetched, err := c.ListRoutes(rt.VpcNo, rt.RouteTableNo); err != nil {
				logFn(fmt.Sprintf("    [경고] 경로 조회 실패 (Table %s, vpcNo=%q): %v", rt.RouteTableName, rt.VpcNo, err))
			} else {
				routes = fetched
			}

			logFn(fmt.Sprintf("    Table %s (vpcNo=%s): 경로 %d개", rt.RouteTableName, rt.VpcNo, len(routes)))
			for _, route := range routes {
				var typeCode string
				switch {
				case natNos[route.TargetNo]:
					typeCode = "NATGW"
				case peeringNos[route.TargetNo]:
					typeCode = "VPCPEERING"
				default:
					continue
				}
				routesToDelete++
				logFn(fmt.Sprintf("      경로 삭제 시도: %s -> %s (%s, target=%s)", rt.RouteTableName, route.DestinationCidrBlock, typeCode, route.TargetName))
				r := route
				r.TargetTypeCode.Code = typeCode
				if err := c.RemoveRoute(rt.VpcNo, rt.RouteTableNo, r); err != nil {
					logFn(fmt.Sprintf("        [실패] %v", err))
				} else {
					logFn("        [성공]")
				}
			}
		}
		if routesToDelete == 0 {
			logFn("    삭제할 NAT/Peering 경로 없음")
		}
	}

	// 13. VPC Peering (Route 삭제 후)
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

	// 14. NAT Gateways (Route 삭제 후)
	// Route removal is asynchronous on NCP, so a NAT delete right after may still
	// see the route and fail (returnCode 1018005). Retry with a short delay.
	for _, nat := range summary.NatGateways {
		logFn(fmt.Sprintf("  NAT Gateway 삭제: %s (%s)", nat.NatGatewayName, nat.NatGatewayInstanceNo))

		maxRetries := 5
		var lastErr error
		for retry := 0; retry < maxRetries; retry++ {
			if retry > 0 {
				logFn(fmt.Sprintf("    경로 반영 대기 후 재시도 %d/%d...", retry+1, maxRetries))
				time.Sleep(10 * time.Second)
			}
			if err := c.DeleteNatGateway(nat.NatGatewayInstanceNo); err != nil {
				lastErr = err
				if retry < maxRetries-1 {
					logFn(fmt.Sprintf("    [대기] %v", err))
				}
			} else {
				lastErr = nil
				break
			}
		}

		if lastErr != nil {
			logFn(fmt.Sprintf("    [실패] %v", lastErr))
			fail++
		} else {
			logFn("    [성공]")
			success++
		}
	}
	if len(summary.NatGateways) > 0 {
		logFn("  NAT Gateway 삭제 완료 대기 중...")
		c.waitForNatGatewaysDeletion(summary.NatGateways, logFn)
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
	var acgsToDelete []AccessControlGroup
	for _, acg := range summary.AccessControlGroups {
		if !acg.IsDefault {
			acgsToDelete = append(acgsToDelete, acg)
		}
	}
	if len(acgsToDelete) > 0 {
		logFn(fmt.Sprintf("  ACG %d개 삭제 중 (Default 제외)...", len(acgsToDelete)))
		for _, acg := range acgsToDelete {
			if err := c.DeleteAccessControlGroup(acg.VpcNo, acg.AccessControlGroupNo); err != nil {
				logFn(fmt.Sprintf("    [실패] ACG(%s) 삭제: %v", acg.AccessControlGroupName, err))
				fail++
			} else {
				logFn(fmt.Sprintf("    [성공] ACG(%s) 삭제", acg.AccessControlGroupName))
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

	// 18-1. Subnet 삭제 완료 대기 (VPC 삭제 전 필수)
	if len(summary.Subnets) > 0 {
		logFn("  Subnet 삭제 완료 대기 중...")
		c.waitForSubnetsDeletion(summary.Subnets, logFn)
	}

	// 19. VPCs (Subnet이 모두 삭제된 후에만 가능)
	for _, vpc := range summary.Vpcs {
		logFn(fmt.Sprintf("  VPC 삭제: %s (%s)", vpc.VpcName, vpc.VpcNo))

		// 재시도 로직: Subnet 삭제가 완료되지 않았을 수 있음
		maxRetries := 3
		var lastErr error

		for retry := 0; retry < maxRetries; retry++ {
			if retry > 0 {
				logFn(fmt.Sprintf("    재시도 %d/%d...", retry+1, maxRetries))
				time.Sleep(10 * time.Second)
			}

			if err := c.DeleteVpc(vpc.VpcNo); err != nil {
				lastErr = err
				if retry < maxRetries-1 {
					logFn(fmt.Sprintf("    [실패] %v (재시도 예정)", err))
				}
			} else {
				logFn("    [성공]")
				success++
				lastErr = nil
				break
			}
		}

		if lastErr != nil {
			logFn(fmt.Sprintf("    [실패] %v", lastErr))
			fail++
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

	// 23. Object Storage Buckets (객체 전부 비운 뒤 버킷 삭제, VPC와 독립)
	for _, b := range summary.Buckets {
		logFn(fmt.Sprintf("  Object Storage 버킷 비우기/삭제: %s", b.Name))
		if err := c.DeleteBucket(b.Name, logFn); err != nil {
			logFn(fmt.Sprintf("    [실패] %v", err))
			fail++
		} else {
			logFn("    [성공]")
			success++
		}
	}

	return success, fail
}

// disableServerProtection turns off 반납 보호 (termination protection) on every
// given server so they can be terminated (by us or by a draining ASG).
func (c *Client) disableServerProtection(servers []ServerInstance, logFn func(string)) {
	for _, s := range servers {
		if err := c.SetServerTerminationProtection(s.ServerInstanceNo, false); err != nil {
			logFn(fmt.Sprintf("    [경고] 반납 보호 해제 실패 (%s): %v", s.ServerName, err))
		}
	}
}

func (c *Client) waitForNatGatewaysDeletion(natGateways []NatGatewayInstance, logFn func(string)) {
	maxWait := 5 * time.Minute
	pollInterval := 10 * time.Second
	deadline := time.Now().Add(maxWait)

	targetNos := make(map[string]bool)
	for _, nat := range natGateways {
		targetNos[nat.NatGatewayInstanceNo] = true
	}

	for time.Now().Before(deadline) {
		remaining, err := c.ListNatGateways()
		if err != nil {
			logFn(fmt.Sprintf("    [경고] NAT Gateway 조회 실패: %v, 재시도...", err))
			time.Sleep(pollInterval)
			continue
		}

		stillExists := 0
		for _, nat := range remaining {
			if targetNos[nat.NatGatewayInstanceNo] {
				stillExists++
			}
		}

		if stillExists == 0 {
			logFn("    NAT Gateway 삭제 완료 확인")
			return
		}

		logFn(fmt.Sprintf("    아직 %d개 NAT Gateway 삭제 중... (%d초 후 재확인)", stillExists, int(pollInterval.Seconds())))
		time.Sleep(pollInterval)
	}

	logFn("    [경고] NAT Gateway 삭제 대기 시간 초과")
}

func (c *Client) waitForSubnetsDeletion(subnets []Subnet, logFn func(string)) {
	maxWait := 5 * time.Minute
	pollInterval := 10 * time.Second
	deadline := time.Now().Add(maxWait)

	targetNos := make(map[string]bool)
	for _, subnet := range subnets {
		targetNos[subnet.SubnetNo] = true
	}

	for time.Now().Before(deadline) {
		remaining, err := c.ListSubnets()
		if err != nil {
			logFn(fmt.Sprintf("    [경고] Subnet 조회 실패: %v, 재시도...", err))
			time.Sleep(pollInterval)
			continue
		}

		stillExists := 0
		for _, s := range remaining {
			if targetNos[s.SubnetNo] {
				stillExists++
			}
		}

		if stillExists == 0 {
			logFn("    Subnet 삭제 완료 확인")
			return
		}

		logFn(fmt.Sprintf("    아직 %d개 Subnet 삭제 중... (%d초 후 재확인)", stillExists, int(pollInterval.Seconds())))
		time.Sleep(pollInterval)
	}

	logFn("    [경고] Subnet 삭제 대기 시간 초과, VPC 삭제를 계속 시도합니다")
}

func (c *Client) waitForServersStopped(serverNos []string, logFn func(string)) {
	maxWait := 3 * time.Minute
	pollInterval := 10 * time.Second
	deadline := time.Now().Add(maxWait)

	targetNos := make(map[string]bool)
	for _, no := range serverNos {
		targetNos[no] = true
	}

	for time.Now().Before(deadline) {
		servers, err := c.ListServers()
		if err != nil {
			logFn(fmt.Sprintf("    [경고] 서버 조회 실패: %v, 재시도...", err))
			time.Sleep(pollInterval)
			continue
		}

		stillRunning := 0
		for _, s := range servers {
			if targetNos[s.ServerInstanceNo] && s.ServerInstanceStatus.Code == "RUN" {
				stillRunning++
			}
		}

		if stillRunning == 0 {
			logFn("    서버 정지 완료 확인")
			return
		}

		logFn(fmt.Sprintf("    아직 %d대 서버 정지 중... (%d초 후 재확인)", stillRunning, int(pollInterval.Seconds())))
		time.Sleep(pollInterval)
	}

	logFn("    [경고] 서버 정지 대기 시간 초과, 반납을 계속 시도합니다")
}

func (c *Client) waitForServersTerminated(servers []ServerInstance, logFn func(string)) {
	maxWait := 5 * time.Minute
	pollInterval := 10 * time.Second
	deadline := time.Now().Add(maxWait)

	targetNos := make(map[string]bool)
	for _, s := range servers {
		targetNos[s.ServerInstanceNo] = true
	}

	for time.Now().Before(deadline) {
		remaining, err := c.ListServers()
		if err != nil {
			logFn(fmt.Sprintf("    [경고] 서버 조회 실패: %v, 재시도...", err))
			time.Sleep(pollInterval)
			continue
		}

		stillExists := 0
		for _, s := range remaining {
			if targetNos[s.ServerInstanceNo] {
				stillExists++
			}
		}

		if stillExists == 0 {
			logFn("    서버 반납 완료 확인")
			return
		}

		logFn(fmt.Sprintf("    아직 %d대 서버 반납 중... (%d초 후 재확인)", stillExists, int(pollInterval.Seconds())))
		time.Sleep(pollInterval)
	}

	logFn("    [경고] 서버 반납 대기 시간 초과")
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

// SetServerTerminationProtection enables/disables 반납 보호 (termination
// protection) for a server so it can be terminated.
func (c *Client) SetServerTerminationProtection(serverInstanceNo string, protect bool) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("serverInstanceNo", serverInstanceNo)
	params.Set("isProtectServerTermination", strconv.FormatBool(protect))
	path := "/setProtectServerTermination?" + params.Encode()
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
