package crdreplicator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"

	configv1alpha1 "github.com/liqotech/liqo/apis/config/v1alpha1"
	netv1alpha1 "github.com/liqotech/liqo/apis/net/v1alpha1"
	"github.com/liqotech/liqo/pkg/consts"
)

func TestDispatcherReconciler_GetConfig(t *testing.T) {
	dispatcher := Controller{}
	// test 1
	// the list of the resources to be replicated is 0, so we expect a 0 length list to be returned by the function
	t1 := configv1alpha1.DispatcherConfig{ResourcesToReplicate: nil}
	// test 2
	// the list of the resources to be replicated contains 2 elements, so we expect  two elements in the list to be returned by the function
	t2 := configv1alpha1.DispatcherConfig{ResourcesToReplicate: []configv1alpha1.Resource{
		{
			Group:    netv1alpha1.GroupVersion.Group,
			Version:  netv1alpha1.GroupVersion.Version,
			Resource: "networkconfigs"},
		{
			Group:    netv1alpha1.GroupVersion.Group,
			Version:  netv1alpha1.GroupVersion.Version,
			Resource: "tunnelendpoints",
		},
	}}
	tests := []struct {
		config           configv1alpha1.DispatcherConfig
		expectedElements int
	}{
		{t1, 0},
		{t2, 2},
	}

	for _, test := range tests {
		cfg := &configv1alpha1.ClusterConfig{
			Spec: configv1alpha1.ClusterConfigSpec{
				DispatcherConfig: test.config,
			},
			Status: configv1alpha1.ClusterConfigStatus{},
		}
		res := dispatcher.getConfig(cfg)
		assert.Equal(t, test.expectedElements, len(res), "length should be equal")
	}
}

func TestDispatcherReconciler_GetRemovedResources(t *testing.T) {
	dispatcher := Controller{
		RegisteredResources: []resourceToReplicate{
			{
				groupVersionResource: schema.GroupVersionResource{
					Group:    netv1alpha1.GroupVersion.Group,
					Version:  netv1alpha1.GroupVersion.Version,
					Resource: "networkconfigs",
				},
				peeringPhase: consts.PeeringPhaseEstablished,
			},
			{
				groupVersionResource: schema.GroupVersionResource{
					Group:    netv1alpha1.GroupVersion.Group,
					Version:  netv1alpha1.GroupVersion.Version,
					Resource: "tunnelendpoints",
				},
				peeringPhase: consts.PeeringPhaseEstablished,
			},
		},
	}
	// test 1
	// the configuration does not change, is the same
	// so we expect expect to get a 0 length list
	t1 := []resourceToReplicate{
		{
			groupVersionResource: schema.GroupVersionResource{
				Group:    netv1alpha1.GroupVersion.Group,
				Version:  netv1alpha1.GroupVersion.Version,
				Resource: "networkconfigs",
			},
			peeringPhase: consts.PeeringPhaseEstablished,
		},
		{
			groupVersionResource: schema.GroupVersionResource{
				Group:    netv1alpha1.GroupVersion.Group,
				Version:  netv1alpha1.GroupVersion.Version,
				Resource: "tunnelendpoints",
			},
			peeringPhase: consts.PeeringPhaseEstablished,
		},
	}
	// test2
	// we remove a resource from the configuration and add a new one to it
	// so we expect to get a list with 1 element
	t2 := []resourceToReplicate{
		{
			groupVersionResource: schema.GroupVersionResource{
				Group:    netv1alpha1.GroupVersion.Group,
				Version:  netv1alpha1.GroupVersion.Version,
				Resource: "networkconfigs",
			},
			peeringPhase: consts.PeeringPhaseEstablished,
		},
		{
			groupVersionResource: schema.GroupVersionResource{
				Group:    netv1alpha1.GroupVersion.Group,
				Version:  netv1alpha1.GroupVersion.Version,
				Resource: "tunnelendpoints-wrong",
			},
			peeringPhase: consts.PeeringPhaseEstablished,
		},
	}

	tests := []struct {
		config           []resourceToReplicate
		expectedElements int
	}{
		{t1, 0},
		{t2, 1},
	}

	for _, test := range tests {
		res := dispatcher.getRemovedResources(test.config)
		assert.Equal(t, test.expectedElements, len(res), "length should be equal")
	}
}

func TestDispatcherReconciler_UpdateConfig(t *testing.T) {
	dispatcher := Controller{}
	// test 1
	// the list of the resources to be replicated is 0, so we expect a 0 length list to be returned by the function
	// and 0 elements removed
	t1 := configv1alpha1.DispatcherConfig{ResourcesToReplicate: nil}
	// test 2
	// the list of the resources to be replicated contains 2 elements, so we expect  two elements in the list to be returned by the function
	// and 0 elements removed
	t2 := configv1alpha1.DispatcherConfig{ResourcesToReplicate: []configv1alpha1.Resource{
		{
			Group:    netv1alpha1.GroupVersion.Group,
			Version:  netv1alpha1.GroupVersion.Version,
			Resource: "networkconfigs"},
		{
			Group:    netv1alpha1.GroupVersion.Group,
			Version:  netv1alpha1.GroupVersion.Version,
			Resource: "tunnelendpoints",
		},
	}}

	// test 3
	// we remove an existing element and add a new one. we expect to have 2 elements in the registeredResources
	// and 1 element removedResources
	t3 := configv1alpha1.DispatcherConfig{ResourcesToReplicate: []configv1alpha1.Resource{
		{
			Group:    netv1alpha1.GroupVersion.Group,
			Version:  netv1alpha1.GroupVersion.Version,
			Resource: "networkconfigs"},
		{
			Group:    netv1alpha1.GroupVersion.Group,
			Version:  netv1alpha1.GroupVersion.Version,
			Resource: "tunnelendpoints-wrong",
		},
	}}
	tests := []struct {
		config                     configv1alpha1.DispatcherConfig
		expectedElementsResources  int
		expectedElementsRemovedRes int
	}{
		{t1, 0, 0},
		{t2, 2, 0},
		{t3, 2, 1},
	}

	for _, test := range tests {
		cfg := &configv1alpha1.ClusterConfig{
			Spec: configv1alpha1.ClusterConfigSpec{
				DispatcherConfig: test.config,
			},
			Status: configv1alpha1.ClusterConfigStatus{},
		}
		dispatcher.UpdateConfig(cfg)
		assert.Equal(t, test.expectedElementsResources, len(dispatcher.RegisteredResources), "length should be equal")
		assert.Equal(t, test.expectedElementsRemovedRes, len(dispatcher.UnregisteredResources), "length should be equal")
	}
}

// we test that if the *rest.config of the custer is not correct the function return the error
func TestDispatcherReconciler_WatchConfiguration(t *testing.T) {
	dispatcher := Controller{}
	// test1
	// the group version is not correct and we expect an error
	config := k8sManagerLocal.GetConfig()
	err := dispatcher.WatchConfiguration(config, nil)
	assert.NotNil(t, err, "error should be not nil")

	// test2
	// the group version is not correct and we expect an error
	err = dispatcher.WatchConfiguration(config, &configv1alpha1.GroupVersion)
	assert.Nil(t, err, "error should be not nil")
}
