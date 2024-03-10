package types

import proto "github.com/hazelcast/hazelcast-go-client"

type WanConsumerConfigHolder struct {
	PersistWanReplicatedData bool
	ClassName                string
	Implementation           proto.Data
	Properties               map[string]proto.Data
}

type WanCustomPublisherConfigHolder struct {
	PublisherId    string
	ClassName      string
	Implementation proto.Data
	Properties     map[string]proto.Data
}

type WanBatchPublisherConfigHolder struct {
	PublisherId               string
	ClassName                 string
	Implementation            proto.Data
	Properties                map[string]proto.Data
	ClusterName               string
	SnapshotEnabled           bool
	InitialPublisherState     byte
	QueueCapacity             int
	BatchSize                 int
	BatchMaxDelayMillis       int
	ResponseTimeoutMillis     int
	QueueFullBehavior         int
	AcknowledgeType           int
	DiscoveryPeriodSeconds    int
	MaxTargetEndpoints        int
	MaxConcurrentInvocations  int
	UseEndpointPrivateAddress bool
	IdleMinParkNs             int64
	IdleMaxParkNs             int64
	TargetEndpoints           string
	AwsConfig                 CloudConfig
	GcpConfig                 CloudConfig
	AzureConfig               CloudConfig
	KubernetesConfig          CloudConfig
	EurekaConfig              CloudConfig
	DiscoveryConfig           DiscoveryConfig
	SyncConfig                WanSyncConfig
	Endpoint                  string
}

type CloudConfig struct {
	Tag         string
	Enabled     bool
	UsePublicIp bool
	Properties  map[string]string
}

type DiscoveryConfig struct {
	DiscoveryStrategyConfigs []DiscoveryStrategyConfig
	DiscoveryServiceProvider proto.Data
	NodeFilter               proto.Data
	NodeFilterClass          string
}

type DiscoveryStrategyConfig struct {
	ClassName  string
	Properties map[string]proto.Data
}

type WanSyncConfig struct {
	ConsistencyCheckStrategy byte
}
