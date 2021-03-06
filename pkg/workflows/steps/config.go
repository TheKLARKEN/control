package steps

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/supergiant/control/pkg/sgerrors"

	"github.com/Azure/go-autorest/autorest"
	"github.com/pkg/errors"

	"github.com/supergiant/control/pkg/clouds"
	"github.com/supergiant/control/pkg/model"
	"github.com/supergiant/control/pkg/profile"
	"github.com/supergiant/control/pkg/runner"
	"github.com/supergiant/control/pkg/storage"
)

type CertificatesConfig struct {
	ServicesCIDR string `json:"servicesCIDR"`
	PublicIP     string `json:"publicIp"`
	PrivateIP    string `json:"privateIp"`

	MasterHost string `json:"masterHost"`
	MasterPort string `json:"masterPort"`
	NodeName   string `json:"nodeName"`

	IsMaster bool `json:"isMaster"`
	// TODO: this shouldn't be a part of SANs
	// https://kubernetes.io/docs/setup/certificates/#all-certificates
	KubernetesSvcIP string `json:"kubernetesSvcIp"`

	StaticAuth profile.StaticAuth `json:"staticAuth"`

	// DEPRECATED: it's a part of staticAuth
	Username string
	// DEPRECATED: it's a part of staticAuth
	Password string

	AdminCert string `json:"adminCert"`
	AdminKey  string `json:"adminKey"`

	ParenCert []byte `json:"parenCert"`
	CACert    string `json:"caCert"`
	CAKey     string `json:"caKey"`
}

type DOConfig struct {
	Name string `json:"name" valid:"required"`
	// These come from UI select
	Region string `json:"region" valid:"required"`
	Size   string `json:"size" valid:"required"`
	Image  string `json:"image" valid:"required"`

	// These come from cloud account
	Fingerprint string `json:"fingerprint" valid:"required"`
	AccessToken string `json:"accessToken" valid:"required"`

	ExternalLoadBalancerID string `json:"externalLoadBalancerId"`
	InternalLoadBalancerID string `json:"internalLoadBalancerId"`
}

type ServiceAccount struct {
	// NOTE(stgleb): This comes from cloud account
	Type      string `json:"type"`
	ProjectID string `json:"project_id"`

	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`

	ClientEmail string `json:"client_email"`
	ClientID    string `json:"client_id"`

	AuthURI  string `json:"auth_uri"`
	TokenURI string `json:"token_uri"`

	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`
}

type GCEConfig struct {
	ServiceAccount

	// This comes from profile
	ImageFamily       string `json:"imageFamily"`
	Region            string `json:"region"`
	AvailabilityZone  string `json:"availabilityZone"`
	Size              string `json:"size"`
	InstanceGroupName string `json:"instanceGroupName"`
	InstanceGroupLink string `json:"instanceGroupLink"`

	// Target pool acts as a balancer for external traffic https://cloud.google.com/load-balancing/docs/target-pools
	TargetPoolName string `json:"targetPoolName"`
	TargetPoolLink string `json:"targetPoolLink"`

	ExternalAddressName string `json:"externalAddressName"`

	ExternalIPAddressLink string `json:"externalIpAddressLink"`

	HealthCheckName string `json:"healthCheckName"`

	ForwardingRuleName string `json:"externalForwardingRuleName"`
}

type AzureConfig struct {
	ClientID       string `json:"clientId"`
	ClientSecret   string `json:"clientSecret"`
	TenantID       string `json:"tenantId"`
	SubscriptionID string `json:"subscriptionId"`

	Location string `json:"location"`

	VMSize string `json:"vmSize"`
	// TODO: cidr validation?
	VNetCIDR string `json:"vNetCIDR"`
}

type PacketConfig struct{}

type OSConfig struct{}

type AWSConfig struct {
	KeyID                  string `json:"access_key"`
	Secret                 string `json:"secret_key"`
	Region                 string `json:"region"`
	AvailabilityZone       string `json:"availabilityZone"`
	KeyPairName            string `json:"keyPairName"`
	VPCID                  string `json:"vpcid"`
	VPCCIDR                string `json:"vpccidr"`
	RouteTableID           string `json:"routeTableId"`
	InternetGatewayID      string `json:"internetGatewayId"`
	NodesSecurityGroupID   string `json:"nodesSecurityGroupID"`
	MastersSecurityGroupID string `json:"mastersSecurityGroupID"`
	MastersInstanceProfile string `json:"mastersInstanceProfile"`
	NodesInstanceProfile   string `json:"nodesInstanceProfile"`
	VolumeSize             string `json:"volumeSize"`
	DeviceName             string `json:"deviceName"`
	EbsOptimized           string `json:"ebsOptimized"`
	ImageID                string `json:"image"`
	InstanceType           string `json:"size"`

	ExternalLoadBalancerName string `json:"externalLoadBalancerName"`
	InternalLoadBalancerName string `json:"internalLoadBalancerName"`

	// Map of availability zone to subnet
	Subnets map[string]string `json:"subnets"`
	// Map az to route table association
	RouteTableAssociationIDs map[string]string `json:"routeTableAssociationIds"`
}

