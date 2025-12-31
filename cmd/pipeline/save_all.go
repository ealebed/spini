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

package pipeline

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ealebed/spini/types"
	"github.com/ealebed/spini/utils"
	spin "github.com/ealebed/spini/utils/spinnaker"
)

// saveAllOptions represents options for save-all command
type saveAllOptions struct {
	*pipelineOptions
	localConfig    bool
	repositoryName string
	branch         string
}

// NewSaveAllCmd returns new save-all pipeline command
func NewSaveAllCmd(pipelineOptions *pipelineOptions) *cobra.Command {
	options := &saveAllOptions{
		pipelineOptions: pipelineOptions,
	}

	cmd := &cobra.Command{
		Use:     "save-all",
		Aliases: []string{"create-all"},
		Short:   "save/update pipeline(s) for all spinnaker's applications from provided GitHub repository",
		Long:    "save/update pipeline(s) for all spinnaker's applications from provided GitHub repository",
		Example: "spini pipeline save [--repo=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return saveAllPipeline(cmd, options)
		},
	}

	cmd.Flags().BoolVar(&options.localConfig, "local", true, "read local configuration.json")
	cmd.Flags().StringVarP(&options.repositoryName, "repo", "r", "", "GitHub repository name to read configuration.json from")
	cmd.Flags().StringVarP(&options.branch, "branch", "b", "master", "branch to read configuration.json from")

	return cmd
}

// saveAllPipeline creates pipelines for all spinnaker's applications from json-formatted file
func saveAllPipeline(_ *cobra.Command, options *saveAllOptions) error {
	var pipeList []*types.Pipeline
	configResponse := utils.LoadConfiguration(options.localConfig, options.Organization, options.repositoryName, options.branch)

	for _, app := range configResponse {
		if !app.SkipAutogeneration {
			pipeList = utils.GeneratePipelines(app, options.Organization, options.GitHubRepositoryName)
		} else {
			fmt.Println("Skip " + app.Application + " due to skip flag")
		}
	}

	if options.DryRun {
		for _, pipeline := range pipeList {
			fmt.Println("[DRY-RUN] \nGenerate json-pipeline(s) for " + pipeline.Application + "\n\t" + pipeline.Name)

			pretty, err := json.MarshalIndent(pipeline, "", " ")
			if err != nil {
				return fmt.Errorf("failed to marshal pipeline: %w", err)
			}
			if err := utils.WriteFileOnDisk(pretty, pipeline.Application+"-"+pipeline.Name+".json"); err != nil {
				return fmt.Errorf("failed to write file: %w", err)
			}
		}
	} else {
		for _, pipeline := range pipeList {
			if err := spin.CreatePipeline(pipeline, options.GateClient); err != nil {
				return fmt.Errorf("failed to create pipeline %s: %w", pipeline.Name, err)
			}
		}
	}

	return nil
}
