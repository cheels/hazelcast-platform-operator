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
	AzureConfigCodecEnabledFieldOffset          = 0
	AzureConfigCodecUsePublicIpFieldOffset      = AzureConfigCodecEnabledFieldOffset + proto.BooleanSizeInBytes
	AzureConfigCodecUsePublicIpInitialFrameSize = AzureConfigCodecUsePublicIpFieldOffset + proto.BooleanSizeInBytes
)

func EncodeAzureConfig(clientMessage *proto.ClientMessage, azureConfig types.CloudConfig) {
	clientMessage.AddFrame(proto.BeginFrame.Copy())
	initialFrame := proto.NewFrame(make([]byte, AzureConfigCodecUsePublicIpInitialFrameSize))
	EncodeBoolean(initialFrame.Content, AzureConfigCodecEnabledFieldOffset, azureConfig.Enabled)
	EncodeBoolean(initialFrame.Content, AzureConfigCodecUsePublicIpFieldOffset, azureConfig.UsePublicIp)
	clientMessage.AddFrame(initialFrame)

	EncodeString(clientMessage, azureConfig.Tag)
	EncodeMapForStringAndString(clientMessage, azureConfig.Properties)

	clientMessage.AddFrame(proto.EndFrame.Copy())
}

func DecodeAzureConfig(frameIterator *proto.ForwardFrameIterator) types.CloudConfig {
	// begin frame
	frameIterator.Next()
	initialFrame := frameIterator.Next()
	enabled := DecodeBoolean(initialFrame.Content, AzureConfigCodecEnabledFieldOffset)
	usePublicIp := DecodeBoolean(initialFrame.Content, AzureConfigCodecUsePublicIpFieldOffset)

	tag := DecodeString(frameIterator)
	properties := DecodeMapForStringAndString(frameIterator)
	FastForwardToEndFrame(frameIterator)

	return types.CloudConfig{
		Tag:         tag,
		Enabled:     enabled,
		UsePublicIp: usePublicIp,
		Properties:  properties,
	}
}