type NetworkConfig struct {
	CIDR            string `json:"cidr"`
	NetworkProvider string `json:"networkProvider"`
}

type PostStartConfig struct {
	IsMaster    bool          `json:"isMaster"`
	Provider    clouds.Name   `json:"provider"`
	Host        string        `json:"host"`
	Port        string        `json:"port"`
	Username    string        `json:"username"`
	RBACEnabled bool          `json:"rbacEnabled"`
	Timeout     time.Duration `json:"timeout"`
}

type TillerConfig struct {
	HelmVersion     string `json:"helmVersion"`
	RBACEnabled     bool   `json:"rbacEnabled"`
	OperatingSystem string `json:"operatingSystem"`
	Arch            string `json:"arch"`
}

type DockerConfig struct {
	Version        string `json:"version"`
	ReleaseVersion string `json:"releaseVersion"`
	Arch           string `json:"arch"`
}

type DownloadK8sBinary struct {
	K8SVersion      string `json:"k8sVersion"`
	Arch            string `json:"arch"`
	OperatingSystem string `json:"operatingSystem"`
}

type ClusterCheckConfig struct {
	MachineCount int
}

type PrometheusConfig struct {
	Port        string `json:"port"`
	RBACEnabled bool   `json:"rbacEnabled"`
}

type KubeadmConfig struct {
	K8SVersion       string `json:"K8SVersion"`
	IsMaster         bool   `json:"isMaster"`
	AdvertiseAddress string `json:"advertiseAddress"`
	IsBootstrap      bool   `json:"IsBootstrap"`
	ServiceCIDR      string `json:"serviceCIDR"`
	CIDR             string `json:"cidr"`
	Token            string `json:"token"`
	Provider         string `json:"provider"`

	MasterPrivateIP string `json:"masterPrivateIp"`
	InternalDNSName string `json:"internalDNSName"`
	ExternalDNSName string `json:"externalDNSName"`
}

type DrainConfig struct {
	PrivateIP string `json:"privateIp"`
}

type Map struct {
	internal map[string]*model.Machine
}

func (m *Map) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &m.internal)
}

func (m *Map) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.internal)
}

func NewMap(m map[string]*model.Machine) Map {
	return Map{
		internal: m,
	}
}

type Config struct {
	Kube model.Kube `json:"kube"`

	TaskID                 string
	Provider               clouds.Name  `json:"provider"`
	IsMaster               bool         `json:"isMaster"`
	IsBootstrap            bool         `json:"IsBootstrap"`
	BootstrapToken         string       `json:"bootstrapToken"`
	ClusterID              string       `json:"clusterId"`
	ClusterName            string       `json:"clusterName"`
	LogBootstrapPrivateKey bool         `json:"logBootstrapPrivateKey"`
	DigitalOceanConfig     DOConfig     `json:"digitalOceanConfig"`
	AWSConfig              AWSConfig    `json:"awsConfig"`
	GCEConfig              GCEConfig    `json:"gceConfig"`
	AzureConfig            AzureConfig  `json:"azureConfig"`
	OSConfig               OSConfig     `json:"osConfig"`
	PacketConfig           PacketConfig `json:"packetConfig"`

	DockerConfig       DockerConfig       `json:"dockerConfig"`
	DownloadK8sBinary  DownloadK8sBinary  `json:"downloadK8sBinary"`
	CertificatesConfig CertificatesConfig `json:"certificatesConfig"`
	NetworkConfig      NetworkConfig      `json:"networkConfig"`
	PostStartConfig    PostStartConfig    `json:"postStartConfig"`
	TillerConfig       TillerConfig       `json:"tillerConfig"`
	PrometheusConfig   PrometheusConfig   `json:"prometheusConfig"`
	DrainConfig        DrainConfig        `json:"drainConfig"`
	KubeadmConfig      KubeadmConfig      `json:"kubeadmConfig"`

	ExternalDNSName string `json:"externalDnsName"`
	InternalDNSName string `json:"internalDnsName"`

	ClusterCheckConfig ClusterCheckConfig `json:"clusterCheckConfig"`

	Node             model.Machine `json:"node"`
	CloudAccountID   string        `json:"cloudAccountId" valid:"required, length(1|32)"`
	CloudAccountName string        `json:"cloudAccountName" valid:"required, length(1|32)"`
	Timeout          time.Duration `json:"timeout"`
	Runner           runner.Runner `json:"-"`

	repository storage.Interface `json:"-"`

	m1      sync.RWMutex
	Masters Map `json:"masters"`

	m2    sync.RWMutex
	Nodes Map `json:"nodes"`

	authorizerMux  sync.RWMutex
	azureAthorizer autorest.Authorizer

	nodeChan      chan model.Machine
	kubeStateChan chan model.KubeState
	configChan    chan *Config
}

