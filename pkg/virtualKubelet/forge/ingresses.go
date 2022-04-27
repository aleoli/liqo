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

package forge

import (
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	corev1apply "k8s.io/client-go/applyconfigurations/core/v1"
	netv1apply "k8s.io/client-go/applyconfigurations/networking/v1"
)

// RemoteIngress forges the apply patch for the reflected ingress, given the local one.
func RemoteIngress(local *netv1.Ingress, targetNamespace string) *netv1apply.IngressApplyConfiguration {
	return netv1apply.Ingress(local.GetName(), targetNamespace).
		WithLabels(local.GetLabels()).WithLabels(ReflectionLabels()).
		WithAnnotations(local.GetAnnotations()).
		WithSpec(RemoteIngressSpec(local.Spec.DeepCopy()))
}

// RemoteIngressSpec forges the apply patch for the specs of the reflected ingress, given the local ones.
// It expects the local object to be a deepcopy, as it is mutated.
func RemoteIngressSpec(local *netv1.IngressSpec) *netv1apply.IngressSpecApplyConfiguration {
	remote := netv1apply.IngressSpec().
		WithDefaultBackend(RemoteIngressBackend(local.DefaultBackend)).
		WithRules(RemoteIngressRules(local.Rules)...).
		WithTLS(RemoteIngressTLS(local.TLS)...)

	if local.IngressClassName != nil {
		remote.WithIngressClassName(*local.IngressClassName)
	}

	return remote
}

// RemoteIngressBackend forges the apply patch for the backend of the reflected ingress, given the local ones.
func RemoteIngressBackend(local *netv1.IngressBackend) *netv1apply.IngressBackendApplyConfiguration {
	if local == nil {
		return nil
	}

	return netv1apply.IngressBackend().
		WithResource(RemoteIngressResource(local.Resource)).
		WithService(RemoteIngressService(local.Service))
}

func RemoteIngressResource(local *corev1.TypedLocalObjectReference) *corev1apply.TypedLocalObjectReferenceApplyConfiguration {
	if local == nil {
		return nil
	}
	res := corev1apply.TypedLocalObjectReference().
		WithKind(local.Kind).
		WithName(local.Name)
	if local.APIGroup != nil {
		res.WithAPIGroup(*local.APIGroup)
	}
	return res
}

func RemoteIngressService(local *netv1.IngressServiceBackend) *netv1apply.IngressServiceBackendApplyConfiguration {
	if local == nil {
		return nil
	}
	return netv1apply.IngressServiceBackend().
		WithName(local.Name).
		WithPort(netv1apply.ServiceBackendPort().
			WithName(local.Port.Name).
			WithNumber(local.Port.Number))
}

func RemoteIngressRules(local []netv1.IngressRule) []*netv1apply.IngressRuleApplyConfiguration {
	remote := make([]*netv1apply.IngressRuleApplyConfiguration, len(local))
	for i := range local {
		remote[i] = netv1apply.IngressRule().
			WithHost(local[i].Host).
			WithHTTP(RemoteIngressHTTP(local[i].HTTP))
	}
	return remote
}

func RemoteIngressHTTP(local *netv1.HTTPIngressRuleValue) *netv1apply.HTTPIngressRuleValueApplyConfiguration {
	if local == nil {
		return nil
	}
	return netv1apply.HTTPIngressRuleValue().
		WithPaths(RemoteIngressPaths(local.Paths)...)
}

func RemoteIngressPaths(local []netv1.HTTPIngressPath) []*netv1apply.HTTPIngressPathApplyConfiguration {
	remote := make([]*netv1apply.HTTPIngressPathApplyConfiguration, len(local))
	for i := range local {
		remote[i] = netv1apply.HTTPIngressPath().
			WithPath(local[i].Path).
			WithBackend(RemoteIngressBackend(&local[i].Backend))
	}
	return remote
}

func RemoteIngressTLS(local []netv1.IngressTLS) []*netv1apply.IngressTLSApplyConfiguration {
	remote := make([]*netv1apply.IngressTLSApplyConfiguration, len(local))
	for i := range local {
		remote[i] = netv1apply.IngressTLS().
			WithHosts(local[i].Hosts...).
			WithSecretName(local[i].SecretName)
	}
	return remote
}
