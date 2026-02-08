package ncp

import "encoding/json"

// --- Common ---
type CommonCode struct {
	Code     string `json:"code"`
	CodeName string `json:"codeName"`
}

// Generic NCP API response wrapper
type apiResponse struct {
	RequestId     string          `json:"requestId"`
	ReturnCode    int             `json:"returnCode"`
	ReturnMessage string          `json:"returnMessage"`
	TotalRows     int             `json:"totalRows"`
	Content       json.RawMessage `json:"-"`
}

// --- Server ---
type ServerInstance struct {
	ServerInstanceNo     string     `json:"serverInstanceNo"`
	ServerName           string     `json:"serverName"`
	ServerInstanceStatus CommonCode `json:"serverInstanceStatus"`
	PublicIp             string     `json:"publicIp"`
	PrivateIp            string     `json:"privateIp"`
	CpuCount             int        `json:"cpuCount"`
	MemorySize           int64      `json:"memorySize"`
}

// --- Block Storage ---
type BlockStorageInstance struct {
	BlockStorageInstanceNo     string     `json:"blockStorageInstanceNo"`
	BlockStorageName           string     `json:"blockStorageName"`
	BlockStorageInstanceStatus CommonCode `json:"blockStorageInstanceStatus"`
	BlockStorageSize           int64      `json:"blockStorageSize"`
	ServerInstanceNo           string     `json:"serverInstanceNo"`
	BlockStorageType           CommonCode `json:"blockStorageType"`
	BlockStorageDiskDetailType CommonCode `json:"blockStorageDiskDetailType"`
}

// --- Public IP ---
type PublicIpInstance struct {
	PublicIpInstanceNo     string     `json:"publicIpInstanceNo"`
	PublicIp               string     `json:"publicIp"`
	PublicIpInstanceStatus CommonCode `json:"publicIpInstanceStatus"`
	ServerInstanceNo       string     `json:"serverInstanceNo"`
	ServerName             string     `json:"serverName"`
}

// --- NAS Volume ---
type NasVolumeInstance struct {
	NasVolumeInstanceNo     string     `json:"nasVolumeInstanceNo"`
	VolumeName              string     `json:"volumeName"`
	NasVolumeInstanceStatus CommonCode `json:"nasVolumeInstanceStatus"`
	VolumeAllotmentProtocol CommonCode `json:"volumeAllotmentProtocolType"`
	VolumeTotalSize         int64      `json:"volumeTotalSize"`
}

// --- Load Balancer ---
type LoadBalancerInstance struct {
	LoadBalancerInstanceNo     string     `json:"loadBalancerInstanceNo"`
	LoadBalancerName           string     `json:"loadBalancerName"`
	LoadBalancerInstanceStatus CommonCode `json:"loadBalancerInstanceStatus"`
	LoadBalancerType           CommonCode `json:"loadBalancerType"`
}

// --- Cloud DB ---
type CloudDBInstance struct {
	CloudDBInstanceNo     string     `json:"cloudDBInstanceNo"`
	CloudDBServiceName    string     `json:"cloudDBServiceName"`
	CloudDBInstanceStatus CommonCode `json:"cloudDBInstanceStatus"`
	DBKindCode            string     `json:"dbKindCode"` // MYSQL, MSSQL, REDIS
}

// --- Cloud DB for PostgreSQL ---
type CloudPostgresqlInstance struct {
	CloudPostgresqlInstanceNo     string     `json:"cloudPostgresqlInstanceNo"`
	CloudPostgresqlServiceName    string     `json:"cloudPostgresqlServiceName"`
	CloudPostgresqlInstanceStatus CommonCode `json:"cloudPostgresqlInstanceStatus"`
}

// --- Cloud DB for MongoDB ---
type CloudMongoDbInstance struct {
	CloudMongoDbInstanceNo     string     `json:"cloudMongoDbInstanceNo"`
	CloudMongoDbServiceName    string     `json:"cloudMongoDbServiceName"`
	CloudMongoDbInstanceStatus CommonCode `json:"cloudMongoDbInstanceStatus"`
}

// --- VPC ---
type Vpc struct {
	VpcNo     string     `json:"vpcNo"`
	VpcName   string     `json:"vpcName"`
	VpcStatus CommonCode `json:"vpcStatus"`
	Ipv4Cidr  string     `json:"ipv4CidrBlock"`
}

