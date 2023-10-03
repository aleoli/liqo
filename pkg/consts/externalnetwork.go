// Copyright 2019-2023 The Liqo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package consts

const (
	// WgServerNameLabel is the label used to indicate the name of the WireGuard server.
	WgServerNameLabel = "liqo.io/wg-server-name"
	// WgClientNameLabel is the label used to indicate the name of the WireGuard client.
	WgClientNameLabel = "liqo.io/wg-client-name"
	// ExternalNetworkLabel is the label added to all components that belong to the external network.
	ExternalNetworkLabel = "liqo.io/external-network"
	// ExternalNetworkLabelValue is the value of the label added to components that belong to the external network.
	ExternalNetworkLabelValue = "true"
)