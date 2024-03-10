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
	DynamicConfigAddWanReplicationConfigCodecRequestMessageType  = int32(0x1B1200)
	DynamicConfigAddWanReplicationConfigCodecResponseMessageType = int32(0x1B1201)

	DynamicConfigAddWanReplicationConfigCodecRequestInitialFrameSize = proto.PartitionIDOffset + proto.IntSizeInBytes
)

// Adds a WAN replication configuration.

func EncodeDynamicConfigAddWanReplicationConfigRequest(name string, consumerConfig types.WanConsumerConfigHolder, customPublisherConfigs []types.WanCustomPublisherConfigHolder, batchPublisherConfigs []types.WanBatchPublisherConfigHolder) *proto.ClientMessage {
	clientMessage := proto.NewClientMessageForEncode()
	clientMessage.SetRetryable(false)

	initialFrame := proto.NewFrameWith(make([]byte, DynamicConfigAddWanReplicationConfigCodecRequestInitialFrameSize), proto.UnfragmentedMessage)
	clientMessage.AddFrame(initialFrame)
	clientMessage.SetMessageType(DynamicConfigAddWanReplicationConfigCodecRequestMessageType)
	clientMessage.SetPartitionId(-1)

	EncodeString(clientMessage, name)
	EncodeNullable(clientMessage, consumerConfig, EncodeWanConsumerConfigHolder)
	EncodeListMultiFrameForWanCustomPublisherConfigHolder(clientMessage, customPublisherConfigs)
	EncodeListMultiFrameForWanBatchPublisherConfigHolder(clientMessage, batchPublisherConfigs)

	return clientMessage
}
