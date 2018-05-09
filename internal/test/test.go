/*
Copyright 2018 The Kubernetes Authors.

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

package test

import (
	"bytes"
	"flag"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
)

// UpdateGolden writes out the golden files with the latest values, rather than failing the test.
var updateGolden = flag.Bool("update", false, "update golden files")

type TestingT interface {
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	HelperT
}

type HelperT interface {
	Helper()
}

func AssertGoldenBytes(t TestingT, actual []byte, filename string) {
	t.Helper()

	if err := compare(actual, path(filename)); err != nil {
		t.Fatalf("%+v", err)
	}
}

func AssertGoldenString(t TestingT, actual, filename string) {
	t.Helper()

	if err := compare([]byte(actual), path(filename)); err != nil {
		t.Fatalf("%+v", err)
	}
}

func path(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join("testdata", filename)
}

func compare(actual []byte, filename string) error {
	if err := update(filename, actual); err != nil {
		return err
	}

	expected, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrapf(err, "unable to read testdata %s", filename)
	}
	if !bytes.Equal(expected, actual) {
		return errors.Errorf("does not match golden file %s\n\nWANT:\n%q\n\nGOT:\n%q\n", filename, expected, actual)
	}
	return nil
}

func update(filename string, in []byte) error {
	if !*updateGolden {
		return nil
	}
	return ioutil.WriteFile(filename, normalize(in), 0666)
}

func normalize(in []byte) []byte {
	return bytes.Replace(in, []byte("\r\n"), []byte("\n"), -1)
}
