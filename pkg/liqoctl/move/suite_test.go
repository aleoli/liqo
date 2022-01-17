// Copyright 2019-2022 The Liqo Authors
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

package move

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMove(t *testing.T) {
	defer GinkgoRecover()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Move Suite")
}

type matchObject struct {
	name      string
	namespace string
}

func (m *matchObject) Match(actual interface{}) (success bool, err error) {
	obj, ok := actual.(metav1.Object)
	if !ok {
		return false, nil
	}

	if obj.GetName() != m.name {
		return false, nil
	}
	if obj.GetNamespace() != m.namespace {
		return false, nil
	}

	return true, nil
}

func (m *matchObject) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n%s\nto be an object with name %s and namespace %s", format.Object(actual, 1), m.name, m.namespace)
}

func (m *matchObject) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n%s\nto not to be an object with name %s and namespace %s", format.Object(actual, 1), m.name, m.namespace)
}
