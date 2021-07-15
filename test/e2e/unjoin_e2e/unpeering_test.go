// Copyright 2019-2021 The Liqo Authors
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

package unjoine2e

import (
	"context"
	"fmt"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"

	liqoconst "github.com/liqotech/liqo/pkg/consts"
	"github.com/liqotech/liqo/test/e2e/testutils/tester"
)

const (
	// clustersRequired is the number of clusters required in this E2E test.
	clustersRequired = 2
	// controllerClientPresence indicates if the test use the controller runtime clients.
	controllerClientPresence = false
	// testName is the name of this E2E test.
	testName = "E2E_UNJOIN"
)

func Test_Unjoin(t *testing.T) {
	util.CheckIfTestIsSkipped(t, clustersRequired, testName)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Liqo E2E Suite")
}

var _ = Describe("Liqo E2E", func() {
	var (
		ctx         = context.Background()
		testContext = tester.GetTester(ctx, true)
	)

	Describe("Assert that Liqo is correctly uninstalled", func() {
		Context("Test Unjoin", func() {
			var PodsUpAndRunningTableEntries []TableEntry
			for index := range testContext.Clusters {
				PodsUpAndRunningTableEntries = append(PodsUpAndRunningTableEntries,
							Entry(strings.Join([]string{"Check Liqo is correctly uninstalled on cluster", fmt.Sprintf("%d", index)}, " "),
								testContext.Clusters[index], testContext.Namespace, ))
			}

			DescribeTable("Liqo Pod to Pod Connectivity Check",
				func(homeCluster tester.ClusterContext, namespace string) {
					err := NoPods(homeCluster.NativeClient, testContext.Namespace)
					Expect(err).ShouldNot(HaveOccurred())
					err = NoJoined(homeCluster.NativeClient)
					Expect(err).ShouldNot(HaveOccurred())
				},
			PodsUpAndRunningTableEntries...)

		},
		)
	})
})

func NoPods(clientset *kubernetes.Clientset, namespace string) error {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Error(err)
		return err
	}
	if len(pods.Items) > 0 {
		return fmt.Errorf("there are still running pods in Liqo namespace")
	}
	return nil
}

func NoJoined(clientset *kubernetes.Clientset) error {
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%v=%v", liqoconst.TypeLabel, liqoconst.TypeNode),
	})
	if err != nil {
		klog.Error(err)
		return err
	}

	if len(nodes.Items) > 0 {
		return fmt.Errorf("there are still virtual nodes in the cluster")
	}
	return nil

}
