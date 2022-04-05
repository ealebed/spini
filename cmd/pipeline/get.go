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

// getOptions represents options for get command
type getOptions struct {
	*pipelineOptions
	applicationName string
	pipelineName    string
}

// NewGetCmd returns new get pipeline command
func NewGetCmd(pipelineOptions *pipelineOptions) *cobra.Command {
	options := &getOptions{
		pipelineOptions: pipelineOptions,
	}

	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"read"},
		Short:   "returns the pipeline with the provided name from the provided spinnaker application",
		Long:    "returns the pipeline with the provided name from the provided spinnaker application",
		Example: "spini pipeline get [--name...] [--pipeline=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getPipeline(cmd, options)
		},
	}

	cmd.Flags().StringVarP(&options.applicationName, "name", "n", "", "Spinnaker application the pipeline belongs to")
	cmd.Flags().StringVarP(&options.pipelineName, "pipeline", "p", "", "name of the pipeline")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("pipeline")

	return cmd
}

// getPipeline returns the pipeline with the provided name from the provided application
func getPipeline(cmd *cobra.Command, options *getOptions) error {
	successPayload, resp, err := options.GateClient.ApplicationControllerApi.GetPipelineConfigUsingGET(
		options.GateClient.Context,
		options.applicationName,
		options.pipelineName)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("encountered an error getting pipeline in application %s with name %s, status code: %d",
			options.applicationName,
			options.pipelineName,
			resp.StatusCode)
	}

	output.JsonOutput(successPayload)

	return nil
}
