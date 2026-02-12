package ncp

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// --- Cloud DB APIs ---

func (c *Client) ListCloudDBInstances() ([]CloudDBInstance, error) {
	path := "/getCloudDBInstanceList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VCloudDBBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}

	var resp struct {
		Response getCloudDBInstanceListResponse `json:"getCloudDBInstanceListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.CloudDBInstanceList, nil
}

func (c *Client) DeleteCloudDBInstance(instanceNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("cloudDBInstanceNo", instanceNo)
	path := "/deleteCloudDBServerInstance?" + params.Encode() // Note: API name is ServerInstance but acts on Instance
	body, status, err := c.doRequestWithBase(VCloudDBBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

// --- Cloud DB for PostgreSQL APIs ---

func (c *Client) ListCloudPostgresqlInstances() ([]CloudPostgresqlInstance, error) {
	path := "/getCloudPostgresqlInstanceList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VPostgreSQLBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}

	var resp struct {
		Response getCloudPostgresqlInstanceListResponse `json:"getCloudPostgresqlInstanceListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.CloudPostgresqlInstanceList, nil
}

func (c *Client) DeleteCloudPostgresqlInstance(instanceNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("cloudPostgresqlInstanceNo", instanceNo)
	path := "/deleteCloudPostgresqlInstance?" + params.Encode()
	body, status, err := c.doRequestWithBase(VPostgreSQLBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

// --- Cloud DB for MongoDB APIs ---

func (c *Client) ListCloudMongoDBInstances() ([]CloudMongoDbInstance, error) {
	path := "/getCloudMongoDbInstanceList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VMongoDBBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}

	var resp struct {
		Response getCloudMongoDbInstanceListResponse `json:"getCloudMongoDbInstanceListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.CloudMongoDbInstanceList, nil
}

func (c *Client) DeleteCloudMongoDBInstance(instanceNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("cloudMongoDbInstanceNo", instanceNo)
	path := "/deleteCloudMongoDbInstance?" + params.Encode()
	body, status, err := c.doRequestWithBase(VMongoDBBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

// --- VPC APIs ---

func (c *Client) ListVpcs() ([]Vpc, error) {
	path := "/getVpcList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VVPCBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}

	var resp struct {
		Response getVpcListResponse `json:"getVpcListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.VpcList, nil
}

func (c *Client) DeleteVpc(vpcNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("vpcNo", vpcNo)
	path := "/deleteVpc?" + params.Encode()
	body, status, err := c.doRequestWithBase(VVPCBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

func (c *Client) ListSubnets() ([]Subnet, error) {
	path := "/getSubnetList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VVPCBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}

	var resp struct {
		Response getSubnetListResponse `json:"getSubnetListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.SubnetList, nil
}

func (c *Client) DeleteSubnet(subnetNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("subnetNo", subnetNo)
	path := "/deleteSubnet?" + params.Encode()
	body, status, err := c.doRequestWithBase(VVPCBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

func (c *Client) ListNatGateways() ([]NatGatewayInstance, error) {
	path := "/getNatGatewayInstanceList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VVPCBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}

	var resp struct {
		Response getNatGatewayInstanceListResponse `json:"getNatGatewayInstanceListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.NatGatewayInstanceList, nil
}

func (c *Client) DeleteNatGateway(instanceNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("natGatewayInstanceNo", instanceNo)
	path := "/deleteNatGatewayInstance?" + params.Encode()
	body, status, err := c.doRequestWithBase(VVPCBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

func (c *Client) ListRouteTables() ([]RouteTable, error) {
	path := "/getRouteTableList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VVPCBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}

	var resp struct {
		Response getRouteTableListResponse `json:"getRouteTableListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.RouteTableList, nil
}

// RemoveRoute removes a specific route from a route table.
func (c *Client) RemoveRoute(routeTableNo string, route Route) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("routeTableNo", routeTableNo)
	params.Set("routeList.1.destinationCidrBlock", route.DestinationCidrBlock)
	params.Set("routeList.1.targetTypeCode", route.TargetTypeCode.Code)
	params.Set("routeList.1.targetNo", route.TargetNo)
	params.Set("routeList.1.targetName", route.TargetName)

	path := "/removeRoute?" + params.Encode()
	body, status, err := c.doRequestWithBase(VVPCBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

// --- ACG APIs ---

func (c *Client) ListAccessControlGroups() ([]AccessControlGroup, error) {
	path := "/getAccessControlGroupList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VServerBaseURL, "GET", path, nil) // ACG is in vserver/v2
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}

	var resp struct {
		Response getAccessControlGroupListResponse `json:"getAccessControlGroupListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.AccessControlGroupList, nil
}

func (c *Client) DeleteAccessControlGroup(acgNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("accessControlGroupNo", acgNo)
	path := "/deleteAccessControlGroup?" + params.Encode()
	body, status, err := c.doRequestWithBase(VServerBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

// --- Auto Scaling APIs ---

func (c *Client) ListAutoScalingGroups() ([]AutoScalingGroup, error) {
	path := "/getAutoScalingGroupList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VAutoScalingBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}

	var resp struct {
		Response getAutoScalingGroupListResponse `json:"getAutoScalingGroupListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.AutoScalingGroupList, nil
}

func (c *Client) DeleteAutoScalingGroup(groupNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("autoScalingGroupNo", groupNo)
	path := "/deleteAutoScalingGroup?" + params.Encode()
	body, status, err := c.doRequestWithBase(VAutoScalingBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

// --- NKS (Kubernetes) APIs ---

func (c *Client) ListNksClusters() ([]NksCluster, error) {
	path := "/clusters"
	body, status, err := c.doRequestWithBase(VNKSBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}

	var resp getNksClusterListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Clusters, nil
}

func (c *Client) DeleteNksCluster(uuid string) error {
	path := "/clusters/" + uuid
	body, status, err := c.doRequestWithBase(VNKSBaseURL, "DELETE", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

// --- Cloud DB for MariaDB APIs ---

func (c *Client) ListCloudMariaDbInstances() ([]CloudMariaDbInstance, error) {
	path := "/getCloudMariaDbInstanceList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VMariaDBBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	var resp struct {
		Response getCloudMariaDbInstanceListResponse `json:"getCloudMariaDbInstanceListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.CloudMariaDbInstanceList, nil
}

func (c *Client) DeleteCloudMariaDbInstance(instanceNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("cloudMariaDbInstanceNo", instanceNo)
	path := "/deleteCloudMariaDbInstance?" + params.Encode()
	body, status, err := c.doRequestWithBase(VMariaDBBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

// --- Cloud DB for MySQL (dedicated) APIs ---

func (c *Client) ListCloudMysqlInstances() ([]CloudMysqlInstance, error) {
	path := "/getCloudMysqlInstanceList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VMySQLBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	var resp struct {
		Response getCloudMysqlInstanceListResponse `json:"getCloudMysqlInstanceListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.CloudMysqlInstanceList, nil
}

func (c *Client) DeleteCloudMysqlInstance(instanceNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("cloudMysqlInstanceNo", instanceNo)
	path := "/deleteCloudMysqlInstance?" + params.Encode()
	body, status, err := c.doRequestWithBase(VMySQLBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

// --- Cloud DB for Redis (dedicated) APIs ---

func (c *Client) ListCloudRedisInstances() ([]CloudRedisInstance, error) {
	path := "/getCloudRedisInstanceList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VRedisBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	var resp struct {
		Response getCloudRedisInstanceListResponse `json:"getCloudRedisInstanceListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.CloudRedisInstanceList, nil
}

func (c *Client) DeleteCloudRedisInstance(instanceNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("cloudRedisInstanceNo", instanceNo)
	path := "/deleteCloudRedisInstance?" + params.Encode()
	body, status, err := c.doRequestWithBase(VRedisBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

// --- Launch Configuration APIs ---

func (c *Client) ListLaunchConfigurations() ([]LaunchConfiguration, error) {
	path := "/getLaunchConfigurationList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VAutoScalingBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	var resp struct {
		Response getLaunchConfigurationListResponse `json:"getLaunchConfigurationListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.LaunchConfigurationList, nil
}

func (c *Client) DeleteLaunchConfiguration(launchConfigurationNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("launchConfigurationNo", launchConfigurationNo)
	path := "/deleteLaunchConfiguration?" + params.Encode()
	body, status, err := c.doRequestWithBase(VAutoScalingBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

// --- Network ACL APIs ---

func (c *Client) ListNetworkAcls() ([]NetworkAcl, error) {
	path := "/getNetworkAclList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VVPCBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	var resp struct {
		Response getNetworkAclListResponse `json:"getNetworkAclListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.NetworkAclList, nil
}

func (c *Client) DeleteNetworkAcl(networkAclNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("networkAclNo", networkAclNo)
	path := "/deleteNetworkAcl?" + params.Encode()
	body, status, err := c.doRequestWithBase(VVPCBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

// --- VPC Peering APIs ---

func (c *Client) ListVpcPeeringInstances() ([]VpcPeeringInstance, error) {
	path := "/getVpcPeeringInstanceList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VVPCBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	var resp struct {
		Response getVpcPeeringInstanceListResponse `json:"getVpcPeeringInstanceListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.VpcPeeringInstanceList, nil
}

func (c *Client) DeleteVpcPeeringInstance(vpcPeeringInstanceNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("vpcPeeringInstanceNo", vpcPeeringInstanceNo)
	path := "/deleteVpcPeeringInstance?" + params.Encode()
	body, status, err := c.doRequestWithBase(VVPCBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

// --- Init Script APIs ---

func (c *Client) ListInitScripts() ([]InitScript, error) {
	path := "/getInitScriptList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VServerBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	var resp struct {
		Response getInitScriptListResponse `json:"getInitScriptListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.InitScriptList, nil
}

func (c *Client) DeleteInitScripts(initScriptNos []string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	for i, no := range initScriptNos {
		params.Set(fmt.Sprintf("initScriptNoList.%d", i+1), no)
	}
	path := "/deleteInitScripts?" + params.Encode()
	body, status, err := c.doRequestWithBase(VServerBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

// --- Login Key APIs ---

func (c *Client) ListLoginKeys() ([]LoginKey, error) {
	path := "/getLoginKeyList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VServerBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	var resp struct {
		Response getLoginKeyListResponse `json:"getLoginKeyListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.LoginKeyList, nil
}

func (c *Client) DeleteLoginKey(keyName string) error {
	return c.DeleteLoginKeys([]string{keyName})
}

func (c *Client) DeleteLoginKeys(keyNames []string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	for i, name := range keyNames {
		params.Set(fmt.Sprintf("keyNameList.%d", i+1), name)
	}
	path := "/deleteLoginKeys?" + params.Encode()
	body, status, err := c.doRequestWithBase(VServerBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

// --- Placement Group APIs ---

func (c *Client) ListPlacementGroups() ([]PlacementGroup, error) {
	path := "/getPlacementGroupList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VServerBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	var resp struct {
		Response getPlacementGroupListResponse `json:"getPlacementGroupListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.PlacementGroupList, nil
}

func (c *Client) DeletePlacementGroup(placementGroupNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("placementGroupNo", placementGroupNo)
	path := "/deletePlacementGroup?" + params.Encode()
	body, status, err := c.doRequestWithBase(VServerBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

// --- Target Group (Load Balancer) APIs ---

func (c *Client) ListTargetGroups() ([]TargetGroup, error) {
	path := "/getTargetGroupList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VLBBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	var resp struct {
		Response getTargetGroupListResponse `json:"getTargetGroupListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.TargetGroupList, nil
}

func (c *Client) DeleteTargetGroup(targetGroupNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("targetGroupNo", targetGroupNo)
	path := "/deleteTargetGroup?" + params.Encode()
	body, status, err := c.doRequestWithBase(VLBBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

// --- Block Storage Snapshot APIs ---

func (c *Client) ListBlockStorageSnapshotInstances() ([]BlockStorageSnapshotInstance, error) {
	path := "/getBlockStorageSnapshotInstanceList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VServerBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	var resp struct {
		Response getBlockStorageSnapshotInstanceListResponse `json:"getBlockStorageSnapshotInstanceListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.BlockStorageSnapshotInstanceList, nil
}

func (c *Client) DeleteBlockStorageSnapshotInstances(instanceNos []string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	for i, no := range instanceNos {
		params.Set(fmt.Sprintf("blockStorageSnapshotInstanceNoList.%d", i+1), no)
	}
	path := "/deleteBlockStorageSnapshotInstances?" + params.Encode()
	body, status, err := c.doRequestWithBase(VServerBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}

// --- NAS Volume Snapshot APIs ---

func (c *Client) ListNasVolumeSnapshots() ([]NasVolumeSnapshot, error) {
	path := "/getNasVolumeSnapshotList?responseFormatType=json"
	body, status, err := c.doRequestWithBase(VNASBaseURL, "GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	var resp struct {
		Response getNasVolumeSnapshotListResponse `json:"getNasVolumeSnapshotListResponse"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	return resp.Response.NasVolumeSnapshotList, nil
}

func (c *Client) DeleteNasVolumeSnapshot(snapshotInstanceNo string) error {
	params := url.Values{}
	params.Set("responseFormatType", "json")
	params.Set("nasVolumeSnapshotInstanceNo", snapshotInstanceNo)
	path := "/deleteNasVolumeSnapshot?" + params.Encode()
	body, status, err := c.doRequestWithBase(VNASBaseURL, "GET", path, nil)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP %d - %s", status, string(body))
	}
	return nil
}
