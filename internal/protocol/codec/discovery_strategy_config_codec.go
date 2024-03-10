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

func EncodeDiscoveryStrategyConfig(clientMessage *proto.ClientMessage, discoveryStrategyConfig types.DiscoveryStrategyConfig) {
	clientMessage.AddFrame(proto.BeginFrame.Copy())

	EncodeString(clientMessage, discoveryStrategyConfig.ClassName)
	EncodeMapForStringAndData(clientMessage, discoveryStrategyConfig.Properties)

	clientMessage.AddFrame(proto.EndFrame.Copy())
}

func DecodeDiscoveryStrategyConfig(frameIterator *proto.ForwardFrameIterator) types.DiscoveryStrategyConfig {
	// begin frame
	frameIterator.Next()

	className := DecodeString(frameIterator)
	properties := DecodeMapForStringAndData(frameIterator)
	FastForwardToEndFrame(frameIterator)

	return types.DiscoveryStrategyConfig{
		ClassName:  className,
		Properties: properties,
	}
}
