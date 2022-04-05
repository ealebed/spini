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
	"github.com/ealebed/spini/cmd"
	"github.com/spf13/cobra"
)

type pipelineOptions struct {
	*cmd.GlobalOptions
}

// NewPipelineCmd create new pipeline command
func NewPipelineCmd(globalOptions *cmd.GlobalOptions) *cobra.Command {
	options := &pipelineOptions{
		GlobalOptions: globalOptions,
	}

	cmd := &cobra.Command{
		Use:     "pipeline",
		Aliases: []string{"pipelines", "pipe", "pi"},
		Short:   "Working with spinnaker pipelines",
		Long:    "Working with spinnaker pipelines",
		Example: "",
	}

	// create subcommands
	cmd.AddCommand(NewDeleteCmd(options))
	cmd.AddCommand(NewDeleteAllCmd(options))
	cmd.AddCommand(NewDisableCmd(options))
	cmd.AddCommand(NewDisableAllCmd(options))
	cmd.AddCommand(NewEnableCmd(options))
	cmd.AddCommand(NewEnableAllCmd(options))
	cmd.AddCommand(NewExecuteCmd(options))
	cmd.AddCommand(NewExecuteAllCmd(options))
	cmd.AddCommand(NewGetCmd(options))
	cmd.AddCommand(NewListCmd(options))
	cmd.AddCommand(NewSaveCmd(options))
	cmd.AddCommand(NewSaveAllCmd(options))

	return cmd
}
