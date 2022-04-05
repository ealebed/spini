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

package application

import (
	"encoding/json"
	"fmt"

	"github.com/ealebed/spini/types"
	"github.com/ealebed/spini/utils"
	spin "github.com/ealebed/spini/utils/spinnaker"
	"github.com/spf13/cobra"
)

// saveAllOptions represents options for save-all command
type saveAllOptions struct {
	*applicationOptions
	localConfig    bool
	repositoryName string
	branch         string
}

// NewSaveAllCmd returns new save-all application command
func NewSaveAllCmd(applicationOptions *applicationOptions) *cobra.Command {
	options := &saveAllOptions{
		applicationOptions: applicationOptions,
	}

	cmd := &cobra.Command{
		Use:     "save-all",
		Aliases: []string{"create-all"},
		Short:   "create or update all spinnaker applications from provided GitHub repository",
		Long:    "create or update all spinnaker applications from provided GitHub repository",
		Example: "spini application save [--repo=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return saveAllApplication(cmd, options)
		},
	}

	cmd.Flags().BoolVar(&options.localConfig, "local", true, "read local configuration.json")
	cmd.Flags().StringVarP(&options.repositoryName, "repo", "r", "", "GitHub repository name to read configuration.json from")
	cmd.Flags().StringVarP(&options.branch, "branch", "b", "master", "branch to read configuration.json from")

	return cmd
}

// saveAllApplication creates spinnaker application from json-formatted files
func saveAllApplication(cmd *cobra.Command, options *saveAllOptions) error {
	var appList []*types.Application
	configResponse := utils.LoadConfiguration(options.localConfig, options.Organization, options.repositoryName, options.branch)

	for _, app := range configResponse {
		if app.SkipAutogeneration {
			fmt.Println("Skip " + app.Application + " due to skip flag")
		} else {
			a := types.NewApplication(app)
			appList = append(appList, a)
		}
	}

	if options.DryRun {
		for _, app := range appList {
			fmt.Println("[DRY_RUN] \nGenerate json config for application: \n", app.Name)

			pretty, _ := json.MarshalIndent(app, "", " ")
			utils.WriteFileOnDisk([]byte(pretty), app.Name+".json")
		}
	} else {
		for _, app := range appList {
			spin.CreateApplication(app, options.GateClient)
		}
	}

	return nil
}
