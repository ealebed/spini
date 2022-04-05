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

package account

import (
	"fmt"
	"net/http"

	"github.com/antihax/optional"
	"github.com/ealebed/spini/pkg/output"
	"github.com/spf13/cobra"
	gate "github.com/spinnaker/spin/gateapi"
)

type listOptions struct {
	*accountOptions
	expand bool
}

// NewListCmd returns new account list command
func NewListCmd(accountOptions *accountOptions) *cobra.Command {
	options := &listOptions{
		accountOptions: accountOptions,
		expand:         false,
	}

	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "returns list of all spinnaker accounts",
		Long:    "returns list of all spinnaker accounts",
		Example: "spini account list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listAccount(cmd, options)
		},
	}

	return cmd
}

// listAccount returns account list from spinnaker
func listAccount(cmd *cobra.Command, options *listOptions) error {
	accountList, resp, err := options.GateClient.CredentialsControllerApi.GetAccountsUsingGET(
		options.GateClient.Context,
		&gate.CredentialsControllerApiGetAccountsUsingGETOpts{Expand: optional.NewBool(options.expand)})

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("encountered an error listing accounts, status code: %d", resp.StatusCode)
	}

	// if options.OutputFormat == "yaml" {
	// 	output.YamlOutput(accountList)
	// } else {
	output.JsonOutput(accountList)
	// }

	return nil
}
