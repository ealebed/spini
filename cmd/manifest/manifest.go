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

package manifest

import (
	"github.com/spf13/cobra"

	"github.com/ealebed/spini/cmd"
)

type manifestOptions struct {
	*cmd.GlobalOptions
}

var (
	// Comma-separated list of files to commit and their location. The local file is separated by its target location by a semi-colon.
	// If the file should be in the same location with the same name, you can just put the file name and omit the repetition.
	// Example: README.md,main.go:github/examples/commitpr/main.go
	str, rmStr []string
)

// NewManifestCmd create new manifest command
func NewManifestCmd(globalOptions *cmd.GlobalOptions) *cobra.Command {
	options := &manifestOptions{
		GlobalOptions: globalOptions,
	}

	cmd := &cobra.Command{ //nolint:gocritic // shadowing cmd is common pattern in cobra
		Use:     "manifest",
		Short:   "Working with github repository as manifests storage",
		Long:    "Working with github repository as manifests storage",
		Example: "",
	}

	// create subcommands
	cmd.AddCommand(NewSaveCmd(options))
	cmd.AddCommand(NewSaveAllCmd(options))
	cmd.AddCommand(NewDeleteCmd(options))

	return cmd
}
