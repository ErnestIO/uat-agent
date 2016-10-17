/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

type networkEvent struct {
	Type                string   `json:"type"`
	Service             string   `json:"service"`
	NetworkType         string   `json:"network_type"`
	NetworkName         string   `json:"network_name"`
	NetworkNetmask      string   `json:"network_netmask"`
	NetworkStartAddress string   `json:"network_start_address"`
	NetworkEndAddress   string   `json:"network_end_address"`
	NetworkGateway      string   `json:"network_gateway"`
	DNS                 []string `json:"network_dns"`
	RouterName          string   `json:"router_name"`
	RouterType          string   `json:"router_type,omitempty"`
	RouterIP            string   `json:"router_ip"`
	ClientName          string   `json:"client_name,omitempty"`
	DatacenterType      string   `json:"datacenter_type,omitempty"`
	DatacenterName      string   `json:"datacenter_name,omitempty"`
	DatacenterUsername  string   `json:"datacenter_username,omitempty"`
	DatacenterPassword  string   `json:"datacenter_password,omitempty"`
	DatacenterRegion    string   `json:"datacenter_region,omitempty"`
	VCloudURL           string   `json:"vcloud_url"`
}

type instanceEvent struct {
	Service            string   `json:"service_id"`
	Type               string   `json:"type"`
	InstanceName       string   `json:"instance_name"`
	InstanceType       string   `json:"instance_type"`
	Resource           resource `json:"instance_resource"`
	RouterName         string   `json:"router_name"`
	RouterType         string   `json:"router_type"`
	RouterIP           string   `json:"router_ip"`
	ClientName         string   `json:"client_name"`
	DatacenterName     string   `json:"datacenter_name"`
	DatacenterPassword string   `json:"datacenter_password"`
	DatacenterRegion   string   `json:"datacenter_region"`
	DatacenterType     string   `json:"datacenter_type"`
	DatacenterUsername string   `json:"datacenter_username"`
	NetworkName        string   `json:"network_name"`
	VCloudURL          string   `json:"vcloud_url"`
}

type fwrule struct {
	SourceIP        string `json:"source_ip"`
	SourcePort      string `json:"source_port"`
	DestinationIP   string `json:"destination_ip"`
	DestinationPort string `json:"destination_port"`
	Protocol        string `json:"protocol"`
}

type firewallEvent struct {
	Service            string   `json:"service_id"`
	Type               string   `json:"type"`
	Name               string   `json:"firewall_name"`
	ClientID           string   `json:"client_id"`
	ClientName         string   `json:"client_name"`
	Datacenter         string   `json:"datacenter_id"`
	DatacenterName     string   `json:"datacenter_name"`
	DatacenterUsername string   `json:"datacenter_username"`
	DatacenterPassword string   `json:"datacenter_password"`
	DatacenterType     string   `json:"datacenter_type"`
	ExternalNetwork    string   `json:"external_network"`
	VCloudURL          string   `json:"vcloud_url"`
	Router             string   `json:"router_id"`
	RouterType         string   `json:"router_type"`
	RouterName         string   `json:"router_name"`
	RouterIP           string   `json:"router_ip"`
	Created            bool     `json:"created"`
	FirewallID         string   `json:"firewall_id"`
	FirewallType       string   `json:"firewall_type"`
	Rules              []fwrule `json:"firewall_rules"`
}

type ntrule struct {
	Type            string `json:"type"`
	OriginIP        string `json:"origin_ip"`
	OriginPort      string `json:"origin_port"`
	TranslationIP   string `json:"translation_ip"`
	TranslationPort string `json:"translation_port"`
	Protocol        string `json:"protocol"`
	Network         string `json:"network"`
}

