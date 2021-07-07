// Licensed to Shingo Omura under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Shingo Omura licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/everpeace/kube-throttler/pkg/apis/schedule/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeThrottles implements ThrottleInterface
type FakeThrottles struct {
	Fake *FakeScheduleV1alpha1
	ns   string
}

var throttlesResource = schema.GroupVersionResource{Group: "schedule.k8s.everpeace.github.co", Version: "v1alpha1", Resource: "throttles"}

var throttlesKind = schema.GroupVersionKind{Group: "schedule.k8s.everpeace.github.co", Version: "v1alpha1", Kind: "Throttle"}

// Get takes name of the throttle, and returns the corresponding throttle object, and an error if there is any.
func (c *FakeThrottles) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Throttle, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(throttlesResource, c.ns, name), &v1alpha1.Throttle{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Throttle), err
}

// List takes label and field selectors, and returns the list of Throttles that match those selectors.
func (c *FakeThrottles) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ThrottleList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(throttlesResource, throttlesKind, c.ns, opts), &v1alpha1.ThrottleList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ThrottleList{ListMeta: obj.(*v1alpha1.ThrottleList).ListMeta}
	for _, item := range obj.(*v1alpha1.ThrottleList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested throttles.
func (c *FakeThrottles) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(throttlesResource, c.ns, opts))

}

// Create takes the representation of a throttle and creates it.  Returns the server's representation of the throttle, and an error, if there is any.
func (c *FakeThrottles) Create(ctx context.Context, throttle *v1alpha1.Throttle, opts v1.CreateOptions) (result *v1alpha1.Throttle, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(throttlesResource, c.ns, throttle), &v1alpha1.Throttle{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Throttle), err
}

// Update takes the representation of a throttle and updates it. Returns the server's representation of the throttle, and an error, if there is any.
func (c *FakeThrottles) Update(ctx context.Context, throttle *v1alpha1.Throttle, opts v1.UpdateOptions) (result *v1alpha1.Throttle, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(throttlesResource, c.ns, throttle), &v1alpha1.Throttle{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Throttle), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeThrottles) UpdateStatus(ctx context.Context, throttle *v1alpha1.Throttle, opts v1.UpdateOptions) (*v1alpha1.Throttle, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(throttlesResource, "status", c.ns, throttle), &v1alpha1.Throttle{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Throttle), err
}

// Delete takes name of the throttle and deletes it. Returns an error if one occurs.
func (c *FakeThrottles) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(throttlesResource, c.ns, name), &v1alpha1.Throttle{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeThrottles) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(throttlesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ThrottleList{})
	return err
}

// Patch applies the patch and returns the patched throttle.
func (c *FakeThrottles) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Throttle, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(throttlesResource, c.ns, name, pt, data, subresources...), &v1alpha1.Throttle{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Throttle), err
}
