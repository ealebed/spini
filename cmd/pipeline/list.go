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

	"github.com/ealebed/spini/pkg/output"
	"github.com/spf13/cobra"
)

// listOptions represents options for list command
type listOptions struct {
	*pipelineOptions
	applicationName string
}

// NewListCmd returns new application list command
func NewListCmd(pipelineOptions *pipelineOptions) *cobra.Command {
	options := &listOptions{
		pipelineOptions: pipelineOptions,
	}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "returns list of all pipelines for the provided spinnaker application",
		Long:    "returns list of all pipelines for the provided spinnaker application",
		Example: "spini pipeline list [--name=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listPipeline(cmd, options)
		},
	}

	cmd.Flags().StringVarP(&options.applicationName, "name", "n", "", "Spinnaker application the pipeline belongs to")
	cmd.MarkFlagRequired("name")

	return cmd
}

// listPipeline returns the pipelines for the provided application
func listPipeline(cmd *cobra.Command, options *listOptions) error {
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
	output.JsonOutput(successPayload)

	return nil
}
