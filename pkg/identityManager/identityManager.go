package identitymanager

import (
	"k8s.io/client-go/kubernetes"

	"github.com/liqotech/liqo/pkg/clusterid"
	tenantnamespace "github.com/liqotech/liqo/pkg/tenantNamespace"
)

type identityManager struct {
	identityProvider

	client           kubernetes.Interface
	localClusterID   clusterid.ClusterID
	namespaceManager tenantnamespace.Manager
}

// NewCertificateIdentityManager gets a new certificate identity manager.
func NewCertificateIdentityManager(client kubernetes.Interface,
	localClusterID clusterid.ClusterID,
	namespaceManager tenantnamespace.Manager) IdentityManager {
	idProvider := &certificateIdentityProvider{
		namespaceManager: namespaceManager,
		client:           client,
	}

	return &identityManager{
		client:           client,
		localClusterID:   localClusterID,
		namespaceManager: namespaceManager,

		identityProvider: idProvider,
	}
}

// NewIAMIdentityManager gets a new identity manager to handle IAM identities.
func NewIAMIdentityManager(client kubernetes.Interface,
	localClusterID clusterid.ClusterID, awsConfig *AwsConfig,
	namespaceManager tenantnamespace.Manager) IdentityManager {
	idProvider := &iamIdentityProvider{
		awsConfig: awsConfig,
		client:    client,
	}

	return &identityManager{
		client:           client,
		localClusterID:   localClusterID,
		namespaceManager: namespaceManager,

		identityProvider: idProvider,
	}
}
