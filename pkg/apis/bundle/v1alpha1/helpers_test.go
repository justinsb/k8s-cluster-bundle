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

package v1alpha1

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestCreateName(t *testing.T) {
	testCases := []struct {
		desc   string
		inName string
		inVer  string
		exp    string
	}{
		{desc: "basic case", inName: "foo", inVer: "0.1.0", exp: "foo-0.1.0"},
		{desc: "empty version", inName: "foo", exp: "foo"},
		{desc: "empty name", inName: "", inVer: "0.1.0", exp: ""},
		{desc: "empty name and version", exp: ""},
	}
	for _, tc := range testCases {
		if n := CreateName(tc.inName, tc.inVer); n != tc.exp {
			t.Errorf("CreateName(%q, %q): got %s, but wanted %s", tc.inName, tc.inVer, n, tc.exp)
		}
	}
}

func TestGetLocalObjectRef(t *testing.T) {
	ref := ComponentReference{"foo", "bar"}
	exp := corev1.LocalObjectReference{"foo-bar"}
	if got := ref.GetLocalObjectRef(); got != exp {
		t.Errorf("GetLocalObjectRef: got %s, but wanted %s", got, exp)
	}
}

func TestGetAllLocalObjectRefs(t *testing.T) {
	cset := ComponentSet{
		Spec: ComponentSetSpec{
			SetName:    "zip",
			Components: []ComponentReference{{"foo", "1.2"}, {"biff", "2.3"}},
		},
	}
	exp := []corev1.LocalObjectReference{{"foo-1.2"}, {"biff-2.3"}}
	if got := cset.GetAllLocalObjectRefs(); !reflect.DeepEqual(got, exp) {
		t.Errorf("GetAllLocalObjectRefs: got %v, but wanted %v", got, exp)
	}
}

func TestMakeComponentReference(t *testing.T) {
	comp := ComponentPackage{
		Spec: ComponentPackageSpec{
			ComponentName: "zip",
			Version:       "1.2.3",
			ObjectFiles: []File{{
				URL: "file://foo.yaml",
			}},
		},
	}
	exp := ComponentReference{"zip", "1.2.3"}
	if got := comp.MakeComponentReference(); got != exp {
		t.Errorf("GetAllLocalObjectRefs: got %v, but wanted %v", got, exp)
	}
}

func TestMakeAndSetName_ComponentSet(t *testing.T) {
	cset := ComponentSet{
		Spec: ComponentSetSpec{
			SetName:    "zip",
			Version:    "1.2.3",
			Components: []ComponentReference{{"foo", "1.2"}, {"biff", "2.3"}},
		},
	}
	exp := "zip-1.2.3"
	cset.MakeAndSetName()
	if got := cset.ObjectMeta.Name; !reflect.DeepEqual(got, exp) {
		t.Errorf("MakeAndSetName: got %s, but wanted %s", got, exp)
	}
}

func TestMakeAndSetName_ComponentPackage(t *testing.T) {
	comp := &ComponentPackage{
		Spec: ComponentPackageSpec{
			ComponentName: "zap",
			Version:       "3.5.3",
		},
	}
	exp := "zap-3.5.3"
	comp.MakeAndSetName()
	if got := comp.ObjectMeta.Name; !reflect.DeepEqual(got, exp) {
		t.Errorf("MakeAndSetName: got %s, but wanted %s", got, exp)
	}
}

func TestMakeAndSetAllNames_Bundle(t *testing.T) {
	comp := &ComponentPackage{
		Spec: ComponentPackageSpec{
			ComponentName: "zap",
			Version:       "3.5.3",
		},
	}

	b := &Bundle{
		SetName:    "zorp",
		Version:    "0.1.0",
		Components: []*ComponentPackage{comp},
	}

	b.MakeAndSetAllNames()
	exp := "zorp-0.1.0"
	if got := b.ObjectMeta.Name; !reflect.DeepEqual(got, exp) {
		t.Errorf("MakeAndSetAllNames: got bundle name %s, but wanted %s", got, exp)
	}

	exp = "zap-3.5.3"
	if got := comp.ObjectMeta.Name; !reflect.DeepEqual(got, exp) {
		t.Errorf("MakeAndSetAllNames: got component name %s, but wanted %s", got, exp)
	}
}

func TestMakeComponentSet(t *testing.T) {
	comp1 := &ComponentPackage{
		Spec: ComponentPackageSpec{
			ComponentName: "zap",
			Version:       "3.5.3",
		},
	}
	comp2 := &ComponentPackage{
		Spec: ComponentPackageSpec{
			ComponentName: "zip",
			Version:       "5.4.3",
		},
	}

	b := &Bundle{
		SetName:    "zorp",
		Version:    "0.1.0",
		Components: []*ComponentPackage{comp1, comp2},
	}

	got := b.MakeComponentSet()
	exp := &ComponentSet{
		Spec: ComponentSetSpec{
			SetName: "zorp",
			Version: "0.1.0",
			Components: []ComponentReference{
				{"zap", "3.5.3"},
				{"zip", "5.4.3"},
			},
		},
	}
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("MakeComponentSet: got %v, but wanted %v", got, exp)
	}
}
