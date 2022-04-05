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
	"github.com/ealebed/spini/cmd"
	"github.com/spf13/cobra"
)

type accountOptions struct {
	*cmd.GlobalOptions
}

// NewAccountCmd create new account command
func NewAccountCmd(globalOptions *cmd.GlobalOptions) *cobra.Command {
	options := &accountOptions{
		GlobalOptions: globalOptions,
	}

	cmd := &cobra.Command{
		Use:     "account",
		Aliases: []string{"acc"},
		Short:   "Working with spinnaker accounts(clusters)",
		Long:    "Working with spinnaker accounts(clusters)",
		Example: "",
	}

	// create subcommands
	cmd.AddCommand(NewGetCmd(options))
	cmd.AddCommand(NewListCmd(options))

	return cmd
}
