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

package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/liqotech/liqo/pkg/liqoctl/status"
)

func newStatusCommand(ctx context.Context) *cobra.Command {
	var params = status.Args{}

	cmd := &cobra.Command{
		Use:           status.UseCommand,
		Short:         status.ShortHelp,
		Long:          status.LongHelp,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return params.Handler(ctx)
		},
	}
	cmd.Flags().StringVarP(&params.Namespace, status.Namespace, "n", "liqo", "Namespace Liqo is running in")
	cmd.Flags().BoolVarP(&params.ShowOnlyLocal, status.ShowOnlyLocal, "l", false, "Shows only local cluster information")
	params.ClusterNameFilter = cmd.Flags().StringSliceP(status.ClusterNameFilter, "N", []string{},
		"show info about clusters specified by name, you can specify more than one cluster separating names with ',' character")
	params.ClusterIDFilter = cmd.Flags().StringSliceP(status.ClusterIDFilter, "I", []string{},
		"show info about clusters specified by ID, you can specify more than one cluster separating IDs with ',' character")
	return cmd
}
