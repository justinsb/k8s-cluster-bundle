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

package converter

import (
	"testing"

	"github.com/GoogleCloudPlatform/k8s-cluster-bundle/pkg/testutil"
)

func TestRealisticDataParse(t *testing.T) {
	b, err := testutil.ReadData("../../", "examples/bundle-example.yaml")
	if err != nil {
		t.Fatalf("Error reading file %v", err)
	}

	dataFiles, err := FromYAML(b).ToBundle()
	if err != nil {
		t.Fatalf("Error calling ToBundle(): %v", err)
	}

	if l := len(dataFiles.ComponentFiles); l == 0 {
		t.Fatalf("found zero files, but expected some")
	}
}

func TestRealisticDataParse_ComponentSet(t *testing.T) {
	b, err := testutil.ReadData("../../", "examples/component-set.yaml")
	if err != nil {
		t.Fatalf("Error reading file %v", err)
	}

	cset, err := FromYAML(b).ToComponentSet()
	if err != nil {
		t.Fatalf("Error calling ToComponentSet(): %v", err)
	}

	if l := len(cset.Spec.Components); l == 0 {
		t.Fatalf("found zero components, but expected some")
	}
}
