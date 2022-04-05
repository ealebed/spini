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

package types

import "github.com/google/go-github/v43/github"

// PullRequestOptions stores additional GitHub specific data for creating PR
type PullRequestOptions struct {
	// Name of the owner (user or org) of the repo to create the commit in
	Organization string
	// Name of repository to create the commit in
	RepositoryName string
	// Name of the author of the commit
	AuthorName string
	// Email of the author of the commit
	AuthorEmail string
	// Name of branch to create the commit in. If it does not already exists, it will be created using the `master` branch
	CommitBranch string
	// Content of the commit message
	CommitMessage string
	// Title of the pull request. If not specified, no pull request will be created
	PRSubject string
	// Text to put in the description of the pull request
	PRDescription string

	Entries   []*github.TreeEntry
	Reference *github.Reference
	Tree      *github.Tree
}
