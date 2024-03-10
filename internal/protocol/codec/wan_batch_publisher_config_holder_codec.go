/*
* Copyright (c) 2008-2024, Hazelcast, Inc. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License")
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package codec

import (
	proto "github.com/hazelcast/hazelcast-go-client"
	"github.com/hazelcast/hazelcast-platform-operator/internal/protocol/types"
)

const (
	WanBatchPublisherConfigHolderCodecSnapshotEnabledFieldOffset           = 0
	WanBatchPublisherConfigHolderCodecInitialPublisherStateFieldOffset     = WanBatchPublisherConfigHolderCodecSnapshotEnabledFieldOffset + proto.BooleanSizeInBytes
	WanBatchPublisherConfigHolderCodecQueueCapacityFieldOffset             = WanBatchPublisherConfigHolderCodecInitialPublisherStateFieldOffset + proto.ByteSizeInBytes
	WanBatchPublisherConfigHolderCodecBatchSizeFieldOffset                 = WanBatchPublisherConfigHolderCodecQueueCapacityFieldOffset + proto.IntSizeInBytes
	WanBatchPublisherConfigHolderCodecBatchMaxDelayMillisFieldOffset       = WanBatchPublisherConfigHolderCodecBatchSizeFieldOffset + proto.IntSizeInBytes
	WanBatchPublisherConfigHolderCodecResponseTimeoutMillisFieldOffset     = WanBatchPublisherConfigHolderCodecBatchMaxDelayMillisFieldOffset + proto.IntSizeInBytes
	WanBatchPublisherConfigHolderCodecQueueFullBehaviorFieldOffset         = WanBatchPublisherConfigHolderCodecResponseTimeoutMillisFieldOffset + proto.IntSizeInBytes
	WanBatchPublisherConfigHolderCodecAcknowledgeTypeFieldOffset           = WanBatchPublisherConfigHolderCodecQueueFullBehaviorFieldOffset + proto.IntSizeInBytes
	WanBatchPublisherConfigHolderCodecDiscoveryPeriodSecondsFieldOffset    = WanBatchPublisherConfigHolderCodecAcknowledgeTypeFieldOffset + proto.IntSizeInBytes
	WanBatchPublisherConfigHolderCodecMaxTargetEndpointsFieldOffset        = WanBatchPublisherConfigHolderCodecDiscoveryPeriodSecondsFieldOffset + proto.IntSizeInBytes
	WanBatchPublisherConfigHolderCodecMaxConcurrentInvocationsFieldOffset  = WanBatchPublisherConfigHolderCodecMaxTargetEndpointsFieldOffset + proto.IntSizeInBytes
	WanBatchPublisherConfigHolderCodecUseEndpointPrivateAddressFieldOffset = WanBatchPublisherConfigHolderCodecMaxConcurrentInvocationsFieldOffset + proto.IntSizeInBytes
	WanBatchPublisherConfigHolderCodecIdleMinParkNsFieldOffset             = WanBatchPublisherConfigHolderCodecUseEndpointPrivateAddressFieldOffset + proto.BooleanSizeInBytes
	WanBatchPublisherConfigHolderCodecIdleMaxParkNsFieldOffset             = WanBatchPublisherConfigHolderCodecIdleMinParkNsFieldOffset + proto.LongSizeInBytes
	WanBatchPublisherConfigHolderCodecIdleMaxParkNsInitialFrameSize        = WanBatchPublisherConfigHolderCodecIdleMaxParkNsFieldOffset + proto.LongSizeInBytes
)

func EncodeWanBatchPublisherConfigHolder(clientMessage *proto.ClientMessage, wanBatchPublisherConfigHolder types.WanBatchPublisherConfigHolder) {
	clientMessage.AddFrame(proto.BeginFrame.Copy())
	initialFrame := proto.NewFrame(make([]byte, WanBatchPublisherConfigHolderCodecIdleMaxParkNsInitialFrameSize))
	EncodeBoolean(initialFrame.Content, WanBatchPublisherConfigHolderCodecSnapshotEnabledFieldOffset, wanBatchPublisherConfigHolder.SnapshotEnabled)
	EncodeByte(initialFrame.Content, WanBatchPublisherConfigHolderCodecInitialPublisherStateFieldOffset, wanBatchPublisherConfigHolder.InitialPublisherState)
	EncodeInt(initialFrame.Content, WanBatchPublisherConfigHolderCodecQueueCapacityFieldOffset, int32(wanBatchPublisherConfigHolder.QueueCapacity))
	EncodeInt(initialFrame.Content, WanBatchPublisherConfigHolderCodecBatchSizeFieldOffset, int32(wanBatchPublisherConfigHolder.BatchSize))
	EncodeInt(initialFrame.Content, WanBatchPublisherConfigHolderCodecBatchMaxDelayMillisFieldOffset, int32(wanBatchPublisherConfigHolder.BatchMaxDelayMillis))
	EncodeInt(initialFrame.Content, WanBatchPublisherConfigHolderCodecResponseTimeoutMillisFieldOffset, int32(wanBatchPublisherConfigHolder.ResponseTimeoutMillis))
	EncodeInt(initialFrame.Content, WanBatchPublisherConfigHolderCodecQueueFullBehaviorFieldOffset, int32(wanBatchPublisherConfigHolder.QueueFullBehavior))
	EncodeInt(initialFrame.Content, WanBatchPublisherConfigHolderCodecAcknowledgeTypeFieldOffset, int32(wanBatchPublisherConfigHolder.AcknowledgeType))
	EncodeInt(initialFrame.Content, WanBatchPublisherConfigHolderCodecDiscoveryPeriodSecondsFieldOffset, int32(wanBatchPublisherConfigHolder.DiscoveryPeriodSeconds))
	EncodeInt(initialFrame.Content, WanBatchPublisherConfigHolderCodecMaxTargetEndpointsFieldOffset, int32(wanBatchPublisherConfigHolder.MaxTargetEndpoints))
	EncodeInt(initialFrame.Content, WanBatchPublisherConfigHolderCodecMaxConcurrentInvocationsFieldOffset, int32(wanBatchPublisherConfigHolder.MaxConcurrentInvocations))
	EncodeBoolean(initialFrame.Content, WanBatchPublisherConfigHolderCodecUseEndpointPrivateAddressFieldOffset, wanBatchPublisherConfigHolder.UseEndpointPrivateAddress)
	EncodeLong(initialFrame.Content, WanBatchPublisherConfigHolderCodecIdleMinParkNsFieldOffset, int64(wanBatchPublisherConfigHolder.IdleMinParkNs))
	EncodeLong(initialFrame.Content, WanBatchPublisherConfigHolderCodecIdleMaxParkNsFieldOffset, int64(wanBatchPublisherConfigHolder.IdleMaxParkNs))
	clientMessage.AddFrame(initialFrame)

	EncodeNullableForString(clientMessage, wanBatchPublisherConfigHolder.PublisherId)
	EncodeNullableForString(clientMessage, wanBatchPublisherConfigHolder.ClassName)
	EncodeNullableForData(clientMessage, wanBatchPublisherConfigHolder.Implementation)
	EncodeMapForStringAndData(clientMessage, wanBatchPublisherConfigHolder.Properties)
	EncodeNullableForString(clientMessage, wanBatchPublisherConfigHolder.ClusterName)
	EncodeString(clientMessage, wanBatchPublisherConfigHolder.TargetEndpoints)
	EncodeAwsConfig(clientMessage, wanBatchPublisherConfigHolder.AwsConfig)
	EncodeGcpConfig(clientMessage, wanBatchPublisherConfigHolder.GcpConfig)
	EncodeAzureConfig(clientMessage, wanBatchPublisherConfigHolder.AzureConfig)
	EncodeKubernetesConfig(clientMessage, wanBatchPublisherConfigHolder.KubernetesConfig)
	EncodeEurekaConfig(clientMessage, wanBatchPublisherConfigHolder.EurekaConfig)
	EncodeDiscoveryConfig(clientMessage, wanBatchPublisherConfigHolder.DiscoveryConfig)
	EncodeWanSyncConfig(clientMessage, wanBatchPublisherConfigHolder.SyncConfig)
	EncodeNullableForString(clientMessage, wanBatchPublisherConfigHolder.Endpoint)

	clientMessage.AddFrame(proto.EndFrame.Copy())
}

func DecodeWanBatchPublisherConfigHolder(frameIterator *proto.ForwardFrameIterator) types.WanBatchPublisherConfigHolder {
	// begin frame
	frameIterator.Next()
	initialFrame := frameIterator.Next()
	snapshotEnabled := DecodeBoolean(initialFrame.Content, WanBatchPublisherConfigHolderCodecSnapshotEnabledFieldOffset)
	initialPublisherState := DecodeByte(initialFrame.Content, WanBatchPublisherConfigHolderCodecInitialPublisherStateFieldOffset)
	queueCapacity := DecodeInt(initialFrame.Content, WanBatchPublisherConfigHolderCodecQueueCapacityFieldOffset)
	batchSize := DecodeInt(initialFrame.Content, WanBatchPublisherConfigHolderCodecBatchSizeFieldOffset)
	batchMaxDelayMillis := DecodeInt(initialFrame.Content, WanBatchPublisherConfigHolderCodecBatchMaxDelayMillisFieldOffset)
	responseTimeoutMillis := DecodeInt(initialFrame.Content, WanBatchPublisherConfigHolderCodecResponseTimeoutMillisFieldOffset)
	queueFullBehavior := DecodeInt(initialFrame.Content, WanBatchPublisherConfigHolderCodecQueueFullBehaviorFieldOffset)
	acknowledgeType := DecodeInt(initialFrame.Content, WanBatchPublisherConfigHolderCodecAcknowledgeTypeFieldOffset)
	discoveryPeriodSeconds := DecodeInt(initialFrame.Content, WanBatchPublisherConfigHolderCodecDiscoveryPeriodSecondsFieldOffset)
	maxTargetEndpoints := DecodeInt(initialFrame.Content, WanBatchPublisherConfigHolderCodecMaxTargetEndpointsFieldOffset)
	maxConcurrentInvocations := DecodeInt(initialFrame.Content, WanBatchPublisherConfigHolderCodecMaxConcurrentInvocationsFieldOffset)
	useEndpointPrivateAddress := DecodeBoolean(initialFrame.Content, WanBatchPublisherConfigHolderCodecUseEndpointPrivateAddressFieldOffset)
	idleMinParkNs := DecodeLong(initialFrame.Content, WanBatchPublisherConfigHolderCodecIdleMinParkNsFieldOffset)
	idleMaxParkNs := DecodeLong(initialFrame.Content, WanBatchPublisherConfigHolderCodecIdleMaxParkNsFieldOffset)

	publisherId := DecodeNullableForString(frameIterator)
	className := DecodeNullableForString(frameIterator)
	implementation := DecodeNullableForData(frameIterator)
	properties := DecodeMapForStringAndData(frameIterator)
	clusterName := DecodeNullableForString(frameIterator)
	targetEndpoints := DecodeString(frameIterator)
	awsConfig := DecodeAwsConfig(frameIterator)
	gcpConfig := DecodeGcpConfig(frameIterator)
	azureConfig := DecodeAzureConfig(frameIterator)
	kubernetesConfig := DecodeKubernetesConfig(frameIterator)
	eurekaConfig := DecodeEurekaConfig(frameIterator)
	discoveryConfig := DecodeDiscoveryConfig(frameIterator)
	syncConfig := DecodeWanSyncConfig(frameIterator)
	endpoint := DecodeNullableForString(frameIterator)
	FastForwardToEndFrame(frameIterator)

	return types.WanBatchPublisherConfigHolder{
		PublisherId:               publisherId,
		ClassName:                 className,
		Implementation:            implementation,
		Properties:                properties,
		ClusterName:               clusterName,
		SnapshotEnabled:           snapshotEnabled,
		InitialPublisherState:     initialPublisherState,
		QueueCapacity:             int(queueCapacity),
		BatchSize:                 int(batchSize),
		BatchMaxDelayMillis:       int(batchMaxDelayMillis),
		ResponseTimeoutMillis:     int(responseTimeoutMillis),
		QueueFullBehavior:         int(queueFullBehavior),
		AcknowledgeType:           int(acknowledgeType),
		DiscoveryPeriodSeconds:    int(discoveryPeriodSeconds),
		MaxTargetEndpoints:        int(maxTargetEndpoints),
		MaxConcurrentInvocations:  int(maxConcurrentInvocations),
		UseEndpointPrivateAddress: useEndpointPrivateAddress,
		IdleMinParkNs:             idleMinParkNs,
		IdleMaxParkNs:             idleMaxParkNs,
		TargetEndpoints:           targetEndpoints,
		AwsConfig:                 awsConfig,
		GcpConfig:                 gcpConfig,
		AzureConfig:               azureConfig,
		KubernetesConfig:          kubernetesConfig,
		EurekaConfig:              eurekaConfig,
		DiscoveryConfig:           discoveryConfig,
		SyncConfig:                syncConfig,
		Endpoint:                  endpoint,
	}
}
