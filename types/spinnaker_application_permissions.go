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

// ApplicationPermissions represents spinnaker application permissions
type ApplicationPermissions struct {
	Execute []string `json:"EXECUTE"`
	Read    []string `json:"READ"`
	Write   []string `json:"WRITE"`
}

// defaultApplicationPermissions new application permissions object with default values (group 'devops')
func defaultApplicationPermissions() *ApplicationPermissions {
	return &ApplicationPermissions{
		Execute: []string{"devops"},
		Read:    []string{"devops"},
		Write:   []string{"devops"},
	}
}

// AppendApplicationPermissions appends team owner permissions to defaults one
func AppendApplicationPermissions(teamOwner string) *ApplicationPermissions {
	applicationPermissions := defaultApplicationPermissions()

	applicationPermissions.Execute = append(applicationPermissions.Execute, teamOwner)
	applicationPermissions.Read = append(applicationPermissions.Read, teamOwner)
	applicationPermissions.Write = append(applicationPermissions.Write, teamOwner)

	return applicationPermissions
}
