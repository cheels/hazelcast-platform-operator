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
	KubernetesConfigCodecEnabledFieldOffset          = 0
	KubernetesConfigCodecUsePublicIpFieldOffset      = KubernetesConfigCodecEnabledFieldOffset + proto.BooleanSizeInBytes
	KubernetesConfigCodecUsePublicIpInitialFrameSize = KubernetesConfigCodecUsePublicIpFieldOffset + proto.BooleanSizeInBytes
)

func EncodeKubernetesConfig(clientMessage *proto.ClientMessage, kubernetesConfig types.CloudConfig) {
	clientMessage.AddFrame(proto.BeginFrame.Copy())
	initialFrame := proto.NewFrame(make([]byte, KubernetesConfigCodecUsePublicIpInitialFrameSize))
	EncodeBoolean(initialFrame.Content, KubernetesConfigCodecEnabledFieldOffset, kubernetesConfig.Enabled)
	EncodeBoolean(initialFrame.Content, KubernetesConfigCodecUsePublicIpFieldOffset, kubernetesConfig.UsePublicIp)
	clientMessage.AddFrame(initialFrame)

	EncodeString(clientMessage, kubernetesConfig.Tag)
	EncodeMapForStringAndString(clientMessage, kubernetesConfig.Properties)

	clientMessage.AddFrame(proto.EndFrame.Copy())
}

func DecodeKubernetesConfig(frameIterator *proto.ForwardFrameIterator) types.CloudConfig {
	// begin frame
	frameIterator.Next()
	initialFrame := frameIterator.Next()
	enabled := DecodeBoolean(initialFrame.Content, KubernetesConfigCodecEnabledFieldOffset)
	usePublicIp := DecodeBoolean(initialFrame.Content, KubernetesConfigCodecUsePublicIpFieldOffset)

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
