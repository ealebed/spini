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

package github

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/ealebed/spini/types"
	"github.com/google/go-github/v44/github"
	"github.com/google/go-querystring/query"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
}

func NewClient() *Client {
	githubToken := os.Getenv("GITHUB_AUTH_TOKEN")
	var c *github.Client

	if githubToken == "" {
		fmt.Printf("Unauthorized: No GitHub token present!\n")
	} else {
		tokenService := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubToken})
		tokenClient := oauth2.NewClient(context.Background(), tokenService)

		c = github.NewClient(tokenClient)
	}

	return &Client{c}
}

// ExecGitConfig check git configuration
func ExecGitConfig(args ...string) (string, error) {
	gitArgs := append([]string{"config", "--get", "--null"}, args...)
	var stdout bytes.Buffer
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdout = &stdout
	cmd.Stderr = ioutil.Discard

	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() == 1 {
				return "", err
			}
		}
		return "", err
	}

	return strings.TrimRight(stdout.String(), "\000"), nil
}

// addOptions adds the parameters in opt as URL query parameters to s. opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// readFileContent gets the file content for the given filename in repository
func (c *Client) readFileContent(org, name, branch, path string) (fileContent *github.RepositoryContent, resp *github.Response, err error) {
	escapedPath := (&url.URL{Path: path}).String()
	u := fmt.Sprintf("repos/%s/%s/contents/%s", org, name, escapedPath)

	// RepositoryContentGetOptions represents an optional ref parameter, which can be a SHA, branch, or tag
	opt := &github.RepositoryContentGetOptions{
		Ref: branch,
	}

	u, err = addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var rawJSON json.RawMessage
	resp, err = c.client.Do(context.Background(), req, &rawJSON)
	if err != nil {
		return nil, resp, err
	}

	fileUnmarshalError := json.Unmarshal(rawJSON, &fileContent)
	if fileUnmarshalError == nil {
		return fileContent, resp, nil
	}

	return nil, resp, fmt.Errorf("unmarshalling failed for file content: %s", fileUnmarshalError)
}

