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
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

// deleteOptions represents options for delete command
type deleteOptions struct {
	*pipelineOptions
	applicationName string
	pipelineName    string
}

// NewDeleteCmd returns new delete pipeline command
func NewDeleteCmd(pipelineOptions *pipelineOptions) *cobra.Command {
	options := &deleteOptions{
		pipelineOptions: pipelineOptions,
	}

	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del", "remove", "rm"},
		Short:   "delete given pipeline from the provided spinnaker application",
		Long:    "delete given pipeline from the provided spinnaker application",
		Example: "spini pipeline delete [--name=...] [--pipeline=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return deletePipeline(cmd, options)
		},
	}

	cmd.Flags().StringVarP(&options.applicationName, "name", "n", "", "Spinnaker application the pipeline belongs to")
	cmd.Flags().StringVarP(&options.pipelineName, "pipeline", "p", "", "name pipeline to delete")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("pipeline")

	return cmd
}

// deletePipeline delete given pipeline in provided application
func deletePipeline(cmd *cobra.Command, options *deleteOptions) error {
	if options.DryRun {
		fmt.Println("[DRY_RUN] \nDelete pipeline " + options.pipelineName + " from application " + options.applicationName)
	} else {
		resp, err := options.GateClient.PipelineControllerApi.DeletePipelineUsingDELETE(
			options.GateClient.Context,
			options.applicationName,
			options.pipelineName)

		if err != nil {
			return fmt.Errorf("encountered an error deleting pipeline '%s' in application '%s': %s", options.pipelineName, options.applicationName, err)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("encountered an error deleting pipeline, status code: %d", resp.StatusCode)
		}

		fmt.Println("Application " + options.applicationName + ":\n \u2714 Pipeline " + options.pipelineName + " deleted")
	}

	return nil
}
