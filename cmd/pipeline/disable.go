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
	"net/http"

	"github.com/ealebed/spini/types"
	"github.com/spf13/cobra"
	gate "github.com/spinnaker/spin/gateapi"
)

// disableOptions represents options for disable command
type disableOptions struct {
	*pipelineOptions
	applicationName string
}

// NewDisableCmd returns new disable pipeline command
func NewDisableCmd(pipelineOptions *pipelineOptions) *cobra.Command {
	options := &disableOptions{
		pipelineOptions: pipelineOptions,
	}

	cmd := &cobra.Command{
		Use:     "disable",
		Aliases: []string{"off"},
		Short:   "disable pipelines in the provided spinnaker application",
		Long:    "disable pipelines in the provided spinnaker application",
		Example: "spini pipeline disable [--name=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return disablePipeline(cmd, options)
		},
	}

	cmd.Flags().StringVarP(&options.applicationName, "name", "n", "", "Spinnaker application the pipeline belongs to")
	cmd.MarkFlagRequired("name")

	return cmd
}

// disablePipeline disable pipelines in selected application
func disablePipeline(cmd *cobra.Command, options *disableOptions) error {
	if options.DryRun {
		fmt.Println("[DRY_RUN] \nDisable pipelines from application " + options.applicationName)
	} else {
		var lp *[]types.Pipeline

		successListPipelines, resp, err := options.GateClient.ApplicationControllerApi.GetPipelineConfigsForApplicationUsingGET(
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

		prettyListStr, _ := json.MarshalIndent(successListPipelines, "", " ")

		json.Unmarshal([]byte(prettyListStr), &lp)
		for _, pipeline := range *lp {
			successPayload, resp, err := options.GateClient.ApplicationControllerApi.GetPipelineConfigUsingGET(
				options.GateClient.Context,
				options.applicationName,
				pipeline.Name)

			if err != nil {
				return fmt.Errorf("encountered an error getting pipeline in application '%s': %s", options.applicationName, err)
			}

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("encountered an error getting pipeline in application %s with name %s, status code: %d",
					options.applicationName,
					pipeline.Name,
					resp.StatusCode)
			}

			var p *types.Pipeline
			prettyPipeStr, _ := json.MarshalIndent(successPayload, "", " ")

			json.Unmarshal([]byte(prettyPipeStr), &p)
			if !p.Disabled {
				p.Disabled = true
			}

			saveResp, saveErr := options.GateClient.PipelineControllerApi.SavePipelineUsingPOST(
				options.GateClient.Context,
				&p,
				&gate.PipelineControllerApiSavePipelineUsingPOSTOpts{})

			if saveErr != nil {
				return fmt.Errorf("encountered an error disabling pipeline definition: %v", saveErr)
			}

			if saveResp.StatusCode != http.StatusOK {
				return fmt.Errorf("encountered an error disabling pipeline, status code: %d", saveResp.StatusCode)
			}

			fmt.Println("Pipeline " + p.Name + " in application " + options.applicationName + " disabled!")
		}
	}

	return nil
}
