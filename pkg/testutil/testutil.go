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

// Package testutil provides utilities for reading testdata from children
// directories.
package testutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// TestPathPrefix returns the empty string or the bazel test path prefix.
func TestPathPrefix(pathToRoot, file string) string {
	path := os.Getenv("TEST_SRCDIR") // For dealing with bazel.
	workspace := os.Getenv("TEST_WORKSPACE")
	if path != "" {
		return filepath.Join(path, workspace, file)
	}
	return filepath.Join(pathToRoot, file)
}

// ReadData reads the test-data from disk.
func ReadData(pathToRoot, file string) ([]byte, error) {
	return ioutil.ReadFile(TestPathPrefix(pathToRoot, file))
}
