package ncp

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

// --- Response types for NCP VPC resource APIs ---

type CommonCode struct {
	Code     string `json:"code"`
	CodeName string `json:"codeName"`
}

type ServerInstance struct {
	ServerInstanceNo     string     `json:"serverInstanceNo"`
	ServerName           string     `json:"serverName"`
	ServerInstanceStatus CommonCode `json:"serverInstanceStatus"`
	PublicIp             string     `json:"publicIp"`
	PrivateIp            string     `json:"privateIp"`
	CpuCount             int        `json:"cpuCount"`
	MemorySize           int64      `json:"memorySize"`
}

type BlockStorageInstance struct {
	BlockStorageInstanceNo     string     `json:"blockStorageInstanceNo"`
	BlockStorageName           string     `json:"blockStorageName"`
	BlockStorageInstanceStatus CommonCode `json:"blockStorageInstanceStatus"`
	BlockStorageSize           int64      `json:"blockStorageSize"`
	ServerInstanceNo           string     `json:"serverInstanceNo"`
	BlockStorageType           CommonCode `json:"blockStorageType"`
	BlockStorageDiskDetailType CommonCode `json:"blockStorageDiskDetailType"`
}

type PublicIpInstance struct {
	PublicIpInstanceNo     string     `json:"publicIpInstanceNo"`
	PublicIp               string     `json:"publicIp"`
	PublicIpInstanceStatus CommonCode `json:"publicIpInstanceStatus"`
	ServerInstanceNo       string     `json:"serverInstanceNo"`
	ServerName             string     `json:"serverName"`
}

type NasVolumeInstance struct {
	NasVolumeInstanceNo     string     `json:"nasVolumeInstanceNo"`
	VolumeName              string     `json:"volumeName"`
	NasVolumeInstanceStatus CommonCode `json:"nasVolumeInstanceStatus"`
	VolumeAllotmentProtocol CommonCode `json:"volumeAllotmentProtocolType"`
	VolumeTotalSize         int64      `json:"volumeTotalSize"`
}

type LoadBalancerInstance struct {
	LoadBalancerInstanceNo     string     `json:"loadBalancerInstanceNo"`
	LoadBalancerName           string     `json:"loadBalancerName"`
	LoadBalancerInstanceStatus CommonCode `json:"loadBalancerInstanceStatus"`
	LoadBalancerType           CommonCode `json:"loadBalancerType"`
}

// Generic NCP API response wrapper
type apiResponse struct {
	RequestId     string          `json:"requestId"`
	ReturnCode    int             `json:"returnCode"`
	ReturnMessage string          `json:"returnMessage"`
	TotalRows     int             `json:"totalRows"`
	Content       json.RawMessage `json:"-"`
}

// ResourceSummary holds a summary of resources found for a root account.
type ResourceSummary struct {
	Servers       []ServerInstance
	BlockStorages []BlockStorageInstance
	PublicIps     []PublicIpInstance
	NasVolumes    []NasVolumeInstance
	LoadBalancers []LoadBalancerInstance
}

// TotalCount returns total number of resources.
func (r *ResourceSummary) TotalCount() int {
	return len(r.Servers) + len(r.BlockStorages) + len(r.PublicIps) + len(r.NasVolumes) + len(r.LoadBalancers)
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

	return summary, errs
}

// --- Delete/Terminate APIs ---

// StopServers stops running servers. Servers must be stopped before termination.
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

// TerminateServers terminates (deletes) stopped servers.
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

// DeleteBlockStorages deletes detached block storage instances.
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

// DisassociatePublicIp disassociates a public IP from a server.
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

// DeletePublicIp releases a public IP.
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

// DeleteNasVolume deletes a NAS volume.
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

// DeleteLoadBalancer deletes a load balancer.
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

// CleanupAllResources deletes all resources for a root account.
// Order: LB → Server(stop→terminate) → Block Storage → Public IP → NAS
// Returns success count, fail count, and errors.
func (c *Client) CleanupAllResources(summary *ResourceSummary, logFn func(string)) (int, int) {
	success, fail := 0, 0

	// 1. Delete Load Balancers
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

	// 2. Stop & Terminate Servers
	if len(summary.Servers) > 0 {
		// Stop running servers first
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

		// Terminate all servers
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
			// Wait for servers to terminate before deleting dependent resources
			logFn("    서버 반납 대기 중 (30초)...")
			time.Sleep(30 * time.Second)
		}
	}

	// 3. Delete Block Storages (non-basic disks)
	var storagesToDelete []string
	for _, bs := range summary.BlockStorages {
		// Skip basic (boot) disks - they are deleted with the server
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

	// 4. Release Public IPs
	for _, ip := range summary.PublicIps {
		logFn(fmt.Sprintf("  공인 IP 해제: %s (%s)", ip.PublicIp, ip.PublicIpInstanceNo))
		// Disassociate first if assigned
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

	// 5. Delete NAS Volumes
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

	return success, fail
}
