//go:build e2e
// +build e2e

/*
Copyright 2021 Frederic Branczyk All rights reserved.

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

package e2e

import (
	"fmt"
	"testing"

	"github.com/brancz/kube-rbac-proxy/test/kubetest"
)

func testTLS(s *kubetest.Suite) kubetest.TestSuite {
	return func(t *testing.T) {
		command := `curl %v --connect-timeout 5 -v -s -k --fail -H "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" https://kube-rbac-proxy.default.svc.cluster.local:8443/metrics`

		for _, tc := range []struct {
			name    string
			tlsFlag string
		}{
			{
				name:    "1.0",
				tlsFlag: "--tlsv1.0",
			},
			{
				name:    "1.1",
				tlsFlag: "--tlsv1.1",
			},
			{
				name:    "1.2",
				tlsFlag: "--tlsv1.2",
			},
			{
				name:    "1.3",
				tlsFlag: "--tlsv1.3",
			},
		} {
			kubetest.Scenario{
				Name: tc.name,

				Given: kubetest.Setups(
					kubetest.CreatedManifests(
						s.KubeClient,
						"basics/clusterRole.yaml",
						"basics/clusterRoleBinding.yaml",
						"basics/deployment.yaml",
						"basics/service.yaml",
						"basics/serviceAccount.yaml",
						// This adds the clients cluster role to succeed
						"basics/clusterRole-client.yaml",
						"basics/clusterRoleBinding-client.yaml",
					),
				),
				When: kubetest.Conditions(
					kubetest.PodsAreReady(
						s.KubeClient,
						1,
						"app=kube-rbac-proxy",
					),
					kubetest.ServiceIsReady(
						s.KubeClient,
						"kube-rbac-proxy",
					),
				),
				Then: kubetest.Checks(
					ClientSucceeds(
						s.KubeClient,
						fmt.Sprintf(command, tc.tlsFlag),
						nil,
					),
				),
			}.Run(t)
		}
	}
}
