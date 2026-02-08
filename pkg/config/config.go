package config

import (
	"encoding/json"
	"os"
)

type ResourceFilter struct {
	Enabled *bool    `json:"enabled,omitempty"`
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

// IsEnabled returns whether this resource filter is enabled.
// nil (not specified) or true → enabled; false → disabled.
func (f *ResourceFilter) IsEnabled() bool {
	return f.Enabled == nil || *f.Enabled
}

// Match checks if a resource name or ID matches the filter.
// Rules:
// 1. If Include is not empty, the resource MUST be in the Include list to be considered.
// 2. If Exclude is not empty, the resource MUST NOT be in the Exclude list.
func (f *ResourceFilter) Match(name, id string) bool {
	if !f.IsEnabled() {
		return false
	}

	// If Include list is provided, default to false (reject), unless matched.
	// If Include list is empty, default to true (allow all).
	allowed := true
	if len(f.Include) > 0 {
		allowed = false
		for _, s := range f.Include {
			if s == name || s == id {
				allowed = true
				break
			}
		}
	}

	if !allowed {
		return false
	}

	// Check Exclude list
	if len(f.Exclude) > 0 {
		for _, s := range f.Exclude {
			if s == name || s == id {
				return false
			}
		}
	}

	return true
}

type Config struct {
	Servers               ResourceFilter `json:"servers"`
	BlockStorages         ResourceFilter `json:"block_storages"`
	BlockStorageSnapshots ResourceFilter `json:"block_storage_snapshots"`
	PublicIps             ResourceFilter `json:"public_ips"`
	NasVolumes            ResourceFilter `json:"nas_volumes"`
	NasVolumeSnapshots    ResourceFilter `json:"nas_volume_snapshots"`
	LoadBalancers         ResourceFilter `json:"load_balancers"`
	TargetGroups          ResourceFilter `json:"target_groups"`
	CloudDBs              ResourceFilter `json:"cloud_dbs"`
	CloudPostgresqls      ResourceFilter `json:"cloud_postgresqls"`
	CloudMongoDBs         ResourceFilter `json:"cloud_mongodbs"`
	CloudMariaDBs         ResourceFilter `json:"cloud_mariadbs"`
	CloudMySQLs           ResourceFilter `json:"cloud_mysqls"`
	CloudRedises          ResourceFilter `json:"cloud_redises"`
	Vpcs                  ResourceFilter `json:"vpcs"`
	Subnets               ResourceFilter `json:"subnets"`
	NatGateways           ResourceFilter `json:"nat_gateways"`
	VpcPeerings           ResourceFilter `json:"vpc_peerings"`
	NetworkAcls           ResourceFilter `json:"network_acls"`
	AccessControlGroups   ResourceFilter `json:"access_control_groups"`
	AutoScalingGroups     ResourceFilter `json:"auto_scaling_groups"`
	LaunchConfigurations  ResourceFilter `json:"launch_configurations"`
	NksClusters           ResourceFilter `json:"nks_clusters"`
	InitScripts           ResourceFilter `json:"init_scripts"`
	LoginKeys             ResourceFilter `json:"login_keys"`
	PlacementGroups       ResourceFilter `json:"placement_groups"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
