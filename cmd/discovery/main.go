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

package main

import (
	"flag"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	discoveryv1alpha1 "github.com/liqotech/liqo/apis/discovery/v1alpha1"
	nettypes "github.com/liqotech/liqo/apis/net/v1alpha1"
	advtypes "github.com/liqotech/liqo/apis/sharing/v1alpha1"
	"github.com/liqotech/liqo/internal/discovery"
	foreignclusteroperator "github.com/liqotech/liqo/internal/discovery/foreign-cluster-operator"
	searchdomainoperator "github.com/liqotech/liqo/internal/discovery/search-domain-operator"
	"github.com/liqotech/liqo/pkg/clusterid"
	"github.com/liqotech/liqo/pkg/consts"
	"github.com/liqotech/liqo/pkg/mapperUtils"
	"github.com/liqotech/liqo/pkg/utils/restcfg"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = discoveryv1alpha1.AddToScheme(scheme)
	_ = advtypes.AddToScheme(scheme)
	_ = nettypes.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	klog.Info("Starting")

	namespace := flag.String("namespace", "default", "Namespace where your configs are stored.")
	requeueAfter := flag.Duration("requeue-after", 30*time.Second,
		"Period after that the PeeringRequests status is synchronized")

	clusterName := flag.String(consts.ClusterNameParameter, "", "A mnemonic name associated with the current cluster")
	authServiceAddressOverride := flag.String(consts.AuthServiceAddressOverrideParameter, "",
		"The address the authentication service is reachable from foreign clusters (automatically retrieved if not set")
	authServicePortOverride := flag.String(consts.AuthServicePortOverrideParameter, "",
		"The port the authentication service is reachable from foreign clusters (automatically retrieved if not set")
	autoJoin := flag.Bool("auto-join-discovered-clusters", true, "Whether to automatically peer with discovered clusters")
	ownerReferencesPermissionEnforcement := flag.Bool("owner-references-permission-enforcement", false,
		"Enable support for the OwnerReferencesPermissionEnforcement admission controller "+
			"https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#ownerreferencespermissionenforcement")

	var mdnsConfig discovery.MDNSConfig
	flag.BoolVar(&mdnsConfig.EnableAdvertisement, "mdns-enable-advertisement", false, "Enable the mDNS advertisement on LANs")
	flag.BoolVar(&mdnsConfig.EnableDiscovery, "mdns-enable-discovery", false, "Enable the mDNS discovery on LANs")
	flag.StringVar(&mdnsConfig.Service, "mdns-service-name", "_liqo_auth._tcp",
		"The name of the service used for mDNS advertisement/discovery on LANs")
	flag.StringVar(&mdnsConfig.Domain, "mdns-domain-name", "local.",
		"The name of the domain used for mDNS advertisement/discovery on LANs")
	flag.DurationVar(&mdnsConfig.TTL, "mdns-ttl", 90*time.Second,
		"The time-to-live before an automatically discovered clusters is deleted if no longer announced")
	flag.DurationVar(&mdnsConfig.ResolveRefreshTime, "mdns-resolve-refresh-time", 10*time.Minute,
		"Period after that mDNS resolve context is refreshed")

	dialTCPTimeout := flag.Duration("dial-tcp-timeout", 500*time.Millisecond,
		"Time to wait for a TCP connection to a remote cluster before to consider it as not reachable")

	restcfg.InitFlags(nil)
	klog.InitFlags(nil)
	flag.Parse()

	klog.Info("Namespace: ", *namespace)
	klog.Info("RequeueAfter: ", *requeueAfter)

	config := restcfg.SetRateLimiter(ctrl.GetConfigOrDie())
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Errorf("Failed to create a new Kubernetes client: %w", err)
		os.Exit(1)
	}

	localClusterID, err := clusterid.NewClusterIDFromClient(clientset)
	if err != nil {
		klog.Error(err.Error())
		os.Exit(1)
	}
	err = localClusterID.SetupClusterID(*namespace)
	if err != nil {
		klog.Error(err.Error())
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(config, ctrl.Options{
		MapperProvider:   mapperUtils.LiqoMapperProvider(scheme),
		Scheme:           scheme,
		LeaderElection:   false,
		LeaderElectionID: "b3156c4e.liqo.io",
	})
	if err != nil {
		klog.Errorf("Unable to create main manager: %w", err)
		os.Exit(1)
	}

	// Create an accessory manager restricted to the given namespace only, to avoid introducing
	// performance overhead and requiring excessively wide permissions when not necessary.
	auxmgr, err := ctrl.NewManager(config, ctrl.Options{
		MapperProvider:     mapperUtils.LiqoMapperProvider(scheme),
		Scheme:             scheme,
		Namespace:          *namespace,
		MetricsBindAddress: "0", // Disable the metrics of the auxiliary manager to prevent conflicts.
	})
	if err != nil {
		klog.Errorf("Unable to create auxiliary (namespaced) manager: %w", err)
		os.Exit(1)
	}

	namespacedClient := client.NewNamespacedClient(auxmgr.GetClient(), *namespace)

	klog.Info("Starting the discovery logic")
	discoveryCtl := discovery.NewDiscoveryCtrl(mgr.GetClient(), namespacedClient, *namespace,
		localClusterID, mdnsConfig, *dialTCPTimeout)
	if err := mgr.Add(discoveryCtl); err != nil {
		klog.Errorf("Unable to add the discovery controller to the manager: %w", err)
		os.Exit(1)
	}

	klog.Info("Starting SearchDomain operator")
	searchdomainoperator.StartOperator(mgr, *requeueAfter, discoveryCtl)

	klog.Info("Starting ForeignCluster operator")
	foreignclusteroperator.StartOperator(mgr, namespacedClient, clientset, *namespace,
		*requeueAfter, localClusterID, *clusterName, *authServiceAddressOverride,
		*authServicePortOverride, *autoJoin, *ownerReferencesPermissionEnforcement)

	if err := mgr.Add(auxmgr); err != nil {
		klog.Errorf("Unable to add the auxiliary manager to the main one: %w", err)
		os.Exit(1)
	}
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		klog.Errorf("Unable to start manager: %w", err)
		os.Exit(1)
	}
}
