package util

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	offv1alpha1 "github.com/liqotech/liqo/apis/offloading/v1alpha1"
	liqoconst "github.com/liqotech/liqo/pkg/consts"
	"github.com/liqotech/liqo/test/e2e/testconsts"
)

// CreateNamespaceOffloading gets the right NamespaceOffloading according to the specified strategy.
func CreateNamespaceOffloading(ctx context.Context, cl client.Client, namespaceName string,
	namespaceMappingStrategy offv1alpha1.NamespaceMappingStrategyType,
	podStrategy offv1alpha1.PodOffloadingStrategyType,
	clusterSelector *corev1.NodeSelector) error {
	namespaceOffloading := &offv1alpha1.NamespaceOffloading{
		ObjectMeta: metav1.ObjectMeta{
			Name:      liqoconst.DefaultNamespaceOffloadingName,
			Namespace: namespaceName,
		},
		Spec: offv1alpha1.NamespaceOffloadingSpec{
			NamespaceMappingStrategy: namespaceMappingStrategy,
			PodOffloadingStrategy:    podStrategy,
			ClusterSelector:          *clusterSelector,
		},
	}
	return cl.Create(ctx, namespaceOffloading)
}

// GetClusterSelector returns a cluster selector for a NamespaceOffloading resource.
func GetClusterSelector() *corev1.NodeSelector {
	return &corev1.NodeSelector{NodeSelectorTerms: []corev1.NodeSelectorTerm{
		{
			MatchExpressions: []corev1.NodeSelectorRequirement{
				{
					Key:      testconsts.RegionKey,
					Operator: corev1.NodeSelectorOpIn,
					Values:   []string{testconsts.RegionB},
				},
				{
					Key:      testconsts.ProviderKey,
					Operator: corev1.NodeSelectorOpIn,
					Values:   []string{testconsts.ProviderAWS},
				},
			},
		},
	}}
}