// NewConfig builds instance of config for provisioning
func NewConfig(clusterName, cloudAccountName string, profile profile.Profile) (*Config, error) {
	return &Config{
		Kube: model.Kube{
			SSHConfig: model.SSHConfig{
				Port:      "22",
				User:      "root",
				Timeout:   30,
				PublicKey: profile.PublicKey,
			},
			ExposedAddresses: profile.ExposedAddresses,
		},
		Provider:    profile.Provider,
		ClusterName: clusterName,
		DigitalOceanConfig: DOConfig{
			Region: profile.Region,
		},
		LogBootstrapPrivateKey: profile.LogBootstrapPrivateKey,
		AWSConfig: AWSConfig{
			Region:                 profile.Region,
			AvailabilityZone:       profile.CloudSpecificSettings[clouds.AwsAZ],
			VPCCIDR:                profile.CloudSpecificSettings[clouds.AwsVpcCIDR],
			VPCID:                  profile.CloudSpecificSettings[clouds.AwsVpcID],
			KeyPairName:            profile.CloudSpecificSettings[clouds.AwsKeyPairName],
			MastersSecurityGroupID: profile.CloudSpecificSettings[clouds.AwsMastersSecGroupID],
			NodesSecurityGroupID:   profile.CloudSpecificSettings[clouds.AwsNodesSecgroupID],
			// TODO(stgleb): Passs this from UI or figure out any better way
			DeviceName: "/dev/sda1",
		},
		GCEConfig: GCEConfig{
			AvailabilityZone: profile.Zone,
			ImageFamily:      "ubuntu-1604-lts",
			Region:           profile.Region,
		},
		AzureConfig: AzureConfig{
			Location: profile.Region,
			VNetCIDR: profile.CloudSpecificSettings[clouds.AzureVNetCIDR],
		},
		OSConfig:     OSConfig{},
		PacketConfig: PacketConfig{},

		DockerConfig: DockerConfig{
			Version:        profile.DockerVersion,
			ReleaseVersion: profile.UbuntuVersion,
			Arch:           profile.Arch,
		},
		DownloadK8sBinary: DownloadK8sBinary{
			K8SVersion:      profile.K8SVersion,
			Arch:            profile.Arch,
			OperatingSystem: profile.OperatingSystem,
		},
		CertificatesConfig: CertificatesConfig{
			ServicesCIDR: profile.K8SServicesCIDR,
			Username:     profile.User,
			Password:     profile.Password,
			StaticAuth:   profile.StaticAuth,
		},
		NetworkConfig: NetworkConfig{
			CIDR:            profile.CIDR,
			NetworkProvider: profile.NetworkProvider,
		},
		PostStartConfig: PostStartConfig{
			Host:        "localhost",
			Port:        "8080",
			Username:    profile.User,
			RBACEnabled: profile.RBACEnabled,
			Timeout:     time.Minute * 30,
			Provider:    profile.Provider,
		},
		TillerConfig: TillerConfig{
			HelmVersion:     profile.HelmVersion,
			OperatingSystem: profile.OperatingSystem,
			Arch:            profile.Arch,
			RBACEnabled:     profile.RBACEnabled,
		},
		ClusterCheckConfig: ClusterCheckConfig{
			MachineCount: len(profile.NodesProfiles) + len(profile.MasterProfiles),
		},
		PrometheusConfig: PrometheusConfig{
			Port:        "30900",
			RBACEnabled: profile.RBACEnabled,
		},
		KubeadmConfig: KubeadmConfig{
			K8SVersion:  profile.K8SVersion,
			IsBootstrap: true,
			CIDR:        profile.CIDR,
			ServiceCIDR: profile.K8SServicesCIDR,
		},

		Masters: Map{
			internal: make(map[string]*model.Machine, len(profile.MasterProfiles)),
		},
		Nodes: Map{
			internal: make(map[string]*model.Machine, len(profile.NodesProfiles)),
		},
		Timeout:          time.Minute * 30,
		CloudAccountName: cloudAccountName,

		nodeChan:      make(chan model.Machine, len(profile.MasterProfiles)+len(profile.NodesProfiles)),
		kubeStateChan: make(chan model.KubeState, 2),
		configChan:    make(chan *Config),
	}, nil
}

