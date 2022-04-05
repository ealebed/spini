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

	"github.com/antihax/optional"
	"github.com/spf13/cobra"
	gate "github.com/spinnaker/spin/gateapi"
)

// executeOptions represents options for execute command
type executeOptions struct {
	*pipelineOptions
	applicationName string
	pipelineName    string
}

// NewExecuteCmd returns new execute pipeline command
func NewExecuteCmd(pipelineOptions *pipelineOptions) *cobra.Command {
	options := &executeOptions{
		pipelineOptions: pipelineOptions,
	}

	cmd := &cobra.Command{
		Use:     "execute",
		Aliases: []string{"exec"},
		Short:   "execute the provided pipeline in the provided spinnaker application",
		Long:    "execute the provided pipeline in the provided spinnaker application",
		Example: "spini pipeline execute [--name=...] [--pipeline=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return executePipeline(cmd, options)
		},
	}

	cmd.Flags().StringVarP(&options.applicationName, "name", "n", "", "Spinnaker application the pipeline lives to")
	cmd.Flags().StringVarP(&options.pipelineName, "pipeline", "p", "", "name pipeline to execute")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("pipeline")

	return cmd
}

func executePipeline(cmd *cobra.Command, options *executeOptions) error {
	if options.DryRun {
		fmt.Println("[DRY_RUN] \nExecute pipeline " + options.pipelineName + " from application " + options.applicationName)
	} else {
		trigger := map[string]interface{}{"type": "manual", "user": "devops"}
		resp, err := options.GateClient.PipelineControllerApi.InvokePipelineConfigUsingPOST1(
			options.GateClient.Context,
			options.applicationName,
			options.pipelineName,
			&gate.PipelineControllerApiInvokePipelineConfigUsingPOST1Opts{Trigger: optional.NewInterface(trigger)})

		if err != nil {
			return fmt.Errorf("execute pipeline failed with response: %v and error: %s", resp, err)
		}

		if resp.StatusCode != http.StatusAccepted {
			return fmt.Errorf("encountered an error executing pipeline, status code: %d", resp.StatusCode)
		}

		fmt.Println("Pipeline execution started")
	}

	return nil
}
