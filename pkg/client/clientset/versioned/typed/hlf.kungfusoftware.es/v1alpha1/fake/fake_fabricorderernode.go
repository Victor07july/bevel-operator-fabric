/*
 * Copyright Kungfusoftware.es. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fake

import (
	"context"
	json "encoding/json"
	"fmt"

	v1alpha1 "github.com/kfsoftware/hlf-operator/api/hlf.kungfusoftware.es/v1alpha1"
	hlfkungfusoftwareesv1alpha1 "github.com/kfsoftware/hlf-operator/pkg/client/applyconfiguration/hlf.kungfusoftware.es/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeFabricOrdererNodes implements FabricOrdererNodeInterface
type FakeFabricOrdererNodes struct {
	Fake *FakeHlfV1alpha1
	ns   string
}

var fabricorderernodesResource = v1alpha1.SchemeGroupVersion.WithResource("fabricorderernodes")

var fabricorderernodesKind = v1alpha1.SchemeGroupVersion.WithKind("FabricOrdererNode")

// Get takes name of the fabricOrdererNode, and returns the corresponding fabricOrdererNode object, and an error if there is any.
func (c *FakeFabricOrdererNodes) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.FabricOrdererNode, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(fabricorderernodesResource, c.ns, name), &v1alpha1.FabricOrdererNode{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricOrdererNode), err
}

// List takes label and field selectors, and returns the list of FabricOrdererNodes that match those selectors.
func (c *FakeFabricOrdererNodes) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.FabricOrdererNodeList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(fabricorderernodesResource, fabricorderernodesKind, c.ns, opts), &v1alpha1.FabricOrdererNodeList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.FabricOrdererNodeList{ListMeta: obj.(*v1alpha1.FabricOrdererNodeList).ListMeta}
	for _, item := range obj.(*v1alpha1.FabricOrdererNodeList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested fabricOrdererNodes.
func (c *FakeFabricOrdererNodes) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(fabricorderernodesResource, c.ns, opts))

}

// Create takes the representation of a fabricOrdererNode and creates it.  Returns the server's representation of the fabricOrdererNode, and an error, if there is any.
func (c *FakeFabricOrdererNodes) Create(ctx context.Context, fabricOrdererNode *v1alpha1.FabricOrdererNode, opts v1.CreateOptions) (result *v1alpha1.FabricOrdererNode, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(fabricorderernodesResource, c.ns, fabricOrdererNode), &v1alpha1.FabricOrdererNode{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricOrdererNode), err
}

// Update takes the representation of a fabricOrdererNode and updates it. Returns the server's representation of the fabricOrdererNode, and an error, if there is any.
func (c *FakeFabricOrdererNodes) Update(ctx context.Context, fabricOrdererNode *v1alpha1.FabricOrdererNode, opts v1.UpdateOptions) (result *v1alpha1.FabricOrdererNode, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(fabricorderernodesResource, c.ns, fabricOrdererNode), &v1alpha1.FabricOrdererNode{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricOrdererNode), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeFabricOrdererNodes) UpdateStatus(ctx context.Context, fabricOrdererNode *v1alpha1.FabricOrdererNode, opts v1.UpdateOptions) (*v1alpha1.FabricOrdererNode, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(fabricorderernodesResource, "status", c.ns, fabricOrdererNode), &v1alpha1.FabricOrdererNode{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricOrdererNode), err
}

// Delete takes name of the fabricOrdererNode and deletes it. Returns an error if one occurs.
func (c *FakeFabricOrdererNodes) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(fabricorderernodesResource, c.ns, name, opts), &v1alpha1.FabricOrdererNode{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeFabricOrdererNodes) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(fabricorderernodesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.FabricOrdererNodeList{})
	return err
}

// Patch applies the patch and returns the patched fabricOrdererNode.
func (c *FakeFabricOrdererNodes) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.FabricOrdererNode, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(fabricorderernodesResource, c.ns, name, pt, data, subresources...), &v1alpha1.FabricOrdererNode{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricOrdererNode), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied fabricOrdererNode.
func (c *FakeFabricOrdererNodes) Apply(ctx context.Context, fabricOrdererNode *hlfkungfusoftwareesv1alpha1.FabricOrdererNodeApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricOrdererNode, err error) {
	if fabricOrdererNode == nil {
		return nil, fmt.Errorf("fabricOrdererNode provided to Apply must not be nil")
	}
	data, err := json.Marshal(fabricOrdererNode)
	if err != nil {
		return nil, err
	}
	name := fabricOrdererNode.Name
	if name == nil {
		return nil, fmt.Errorf("fabricOrdererNode.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(fabricorderernodesResource, c.ns, *name, types.ApplyPatchType, data), &v1alpha1.FabricOrdererNode{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricOrdererNode), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeFabricOrdererNodes) ApplyStatus(ctx context.Context, fabricOrdererNode *hlfkungfusoftwareesv1alpha1.FabricOrdererNodeApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.FabricOrdererNode, err error) {
	if fabricOrdererNode == nil {
		return nil, fmt.Errorf("fabricOrdererNode provided to Apply must not be nil")
	}
	data, err := json.Marshal(fabricOrdererNode)
	if err != nil {
		return nil, err
	}
	name := fabricOrdererNode.Name
	if name == nil {
		return nil, fmt.Errorf("fabricOrdererNode.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(fabricorderernodesResource, c.ns, *name, types.ApplyPatchType, data, "status"), &v1alpha1.FabricOrdererNode{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.FabricOrdererNode), err
}
