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

	"github.com/spf13/cobra"
	gate "github.com/spinnaker/spin/gateapi"

	"github.com/ealebed/spini/pkg/output"
)

// getOptions represents options for get command
type getOptions struct {
	*accountOptions
	accountName string
}

// NewGetCmd returns new get account command
func NewGetCmd(accountOptions *accountOptions) *cobra.Command {
	options := &getOptions{
		accountOptions: accountOptions,
	}

	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"read"},
		Short:   "returns the specified spinnaker account information",
		Long:    "returns the specified spinnaker account information",
		Example: "spini account get [--account=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getAccount(cmd, options)
		},
	}

	cmd.Flags().StringVarP(&options.accountName, "account", "a", "", "spinnaker account(cluster) for getting info")
	if err := cmd.MarkFlagRequired("account"); err != nil {
		return nil
	}

	return cmd
}

// getAccount returns actual attributes for specified account
func getAccount(_ *cobra.Command, options *getOptions) error {
	account, resp, err := options.GateClient.CredentialsControllerApi.GetAccountUsingGET(
		options.GateClient.Context,
		options.accountName,
		&gate.CredentialsControllerApiGetAccountUsingGETOpts{})

	if resp != nil {
		defer resp.Body.Close() //nolint:errcheck // acceptable to ignore close errors in defer
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("account '%s' not found", options.accountName)
		} else if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("encountered an error getting account, status code: %d", resp.StatusCode)
		}
	}

	if err != nil {
		return err
	}

	output.JsonOutput(account)

	return nil
}
