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

func EncodeWanCustomPublisherConfigHolder(clientMessage *proto.ClientMessage, wanCustomPublisherConfigHolder types.WanCustomPublisherConfigHolder) {
	clientMessage.AddFrame(proto.BeginFrame.Copy())

	EncodeNullableForString(clientMessage, wanCustomPublisherConfigHolder.PublisherId)
	EncodeNullableForString(clientMessage, wanCustomPublisherConfigHolder.ClassName)
	EncodeNullableForData(clientMessage, wanCustomPublisherConfigHolder.Implementation)
	EncodeMapForStringAndData(clientMessage, wanCustomPublisherConfigHolder.Properties)

	clientMessage.AddFrame(proto.EndFrame.Copy())
}

func DecodeWanCustomPublisherConfigHolder(frameIterator *proto.ForwardFrameIterator) types.WanCustomPublisherConfigHolder {
	// begin frame
	frameIterator.Next()

	publisherId := DecodeNullableForString(frameIterator)
	className := DecodeNullableForString(frameIterator)
	implementation := DecodeNullableForData(frameIterator)
	properties := DecodeMapForStringAndData(frameIterator)
	FastForwardToEndFrame(frameIterator)

	return types.WanCustomPublisherConfigHolder{
		PublisherId:    publisherId,
		ClassName:      className,
		Implementation: implementation,
		Properties:     properties,
	}
}
