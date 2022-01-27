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

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/liqotech/liqo/pkg/liqoctl/common"
	"github.com/liqotech/liqo/pkg/utils"
)

// Args encapsulates arguments required to move a resource.
type Args struct {
	VolumeName string
	Namespace  string
	TargetNode string

	ResticPassword string
}

// HandleMoveVolumeCommand handles the move volume command,
// configuring all the resources required to move a liqo volume.
func HandleMoveVolumeCommand(ctx context.Context, t *Args) error {
	restConfig, err := common.GetLiqoctlRestConf()
	if err != nil {
		return err
	}

	fmt.Println("* Initializing... 🔌 ")
	k8sClient, err := client.New(restConfig, client.Options{})
	if err != nil {
		return err
	}

	if t.ResticPassword == "" {
		t.ResticPassword = utils.RandomString(16)
	}

	fmt.Println("* Processing Volume Moving... 💾 ")
	return processMoveVolume(ctx, t, k8sClient)
}

func processMoveVolume(ctx context.Context, t *Args, k8sClient client.Client) error {
	var pvc corev1.PersistentVolumeClaim
	if err := k8sClient.Get(ctx, client.ObjectKey{Namespace: t.Namespace, Name: t.VolumeName}, &pvc); err != nil {
		return err
	}

	mounter, err := getMounter(ctx, k8sClient, &pvc)
	if err != nil {
		return err
	}
	if mounter != nil {
		return fmt.Errorf(
			"the volume to move (%s/%s) must not to be mounted by any pod, but found mounter pod %s/%s",
			t.Namespace, t.VolumeName, mounter.Namespace, mounter.Name)
	}

	var targetNode corev1.Node
	if err := k8sClient.Get(ctx, client.ObjectKey{Name: t.TargetNode}, &targetNode); err != nil {
		return err
	}
	targetIsLocal := !utils.IsVirtualNode(&targetNode)

	originIsLocal, originNode, err := isLocalVolume(ctx, k8sClient, &pvc)
	if err != nil {
		return err
	}

	if err = offloadLiqoStorageNamespace(ctx, k8sClient, originNode, &targetNode); err != nil {
		return err
	}

	defer func() {
		if err := repatriateLiqoStorageNamespace(ctx, k8sClient); err != nil {
			fmt.Printf("Error while repatriating liqo-storage namespace: %v", err)
		}
	}()

	if err := ensureResticRepository(ctx, k8sClient, &pvc); err != nil {
		return err
	}

	defer func() {
		if err := deleteResticRepository(ctx, k8sClient); err != nil {
			fmt.Printf("Error deleting restic repository: %v\n", err)
		}
	}()

	if err := waitForResticRepository(ctx, k8sClient); err != nil {
		return err
	}

	fmt.Print("* Taking a snapshot... 🌅 \n")
	originResticRepositoryURL, err := getResticRepositoryURL(ctx, k8sClient, originIsLocal)
	if err != nil {
		return err
	}
	if err = takeSnapshot(ctx, k8sClient, &pvc,
		originResticRepositoryURL, t.ResticPassword); err != nil {
		return err
	}

	fmt.Print("* Moving the volume... 🚚 \n")
	newPvc, err := recreatePvc(ctx, k8sClient, &pvc)
	if err != nil {
		return err
	}

	targetResticRepositoryURL, err := getResticRepositoryURL(ctx, k8sClient, targetIsLocal)
	if err != nil {
		return err
	}
	if err = restoreSnapshot(ctx, k8sClient,
		&pvc, newPvc, t.TargetNode,
		targetResticRepositoryURL, t.ResticPassword); err != nil {
		return err
	}
	fmt.Print("* Restore completed... 🚀 \n")

	return nil
}

func getResticRepositoryURL(ctx context.Context, cl client.Client, isLocal bool) (string, error) {
	var namespace string
	if isLocal {
		namespace = liqoStorageNamespace
	} else {
		var err error
		namespace, err = getRemoteStorageNamespaceName(ctx, cl, nil)
		if err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("rest:http://%s.%s.svc.cluster.local:8000/", resticRegistry, namespace), nil
}
