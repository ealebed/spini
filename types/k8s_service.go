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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewService return k8s service object
func NewService(application, stage, namespace string, ports []Port) *apiv1.Service {
	servicePorts := []apiv1.ServicePort{}

	if stage != "production" {
		application = application + "-" + stage
	}

	for _, port := range ports {
		servicePort := apiv1.ServicePort{
			Name: port.Name,
			Port: port.ContainerPort,
		}
		servicePorts = append(servicePorts, servicePort)
	}

	return &apiv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      application,
			Namespace: namespace,
		},
		Spec: apiv1.ServiceSpec{
			ClusterIP: "None",
			Ports:     servicePorts,
			Selector: map[string]string{
				"app": application,
			},
		},
	}
}
