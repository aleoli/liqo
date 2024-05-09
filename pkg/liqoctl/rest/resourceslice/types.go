// Copyright 2019-2024 The Liqo Authors
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

package resourceslice

import (
	"github.com/liqotech/liqo/pkg/liqoctl/rest"
	tenantnamespace "github.com/liqotech/liqo/pkg/tenantNamespace"
)

// Options encapsulates the arguments of the resourceslice command.
type Options struct {
	createOptions *rest.CreateOptions

	namespaceManager tenantnamespace.Manager

	remoteClusterID string
	class           string

	cpu    string
	memory string
	pods   string

	disableVirtualNodeCreation bool
}

var _ rest.API = &Options{}

// ResourceSlice returns the rest API for the resourceslice command.
func ResourceSlice() rest.API {
	return &Options{}
}

// APIOptions returns the APIOptions for the identity API.
func (o *Options) APIOptions() *rest.APIOptions {
	return &rest.APIOptions{
		EnableCreate: true,
	}
}