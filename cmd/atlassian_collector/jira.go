package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
)

// Define a struct to match the JSON structure of the response for easier parsing.
type JiraResponse struct {
	Total  int `json:"total"`
	Issues []struct {
		Key    string `json:"key"`
		Fields struct {
			Summary  string `json:"summary"`
			Assignee struct {
				AccountID    string `json:"accountId"`
				DisplayName  string `json:"displayName"`
				EmailAddress string `json:"emailAddress"`
				// Add more fields as needed
			} `json:"assignee"`
			Created string `json:"created"`
			Updated string `json:"updated"`
		} `json:"fields"`
	} `json:"issues"`
}

// Define structures to model the JSON response for issue changelog
type IssueChangelogResponse struct {
	Changelog struct {
		Histories []struct {
			Created string `json:"created"`
			Author  struct {
				AccountID    string `json:"accountId"`
				EmailAddress string `json:"emailAddress"`
				DisplayName  string `json:"displayName"`
			} `json:"author"`
			Items []struct {
				Field      string `json:"field"`
				From       string `json:"from"`       // Account ID
				FromString string `json:"fromString"` // Display Name
				ToString   string `json:"toString"`   // Display Name
				To         string `json:"to"`         // Account ID
			} `json:"items"`
		} `json:"histories"`
	} `json:"changelog"`
}

// Define a struct to model the JSON structure of comments, including the creation datetime
type JiraCommentResponse struct {
	Comments []struct {
		Author struct {
			AccountID    string `json:"accountId"`
			EmailAddress string `json:"emailAddress"`
			DisplayName  string `json:"displayName"`
		} `json:"author"`
		//Body    string `json:"body"`
		Created string `json:"created"` // Added to capture the datetime the comment was created
	} `json:"comments"`
}

func loadIssues() {
	startAt := 0
	maxResults := 50 // Adjust as needed//

	//jql := url.QueryEscape(fmt.Sprintf("updated >= '%s' AND updated <= '%s'", fromDate.Format("2006-01-02"), toDate.Format("2006-01-02")))
	jql := url.QueryEscape(fmt.Sprintf("updated >= '%s'", fromDate.Format("2006-01-02")))
	for {
		url := fmt.Sprintf("https://%s/rest/api/3/search?jql=%s&startAt=%d&maxResults=%d", atlasianDomain, jql, startAt, maxResults)
		body := callAPI(url, jiraUserName, jiraToken)
		if body == nil {
			logError("not data returned loadIssues")
			return
		}

		// Parse the JSON response
		var jiraResp JiraResponse
		if err := json.Unmarshal(body, &jiraResp); err != nil {
			logError("Error parsing JSON: %v", err)
			logError(string(body))
			return
		}

		// Process each issue
		for _, issue := range jiraResp.Issues {
			totalTickets++

			thisStatus := inactive
			if wasInRangeWeek(issue.Fields.Created) == 0 {
				thisStatus = assigned
			}
			userId := addFindUser(issue.Fields.Assignee.AccountID, issue.Fields.Assignee.DisplayName, issue.Fields.Assignee.EmailAddress, "")
			if thisStatus == assigned || wasInRangeWeek(issue.Fields.Updated) == 0 {
				addTicketActivity(userId, issue.Key, issue.Fields.Summary, thisStatus)
			}
			getChangeLog(issue.Key, issue.Fields.Summary)
			getComments(issue.Key, issue.Fields.Summary)
		}

		// Check if there are more issues to fetch
		startAt += len(jiraResp.Issues)
		log.Printf("On %d of %d", startAt, jiraResp.Total)
		if startAt >= jiraResp.Total {
			break
		}
	}

	log.Printf("Completed fetching all issues.")
}

func getChangeLog(issueKey string, summary string) {
	// Prepare the URL and request
	url := fmt.Sprintf("https://%s/rest/api/3/issue/%s?expand=changelog", atlasianDomain, issueKey)
	body := callAPI(url, jiraUserName, jiraToken)
	if body == nil {
		logError("not data returned getChangeLog")
		return
	}

	var response IssueChangelogResponse
	if err := json.Unmarshal(body, &response); err != nil {
		logError("Error parsing JSON: %v", err)
		return
	}

	// Process the changelog
	for _, history := range response.Changelog.Histories {
		if wasInRangeWeek(history.Created) == 0 {
			for _, item := range history.Items {
				if item.Field == "assignee" { // Change this to your field of interest
					totalAssigneeChanges++

					userId := addFindUser(item.From, item.FromString, "", "")
					addTicketActivity(userId, issueKey, summary, assigned)

					userId = addFindUser(item.To, item.ToString, "", "")
					addTicketActivity(userId, issueKey, summary, assigned)

				} else if item.Field == "status" {
					totalStatusChanges++
					_, exists1 := trackStatus[strings.ToLower(item.ToString)]
					_, exists2 := trackStatus[strings.ToLower(item.FromString)]
					if exists1 || exists2 {
						userId := addFindUser(history.Author.AccountID, history.Author.DisplayName, history.Author.EmailAddress, "")
						if exists1 {
							addTicketActivity(userId, issueKey, summary, item.ToString)
						}
						if exists2 {
							addTicketActivity(userId, issueKey, summary, item.FromString)
						}
					}
				}
			}
		}
	}
}

func getComments(issueKey string, summary string) {
	url := fmt.Sprintf("https://%s/rest/api/3/issue/%s/comment", atlasianDomain, issueKey)

	body := callAPI(url, jiraUserName, jiraToken)
	if body == nil {
		logError("not data returned getComments")
		return
	}

	var comments JiraCommentResponse
	if err := json.Unmarshal(body, &comments); err != nil {
		logError("Error parsing JSON: %v", err)
		logError(string(body))
		return
	}

	for _, comment := range comments.Comments {
		//fmt.Printf("Author: %s (%s)\nComment: %s\nCreated: %s\n\n", comment.Author.DisplayName, comment.Author.EmailAddress, comment.Body, comment.Created)
		if wasInRangeWeek(comment.Created) == 0 {
			totalComments++
			userId := addFindUser(comment.Author.AccountID, comment.Author.DisplayName, comment.Author.EmailAddress, "")
			addTicketActivity(userId, issueKey, summary, commented)
		}
	}
}
