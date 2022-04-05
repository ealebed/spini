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

package assembler

import (
	"github.com/ealebed/spini/cmd"
	"github.com/ealebed/spini/cmd/account"
	"github.com/ealebed/spini/cmd/application"
	"github.com/ealebed/spini/cmd/manifest"
	"github.com/ealebed/spini/cmd/pipeline"
	"github.com/spf13/cobra"
)

// AddSubCommands adds all the subcommands to the rootCmd.
// globalOptions are passed through to the subcommands.
func AddSubCommands(rootCmd *cobra.Command, globalOptions *cmd.GlobalOptions) {
	rootCmd.AddCommand(account.NewAccountCmd(globalOptions))
	rootCmd.AddCommand(application.NewApplicationCmd(globalOptions))
	rootCmd.AddCommand(pipeline.NewPipelineCmd(globalOptions))
	rootCmd.AddCommand(manifest.NewManifestCmd(globalOptions))
}
