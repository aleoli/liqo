// Package roles defines the ClusterRole containing the permissions required by the virtual kubelet in the local cluster.
package roles

// +kubebuilder:rbac:groups="",resources=configmaps;pods;services;namespaces;secrets;nodes,verbs=get;update;patch;list;watch;delete;create
// +kubebuilder:rbac:groups="",resources=pods/status;services/status;namespaces/status;nodes/status,verbs=get;update;patch;list;watch;delete;create
// +kubebuilder:rbac:groups=apps,resources=replicasets,verbs=get;update;patch;list;watch;delete;create
// +kubebuilder:rbac:groups=apps,resources=replicasets/status,verbs=get;update;patch;list;watch;delete;create
// +kubebuilder:rbac:groups=discovery.k8s.io,resources=endpointslices,verbs=get;update;patch;list;watch;delete;create
// +kubebuilder:rbac:groups=certificates.k8s.io,resources=certificatesigningrequests,verbs=create;get;list;watch
// +kubebuilder:rbac:groups=certificates.k8s.io,resources=certificatesigningrequests/approval,verbs=update

// +kubebuilder:rbac:groups=virtualkubelet.liqo.io,resources=namespacenattingtables,verbs=get;update;patch;list;watch;delete;create
// +kubebuilder:rbac:groups=net.liqo.io,resources=tunnelendpoints,verbs=get;list;watch
// +kubebuilder:rbac:groups=sharing.liqo.io,resources=advertisements,verbs=get;list;watch;update;delete