type natEvent struct {
	Service            string   `json:"service_id"`
	Type               string   `json:"type"`
	NatName            string   `json:"nat_name"`
	NatRules           []ntrule `json:"nat_rules"`
	RouterName         string   `json:"router_name"`
	RouterType         string   `json:"router_type"`
	RouterIP           string   `json:"router_ip"`
	ClientName         string   `json:"client_name"`
	DatacenterName     string   `json:"datacenter_name"`
	DatacenterUsername string   `json:"datacenter_username"`
	DatacenterPassword string   `json:"datacenter_password"`
	DatacenterRegion   string   `json:"datacenter_region"`
	DatacenterType     string   `json:"datacenter_type"`
	ExternalNetwork    string   `json:"external_network"`
	VCloudURL          string   `json:"vcloud_url"`
	Status             string   `json:"status"`
}

type disk struct {
	ID   int `json:"id"`
	Size int `json:"size"`
}

type resource struct {
	CPU     int    `json:"cpus"`
	RAM     int    `json:"ram"`
	IP      string `json:"ip"`
	Catalog string `json:"reference_catalog"`
	Image   string `json:"reference_image"`
	Disks   []disk `json:"disks"`
}

type routerEvent struct {
	Service            string `json:"service_id"`
	Type               string `json:"type"`
	RouterName         string `json:"router_name"`
	RouterType         string `json:"router_type"`
	ClientName         string `json:"client_name"`
	DatacenterName     string `json:"datacenter_name"`
	DatacenterUsername string `json:"datacenter_username"`
	DatacenterPassword string `json:"datacenter_password"`
	DatacenterRegion   string `json:"datacenter_region"`
	DatacenterType     string `json:"datacenter_type"`
	ExternalNetwork    string `json:"external_network"`
	VCloudURL          string `json:"vcloud_url"`
	VseURL             string `json:"vse_url"`
	Status             string `json:"status"`
}

type report struct {
	Code     int    `json:"return_code"`
	Instance string `json:"instance"`
	StdErr   string `json:"stderr"`
	StdOut   string `json:"stdout"`
}

type serviceOptions struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type executionEvent struct {
	Type             string         `json:"type"`
	Created          bool           `json:"created"`
	Name             string         `json:"execution_name"`
	DatacenterName   string         `json:"datacenter_name"`
	ClientName       string         `json:"client_name"`
	Service          string         `json:"service_id"`
	ServiceName      string         `json:"service_name"`
	ServiceType      string         `json:"service_type"`
	ServiceOptions   serviceOptions `json:"service_options"`
	ExecutionType    string         `json:"execution_type"`
	ExecutionPayload string         `json:"execution_payload"`
	ExecutionTarget  string         `json:"execution_target"`
	ExecutionResults struct {
		Reports []report `json:"reports,omitempty"`
	} `json:"execution_results,omitempty"`
	ExecutionMatchedInstances []string `json:"execution_matched_instances,omitempty"`
	ExecutionStatus           string   `json:"execution_status,omitempty"`
}

type awsInstanceEvent struct {
	Uuid                  string   `json:"_uuid"`
	BatchID               string   `json:"_batch_id"`
	Type                  string   `json:"_type"`
	DatacenterRegion      string   `json:"datacenter_region,omitempty"`
	DatacenterAccessToken string   `json:"datacenter_access_token,omitempty"`
	DatacenterAccessKey   string   `json:"datacenter_access_key,omitempty"`
	DatacenterVpcID       string   `json:"datacenter_vpc_id,omitempty"`
	NetworkAWSID          string   `json:"network_aws_id"`
	SecurityGroupAWSIDs   []string `json:"security_group_aws_ids"`
	InstanceName          string   `json:"instance_name"`
	InstanceImage         string   `json:"instance_image"`
	InstanceType          string   `json:"instance_type"`
	Status                string   `json:"status"`
	ErrorCode             string   `json:"error_code"`
	ErrorMessage          string   `json:"error_message"`
}