func NewConfigFromKube(profile *profile.Profile, k *model.Kube) (*Config, error) {
	if k == nil {
		return nil, errors.Wrapf(sgerrors.ErrNilEntity, "kube must not be nil")
	}

	cfg := &Config{
		ClusterID:      k.ID,
		Provider:       profile.Provider,
		ClusterName:    k.Name,
		BootstrapToken: k.BootstrapToken,
		DigitalOceanConfig: DOConfig{
			Region: profile.Region,
		},
		Kube:                   *k,
		ExternalDNSName:        k.ExternalDNSName,
		InternalDNSName:        k.InternalDNSName,
		LogBootstrapPrivateKey: profile.LogBootstrapPrivateKey,
		AWSConfig: AWSConfig{
			Region:                   profile.Region,
			AvailabilityZone:         k.CloudSpec[clouds.AwsAZ],
			VPCCIDR:                  k.CloudSpec[clouds.AwsVpcCIDR],
			VPCID:                    k.CloudSpec[clouds.AwsVpcID],
			KeyPairName:              k.CloudSpec[clouds.AwsKeyPairName],
			Subnets:                  k.Subnets,
			MastersSecurityGroupID:   k.CloudSpec[clouds.AwsMastersSecGroupID],
			NodesSecurityGroupID:     k.CloudSpec[clouds.AwsNodesSecgroupID],
			ImageID:                  k.CloudSpec[clouds.AwsImageID],
			ExternalLoadBalancerName: k.CloudSpec[clouds.AwsExternalLoadBalancerName],
			InternalLoadBalancerName: k.CloudSpec[clouds.AwsInternalLoadBalancerName],
			// TODO(stgleb): Passs this from UI or figure out any better way
			DeviceName: "/dev/sda1",
		},
		GCEConfig: GCEConfig{
			AvailabilityZone: profile.Zone,
		},
		AzureConfig: AzureConfig{
			Location: profile.Region,
		},
		OSConfig:     OSConfig{},
		PacketConfig: PacketConfig{},

		DockerConfig: DockerConfig{
			Version:        profile.DockerVersion,
			ReleaseVersion: profile.UbuntuVersion,
			Arch:           profile.Arch,
		},
		DownloadK8sBinary: DownloadK8sBinary{
			K8SVersion:      profile.K8SVersion,
			Arch:            profile.Arch,
			OperatingSystem: profile.OperatingSystem,
		},
		CertificatesConfig: CertificatesConfig{
			ServicesCIDR: profile.K8SServicesCIDR,
			Username:     profile.User,
			Password:     profile.Password,
			StaticAuth:   profile.StaticAuth,
			CAKey:        k.Auth.CAKey,
			CACert:       k.Auth.CACert,
			AdminCert:    k.Auth.AdminCert,
			AdminKey:     k.Auth.AdminKey,
		},
		NetworkConfig: NetworkConfig{
			NetworkProvider: profile.NetworkProvider,
			CIDR:            profile.CIDR,
		},

		PostStartConfig: PostStartConfig{
			Host:        "localhost",
			Port:        "8080",
			Username:    profile.User,
			RBACEnabled: profile.RBACEnabled,
			Timeout:     time.Minute * 30,
			Provider:    k.Provider,
		},
		TillerConfig: TillerConfig{
			HelmVersion:     profile.HelmVersion,
			OperatingSystem: profile.OperatingSystem,
			Arch:            profile.Arch,
			RBACEnabled:     profile.RBACEnabled,
		},
		ClusterCheckConfig: ClusterCheckConfig{
			MachineCount: len(profile.NodesProfiles) + len(profile.MasterProfiles),
		},
		PrometheusConfig: PrometheusConfig{
			Port:        "30900",
			RBACEnabled: profile.RBACEnabled,
		},
		KubeadmConfig: KubeadmConfig{
			K8SVersion:  profile.K8SVersion,
			IsBootstrap: true,
			Token:       k.BootstrapToken,
			CIDR:        profile.CIDR,
			ServiceCIDR: profile.K8SServicesCIDR,
		},
		Masters: Map{
			internal: make(map[string]*model.Machine, len(profile.MasterProfiles)),
		},
		Nodes: Map{
			internal: make(map[string]*model.Machine, len(profile.NodesProfiles)),
		},
		Timeout:          time.Minute * 30,
		CloudAccountName: k.AccountName,
		nodeChan:         make(chan model.Machine, len(profile.MasterProfiles)+len(profile.NodesProfiles)),
		kubeStateChan:    make(chan model.KubeState, 5),
		configChan:       make(chan *Config),
	}

	// Restore all masters and workers from kube
	for index := range k.Masters {
		cfg.AddMaster(k.Masters[index])
	}

	for index := range k.Nodes {
		cfg.AddNode(k.Nodes[index])
	}

	cfg.Kube = *k

	cfg.Kube.SSHConfig = model.SSHConfig{
		Port:      "22",
		User:      "root",
		Timeout:   10,
		PublicKey: profile.PublicKey,
	}

	return cfg, nil
}

