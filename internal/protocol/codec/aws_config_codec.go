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
	AwsConfigCodecEnabledFieldOffset          = 0
	AwsConfigCodecUsePublicIpFieldOffset      = AwsConfigCodecEnabledFieldOffset + proto.BooleanSizeInBytes
	AwsConfigCodecUsePublicIpInitialFrameSize = AwsConfigCodecUsePublicIpFieldOffset + proto.BooleanSizeInBytes
)

func EncodeAwsConfig(clientMessage *proto.ClientMessage, awsConfig types.CloudConfig) {
	clientMessage.AddFrame(proto.BeginFrame.Copy())
	initialFrame := proto.NewFrame(make([]byte, AwsConfigCodecUsePublicIpInitialFrameSize))
	EncodeBoolean(initialFrame.Content, AwsConfigCodecEnabledFieldOffset, awsConfig.Enabled)
	EncodeBoolean(initialFrame.Content, AwsConfigCodecUsePublicIpFieldOffset, awsConfig.UsePublicIp)
	clientMessage.AddFrame(initialFrame)

	EncodeString(clientMessage, awsConfig.Tag)
	EncodeMapForStringAndString(clientMessage, awsConfig.Properties)

	clientMessage.AddFrame(proto.EndFrame.Copy())
}

func DecodeAwsConfig(frameIterator *proto.ForwardFrameIterator) types.CloudConfig {
	// begin frame
	frameIterator.Next()
	initialFrame := frameIterator.Next()
	enabled := DecodeBoolean(initialFrame.Content, AwsConfigCodecEnabledFieldOffset)
	usePublicIp := DecodeBoolean(initialFrame.Content, AwsConfigCodecUsePublicIpFieldOffset)

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
