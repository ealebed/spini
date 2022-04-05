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
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/antihax/optional"
	"github.com/ealebed/spini/types"
	"github.com/spf13/cobra"
	gate "github.com/spinnaker/spin/gateapi"
)

// executeAllOptions represents the pipeline execute-all command
type executeAllOptions = struct {
	*pipelineOptions
	applicationName string
	accountName     string
}

// NewExecuteAllCmd returns new delete all pipeline command
func NewExecuteAllCmd(pipelineOptions *pipelineOptions) *cobra.Command {
	options := &executeAllOptions{
		pipelineOptions: pipelineOptions,
	}

	cmd := &cobra.Command{
		Use:     "execute-all",
		Short:   "execute all pipelines in the provided spinnaker application or account(cluster)",
		Long:    "execute all pipelines in the provided spinnaker application or account(cluster)",
		Example: "spini pipeline execute-all [--name=...] / [--account=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeAllPipelines(cmd, options)
		},
	}
	cmd.Flags().StringVarP(&options.applicationName, "name", "n", "", "Spinnaker application the pipelines belongs to")
	cmd.Flags().StringVarP(&options.accountName, "account", "a", "", "Spinnaker account(cluster) the pipelines belongs to")

	return cmd
}

// executeAllPipelines initiates execution of all pipelines in the provided application or account(cluster)
func executeAllPipelines(cmd *cobra.Command, options *executeAllOptions) error {
	var message string

	if options.accountName == "" && options.applicationName == "" {
		return errors.New("you should provide application or account(cluster) name")
	} else if options.accountName != "" && options.applicationName != "" {
		return errors.New("you should provide only one option: application or account(cluster) name")
	} else if options.accountName != "" && options.applicationName == "" {
		message = "[DRY-RUN] \nExecute all pipelines from account(cluster) " + options.accountName
	} else {
		message = "[DRY-RUN] \nExecute all pipelines from application " + options.applicationName
	}

	if options.DryRun {
		fmt.Println(message)
	} else {
		trigger := map[string]interface{}{"type": "manual"}
		if options.applicationName != "" {
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
				if !pipeline.Disabled && strings.Contains(pipeline.Name, "deploy-") {
					resp, err := options.GateClient.PipelineControllerApi.InvokePipelineConfigUsingPOST1(
						options.GateClient.Context,
						options.applicationName,
						pipeline.Name,
						&gate.PipelineControllerApiInvokePipelineConfigUsingPOST1Opts{
							Trigger: optional.NewInterface(trigger),
						})

					if err != nil {
						return fmt.Errorf("execute pipeline failed with response: %v and error: %s", resp, err)
					}

					if resp.StatusCode != http.StatusAccepted {
						return fmt.Errorf("encountered an error executing pipeline, status code: %d", resp.StatusCode)
					}

					fmt.Println("Pipeline " + pipeline.Name + " execution for application " + options.applicationName + " started!")
				}
			}
		}

		if options.accountName != "" {
			appList, resp, err := options.GateClient.ApplicationControllerApi.GetAllApplicationsUsingGET(
				options.GateClient.Context,
				&gate.ApplicationControllerApiGetAllApplicationsUsingGETOpts{
					Account: optional.NewString(options.accountName),
				})

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("encountered an error listing application, status code: %d", resp.StatusCode)
			}

			if err != nil {
				return fmt.Errorf("encountered an error listing application: %s", err)
			}

			for _, application := range appList {
				var lp *[]types.Pipeline
				successListPipelines, resp, err := options.GateClient.ApplicationControllerApi.GetPipelineConfigsForApplicationUsingGET(
					options.GateClient.Context,
					application.(map[string]interface{})["name"].(string))
				if err != nil {
					return fmt.Errorf("encountered an error listing pipelines for account(cluster) '%s': %s",
						application.(map[string]interface{})["name"].(string),
						err)
				}

				if resp.StatusCode != http.StatusOK {
					return fmt.Errorf("encountered an error listing pipelines for account(cluster) %s, status code: %d",
						application.(map[string]interface{})["name"].(string),
						resp.StatusCode)
				}

				prettyListStr, _ := json.MarshalIndent(successListPipelines, "", " ")

				json.Unmarshal([]byte(prettyListStr), &lp)
				for _, pipeline := range *lp {
					if strings.Contains(pipeline.Name, options.accountName) {
						if !pipeline.Disabled && strings.Contains(pipeline.Name, "deploy-") {
							resp, err := options.GateClient.PipelineControllerApi.InvokePipelineConfigUsingPOST1(
								options.GateClient.Context,
								application.(map[string]interface{})["name"].(string),
								pipeline.Name,
								&gate.PipelineControllerApiInvokePipelineConfigUsingPOST1Opts{
									Trigger: optional.NewInterface(trigger),
								})

							if err != nil {
								return fmt.Errorf("execute pipeline failed with response: %v and error: %s", resp, err)
							}

							if resp.StatusCode != http.StatusAccepted {
								return fmt.Errorf("encountered an error executing pipeline, status code: %d", resp.StatusCode)
							}

							fmt.Println("Pipeline " + pipeline.Name + " execution for application " + application.(map[string]interface{})["name"].(string) + " in cluster " + options.accountName + " started!")
						}
					}
				}
			}
		}
	}

	return nil
}
