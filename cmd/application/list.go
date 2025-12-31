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
	gate "github.com/spinnaker/spin/gateapi"

	"github.com/ealebed/spini/pkg/output"
)

// listOptions represents options for list command
type listOptions struct {
	*applicationOptions
	accountName string
}

// NewListCmd returns new application list command
func NewListCmd(applicationOptions *applicationOptions) *cobra.Command {
	options := listOptions{
		applicationOptions: applicationOptions,
	}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "returns list of all spinnaker applications",
		Long:    "returns list of all spinnaker applications, optionally listed by account(cluster) name",
		Example: "spini application list [--account=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listApplication(cmd, options)
		},
	}

	cmd.PersistentFlags().StringVarP(&options.accountName, "account", "a", "", "Spinnaker account(cluster) the application belongs to")

	return cmd
}

// listApplication returns application list from spinnaker
func listApplication(_ *cobra.Command, options listOptions) error {
	appList, resp, err := options.GateClient.ApplicationControllerApi.GetAllApplicationsUsingGET(
		options.GateClient.Context,
		&gate.ApplicationControllerApiGetAllApplicationsUsingGETOpts{Account: optional.NewString(options.accountName)})

	if resp != nil {
		defer resp.Body.Close() //nolint:errcheck // acceptable to ignore close errors in defer
	}

	if err != nil {
		return err
	}

	if resp != nil && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("encountered an error listing application, status code: %d", resp.StatusCode)
	}

	output.JsonOutput(appList)

	return nil
}
