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
	"github.com/ealebed/spini/cmd"
	"github.com/spf13/cobra"
)

type applicationOptions struct {
	*cmd.GlobalOptions
}

// NewApplicationCmd create new application command
func NewApplicationCmd(globalOptions *cmd.GlobalOptions) *cobra.Command {
	options := &applicationOptions{
		GlobalOptions: globalOptions,
	}

	cmd := &cobra.Command{
		Use:     "application",
		Aliases: []string{"app"},
		Short:   "Working with spinnaker applications",
		Long:    "Working with spinnaker applications",
		Example: "",
	}

	// create subcommands
	cmd.AddCommand(NewDeleteCmd(options))
	cmd.AddCommand(NewGetCmd(options))
	cmd.AddCommand(NewListCmd(options))
	cmd.AddCommand(NewSaveCmd(options))
	cmd.AddCommand(NewSaveAllCmd(options))

	return cmd
}
