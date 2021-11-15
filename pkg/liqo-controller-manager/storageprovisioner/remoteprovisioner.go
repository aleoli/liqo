// Copyright 2019-2021 The Liqo Authors
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

package storageprovisioner

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1apply "k8s.io/client-go/applyconfigurations/core/v1"
	corev1clients "k8s.io/client-go/kubernetes/typed/core/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"sigs.k8s.io/sig-storage-lib-external-provisioner/v7/controller"

	"github.com/liqotech/liqo/pkg/virtualKubelet/forge"
)

// ProvisionRemotePVC ensures the existence of a remote PVC and returns a virtual PV for that remote storage device.
func ProvisionRemotePVC(ctx context.Context,
	options controller.ProvisionOptions,
	remoteNamespace, remoteRealStorageClass string,
	remotePvcLister corev1listers.PersistentVolumeClaimNamespaceLister,
	remotePvcClient corev1clients.PersistentVolumeClaimInterface) (*corev1.PersistentVolume, controller.ProvisioningState, error) {
	virtualPvc := options.PVC

	mutation := remotePersistentVolumeClaim(virtualPvc, remoteRealStorageClass, remoteNamespace)
	_, err := remotePvcClient.Apply(ctx, mutation, forge.ApplyOptions())
	if err != nil {
		return nil, controller.ProvisioningInBackground, err
	}

	pv := &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: options.PVName,
		},
		Spec: corev1.PersistentVolumeSpec{
			StorageClassName: options.StorageClass.Name,
			AccessModes:      options.PVC.Spec.AccessModes,
			Capacity: corev1.ResourceList{
				corev1.ResourceStorage: options.PVC.Spec.Resources.Requests[corev1.ResourceStorage],
			},
			PersistentVolumeSource: corev1.PersistentVolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/tmp/liqo-placeholder",
				},
			},
			NodeAffinity: &corev1.VolumeNodeAffinity{
				Required: &corev1.NodeSelector{
					NodeSelectorTerms: []corev1.NodeSelectorTerm{
						{
							MatchExpressions: []corev1.NodeSelectorRequirement{
								{
									Key:      corev1.LabelHostname,
									Operator: corev1.NodeSelectorOpIn,
									Values:   []string{options.SelectedNode.Name},
								},
							},
						},
					},
				},
			},
		},
	}

	if options.StorageClass.ReclaimPolicy != nil {
		pv.Spec.PersistentVolumeReclaimPolicy = *options.StorageClass.ReclaimPolicy
	}

	return pv, controller.ProvisioningFinished, nil
}

// remotePersistentVolumeClaim forges the apply patch for the reflected PersistentVolumeClaim, given the local one.
func remotePersistentVolumeClaim(virtualPvc *corev1.PersistentVolumeClaim,
	storageClass, namespace string) *v1apply.PersistentVolumeClaimApplyConfiguration {
	return v1apply.PersistentVolumeClaim(virtualPvc.Name, namespace).
		WithLabels(virtualPvc.GetLabels()).
		WithLabels(forge.ReflectionLabels()).
		WithSpec(remotePersistentVolumeClaimSpec(virtualPvc, storageClass))
}

func remotePersistentVolumeClaimSpec(virtualPvc *corev1.PersistentVolumeClaim,
	storageClass string) *v1apply.PersistentVolumeClaimSpecApplyConfiguration {
	res := v1apply.PersistentVolumeClaimSpec().
		WithAccessModes(virtualPvc.Spec.AccessModes...).
		WithVolumeMode(func() corev1.PersistentVolumeMode {
			if virtualPvc.Spec.VolumeMode != nil {
				return *virtualPvc.Spec.VolumeMode
			}
			return corev1.PersistentVolumeFilesystem
		}()).
		WithResources(persistenVolumeClaimResources(virtualPvc.Spec.Resources))

	if storageClass != "" {
		res.WithStorageClassName(storageClass)
	}

	return res
}

func persistenVolumeClaimResources(resources corev1.ResourceRequirements) *v1apply.ResourceRequirementsApplyConfiguration {
	return v1apply.ResourceRequirements().
		WithLimits(resources.Limits).
		WithRequests(resources.Requests)
}
