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

package cmd

import (
	"io"

	"github.com/spf13/cobra"
	gateclient "github.com/spinnaker/spin/cmd/gateclient"
	"github.com/spinnaker/spin/cmd/output"

	"github.com/ealebed/spini/cmd/version"
	git "github.com/ealebed/spini/utils/github"
)

type GlobalOptions struct {
	configPath           string
	gateEndpoint         string
	Organization         string
	GitHubUser           string
	GitHubEmail          string
	GitHubRepositoryName string
	OutputFormat         string
	DryRun               bool

	GateClient *gateclient.GatewayClient
}

func NewCmdRoot(outWriter, errWriter io.Writer) (*cobra.Command, *GlobalOptions) {
	options := &GlobalOptions{}

	cmd := &cobra.Command{
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version.String(),
	}

	cmd.SetOut(outWriter)
	cmd.SetErr(errWriter)

	// GateClient Flags
	cmd.PersistentFlags().StringVar(&options.configPath, "config", "", "path to config file (default $HOME/.spin/config)")
	cmd.PersistentFlags().StringVar(&options.gateEndpoint, "gate-endpoint", "", "Gate (API server) endpoint (default http://localhost:8084)")

	// TODO: configure colored/formatted output
	// cmd.PersistentFlags().StringVar(&options.OutputFormat, "output", "", "configure output formatting")

	// Other flags
	cmd.PersistentFlags().BoolVar(&options.DryRun, "dry-run", true, "print output / save generated files without real changing system configuration")
	cmd.PersistentFlags().StringVar(&options.Organization, "org", "ealebed", "source owner organization")

	// Initialize GateClient
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		ui := output.NewUI(false, false, nil, outWriter, errWriter)
		gateClient, err := gateclient.NewGateClient(ui, options.gateEndpoint, "", options.configPath, false, false, 0)
		if err != nil {
			return err
		}
		options.GateClient = gateClient

		options.GitHubUser, err = git.ExecGitConfig("user.name")
		if err != nil {
			return err
		}

		options.GitHubEmail, err = git.ExecGitConfig("user.email")
		if err != nil {
			return err
		}

		options.GitHubRepositoryName = "test-k8s"

		return nil
	}

	return cmd, options
}