type Subnet struct {
	SubnetNo     string     `json:"subnetNo"`
	SubnetName   string     `json:"subnetName"`
	SubnetStatus CommonCode `json:"subnetStatus"`
	VpcNo        string     `json:"vpcNo"`
}

type NatGatewayInstance struct {
	NatGatewayInstanceNo     string     `json:"natGatewayInstanceNo"`
	NatGatewayName           string     `json:"natGatewayName"`
	NatGatewayInstanceStatus CommonCode `json:"natGatewayInstanceStatus"`
	VpcNo                    string     `json:"vpcNo"`
}

// --- Route Table ---
type Route struct {
	DestinationCidrBlock string     `json:"destinationCidrBlock"`
	TargetName           string     `json:"targetName"`
	TargetNo             string     `json:"targetNo"`
	TargetTypeCode       CommonCode `json:"targetTypeCode"` // NATGW, LOCAL, VPCPEERING, VGW
}

type RouteTable struct {
	RouteTableNo     string     `json:"routeTableNo"`
	RouteTableName   string     `json:"routeTableName"`
	RouteTableStatus CommonCode `json:"routeTableStatus"`
	VpcNo            string     `json:"vpcNo"`
	IsDefault        bool       `json:"isDefault"`
	RouteList        []Route    `json:"routeList"`
}

// --- Access Control Group (ACG) ---
type AccessControlGroup struct {
	AccessControlGroupNo     string     `json:"accessControlGroupNo"`
	AccessControlGroupName   string     `json:"accessControlGroupName"`
	AccessControlGroupStatus CommonCode `json:"accessControlGroupStatus"`
	VpcNo                    string     `json:"vpcNo"`
	IsDefault                bool       `json:"isDefault"`
}

// --- Auto Scaling ---
type AutoScalingGroup struct {
	AutoScalingGroupName string     `json:"autoScalingGroupName"`
	AutoScalingGroupNo   string     `json:"autoScalingGroupNo"`
	InAutoScalingGroupNo string     `json:"inAutoScalingGroupNo"` // Sometimes used
	HealthCheckTypeCode  CommonCode `json:"healthCheckTypeCode"`
}