type awsNetworkEvent struct {
	Uuid                  string `json:"_uuid"`
	BatchID               string `json:"_batch_id"`
	Type                  string `json:"_type"`
	Service               string `json:"service"`
	DatacenterRegion      string `json:"datacenter_region,omitempty"`
	DatacenterAccessToken string `json:"datacenter_access_token,omitempty"`
	DatacenterAccessKey   string `json:"datacenter_access_key,omitempty"`
	DatacenterVpcID       string `json:"datacenter_vpc_id,omitempty"`
	NetworkType           string `json:"network_type"`
	NetworkSubnet         string `json:"network_subnet"`
	NetworkAWSID          string `json:"network_aws_id"`
	NetworkIsPublic       bool   `json:"network_is_public"`
}

type awsFirewallRule struct {
	IP       string `json:"ip"`
	From     int    `json:"from_port"`
	To       int    `json:"to_port"`
	Protocol string `json:"protocol"`
}

type awsFirewallEvent struct {
	Uuid                  string `json:"_uuid"`
	BatchID               string `json:"_batch_id"`
	Type                  string `json:"_type"`
	DatacenterRegion      string `json:"datacenter_region"`
	DatacenterAccessToken string `json:"datacenter_access_token"`
	DatacenterAccessKey   string `json:"datacenter_access_key"`
	DatacenterVPCID       string `json:"datacenter_vpc_id"`
	SecurityGroupName     string `json:"security_group_name"`
	SecurityGroupRules    struct {
		Ingress []awsFirewallRule `json:"ingress"`
		Egress  []awsFirewallRule `json:"egress"`
	} `json:"security_group_rules"`
	Status       string `json:"status"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

type awsNatEvent struct {
	Uuid                  string   `json:"_uuid"`
	BatchID               string   `json:"_batch_id"`
	Type                  string   `json:"_type"`
	DatacenterRegion      string   `json:"datacenter_region"`
	DatacenterAccessToken string   `json:"datacenter_access_token"`
	DatacenterAccessKey   string   `json:"datacenter_access_key"`
	DatacenterVPCID       string   `json:"datacenter_vpc_id"`
	NatGatewayAWSID       string   `json:"nat_gateway_aws_id"`
	PublicNetwork         string   `json:"public_network"`
	PublicNetworkAWSID    string   `json:"public_network_aws_id"`
	RoutedNetworks        []string `json:"routed_networks"`
	RoutedNetworkAWSIDs   []string `json:"routed_networks_aws_ids"`
	Status                string   `json:"status"`
	ErrorCode             string   `json:"error_code"`
	ErrorMessage          string   `json:"error_message"`
}

type awsELBListener struct {
	FromPort int    `json:"from_port"`
	ToPort   int    `json:"to_port"`
	Protocol string `json:"protocol"`
	SSLCert  string `json:"ssl_cert"`
}

type awsELBEvent struct {
	Uuid                string           `json:"_uuid"`
	BatchID             string           `json:"_batch_id"`
	Type                string           `json:"_type"`
	Name                string           `json:"name"`
	IsPrivate           bool             `json:"is_private"`
	DNSName             string           `json:"dns_name"`
	Listeners           []awsELBListener `json:"listeners"`
	NetworkAWSIDs       []string         `json:"network_aws_ids"`
	Instances           []string         `json:"instances"`
	InstanceNames       []string         `json:"instance_names"`
	InstanceAWSIDs      []string         `json:"instance_aws_ids"`
	SecurityGroups      []string         `json:"security_groups"`
	SecurityGroupAWSIDs []string         `json:"security_group_aws_ids"`
	DatacenterType      string           `json:"datacenter_type,omitempty"`
	DatacenterName      string           `json:"datacenter_name,omitempty"`
	DatacenterRegion    string           `json:"datacenter_region"`
	DatacenterToken     string           `json:"datacenter_token"`
	DatacenterSecret    string           `json:"datacenter_secret"`
	VpcID               string           `json:"vpc_id"`
	Service             string           `json:"service"`
	Status              string           `json:"status"`
	ErrorCode           string           `json:"error_code"`
	ErrorMessage        string           `json:"error_message"`
}
