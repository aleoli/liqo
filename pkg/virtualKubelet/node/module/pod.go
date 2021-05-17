// Copyright © 2017 The virtual-kubelet authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package module

import (
	"context"

	"github.com/google/go-cmp/cmp"
	pkgerrors "github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
)

const (
	podStatusReasonProviderFailed = "ProviderFailed"
)

func (pc *PodController) createOrUpdatePod(ctx context.Context, pod *corev1.Pod) error {
	// We do this so we don't mutate the pod from the informer cache
	pod = pod.DeepCopy()
	if err := populateEnvironmentVariables(ctx, pod, pc.resourceManager, pc.recorder); err != nil {
		return err
	}

	// We have to use a  different pod that we pass to the provider than the one that gets used in handleProviderError
	// because the provider  may manipulate the pod in a separate goroutine while we were doing work
	podForProvider := pod.DeepCopy()

	// Check if the pod is already known by the provider.
	// NOTE: Some providers return a non-nil error in their GetPod implementation when the pod is not found while some other don't.
	// Hence, we ignore the error and just act upon the pod if it is non-nil (meaning that the provider still knows about the pod).
	if podFromProvider, _ := pc.provider.GetPod(ctx, pod.Namespace, pod.Name); podFromProvider != nil {
		if !podsEqual(podFromProvider, podForProvider) {
			if origErr := pc.provider.UpdatePod(ctx, podForProvider); origErr != nil {
				pc.handleProviderError(ctx, origErr, pod)
				return origErr
			}
		}
	} else {
		if origErr := pc.provider.CreatePod(ctx, podForProvider); origErr != nil {
			pc.handleProviderError(ctx, origErr, pod)
			return origErr
		}
	}
	return nil
}

// podsEqual checks if two pods are equal according to the fields we know that are allowed
// to be modified after startup time.
func podsEqual(pod1, pod2 *corev1.Pod) bool {
	// Pod Update Only Permits update of:
	// - `spec.containers[*].image`
	// - `spec.initContainers[*].image`
	// - `spec.activeDeadlineSeconds`
	// - `spec.tolerations` (only additions to existing tolerations)
	// - `objectmeta.labels`
	// - `objectmeta.annotations`
	// compare the values of the pods to see if the values actually changed

	var (
		containers     = true
		initContainers = true
	)

	if len(pod1.Annotations) == 0 {
		pod1.Annotations = nil
	}
	if len(pod2.Annotations) == 0 {
		pod2.Annotations = nil
	}
	if len(pod1.Labels) == 0 {
		pod1.Labels = nil
	}
	if len(pod2.Labels) == 0 {
		pod2.Labels = nil
	}

	// since the only mutable fields in pods containers and initContainers are the images,
	// we check only them
	for i := range pod1.Spec.Containers {
		if pod2.Spec.Containers[i].Image != pod1.Spec.Containers[i].Image {
			containers = false
			break
		}
	}
	for i := range pod1.Spec.InitContainers {
		if pod2.Spec.InitContainers[i].Image != pod1.Spec.InitContainers[i].Image {
			initContainers = false
			break
		}
	}

	deadline := cmp.Equal(pod1.Spec.ActiveDeadlineSeconds, pod2.Spec.ActiveDeadlineSeconds)
	tolerations := cmp.Equal(pod1.Spec.Tolerations, pod2.Spec.Tolerations)
	labels := cmp.Equal(pod1.ObjectMeta.Labels, pod2.Labels)
	annotations := cmp.Equal(pod1.ObjectMeta.Annotations, pod2.Annotations)

	return containers && initContainers && deadline && tolerations && labels && annotations
}

func (pc *PodController) handleProviderError(ctx context.Context, origErr error, pod *corev1.Pod) {
	// For now this switch case keeps in consideration only the error
	// of type notFound and handles it properly (by deleting the local pod)
	// if in further investigations we notice different error types,
	// the related cases should be added here
	switch errors.ReasonForError(origErr) {
	case metav1.StatusReasonNotFound:
		err := pc.client.Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
		if err != nil {
			err = pkgerrors.Wrapf(err, "setting provider failed, cannot delete local pod %s/%s", pod.Namespace, pod.Name)
			klog.Error(err)
			pc.setProviderFailed(ctx, origErr, pod)
		}

	case metav1.StatusReasonServiceUnavailable, metav1.StatusReasonAlreadyExists:
		// if the pod creation/update process ends with one of these errors we have not to set the pod in the provider
		// failed status
		// StatusReasonAlreadyExists: if the cache is not aligned we can have double pod creation, if that happens
		// simply ignore second creation
		// StatusReasonServiceUnavailable: it can be generated by 2 different things, with both it will retried very soon
		// 1. the cache for the remote cluster is not started yet
		// 2. no secret referencing the same ServiceAccount of the pod is available in the remote cluster, probably
		//    because the secret reflection is not completed yet
		klog.V(4).Info(origErr.Error())

	default:
		pc.setProviderFailed(ctx, origErr, pod)
	}
}

