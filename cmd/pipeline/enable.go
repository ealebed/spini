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

	"github.com/spf13/cobra"
	gate "github.com/spinnaker/spin/gateapi"

	"github.com/ealebed/spini/types"
)

// enableOptions represents options for enable command
type enableOptions struct {
	*pipelineOptions
	applicationName string
}

// NewEnableCmd returns new enable pipeline command
func NewEnableCmd(pipelineOptions *pipelineOptions) *cobra.Command { //nolint:dupl // similar structure to other command constructors is acceptable
	options := &enableOptions{
		pipelineOptions: pipelineOptions,
	}

	cmd := &cobra.Command{
		Use:     "enable",
		Aliases: []string{"on"},
		Short:   "enable pipelines in the provided spinnaker application",
		Long:    "enable pipelines in the provided spinnaker application",
		Example: "spini pipeline enable [--name=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return enablePipeline(cmd, options)
		},
	}

	cmd.Flags().StringVarP(&options.applicationName, "name", "n", "", "Spinnaker application the pipeline belongs to")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		return nil
	}

	return cmd
}

// enablePipeline enable all pipelines in selected application
func enablePipeline(_ *cobra.Command, options *enableOptions) error { //nolint:gocyclo // complex business logic requires multiple conditionals
	if options.DryRun {
		fmt.Println("[DRY_RUN] \nDisable pipelines from application " + options.applicationName)
	} else {
		var lp *[]types.Pipeline

		successListPipelines, resp, err := options.GateClient.ApplicationControllerApi.GetPipelineConfigsForApplicationUsingGET(
			options.GateClient.Context,
			options.applicationName)

		if resp != nil {
			defer resp.Body.Close() //nolint:errcheck // acceptable to ignore close errors in defer
		}

		if err != nil {
			return fmt.Errorf("encountered an error listing pipelines for application '%s': %s", options.applicationName, err)
		}

		if resp != nil && resp.StatusCode != http.StatusOK {
			return fmt.Errorf("encountered an error listing pipelines for application %s, status code: %d",
				options.applicationName,
				resp.StatusCode)
		}

		prettyListStr, err := json.MarshalIndent(successListPipelines, "", " ")
		if err != nil {
			return fmt.Errorf("failed to marshal pipeline list: %w", err)
		}

		if err := json.Unmarshal(prettyListStr, &lp); err != nil {
			return fmt.Errorf("failed to unmarshal pipeline list: %w", err)
		}
		for _, pipeline := range *lp {
			successPayload, resp, err := options.GateClient.ApplicationControllerApi.GetPipelineConfigUsingGET(
				options.GateClient.Context,
				options.applicationName,
				pipeline.Name)

			if resp != nil {
				defer resp.Body.Close() //nolint:errcheck,gocritic // acceptable to ignore close errors in defer, defer in loop is intentional
			}

			if err != nil {
				return fmt.Errorf("encountered an error getting pipeline in application '%s': %s", options.applicationName, err)
			}

			if resp != nil && resp.StatusCode != http.StatusOK {
				return fmt.Errorf("encountered an error getting pipeline in application %s with name %s, status code: %d",
					options.applicationName,
					pipeline.Name,
					resp.StatusCode)
			}

			var p *types.Pipeline
			prettyPipeStr, err := json.MarshalIndent(successPayload, "", " ")
			if err != nil {
				return fmt.Errorf("failed to marshal pipeline: %w", err)
			}

			if err := json.Unmarshal(prettyPipeStr, &p); err != nil {
				return fmt.Errorf("failed to unmarshal pipeline: %w", err)
			}
			if p.Disabled {
				p.Disabled = false
			}

			saveResp, saveErr := options.GateClient.PipelineControllerApi.SavePipelineUsingPOST(
				options.GateClient.Context,
				&p,
				&gate.PipelineControllerApiSavePipelineUsingPOSTOpts{})

			if saveResp != nil {
				defer saveResp.Body.Close() //nolint:errcheck,gocritic // acceptable to ignore close errors in defer, defer in loop is intentional
			}

			if saveErr != nil {
				return fmt.Errorf("encountered an error enabling pipeline definition: %v", saveErr)
			}

			if saveResp != nil && saveResp.StatusCode != http.StatusOK {
				return fmt.Errorf("encountered an error enabling pipeline, status code: %d", saveResp.StatusCode)
			}

			fmt.Println("Pipeline " + p.Name + " in application " + options.applicationName + " enabled!")
		}
	}

	return nil
}