// AddMaster to map of master, map is used because it is reference and can be shared among
// goroutines that run multiple tasks of cluster deployment
func (c *Config) AddMaster(n *model.Machine) {
	c.m1.Lock()
	defer c.m1.Unlock()
	c.Masters.internal[n.ID] = n
}

// AddNode to map of nodes in cluster
func (c *Config) AddNode(n *model.Machine) {
	c.m2.Lock()
	defer c.m2.Unlock()
	c.Nodes.internal[n.ID] = n
}

// GetMaster returns first master in master map or nil
func (c *Config) GetMaster() *model.Machine {
	// non-blocking fast path for master nodes
	if c.IsMaster && c.Node.State == model.MachineStateActive {
		return &c.Node
	}

	c.m1.RLock()
	defer c.m1.RUnlock()

	if len(c.Masters.internal) == 0 {
		return nil
	}

	for key := range c.Masters.internal {
		// Skip inactive nodes for selecting
		if c.Masters.internal[key] != nil && c.Masters.internal[key].State == model.MachineStateActive {
			return c.Masters.internal[key]
		}
	}

	return nil
}

func (c *Config) GetMasters() map[string]*model.Machine {
	c.m1.RLock()
	defer c.m1.RUnlock()

	m := make(map[string]*model.Machine, len(c.Masters.internal))

	for key := range c.Masters.internal {
		m[c.Masters.internal[key].Name] = c.Masters.internal[key]
	}

	return m
}

func (c *Config) GetNodes() map[string]*model.Machine {
	c.m2.RLock()
	defer c.m2.RUnlock()

	m := make(map[string]*model.Machine, len(c.Nodes.internal))

	for key := range c.Nodes.internal {
		m[c.Nodes.internal[key].Name] = c.Nodes.internal[key]
	}

	return m
}

// GetMaster returns first master in master map or nil
func (c *Config) GetNode() *model.Machine {
	c.m2.RLock()
	defer c.m2.RUnlock()

	if len(c.Nodes.internal) == 0 {
		return nil
	}

	for key := range c.Nodes.internal {
		// Skip inactive nodes for selecting
		if c.Nodes.internal[key] != nil && c.Nodes.internal[key].State == model.MachineStateActive {
			return c.Nodes.internal[key]
		}
	}

	return nil
}

func (c *Config) NodeChan() chan model.Machine {
	return c.nodeChan
}

func (c *Config) KubeStateChan() chan model.KubeState {
	return c.kubeStateChan
}

func (c *Config) ConfigChan() chan *Config {
	return c.configChan
}

func (c *Config) SetNodeChan(nodeChan chan model.Machine) {
	c.nodeChan = nodeChan
}

func (c *Config) SetKubeStateChan(kubeStateChan chan model.KubeState) {
	c.kubeStateChan = kubeStateChan
}

func (c *Config) SetConfigChan(configChan chan *Config) {
	c.configChan = configChan
}

func (c *Config) SetAzureAuthorizer(a autorest.Authorizer) {
	c.authorizerMux.Lock()
	defer c.authorizerMux.Unlock()

	c.azureAthorizer = a
}

func (c *Config) GetAzureAuthorizer() autorest.Authorizer {
	c.authorizerMux.RLock()
	defer c.authorizerMux.RUnlock()

	return c.azureAthorizer
}
