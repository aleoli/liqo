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

package publickey

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	networkingv1alpha1 "github.com/liqotech/liqo/apis/networking/v1alpha1"
	"github.com/liqotech/liqo/pkg/consts"
	liqoutils "github.com/liqotech/liqo/pkg/utils"
)

// ForgePublicKey forges a PublicKey.
func ForgePublicKey(name, namespace, remoteClusterID string, key []byte) (*networkingv1alpha1.PublicKey, error) {
	pubKey := &networkingv1alpha1.PublicKey{
		TypeMeta: metav1.TypeMeta{
			Kind:       networkingv1alpha1.PublicKeyKind,
			APIVersion: networkingv1alpha1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				consts.RemoteClusterID: remoteClusterID,
			},
		},
	}
	err := MutatePublicKey(pubKey, remoteClusterID, key)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}

// MutatePublicKey mutates a PublicKey.
func MutatePublicKey(pubKey *networkingv1alpha1.PublicKey, remoteClusterID string, key []byte) error {
	pubKey.Kind = networkingv1alpha1.PublicKeyKind
	pubKey.APIVersion = networkingv1alpha1.GroupVersion.String()

	if pubKey.Labels == nil {
		pubKey.Labels = make(map[string]string)
	}

	pubKey.Labels[consts.RemoteClusterID] = remoteClusterID
	pubKey.Labels[consts.GatewayResourceLabel] = consts.GatewayResourceLabelValue

	pubKey.Spec.PublicKey = key

	return nil
}

// ExtractKeyFromSecretRef extracts the public key data of a secret from a secret reference.
func ExtractKeyFromSecretRef(ctx context.Context, cl client.Client, secretRef *corev1.ObjectReference) ([]byte, error) {
	var secret corev1.Secret
	if err := cl.Get(ctx, client.ObjectKey{Name: secretRef.Name, Namespace: secretRef.Namespace}, &secret); err != nil {
		return nil, err
	}
	key, ok := secret.Data[consts.PublicKeyField]
	if !ok {
		return nil, fmt.Errorf("secret %q does not contain %s field", client.ObjectKeyFromObject(&secret), consts.PublicKeyField)
	}
	return key, nil
}

// ForgePublicKeyForRemoteCluster forges a PublicKey to be applied on a remote cluster.
func ForgePublicKeyForRemoteCluster(ctx context.Context, cl client.Client,
	liqoNamespace, namespace, gatewayName, gatewayType string) (*networkingv1alpha1.PublicKey, error) {
	clusterIdentity, err := liqoutils.GetClusterIdentityWithControllerClient(ctx, cl, liqoNamespace)
	if err != nil {
		return nil, fmt.Errorf("unable to get cluster identity: %w", err)
	}

	pubKey := &networkingv1alpha1.PublicKey{
		TypeMeta: metav1.TypeMeta{
			Kind:       networkingv1alpha1.PublicKeyKind,
			APIVersion: networkingv1alpha1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterIdentity.ClusterName,
			Labels: map[string]string{
				consts.RemoteClusterID:      clusterIdentity.ClusterID,
				consts.GatewayResourceLabel: consts.GatewayResourceLabelValue,
			},
		},
	}

	if namespace != "" && namespace != corev1.NamespaceDefault {
		pubKey.Namespace = namespace
	}

	// Get public keys of the gateway form the secret reference.
	secretRef, err := GetGatewaySecretReference(ctx, cl, namespace, gatewayName, gatewayType)
	if err != nil {
		return nil, err
	}
	key, err := ExtractKeyFromSecretRef(ctx, cl, secretRef)
	if err != nil {
		return nil, err
	}
	pubKey.Spec = networkingv1alpha1.PublicKeySpec{
		PublicKey: key,
	}

	return pubKey, nil
}

// GetGatewaySecretReference returns the secret reference of a gateway.
func GetGatewaySecretReference(ctx context.Context, cl client.Client, namespace, gatewayName, gatewayType string) (*corev1.ObjectReference, error) {
	switch gatewayType {
	case consts.GatewayTypeServer:
		var gwServer networkingv1alpha1.GatewayServer
		if err := cl.Get(ctx, client.ObjectKey{Namespace: namespace, Name: gatewayName}, &gwServer); err != nil {
			return nil, err
		}
		return gwServer.Status.SecretRef, nil
	case consts.GatewayTypeClient:
		var gwClient networkingv1alpha1.GatewayClient
		if err := cl.Get(ctx, client.ObjectKey{Namespace: namespace, Name: gatewayName}, &gwClient); err != nil {
			return nil, err
		}
		return gwClient.Status.SecretRef, nil
	default:
		return nil, fmt.Errorf("unable to forge PublicKey: invalid gateway type %q", gatewayType)
	}
}
