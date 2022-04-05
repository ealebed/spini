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
	"os"

	"github.com/ealebed/spini/types"
	"github.com/ealebed/spini/utils"
	spin "github.com/ealebed/spini/utils/spinnaker"
	"github.com/spf13/cobra"
)

// saveOptions represents options for save command
type saveOptions struct {
	*pipelineOptions
	applicationName string
	localConfig     bool
	repositoryName  string
	branch          string
}

// NewSaveCmd returns new save pipeline command
func NewSaveCmd(pipelineOptions *pipelineOptions) *cobra.Command {
	options := &saveOptions{
		pipelineOptions: pipelineOptions,
	}

	cmd := &cobra.Command{
		Use:     "save",
		Aliases: []string{"create"},
		Short:   "save/update pipeline(s) for the provided spinnaker application",
		Long:    "save/update pipeline(s) for the provided spinnaker application",
		Example: "spini pipeline save [--name=...] [--repo=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return savePipeline(cmd, options)
		},
	}

	cmd.Flags().StringVarP(&options.applicationName, "name", "n", "", "spinnaker application the pipeline belongs to")
	cmd.Flags().BoolVar(&options.localConfig, "local", true, "read local configuration.json")
	cmd.Flags().StringVarP(&options.repositoryName, "repo", "r", "", "GitHub repository name to read configuration.json from")
	cmd.Flags().StringVarP(&options.branch, "branch", "b", "master", "branch to read configuration.json from")

	cmd.MarkFlagRequired("name")

	return cmd
}

// savePipeline creates pipeline on spinnaker application from json-formatted file
func savePipeline(cmd *cobra.Command, options *saveOptions) error {
	var pipeList []*types.Pipeline
	configResponse := utils.LoadConfiguration(options.localConfig, options.Organization, options.repositoryName, options.branch)

	for _, app := range configResponse {
		if app.Application == options.applicationName {
			if !app.SkipAutogeneration {
				pipeList = utils.GeneratePipelines(app, options.Organization, options.GitHubRepositoryName)
			} else {
				fmt.Println("Skip " + app.Application + " due to skip flag")
				os.Exit(0)
			}
		}
	}

	if options.DryRun {
		for _, pipe := range pipeList {
			fmt.Println("[DRY_RUN] \nGenerate pipeline: ", pipe.Name)

			pretty, _ := json.MarshalIndent(pipe, "", " ")
			utils.WriteFileOnDisk([]byte(pretty), pipe.Name+".json")
		}
	} else {
		for _, pipeline := range pipeList {
			spin.CreatePipeline(pipeline, options.GateClient)
		}
	}

	return nil
}
