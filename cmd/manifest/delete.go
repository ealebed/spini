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
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ealebed/spini/types"
	"github.com/ealebed/spini/utils"
)

// deleteOptions represents options for delete command
type deleteOptions struct {
	*manifestOptions
	applicationName string
	localConfig     bool
	repositoryName  string
	branch          string
}

// NewDeleteCmd returns new delete manifest command
func NewDeleteCmd(manifestOptions *manifestOptions) *cobra.Command {
	options := &deleteOptions{
		manifestOptions: manifestOptions,
	}

	cmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del"},
		Short:   "delete yaml manifest(s) for provided application",
		Long:    "delete yaml manifest(s) for provided application",
		Example: "spini manifest delete [--name=...] [--repo=...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteManifest(cmd, options)
		},
	}

	cmd.Flags().StringVarP(&options.applicationName, "name", "n", "", "application name for deleting manifest")
	cmd.Flags().BoolVar(&options.localConfig, "local", true, "read local configuration.json")
	cmd.Flags().StringVarP(&options.repositoryName, "repo", "r", "", "GitHub repository name to read configuration.json from")
	cmd.Flags().StringVarP(&options.branch, "branch", "b", "master", "branch to read configuration.json from")

	if err := cmd.MarkFlagRequired("name"); err != nil {
		return nil
	}
	if err := cmd.MarkFlagRequired("repo"); err != nil {
		return nil
	}

	return cmd
}

// deleteManifest deletes manifest in github repository
func deleteManifest(_ *cobra.Command, options *deleteOptions) error {
	var filename string
	var str []string

	configResponse := utils.LoadConfiguration(options.localConfig, options.Organization, options.repositoryName, options.branch)

	for _, app := range configResponse {
		if app.Application == options.applicationName {
			for _, profile := range *app.Profiles {
				for _, tier := range *profile.Datacenters {
					if profile.ProfileName == "production" {
						filename = options.applicationName + ".yaml"
					} else {
						filename = options.applicationName + "-" + profile.ProfileName + ".yaml"
					}
					str = append(str, "datacenters/"+tier.TierName+"/"+app.Namespace+"/"+filename)
				}
			}
		}
	}
	sourceFiles := strings.Join(str, ",")

	if options.DryRun {
		fmt.Println("[DRY_RUN] Delete yaml-manifest(s):\n", strings.Join(str, "\n"))
	} else {
		PROptions := &types.PullRequestOptions{
			Organization:   options.Organization,
			RepositoryName: options.GitHubRepositoryName,
			AuthorName:     options.GitHubUser,
			AuthorEmail:    options.GitHubEmail,
			PRSubject:      "Delete manifest(s) " + options.applicationName,
			PRDescription:  "Delete *.yaml manifest(s) for application " + options.applicationName,
			CommitMessage:  "Delete *.yaml manifest(s) for application " + options.applicationName,
			CommitBranch:   "update_" + time.Now().Format("2006-01-02-1504"),
		}

		if err := utils.CreatePullRequest(sourceFiles, PROptions); err != nil {
			return fmt.Errorf("failed to create pull request: %w", err)
		}

		fmt.Println("\nManifest(s) deletion succeeded")
	}

	return nil
}
