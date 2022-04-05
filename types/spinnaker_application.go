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

// CustomBanner represents custom banners config for spinnaker application
type CustomBanner struct {
	BackgroundColor string `json:"backgroundColor"`
	Enabled         bool   `json:"enabled"`
	Text            string `json:"text"`
	TextColor       string `json:"textColor"`
}

// DataSources represents data sources config for spinnaker application
type DataSources struct {
	Disabled []string `json:"disabled"`
	Enabled  []string `json:"enabled"`
}

// Application represents full config for spinnaker application
type Application struct {
	Accounts       string                  `json:"accounts,omitempty"`
	CloudProviders string                  `json:"cloudProviders"`
	CreateTs       string                  `json:"createTs,omitempty"`
	CustomBanners  []*CustomBanner         `json:"customBanners,omitempty"`
	DataSources    *DataSources            `json:"dataSources,omitempty"`
	Description    string                  `json:"description,omitempty"`
	Email          string                  `json:"email,omitempty"`
	LastModifiedBy string                  `json:"lastModifiedBy,omitempty"`
	Name           string                  `json:"name"`
	Permissions    *ApplicationPermissions `json:"permissions,omitempty"`
	TrafficGuards  []string                `json:"trafficGuards,omitempty"`
	UpdateTs       string                  `json:"updateTs,omitempty"`
	User           string                  `json:"user,omitempty"`

	// DesiredCount                   string            `json:"desiredCount"`
	// EnableRestartRunningExecutions bool              `json:"enableRestartRunningExecutions"`
	// IamRole                        string            `json:"iamRole"`
	// InstancePort                   int               `json:"instancePort"`
	// PlatformHealthOnly             bool              `json:"platformHealthOnly"`
	// PlatformHealthOnlyShowOverride bool              `json:"platformHealthOnlyShowOverride"`
	// ProviderSettings               *ProviderSettings `json:"providerSettings"`
	// RepoProjectKey                 string            `json:"repoProjectKey"`
	// RepoSlug                       string            `json:"repoSlug"`
	// RepoType                       string            `json:"repoType"`
	// TaskDefinition                 string            `json:"taskDefinition"`
}

// defaultCustomBanner return CustomBanner object with default values
func defaultCustomBanner() *CustomBanner {
	return &CustomBanner{
		BackgroundColor: "var(--color-accessory-light)",
		Enabled:         true,
		Text:            "Default Custom Banner Text",
		TextColor:       "var(--color-text-primary)",
	}
}

// defaultApplication return application object with default values
func defaultApplication() *Application {
	return &Application{
		CloudProviders: "kubernetes",
		CustomBanners:  []*CustomBanner{defaultCustomBanner()},
		DataSources: &DataSources{
			Disabled: []string{"securityGroups"},
			Enabled:  []string{},
		},
		TrafficGuards: []string{},
	}
}

func NewApplication(app *Configuration) *Application {
	application := defaultApplication()

	application.CustomBanners[len(application.CustomBanners)-1].Text = app.Application
	if app.Type != "service" {
		application.DataSources.Disabled = append(application.DataSources.Disabled, "loadBalancers")
	}
	application.Description = app.Application
	application.Email = app.OwnerEmail
	application.Name = app.Application
	application.Permissions = AppendApplicationPermissions(app.Owners)
	application.User = app.OwnerEmail

	return application
}
