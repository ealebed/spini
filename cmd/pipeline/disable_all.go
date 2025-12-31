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
	"strings"

	"github.com/antihax/optional"
	"github.com/spf13/cobra"
	gate "github.com/spinnaker/spin/gateapi"

	"github.com/ealebed/spini/types"
)

// disableAllOptions represents options for disable all command
type disableAllOptions struct {
	*pipelineOptions
	accountName string
}

// NewDisableAllCmd returns new disable all pipeline command
func NewDisableAllCmd(pipelineOptions *pipelineOptions) *cobra.Command {
	options := &disableAllOptions{
		pipelineOptions: pipelineOptions,
	}

	cmd := &cobra.Command{
		Use:     "disable-all",
		Short:   "disable all pipelines in the provided spinnaker account(cluster)",
		Long:    "disable all pipelines in the provided spinnaker account(cluster)",
		Example: "spini pipeline disable-all [--account=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return disableAllPipeline(cmd, options)
		},
	}

	cmd.Flags().StringVarP(&options.accountName, "account", "a", "", "Spinnaker account(cluster) the pipelines belongs to")
	if err := cmd.MarkFlagRequired("account"); err != nil {
		return nil
	}

	return cmd
}

// disableAllPipeline disable all pipeline in selected account(cluster)
func disableAllPipeline(_ *cobra.Command, options *disableAllOptions) error { //nolint:gocyclo // complex business logic requires multiple conditionals
	if options.DryRun {
		fmt.Println("[DRY_RUN] \nDisable pipelines from account(cluster) " + options.accountName)
	} else {
		appList, resp, err := options.GateClient.ApplicationControllerApi.GetAllApplicationsUsingGET(
			options.GateClient.Context,
			&gate.ApplicationControllerApiGetAllApplicationsUsingGETOpts{
				Account: optional.NewString(options.accountName),
			})

		if resp != nil {
			defer resp.Body.Close() //nolint:errcheck // acceptable to ignore close errors in defer
		}

		if err != nil {
			return fmt.Errorf("encountered an error listing application: %s", err)
		}

		if resp != nil && resp.StatusCode != http.StatusOK {
			return fmt.Errorf("encountered an error listing application, status code: %d", resp.StatusCode)
		}

		for _, application := range appList {
			var lp *[]types.Pipeline

			successListPipelines, resp, err := options.GateClient.ApplicationControllerApi.GetPipelineConfigsForApplicationUsingGET(
				options.GateClient.Context,
				application.(map[string]interface{})["name"].(string))

			if resp != nil {
				defer resp.Body.Close() //nolint:errcheck,gocritic // acceptable to ignore close errors in defer, defer in loop is intentional
			}

			if err != nil {
				return fmt.Errorf("encountered an error listing pipelines for account(cluster) '%s': %s",
					application.(map[string]interface{})["name"].(string),
					err)
			}

			if resp != nil && resp.StatusCode != http.StatusOK {
				return fmt.Errorf("encountered an error listing pipelines for account(cluster) %s, status code: %d",
					application.(map[string]interface{})["name"].(string),
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
				if !strings.Contains(pipeline.Name, options.accountName) {
					continue
				}
				successPayload, resp, err := options.GateClient.ApplicationControllerApi.GetPipelineConfigUsingGET(
					options.GateClient.Context,
					application.(map[string]interface{})["name"].(string),
					pipeline.Name)

				if resp != nil {
					defer resp.Body.Close() //nolint:errcheck,gocritic // acceptable to ignore close errors in defer, defer in loop is intentional
				}

				if err != nil {
					return fmt.Errorf("encountered an error getting pipelines in account(cluster) '%s': %s",
						application.(map[string]interface{})["name"].(string),
						err)
				}

				if resp != nil && resp.StatusCode != http.StatusOK {
					return fmt.Errorf("encountered an error getting pipelines in account(cluster) %s with name %s, status code: %d",
						application.(map[string]interface{})["name"].(string),
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
				if !p.Disabled {
					p.Disabled = true
				}

				saveResp, saveErr := options.GateClient.PipelineControllerApi.SavePipelineUsingPOST(
					options.GateClient.Context,
					&p,
					&gate.PipelineControllerApiSavePipelineUsingPOSTOpts{})

				if saveResp != nil {
					defer saveResp.Body.Close() //nolint:errcheck,gocritic // acceptable to ignore close errors in defer, defer in loop is intentional
				}

				if saveErr != nil {
					return fmt.Errorf("encountered an error disabling pipeline definition: %v", saveErr)
				}

				if saveResp != nil && saveResp.StatusCode != http.StatusOK {
					return fmt.Errorf("encountered an error disabling pipeline, status code: %d", saveResp.StatusCode)
				}

				fmt.Println("Pipeline " + p.Name + " in application " + application.(map[string]interface{})["name"].(string) + " disabled!")
			}
		}
	}

	return nil
}
