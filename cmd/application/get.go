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

// getOptions represents options for get command
type getOptions struct {
	*applicationOptions
	applicationName string
	expand          bool
}

// NewGetCmd returns new get application command
func NewGetCmd(applicationOptions *applicationOptions) *cobra.Command {
	options := &getOptions{
		applicationOptions: applicationOptions,
	}

	cmd := &cobra.Command{
		Use:     "get",
		Aliases: []string{"read"},
		Short:   "returns the specified spinnaker application information",
		Long:    "returns the specified spinnaker application information",
		Example: "spini application get [--name=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getApplication(cmd, options)
		},
	}

	cmd.Flags().StringVarP(&options.applicationName, "name", "n", "", "spinnaker application name for getting info")
	// Note that false here means defaults to false, and flips to true if the flag is present.
	cmd.PersistentFlags().BoolVarP(&options.expand, "expand", "x", false, "expand app payload to include clusters")

	if err := cmd.MarkFlagRequired("name"); err != nil {
		return nil
	}

	return cmd
}

// getApplication returns actual attributes for specified application
func getApplication(_ *cobra.Command, options *getOptions) error {
	app, resp, err := options.GateClient.ApplicationControllerApi.GetApplicationUsingGET(
		options.GateClient.Context,
		options.applicationName,
		&gate.ApplicationControllerApiGetApplicationUsingGETOpts{Expand: optional.NewBool(options.expand)})

	if resp != nil {
		defer resp.Body.Close() //nolint:errcheck // acceptable to ignore close errors in defer
	}

	if err != nil {
		return err
	}

	if resp != nil {
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("application '%s' not found", options.applicationName)
		} else if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("encountered an error getting application, status code: %d", resp.StatusCode)
		}
	}

	if options.expand {
		// NOTE: expand returns the actual attributes as well as the app's cluster details, nested in
		// their own fields. This means that the expanded output can't be submitted as input to `save`.
		output.JsonOutput(app)
	} else {
		// NOTE: app GET wraps the actual app attributes in an 'attributes' field.
		output.JsonOutput(app["attributes"])
	}

	return nil
}
