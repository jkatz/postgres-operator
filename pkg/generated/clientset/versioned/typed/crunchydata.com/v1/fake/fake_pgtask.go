/*
Copyright 2020 Crunchy Data Solutions, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	crunchydatacomv1 "github.com/crunchydata/postgres-operator/v4/pkg/apis/crunchydata.com/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakePgtasks implements PgtaskInterface
type FakePgtasks struct {
	Fake *FakeCrunchydataV1
	ns   string
}

var pgtasksResource = schema.GroupVersionResource{Group: "crunchydata.com", Version: "v1", Resource: "pgtasks"}

var pgtasksKind = schema.GroupVersionKind{Group: "crunchydata.com", Version: "v1", Kind: "Pgtask"}

// Get takes name of the pgtask, and returns the corresponding pgtask object, and an error if there is any.
func (c *FakePgtasks) Get(name string, options v1.GetOptions) (result *crunchydatacomv1.Pgtask, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(pgtasksResource, c.ns, name), &crunchydatacomv1.Pgtask{})

	if obj == nil {
		return nil, err
	}
	return obj.(*crunchydatacomv1.Pgtask), err
}

// List takes label and field selectors, and returns the list of Pgtasks that match those selectors.
func (c *FakePgtasks) List(opts v1.ListOptions) (result *crunchydatacomv1.PgtaskList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(pgtasksResource, pgtasksKind, c.ns, opts), &crunchydatacomv1.PgtaskList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &crunchydatacomv1.PgtaskList{ListMeta: obj.(*crunchydatacomv1.PgtaskList).ListMeta}
	for _, item := range obj.(*crunchydatacomv1.PgtaskList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested pgtasks.
func (c *FakePgtasks) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(pgtasksResource, c.ns, opts))

}

// Create takes the representation of a pgtask and creates it.  Returns the server's representation of the pgtask, and an error, if there is any.
func (c *FakePgtasks) Create(pgtask *crunchydatacomv1.Pgtask) (result *crunchydatacomv1.Pgtask, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(pgtasksResource, c.ns, pgtask), &crunchydatacomv1.Pgtask{})

	if obj == nil {
		return nil, err
	}
	return obj.(*crunchydatacomv1.Pgtask), err
}

// Update takes the representation of a pgtask and updates it. Returns the server's representation of the pgtask, and an error, if there is any.
func (c *FakePgtasks) Update(pgtask *crunchydatacomv1.Pgtask) (result *crunchydatacomv1.Pgtask, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(pgtasksResource, c.ns, pgtask), &crunchydatacomv1.Pgtask{})

	if obj == nil {
		return nil, err
	}
	return obj.(*crunchydatacomv1.Pgtask), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakePgtasks) UpdateStatus(pgtask *crunchydatacomv1.Pgtask) (*crunchydatacomv1.Pgtask, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(pgtasksResource, "status", c.ns, pgtask), &crunchydatacomv1.Pgtask{})

	if obj == nil {
		return nil, err
	}
	return obj.(*crunchydatacomv1.Pgtask), err
}

// Delete takes name of the pgtask and deletes it. Returns an error if one occurs.
func (c *FakePgtasks) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(pgtasksResource, c.ns, name), &crunchydatacomv1.Pgtask{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakePgtasks) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(pgtasksResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &crunchydatacomv1.PgtaskList{})
	return err
}

// Patch applies the patch and returns the patched pgtask.
func (c *FakePgtasks) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *crunchydatacomv1.Pgtask, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(pgtasksResource, c.ns, name, pt, data, subresources...), &crunchydatacomv1.Pgtask{})

	if obj == nil {
		return nil, err
	}
	return obj.(*crunchydatacomv1.Pgtask), err
}
