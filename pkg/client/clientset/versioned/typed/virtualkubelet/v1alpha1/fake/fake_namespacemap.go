// Copyright 2019-2024 The Liqo Authors
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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"

	v1alpha1 "github.com/liqotech/liqo/apis/virtualkubelet/v1alpha1"
)

// FakeNamespaceMaps implements NamespaceMapInterface
type FakeNamespaceMaps struct {
	Fake *FakeVirtualkubeletV1alpha1
	ns   string
}

var namespacemapsResource = v1alpha1.SchemeGroupVersion.WithResource("namespacemaps")

var namespacemapsKind = v1alpha1.SchemeGroupVersion.WithKind("NamespaceMap")

// Get takes name of the namespaceMap, and returns the corresponding namespaceMap object, and an error if there is any.
func (c *FakeNamespaceMaps) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.NamespaceMap, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(namespacemapsResource, c.ns, name), &v1alpha1.NamespaceMap{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NamespaceMap), err
}

// List takes label and field selectors, and returns the list of NamespaceMaps that match those selectors.
func (c *FakeNamespaceMaps) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.NamespaceMapList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(namespacemapsResource, namespacemapsKind, c.ns, opts), &v1alpha1.NamespaceMapList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.NamespaceMapList{ListMeta: obj.(*v1alpha1.NamespaceMapList).ListMeta}
	for _, item := range obj.(*v1alpha1.NamespaceMapList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested namespaceMaps.
func (c *FakeNamespaceMaps) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(namespacemapsResource, c.ns, opts))

}

// Create takes the representation of a namespaceMap and creates it.  Returns the server's representation of the namespaceMap, and an error, if there is any.
func (c *FakeNamespaceMaps) Create(ctx context.Context, namespaceMap *v1alpha1.NamespaceMap, opts v1.CreateOptions) (result *v1alpha1.NamespaceMap, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(namespacemapsResource, c.ns, namespaceMap), &v1alpha1.NamespaceMap{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NamespaceMap), err
}

// Update takes the representation of a namespaceMap and updates it. Returns the server's representation of the namespaceMap, and an error, if there is any.
func (c *FakeNamespaceMaps) Update(ctx context.Context, namespaceMap *v1alpha1.NamespaceMap, opts v1.UpdateOptions) (result *v1alpha1.NamespaceMap, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(namespacemapsResource, c.ns, namespaceMap), &v1alpha1.NamespaceMap{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NamespaceMap), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeNamespaceMaps) UpdateStatus(ctx context.Context, namespaceMap *v1alpha1.NamespaceMap, opts v1.UpdateOptions) (*v1alpha1.NamespaceMap, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(namespacemapsResource, "status", c.ns, namespaceMap), &v1alpha1.NamespaceMap{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NamespaceMap), err
}

// Delete takes name of the namespaceMap and deletes it. Returns an error if one occurs.
func (c *FakeNamespaceMaps) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(namespacemapsResource, c.ns, name, opts), &v1alpha1.NamespaceMap{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeNamespaceMaps) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(namespacemapsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.NamespaceMapList{})
	return err
}

// Patch applies the patch and returns the patched namespaceMap.
func (c *FakeNamespaceMaps) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.NamespaceMap, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(namespacemapsResource, c.ns, name, pt, data, subresources...), &v1alpha1.NamespaceMap{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NamespaceMap), err
}
