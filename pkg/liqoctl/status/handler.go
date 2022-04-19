// Copyright 2019-2022 The Liqo Authors
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

package status

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/liqotech/liqo/pkg/liqoctl/common"
)

// Args flags of the status command.
type Args struct {
	ShowOnlyLocal     bool
	Namespace         string
	ClusterNameFilter *[]string
	ClusterIDFilter   *[]string
	CheckerSelector   *[]string
}

// Handler implements the logic of the status command.
func (a *Args) Handler(ctx context.Context) error {
	restConfig, err := common.GetLiqoctlRestConf()
	if err != nil {
		return err
	}

	client, err := client.New(restConfig, client.Options{})
	if err != nil {
		return err
	}

	collector := newK8sStatusCollector(ctx, client, *a)

	return collector.collectStatus(ctx)
}
