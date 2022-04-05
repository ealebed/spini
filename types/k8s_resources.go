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
	"k8s.io/apimachinery/pkg/api/resource"
)

// newResourceRequirements return k8s resource requirements object
func newResourceRequirements(tier *Datacenter) apiv1.ResourceRequirements {
	if tier.Resources.Limits == nil {
		tier.Resources.Limits = tier.Resources.Requests
	}

	return apiv1.ResourceRequirements{
		Limits: apiv1.ResourceList{
			apiv1.ResourceCPU:    resource.MustParse(tier.Resources.Limits.CPU),
			apiv1.ResourceMemory: resource.MustParse(tier.Resources.Limits.Memory),
		},
		Requests: apiv1.ResourceList{
			apiv1.ResourceCPU:    resource.MustParse(tier.Resources.Requests.CPU),
			apiv1.ResourceMemory: resource.MustParse(tier.Resources.Requests.Memory),
		},
	}
}
