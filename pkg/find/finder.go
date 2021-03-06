// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package find

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	bundle "github.com/GoogleCloudPlatform/k8s-cluster-bundle/pkg/apis/bundle/v1alpha1"
	"github.com/GoogleCloudPlatform/k8s-cluster-bundle/pkg/core"
)

// ComponentFinder is a wrapper which allows for efficient searching through
// component data. The data is intended to be readonly; if modifications are
// made to the data, subsequent lookups will fail.
type ComponentFinder struct {
	nameCompLookup map[string][]*bundle.ComponentPackage
	keyCompLookup  map[bundle.ComponentReference]*bundle.ComponentPackage
	data           []*bundle.ComponentPackage
}

// NewComponentFinder creates a new ComponentFinder or returns an error.
func NewComponentFinder(data []*bundle.ComponentPackage) *ComponentFinder {
	nlup := make(map[string][]*bundle.ComponentPackage)
	klup := make(map[bundle.ComponentReference]*bundle.ComponentPackage)
	for _, comp := range data {
		name := comp.Spec.ComponentName
		klup[comp.MakeComponentReference()] = comp
		if list := nlup[name]; list == nil {
			nlup[name] = []*bundle.ComponentPackage{comp}
		} else {
			nlup[name] = append(nlup[name], comp)
		}
	}
	return &ComponentFinder{
		nameCompLookup: nlup,
		keyCompLookup:  klup,
		data:           data,
	}
}

// ComponentPackage returns the component package that matches a reference,
// returning nil if no match is found.
func (f *ComponentFinder) Component(ref bundle.ComponentReference) *bundle.ComponentPackage {
	return f.keyCompLookup[ref]
}

// ComponentsFromName returns thes components that matches a string-name.
func (f *ComponentFinder) ComponentsFromName(name string) []*bundle.ComponentPackage {
	return f.nameCompLookup[name]
}

// ComponentPackage returns the single component package that matches a
// string-name. If no component is found, nil is returne. If there are two
// components that match the name, the method panics.
func (f *ComponentFinder) UniqueComponentFromName(name string) *bundle.ComponentPackage {
	comps := f.ComponentsFromName(name)
	if len(comps) == 0 {
		return nil
	} else if len(comps) > 1 {
		panic(fmt.Sprintf("duplicate component found for name %q", name))
	}
	return comps[0]
}

// Objects returns ComponentPackage's Cluster objects (given some object
// ref) or nil.
func (f *ComponentFinder) Objects(cref bundle.ComponentReference, ref core.ObjectRef) []*unstructured.Unstructured {
	comp := f.Component(cref)
	if comp == nil {
		return nil
	}
	return NewObjectFinder(comp).Objects(ref)
}

// ObjectsFromUniqueComponent gets the objects for a component, which
// has the same behavior as Objects, except that the component name is
// assumed to be unique (and so panics if that assumption does not hold).
func (f *ComponentFinder) ObjectsFromUniqueComponent(name string, ref core.ObjectRef) []*unstructured.Unstructured {
	comp := f.UniqueComponentFromName(name)
	if comp == nil {
		return nil
	}
	return NewObjectFinder(comp).Objects(ref)
}

// ObjectFinder finds objects within components
type ObjectFinder struct {
	component *bundle.ComponentPackage
}

// NewObjectFinder returns an ObjectFinder instance.
func NewObjectFinder(component *bundle.ComponentPackage) *ObjectFinder {
	return &ObjectFinder{component}
}

// Objects finds cluster objects matching a certain ObjectRef key. If
// the ObjectRef is partially filled out, then only those fields will be used
// for searching and the partial matches will be returned.
func (c *ObjectFinder) Objects(ref core.ObjectRef) []*unstructured.Unstructured {
	var out []*unstructured.Unstructured
	for _, o := range c.component.Spec.Objects {
		var key core.ObjectRef
		if ref.Name != "" {
			// Doing a search based on name
			key.Name = o.GetName()
		}
		if ref.APIVersion != "" {
			// Doing a search based on API version
			key.APIVersion = o.GetAPIVersion()
		}
		if ref.Kind != "" {
			// Doing a search based on kind
			key.Kind = o.GetKind()
		}
		if key == ref {
			out = append(out, o)
		}
	}
	return out
}
