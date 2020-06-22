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

package v1

import (
	"time"

	v1 "github.com/crunchydata/postgres-operator/v4/pkg/apis/crunchydata.com/v1"
	scheme "github.com/crunchydata/postgres-operator/v4/pkg/generated/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// PgpoliciesGetter has a method to return a PgpolicyInterface.
// A group's client should implement this interface.
type PgpoliciesGetter interface {
	Pgpolicies(namespace string) PgpolicyInterface
}

// PgpolicyInterface has methods to work with Pgpolicy resources.
type PgpolicyInterface interface {
	Create(*v1.Pgpolicy) (*v1.Pgpolicy, error)
	Update(*v1.Pgpolicy) (*v1.Pgpolicy, error)
	UpdateStatus(*v1.Pgpolicy) (*v1.Pgpolicy, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.Pgpolicy, error)
	List(opts metav1.ListOptions) (*v1.PgpolicyList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Pgpolicy, err error)
	PgpolicyExpansion
}

// pgpolicies implements PgpolicyInterface
type pgpolicies struct {
	client rest.Interface
	ns     string
}

// newPgpolicies returns a Pgpolicies
func newPgpolicies(c *CrunchydataV1Client, namespace string) *pgpolicies {
	return &pgpolicies{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the pgpolicy, and returns the corresponding pgpolicy object, and an error if there is any.
func (c *pgpolicies) Get(name string, options metav1.GetOptions) (result *v1.Pgpolicy, err error) {
	result = &v1.Pgpolicy{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("pgpolicies").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Pgpolicies that match those selectors.
func (c *pgpolicies) List(opts metav1.ListOptions) (result *v1.PgpolicyList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.PgpolicyList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("pgpolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested pgpolicies.
func (c *pgpolicies) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("pgpolicies").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a pgpolicy and creates it.  Returns the server's representation of the pgpolicy, and an error, if there is any.
func (c *pgpolicies) Create(pgpolicy *v1.Pgpolicy) (result *v1.Pgpolicy, err error) {
	result = &v1.Pgpolicy{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("pgpolicies").
		Body(pgpolicy).
		Do().
		Into(result)
	return
}

// Update takes the representation of a pgpolicy and updates it. Returns the server's representation of the pgpolicy, and an error, if there is any.
func (c *pgpolicies) Update(pgpolicy *v1.Pgpolicy) (result *v1.Pgpolicy, err error) {
	result = &v1.Pgpolicy{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("pgpolicies").
		Name(pgpolicy.Name).
		Body(pgpolicy).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *pgpolicies) UpdateStatus(pgpolicy *v1.Pgpolicy) (result *v1.Pgpolicy, err error) {
	result = &v1.Pgpolicy{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("pgpolicies").
		Name(pgpolicy.Name).
		SubResource("status").
		Body(pgpolicy).
		Do().
		Into(result)
	return
}

// Delete takes name of the pgpolicy and deletes it. Returns an error if one occurs.
func (c *pgpolicies) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("pgpolicies").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *pgpolicies) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("pgpolicies").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched pgpolicy.
func (c *pgpolicies) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Pgpolicy, err error) {
	result = &v1.Pgpolicy{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("pgpolicies").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
