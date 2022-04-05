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

package application

import (
	"fmt"
	"net/http"

	"github.com/antihax/optional"
	"github.com/spf13/cobra"
	orcaTasks "github.com/spinnaker/spin/cmd/orca-tasks"
	gate "github.com/spinnaker/spin/gateapi"
)

// deleteOptions represents options for delete command
type deleteOptions struct {
	*applicationOptions
	applicationName string
}

// NewDeleteCmd returns new delete application command
func NewDeleteCmd(applicationOptions *applicationOptions) *cobra.Command {
	options := &deleteOptions{
		applicationOptions: applicationOptions,
	}

	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del", "remove", "rm"},
		Short:   "delete the provided application",
		Long:    "delete the provided application `--name`: Name of the Spinnaker application to delete",
		Example: "spini application delete [--name=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteApplication(cmd, options)
		},
	}

	cmd.Flags().StringVarP(&options.applicationName, "name", "n", "", "spinnaker application name to delete")
	cmd.MarkFlagRequired("name")

	return cmd
}

// deleteApplication delete the provided application
func deleteApplication(cmd *cobra.Command, options *deleteOptions) error {
	if options.DryRun {
		fmt.Println("[DRY_RUN] \nDelete application: " + options.applicationName)
	} else {
		appSpec := map[string]interface{}{
			"type": "deleteApplication",
			"application": map[string]interface{}{
				"name": options.applicationName,
			},
			"user": "devops",
		}

		_, resp, err := options.GateClient.ApplicationControllerApi.GetApplicationUsingGET(
			options.GateClient.Context,
			options.applicationName,
			&gate.ApplicationControllerApiGetApplicationUsingGETOpts{Expand: optional.NewBool(false)})

		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("application '%s' does not exist, exiting", options.applicationName)
		}

		if err != nil {
			return fmt.Errorf("encountered an error checking application existence, status code: %d", resp.StatusCode)
		}

		deleteAppTask := map[string]interface{}{
			"job":         []interface{}{appSpec},
			"application": options.applicationName,
			"description": fmt.Sprintf("Delete Application: %s", options.applicationName),
		}

		taskRef, resp, err := options.GateClient.TaskControllerApi.TaskUsingPOST1(
			options.GateClient.Context,
			deleteAppTask)

		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("encountered an error deleting application, status code: %d", resp.StatusCode)
		}

		err = orcaTasks.WaitForSuccessfulTask(options.GateClient, taskRef, 5)
		if err != nil {
			return err
		}

		fmt.Println("\u2714 Application " + options.applicationName + " deleted")
	}

	return nil
}