// ReadJSONFromRepoToStruct returns struct with configuration values from selected repository
func (c *Client) ReadJSONFromRepoToStruct(org, repoName, branch string) []*types.Configuration {
	conf := make([]*types.Configuration, 0)

	fileContentToEncode, _, err := c.readFileContent(org, repoName, branch, "configuration.json")
	if err != nil {
		fmt.Printf("Could not get file 'configuration.json' from repository %s due to %s", repoName, err)
		os.Exit(1)
	}

	DownloadURL := fileContentToEncode.GetDownloadURL()
	resp, err := http.Get(DownloadURL)
	if err != nil {
		fmt.Printf("Failed to download file 'configuration.json' due to %s", err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respByte := buf.Bytes()

	if err := json.Unmarshal(respByte, &conf); err != nil {
		fmt.Printf("Failed to unmarshal JSON file 'configuration.json' due to %s", err)
		os.Exit(1)
	}

	return conf
}

// NewPullRequest prepare additional GitHub specific data for creating pr
func (c *Client) NewPullRequest(pro *types.PullRequestOptions) (err error) {
	ref, err := c.getRef(pro.Organization, pro.RepositoryName, pro.CommitBranch)
	if err != nil {
		fmt.Printf("Unable to get/create the commit reference: %s", err)
		os.Exit(1)
	}
	if ref == nil {
		fmt.Printf("No error where returned but the reference is nil")
		os.Exit(1)
	}
	tree, err := c.getTree(ref, pro.Organization, pro.RepositoryName, pro.Entries)
	if err != nil {
		fmt.Printf("Unable to create the tree based on the provided files: %s", err)
		os.Exit(1)
	}

	pro.Reference = ref
	pro.Tree = tree

	if err := c.pushCommit(pro); err != nil {
		fmt.Printf("Unable to create the commit: %s", err)
		os.Exit(1)
	}

	if err := c.createPR(pro); err != nil {
		fmt.Printf("Error while creating the pull request: %s", err)
		os.Exit(1)
	}

	return nil
}

// getRef returns the commit branch reference object if it exists or creates it from the base branch before returning it.
func (c *Client) getRef(org, name, commitBranch string) (ref *github.Reference, err error) {
	if ref, _, err = c.client.Git.GetRef(context.Background(), org, name, "refs/heads/"+commitBranch); err == nil {
		return ref, nil
	}

	// We consider that an error means the branch has not been found and needs to be created.
	if commitBranch == "master" {
		return nil, errors.New("the commit branch does not exist, base-branch is master")
	}

	var baseRef *github.Reference
	if baseRef, _, err = c.client.Git.GetRef(context.Background(), org, name, "refs/heads/master"); err != nil {
		return nil, err
	}
	newRef := &github.Reference{Ref: github.String("refs/heads/" + commitBranch), Object: &github.GitObject{SHA: baseRef.Object.SHA}}
	ref, _, err = c.client.Git.CreateRef(context.Background(), org, name, newRef)
	return ref, err
}

// getTree generates the tree to commit based on the given files and the commit of the ref you got in getRef.
func (c *Client) getTree(ref *github.Reference, org, name string, entries []*github.TreeEntry) (tree *github.Tree, err error) {
	tree, _, err = c.client.Git.CreateTree(context.Background(), org, name, *ref.Object.SHA, entries)
	return tree, err
}

// pushCommit creates the commit in the given reference using the given tree.
func (c *Client) pushCommit(pro *types.PullRequestOptions) (err error) {
	// Get the parent commit to attach the commit to.
	parent, _, err := c.client.Repositories.GetCommit(
		context.Background(),
		pro.Organization,
		pro.RepositoryName,
		*pro.Reference.Object.SHA,
		nil)
	if err != nil {
		return err
	}
	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA

	// Create the commit using the tree.
	date := time.Now()
	author := &github.CommitAuthor{Date: &date, Name: &pro.AuthorName, Email: &pro.AuthorEmail}
	commit := &github.Commit{Author: author, Message: &pro.CommitMessage, Tree: pro.Tree, Parents: []*github.Commit{parent.Commit}}

	newCommit, _, err := c.client.Git.CreateCommit(context.Background(), pro.Organization, pro.RepositoryName, commit)
	if err != nil {
		return err
	}

	gotComparsion, _, err := c.client.Repositories.CompareCommits(
		context.Background(),
		pro.Organization,
		pro.RepositoryName,
		parent.Commit.GetSHA(),
		newCommit.GetSHA(),
		nil)
	if err != nil {
		fmt.Printf("Repositories.CompareCommits returned error: %v", err)
		os.Exit(1)
	}

	if len(gotComparsion.Files) == 0 {
		// Delete `fake` branches (references) if there are no real changes in commit
		// https://developer.github.com/v3/git/refs/#delete-a-reference
		_, err = c.client.Git.DeleteRef(
			context.Background(),
			pro.Organization,
			pro.RepositoryName,
			pro.Reference.GetRef())
		if err != nil {
			fmt.Printf("Git.DeleteRef returned error: %v\n", err)
		}

		fmt.Printf("No files changed, skip PR creation! \nSee diff: %s", gotComparsion.GetDiffURL())
		os.Exit(0)
	} else {
		// Attach the commit to the master branch.
		pro.Reference.Object.SHA = newCommit.SHA
		_, _, err = c.client.Git.UpdateRef(
			context.Background(),
			pro.Organization,
			pro.RepositoryName,
			pro.Reference,
			false)
	}

	return err
}

// createPR creates a pull request. Based on: https://godoc.org/github.com/google/go-github/github#example-PullRequestsService-Create
// Also, add reviewers to created PR
func (c *Client) createPR(pro *types.PullRequestOptions) (err error) {
	reviewers := &types.Assignees{}

	// we shouldn't add commit author name to reviewers list before assign
	if pro.AuthorName == "Yevhen Lebid" {
		reviewers.Add("babinin")
	} else {
		reviewers.Add("ealebed")
	}

	newPR := &github.NewPullRequest{
		Title:               &pro.PRSubject,
		Head:                &pro.CommitBranch,
		Base:                github.String("master"),
		Body:                &pro.PRDescription,
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := c.client.PullRequests.Create(
		context.Background(),
		pro.Organization,
		pro.RepositoryName,
		newPR)
	if err != nil {
		fmt.Printf("Unable to create a PR: %s", err)
	}
	fmt.Printf("PR successfully created: %s\n", pr.GetHTMLURL())

	_, _, error := c.client.PullRequests.RequestReviewers(
		context.Background(),
		pro.Organization,
		pro.RepositoryName,
		pr.GetNumber(),
		github.ReviewersRequest{
			Reviewers: reviewers.List(),
		})
	if error != nil {
		fmt.Printf("Unable to add reviewers to created PR: %s", err)
	}

	fmt.Printf("Reviewers %v successfully added to PR!", reviewers.List())

	return nil
}
