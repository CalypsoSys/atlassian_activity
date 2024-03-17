package main

import (
	"encoding/json"
	"fmt"
	"net/url"

	"log"
)

// Define a struct for parsing the JSON structure of the response
type WorkspacesResponse struct {
	Values []struct {
		Slug string `json:"slug"`
		Name string `json:"name"`
		UUID string `json:"uuid"`
	} `json:"values"`
	Next string `json:"next"` // For pagination
}

// Define a struct to model the JSON structure of the response for listing repositories
type BitbucketRepositoriesResponse struct {
	Values []struct {
		Name        string `json:"name"`
		Slug        string `json:"slug"`
		FullName    string `json:"full_name"`
		Description string `json:"description"`
	} `json:"values"`
	Next string `json:"next"` // For pagination
}

// BitbucketPullRequestsResponse reflects the expected structure of the API response
type BitbucketPullRequestsResponse struct {
	Values []struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		State       string `json:"state"`
		CreatedOn   string `json:"created_on"`
		Author      struct {
			DisplayName string `json:"display_name"`
			AccountID   string `json:"account_id"`
			NickName    string `json:"nickname"`
		} `json:"author"`
	} `json:"values"`
	Next string `json:"next"` // For pagination
}

// Define structures reflecting the part of the Bitbucket API response you're interested in
type BitbucketCommitsResponse struct {
	Values []struct {
		//Hash   string `json:"hash"`
		Author struct {
			Raw  string `json:"raw"`
			User struct {
				DisplayName string `json:"display_name"`
				AccountID   string `json:"account_id"`
				NickName    string `json:"nickname"`
			} `json:"user"`
		} `json:"author"`
		Date string `json:"date"` // Date in RFC3339 format
	} `json:"values"`
	Next string `json:"next"` // For pagination
}

func bitbucketActivity() {
	url := "https://api.bitbucket.org/2.0/workspaces"

	for url != "" {
		body := callAPI(url, bitbucketUserName, bitbucketToken)
		if body == nil {
			logError("not data returned bitbucketActivity")
			return
		}

		var workspaces WorkspacesResponse
		if err := json.Unmarshal(body, &workspaces); err != nil {
			logError("Error parsing JSON: %v", err)
			logError(string(body))
			return
		}

		// Iterate through the workspaces and print details
		for _, workspace := range workspaces.Values {
			totalWorkSpaces++
			//fmt.Printf("UUID: %s\nName: %s\nSlug: %s\n\n", workspace.UUID, workspace.Name, workspace.Slug)
			log.Printf("On Workspace: %s", workspace.Name)
			getRepositories(workspace.Slug)
		}

		url = workspaces.Next
	}
}

func getRepositories(workspace string) {
	url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s", workspace)

	for url != "" {
		body := callAPI(url, bitbucketUserName, bitbucketToken)
		if body == nil {
			logError("not data returned getRepositories")
			return
		}

		var repos BitbucketRepositoriesResponse
		if err := json.Unmarshal(body, &repos); err != nil {
			logError("Error parsing JSON: %v", err)
			logError(string(body))
			return
		}

		// Print out the fetched repositories
		for _, repo := range repos.Values {
			totalRepos++
			//fmt.Printf("Name: %s\nFull Name: %s\nDescription: %s\n\n", repo.Name, repo.FullName, repo.Description)
			log.Printf("On Repository: %s", repo.Name)
			getPullRequests(repo.FullName)
			getCommits(repo.FullName)
		}

		url = repos.Next
	}
}

func getPullRequests(repoName string) {
	query := fmt.Sprintf("created_on >= %s AND created_on <= %s", fromDate.Format("2006-01-02"), toDate.Format("2006-01-02"))
	encodedQuery := url.QueryEscape(query)

	url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/pullrequests?q=%s", repoName, encodedQuery)

	for url != "" {
		body := callAPI(url, bitbucketUserName, bitbucketToken)
		if body == nil {
			logError("not data returned getPullRequests")
			return
		}

		var pullRequests BitbucketPullRequestsResponse
		if err := json.Unmarshal(body, &pullRequests); err != nil {
			logError("Error parsing JSON: %v", err)
			logError(string(body))
			return
		}

		// Print out the fetched repositories
		for _, pr := range pullRequests.Values {
			totalPRs++
			//fmt.Printf("Name: %s\nNick Name: %s\n\n", pr.Author.DisplayName, pr.Author.NickName)
			userId := addFindUser(pr.Author.AccountID, pr.Author.DisplayName, "", pr.Author.NickName)
			addPullRequest(userId, repoName, pr.Title)
		}

		url = pullRequests.Next
	}
}

func getCommits(repoName string) {
	url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/commits", repoName)

	for url != "" {
		body := callAPI(url, bitbucketUserName, bitbucketToken)
		if body == nil {
			logError("not data returned getCommits")
			return
		}

		var commits BitbucketCommitsResponse
		if err := json.Unmarshal(body, &commits); err != nil {
			logError("Error parsing JSON: %v", err)
			logError(string(body))
			return
		}

		added := false
		for _, commit := range commits.Values {
			dateMatch := wasInRangeWeek(commit.Date)
			if dateMatch == 0 {
				totalCommits++
				//fmt.Printf("Name: %s\nFull Name: %s\n\n", commit.Author, commit.Author.Raw)
				userId := addFindUser(commit.Author.User.AccountID, commit.Author.User.DisplayName, commit.Author.Raw, "")
				addCommit(userId, repoName)
				added = true
			} else if dateMatch == 1 {
				added = true
			}
		}

		if !added {
			url = ""
		} else {
			url = commits.Next
		}
	}
}
