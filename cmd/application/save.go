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
	"os"

	"github.com/spf13/cobra"

	"github.com/ealebed/spini/types"
	"github.com/ealebed/spini/utils"
	spin "github.com/ealebed/spini/utils/spinnaker"
)

// saveOptions represents options for save command
type saveOptions struct {
	*applicationOptions
	applicationName string
	localConfig     bool
	repositoryName  string
	branch          string
}

// NewSaveCmd returns new save application command
func NewSaveCmd(applicationOptions *applicationOptions) *cobra.Command {
	options := &saveOptions{
		applicationOptions: applicationOptions,
	}

	cmd := &cobra.Command{
		Use:     "save",
		Aliases: []string{"create"},
		Short:   "save the provided spinnaker application",
		Long:    "save the provided spinnaker application",
		Example: "spini application save [--name=...] [--repo=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return saveApplication(cmd, options)
		},
	}

	cmd.Flags().StringVarP(&options.applicationName, "name", "n", "", "spinnaker application name for creating")
	cmd.Flags().BoolVar(&options.localConfig, "local", true, "read local configuration.json")
	cmd.Flags().StringVarP(&options.repositoryName, "repo", "r", "", "GitHub repository name to read configuration.json from")
	cmd.Flags().StringVarP(&options.branch, "branch", "b", "master", "branch to read configuration.json from")

	if err := cmd.MarkFlagRequired("name"); err != nil {
		return nil
	}

	return cmd
}

// saveApplication creates application on spinnaker from json-formatted file
func saveApplication(_ *cobra.Command, options *saveOptions) error {
	var a *types.Application
	configResponse := utils.LoadConfiguration(options.localConfig, options.Organization, options.repositoryName, options.branch)

	for _, app := range configResponse {
		if app.Application == options.applicationName {
			if app.SkipAutogeneration {
				fmt.Println("Skip " + app.Application + " due to skip flag")
				os.Exit(0)
			} else {
				a = types.NewApplication(app)
			}
		}
	}

	if options.DryRun {
		fmt.Println("[DRY_RUN] \nGenerate json config for application: " + options.applicationName)

		pretty, err := json.MarshalIndent(a, "", " ")
		if err != nil {
			return fmt.Errorf("failed to marshal application: %w", err)
		}
		if err := utils.WriteFileOnDisk(pretty, a.Name+".json"); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
	} else {
		if err := spin.CreateApplication(a, options.GateClient); err != nil {
			return fmt.Errorf("failed to create application: %w", err)
		}
	}

	return nil
}
