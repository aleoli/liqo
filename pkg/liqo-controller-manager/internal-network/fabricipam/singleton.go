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

package fabricipam

import (
	"context"
	"sync"

	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/liqotech/liqo/pkg/consts"
	"github.com/liqotech/liqo/pkg/utils/getters"
)

var (
	fabricIPAM *IPAM
	ready      bool
	mutex      sync.Mutex
)

// Get retrieve and init the IPAM singleton.
func Get(ctx context.Context, cl client.Client) (*IPAM, error) {
	mutex.Lock()
	defer mutex.Unlock()
	if ready {
		return fabricIPAM, nil
	}

	network, err := getters.GetUniqueNetworkByLabel(ctx, cl, labels.SelectorFromSet(map[string]string{
		consts.NetworkTypeLabelKey: string(consts.NetworkTypeInternalCIDR),
	}))
	if err != nil {
		return nil, err
	}
	if network.Status.CIDR.String() == "" {
		return nil, nil
	}

	fabricIPAM, err = newIPAM(network.Status.CIDR.String())
	if err != nil {
		return nil, err
	}

	if err := Init(ctx, cl, fabricIPAM); err != nil {
		return nil, err
	}

	ready = true

	return fabricIPAM, nil
}
