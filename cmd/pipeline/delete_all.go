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

// deleteAllOptions represents options for delete all command
type deleteAllOptions struct {
	*pipelineOptions
	applicationName string
}

// NewDeleteAllCmd returns new delete all pipeline command
func NewDeleteAllCmd(pipelineOptions *pipelineOptions) *cobra.Command {
	options := &deleteAllOptions{
		pipelineOptions: pipelineOptions,
	}

	cmd := &cobra.Command{
		Use:     "delete-all",
		Aliases: []string{"prune"},
		Short:   "delete all pipelines in the provided spinnaker application",
		Long:    "delete all pipelines in the provided spinnaker application",
		Example: "spini pipeline delete-all [--name=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteAllPipelines(cmd, options)
		},
	}

	cmd.Flags().StringVarP(&options.applicationName, "name", "n", "", "Spinnaker application the pipeline belongs to")
	cmd.MarkFlagRequired("name")

	return cmd
}

// deletePipeline delete the provided pipeline in selected application
func deleteAllPipelines(cmd *cobra.Command, options *deleteAllOptions) error {
	successPayload, resp, err := options.GateClient.ApplicationControllerApi.GetPipelineConfigsForApplicationUsingGET(
		options.GateClient.Context,
		options.applicationName)

	if err != nil {
		return fmt.Errorf("encountered an error listing pipelines for application '%s': %s", options.applicationName, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("encountered an error listing pipelines for application %s, status code: %d",
			options.applicationName,
			resp.StatusCode)
	}

	for index := range successPayload {
		object := successPayload[index].(map[string]interface{})

		if options.DryRun {
			fmt.Println("[DRY_RUN] \nDelete pipeline " + object["name"].(string) + " from application " + options.applicationName)
		} else {
			resp, err := options.GateClient.PipelineControllerApi.DeletePipelineUsingDELETE(
				options.GateClient.Context,
				options.applicationName,
				object["name"].(string))

			if err != nil {
				return fmt.Errorf("encountered an error deleting pipeline '%s' in application '%s': %s",
					object["name"].(string),
					options.applicationName,
					err)
			}

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("encountered an error deleting pipeline, status code: %d", resp.StatusCode)
			}

			fmt.Println("Application " + options.applicationName + ":\n \u2714 Pipeline " + object["name"].(string) + " deleted")
		}
	}

	return nil
}