func (pc *PodController) setProviderFailed(ctx context.Context, origErr error, pod *corev1.Pod) {
	podPhase := corev1.PodPending
	if pod.Spec.RestartPolicy == corev1.RestartPolicyNever {
		podPhase = corev1.PodFailed
	}

	pod.ResourceVersion = "" // Blank out resource version to prevent object has been modified error
	pod.Status.Phase = podPhase
	pod.Status.Reason = podStatusReasonProviderFailed
	pod.Status.Message = origErr.Error()

	_, err := pc.client.Pods(pod.Namespace).UpdateStatus(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		klog.Error("Failed to update pod status")
	} else {
		klog.Info("Updated k8s pod status")
	}
}

func (pc *PodController) deletePod(ctx context.Context, pod *corev1.Pod) error {
	err := pc.provider.DeletePod(ctx, pod.DeepCopy())
	if err != nil {
		return err
	}

	klog.Info("Deleted pod from provider")

	return nil
}

func shouldSkipPodStatusUpdate(pod *corev1.Pod) bool {
	return pod.Status.Phase == corev1.PodSucceeded ||
		pod.Status.Phase == corev1.PodFailed ||
		pod.Status.Reason == podStatusReasonProviderFailed
}

func (pc *PodController) updatePodStatus(ctx context.Context, podFromKubernetes *corev1.Pod, key string) error {
	if shouldSkipPodStatusUpdate(podFromKubernetes) {
		return nil
	}

	obj, ok := pc.knownPods.Load(key)
	if !ok {
		// This means there was a race and the pod has been deleted from K8s
		return nil
	}
	kPod := obj.(*knownPod)
	kPod.Lock()
	podFromProvider := kPod.lastPodStatusReceivedFromProvider.DeepCopy()
	kPod.Unlock()
	// We need to do this because the other parts of the pod can be updated elsewhere. Since we're only updating
	// the pod status, and we should be the sole writers of the pod status, we can blind overwrite it. Therefore
	// we need to copy the pod and set ResourceVersion to 0.
	podFromProvider.ResourceVersion = "0"

	if _, err := pc.client.Pods(podFromKubernetes.Namespace).UpdateStatus(context.TODO(), podFromProvider, metav1.UpdateOptions{}); err != nil {
		return pkgerrors.Wrap(err, "error while updating pod status in kubernetes")
	}

	klog.V(4).Infof("Updated pod status in kubernetes\tnew phase:%s\told phase:%s",
		string(podFromProvider.Status.Phase),
		string(podFromKubernetes.Status.Phase))

	return nil
}

// enqueuePodStatusUpdate updates our pod status map, and marks the pod as dirty in the workqueue. The pod must be DeepCopy'd
// prior to enqueuePodStatusUpdate.
func (pc *PodController) enqueuePodStatusUpdate(ctx context.Context, pod *corev1.Pod) {
	if key, err := cache.MetaNamespaceKeyFunc(pod); err != nil {
		klog.Error(err, "Error getting pod meta namespace key")
	} else if obj, ok := pc.knownPods.Load(key); ok {
		kpod := obj.(*knownPod)
		kpod.Lock()
		kpod.lastPodStatusReceivedFromProvider = pod
		kpod.Unlock()
		pc.syncPodStatusFromProvider.Enqueue(key)
	}
}

func (pc *PodController) syncPodStatusFromProviderHandler(ctx context.Context, key string) (retErr error) {
	defer func() {
		if retErr != nil {
			klog.Error("Error processing pod status update")
		}
	}()

	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return pkgerrors.Wrap(err, "error splitting cache key")
	}

	pod, err := pc.podsLister.Pods(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			klog.Error(err, "Skipping pod status update for pod missing in Kubernetes")
			return nil
		}
		return pkgerrors.Wrap(err, "error looking up pod")
	}

	return pc.updatePodStatus(ctx, pod, key)
}

func (pc *PodController) deletePodsFromKubernetesHandler(ctx context.Context, key string) (retErr error) {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)

	if err != nil {
		// Log the error as a warning, but do not requeue the key as it is invalid.
		klog.Info(pkgerrors.Wrapf(err, "invalid resource key: %q", key))
		return nil
	}

	defer func() {
		if retErr == nil {
			if w, ok := pc.provider.(syncWrapper); ok {
				w._deletePodKey(ctx, key)
			}
		}
	}()

	// If the pod has been deleted from API server, we don't need to do anything.
	k8sPod, err := pc.podsLister.Pods(namespace).Get(name)
	if errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}

	if running(&k8sPod.Status) {
		klog.Error("Force deleting pod in running state")
	}

	// We don't check with the provider before doing this delete. At this point, even if an outstanding pod status update
	// was in progress,
	err = pc.client.Pods(namespace).Delete(context.TODO(), name, *metav1.NewDeleteOptions(0))
	if errors.IsNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}
