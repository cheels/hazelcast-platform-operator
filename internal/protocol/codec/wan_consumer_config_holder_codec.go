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
	WanConsumerConfigHolderCodecPersistWanReplicatedDataFieldOffset      = 0
	WanConsumerConfigHolderCodecPersistWanReplicatedDataInitialFrameSize = WanConsumerConfigHolderCodecPersistWanReplicatedDataFieldOffset + proto.BooleanSizeInBytes
)

func EncodeWanConsumerConfigHolder(clientMessage *proto.ClientMessage, value interface{}) {
	wanConsumerConfigHolder := value.(types.WanConsumerConfigHolder)

	clientMessage.AddFrame(proto.BeginFrame.Copy())
	initialFrame := proto.NewFrame(make([]byte, WanConsumerConfigHolderCodecPersistWanReplicatedDataInitialFrameSize))
	EncodeBoolean(initialFrame.Content, WanConsumerConfigHolderCodecPersistWanReplicatedDataFieldOffset, wanConsumerConfigHolder.PersistWanReplicatedData)
	clientMessage.AddFrame(initialFrame)

	EncodeNullableForString(clientMessage, wanConsumerConfigHolder.ClassName)
	EncodeNullableForData(clientMessage, wanConsumerConfigHolder.Implementation)
	EncodeMapForStringAndData(clientMessage, wanConsumerConfigHolder.Properties)

	clientMessage.AddFrame(proto.EndFrame.Copy())
}

func DecodeWanConsumerConfigHolder(frameIterator *proto.ForwardFrameIterator) types.WanConsumerConfigHolder {
	// begin frame
	frameIterator.Next()
	initialFrame := frameIterator.Next()
	persistWanReplicatedData := DecodeBoolean(initialFrame.Content, WanConsumerConfigHolderCodecPersistWanReplicatedDataFieldOffset)

	className := DecodeNullableForString(frameIterator)
	implementation := DecodeNullableForData(frameIterator)
	properties := DecodeMapForStringAndData(frameIterator)
	FastForwardToEndFrame(frameIterator)

	return types.WanConsumerConfigHolder{
		PersistWanReplicatedData: persistWanReplicatedData,
		ClassName:                className,
		Implementation:           implementation,
		Properties:               properties,
	}
}
