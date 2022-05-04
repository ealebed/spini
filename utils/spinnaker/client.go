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

package spin

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ealebed/spini/types"
	"github.com/spinnaker/spin/cmd/gateclient"
	orcaTasks "github.com/spinnaker/spin/cmd/orca-tasks"
	gate "github.com/spinnaker/spin/gateapi"
)

// CreateApplication parse json file with application config and POST creating application task to ORCA endpoint
func CreateApplication(application *types.Application, GateClient *gateclient.GatewayClient) error {
	createAppTask := map[string]interface{}{
		"job": []interface{}{map[string]interface{}{
			"type":        "createApplication",
			"application": application,
			"user":        "devops"},
		},
		"application": application.Name,
		"description": fmt.Sprintf("Create Application: " + application.Name),
	}

	ref, _, err := GateClient.TaskControllerApi.TaskUsingPOST1(
		GateClient.Context,
		createAppTask)
	if err != nil {
		return err
	}

	err = orcaTasks.WaitForSuccessfulTask(GateClient, ref, 5)
	if err != nil {
		return err
	}

	fmt.Println("\u2714 Application " + application.Name + " save succeeded")

	return nil
}

// CreatePipeline parse json file with pipeline config and POST creating pipeline task to ORCA endpoint
func CreatePipeline(pipeline *types.Pipeline, GateClient *gateclient.GatewayClient) error {
	var pipe *types.Pipeline

	foundPipeline, queryResp, _ := GateClient.ApplicationControllerApi.GetPipelineConfigUsingGET(
		GateClient.Context,
		pipeline.Application,
		pipeline.Name)

	if queryResp.StatusCode == http.StatusOK && len(foundPipeline) > 0 {
		fmt.Println("Pipeline " + pipeline.Name + " exists, let's update it!")

		prettyStr, _ := json.MarshalIndent(foundPipeline, "", " ")
		json.Unmarshal([]byte(prettyStr), &pipe)

		for _, triggerExists := range pipe.Triggers {
			// let's use Spinnaker's known service-account in triggers
			for _, triggerCreated := range pipeline.Triggers {
				triggerCreated.RunAsUser = triggerExists.RunAsUser
			}
			// let's use Spinnaker's known dependent Pipeline ID in triggers 'pipeline' type
			if triggerExists.Type == "pipeline" {
				for _, triggerCreated := range pipeline.Triggers {
					if triggerCreated.Type == "pipeline" {
						triggerCreated.Pipeline = triggerExists.Pipeline
					}
				}
			}
		}

		// let's use Spinnaker's known Pipeline ID and index
		pipeline.ID = pipe.ID
		pipeline.Index = pipe.Index
	} else {
		fmt.Println("Pipeline " + pipeline.Name + " doesn't exists, let's create a new one!")
	}

	saveResp, saveErr := GateClient.PipelineControllerApi.SavePipelineUsingPOST(
		GateClient.Context,
		pipeline,
		&gate.PipelineControllerApiSavePipelineUsingPOSTOpts{})

	if saveErr != nil {
		return fmt.Errorf("encountered an error saving pipeline definition: %v", saveErr)
	} else if saveResp.StatusCode != http.StatusOK {
		return fmt.Errorf("encountered an error saving pipeline, status code: %d", saveResp.StatusCode)
	}

	fmt.Println("Application " + pipeline.Application + ":\n \u2714 Pipeline " + pipeline.Name + " save succeeded")

	return nil
}