// --- NKS (Kubernetes) ---
type NksCluster struct {
	Uuid   string `json:"uuid"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// --- Cloud DB for MariaDB ---
type CloudMariaDbInstance struct {
	CloudMariaDbInstanceNo     string     `json:"cloudMariaDbInstanceNo"`
	CloudMariaDbServiceName    string     `json:"cloudMariaDbServiceName"`
	CloudMariaDbInstanceStatus CommonCode `json:"cloudMariaDbInstanceStatus"`
}

// --- Cloud DB for MySQL (dedicated) ---
type CloudMysqlInstance struct {
	CloudMysqlInstanceNo     string     `json:"cloudMysqlInstanceNo"`
	CloudMysqlServiceName    string     `json:"cloudMysqlServiceName"`
	CloudMysqlInstanceStatus CommonCode `json:"cloudMysqlInstanceStatus"`
}

// --- Cloud DB for Redis (dedicated) ---
type CloudRedisInstance struct {
	CloudRedisInstanceNo     string     `json:"cloudRedisInstanceNo"`
	CloudRedisServiceName    string     `json:"cloudRedisServiceName"`
	CloudRedisInstanceStatus CommonCode `json:"cloudRedisInstanceStatus"`
}

// --- Launch Configuration (Auto Scaling) ---
type LaunchConfiguration struct {
	LaunchConfigurationNo   string `json:"launchConfigurationNo"`
	LaunchConfigurationName string `json:"launchConfigurationName"`
}

// --- Network ACL ---
type NetworkAcl struct {
	NetworkAclNo     string     `json:"networkAclNo"`
	NetworkAclName   string     `json:"networkAclName"`
	NetworkAclStatus CommonCode `json:"networkAclStatus"`
	VpcNo            string     `json:"vpcNo"`
	IsDefault        bool       `json:"isDefault"`
}

// --- VPC Peering ---
type VpcPeeringInstance struct {
	VpcPeeringInstanceNo     string     `json:"vpcPeeringInstanceNo"`
	VpcPeeringName           string     `json:"vpcPeeringName"`
	VpcPeeringInstanceStatus CommonCode `json:"vpcPeeringInstanceStatus"`
}

// --- Init Script ---
type InitScript struct {
	InitScriptNo   string `json:"initScriptNo"`
	InitScriptName string `json:"initScriptName"`
}

// --- Login Key ---
type LoginKey struct {
	KeyName string `json:"keyName"`
}

// --- Placement Group ---
type PlacementGroup struct {
	PlacementGroupNo   string `json:"placementGroupNo"`
	PlacementGroupName string `json:"placementGroupName"`
}

// --- Target Group (Load Balancer) ---
type TargetGroup struct {
	TargetGroupNo     string     `json:"targetGroupNo"`
	TargetGroupName   string     `json:"targetGroupName"`
	TargetGroupStatus CommonCode `json:"targetGroupStatus"`
}

// --- Block Storage Snapshot ---
type BlockStorageSnapshotInstance struct {
	BlockStorageSnapshotInstanceNo     string     `json:"blockStorageSnapshotInstanceNo"`
	BlockStorageSnapshotName           string     `json:"blockStorageSnapshotName"`
	BlockStorageSnapshotInstanceStatus CommonCode `json:"blockStorageSnapshotInstanceStatus"`
}

// --- NAS Volume Snapshot ---
type NasVolumeSnapshot struct {
	NasVolumeSnapshotInstanceNo string `json:"nasVolumeSnapshotInstanceNo"`
	NasVolumeSnapshotName       string `json:"nasVolumeSnapshotName"`
	NasVolumeInstanceNo         string `json:"nasVolumeInstanceNo"`
}

// ListResponse wrappers for new services
type getCloudDBInstanceListResponse struct {
	CloudDBInstanceList []CloudDBInstance `json:"cloudDBInstanceList"`
}

type getCloudPostgresqlInstanceListResponse struct {
	CloudPostgresqlInstanceList []CloudPostgresqlInstance `json:"cloudPostgresqlInstanceList"`
}

type getCloudMongoDbInstanceListResponse struct {
	CloudMongoDbInstanceList []CloudMongoDbInstance `json:"cloudMongoDbInstanceList"`
}

type getVpcListResponse struct {
	VpcList []Vpc `json:"vpcList"`
}

type getSubnetListResponse struct {
	SubnetList []Subnet `json:"subnetList"`
}

type getNatGatewayInstanceListResponse struct {
	NatGatewayInstanceList []NatGatewayInstance `json:"natGatewayInstanceList"`
}

type getRouteTableListResponse struct {
	RouteTableList []RouteTable `json:"routeTableList"`
}

type getAccessControlGroupListResponse struct {
	AccessControlGroupList []AccessControlGroup `json:"accessControlGroupList"`
}

type getAutoScalingGroupListResponse struct {
	AutoScalingGroupList []AutoScalingGroup `json:"autoScalingGroupList"`
}

type getNksClusterListResponse struct {
	Clusters []NksCluster `json:"clusters"`
}

type getCloudMariaDbInstanceListResponse struct {
	CloudMariaDbInstanceList []CloudMariaDbInstance `json:"cloudMariaDbInstanceList"`
}

type getCloudMysqlInstanceListResponse struct {
	CloudMysqlInstanceList []CloudMysqlInstance `json:"cloudMysqlInstanceList"`
}

type getCloudRedisInstanceListResponse struct {
	CloudRedisInstanceList []CloudRedisInstance `json:"cloudRedisInstanceList"`
}

type getLaunchConfigurationListResponse struct {
	LaunchConfigurationList []LaunchConfiguration `json:"launchConfigurationList"`
}

type getNetworkAclListResponse struct {
	NetworkAclList []NetworkAcl `json:"networkAclList"`
}

type getVpcPeeringInstanceListResponse struct {
	VpcPeeringInstanceList []VpcPeeringInstance `json:"vpcPeeringInstanceList"`
}

type getInitScriptListResponse struct {
	InitScriptList []InitScript `json:"initScriptList"`
}

type getLoginKeyListResponse struct {
	LoginKeyList []LoginKey `json:"loginKeyList"`
}

type getPlacementGroupListResponse struct {
	PlacementGroupList []PlacementGroup `json:"placementGroupList"`
}

type getTargetGroupListResponse struct {
	TargetGroupList []TargetGroup `json:"targetGroupList"`
}

type getBlockStorageSnapshotInstanceListResponse struct {
	BlockStorageSnapshotInstanceList []BlockStorageSnapshotInstance `json:"blockStorageSnapshotInstanceList"`
}

type getNasVolumeSnapshotListResponse struct {
	NasVolumeSnapshotList []NasVolumeSnapshot `json:"nasVolumeSnapshotList"`
}
