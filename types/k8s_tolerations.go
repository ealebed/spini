/*
Copyright © 2022 Yevhen Lebid ealebed@gmail.com

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

package types

import (
	apiv1 "k8s.io/api/core/v1"
)

// newToleration return k8s toleration objects list
func newToleration(nodepool string) []apiv1.Toleration {
	listTolerations := []apiv1.Toleration{
		{
			Effect:            apiv1.TaintEffect("NoExecute"),
			Key:               "node.kubernetes.io/not-ready",
			Operator:          apiv1.TolerationOperator("Exists"),
			TolerationSeconds: int64Ptr(20),
		},
		{
			Effect:            apiv1.TaintEffect("NoExecute"),
			Key:               "node.kubernetes.io/unreachable",
			Operator:          apiv1.TolerationOperator("Exists"),
			TolerationSeconds: int64Ptr(20),
		},
	}

	listTolerations = append(listTolerations, apiv1.Toleration{
		Effect:   apiv1.TaintEffect("NoExecute"),
		Key:      "dedicated",
		Operator: apiv1.TolerationOperator("Equal"),
		Value:    nodepool,
	})

	return listTolerations
}
