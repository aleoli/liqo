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

package move

const (
	// LiqoctlMoveShortHelp contains the short help string for liqoctl move command.
	LiqoctlMoveShortHelp = "Move liqo volumes to other clusters"
	// LiqoctlMoveLongHelp contains the Long help string for liqoctl move command.
	LiqoctlMoveLongHelp = `Move liqo volumes to other clusters`
	// UseCommand contains the verb of the move command.
	UseCommand = "move"
	// VolumeResourceName contains the name of the resource moved in liqoctl move.
	VolumeResourceName = "volume"
	// NamespaceFlagName contains the namespace where the volume is stored.
	NamespaceFlagName = "namespace"
	// TargetNodeFlagName contains the node where the volume is moved to.
	TargetNodeFlagName = "node"
	// ResticPasswordFlagName contains the restic password to be used to for the restic repository.
	ResticPasswordFlagName = "restic-password"

	liqoStorageNamespace = "liqo-storage"
	resticRegistry       = "restic-registry"
	resticServerImage    = "restic/rest-server:0.10.0"
	resticImage          = "restic/restic:0.12.1"
)
