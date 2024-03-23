package main

import (
	"atlassian_activity/internal/results"
	"bufio"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
)

func outputTextFile(name string, report *results.Report) {
	file, _ := os.Create(name)
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString(fmt.Sprintln(report.Title))
	writer.WriteString(fmt.Sprintf("Domain: %s\n", report.Domain))
	writer.WriteString(fmt.Sprintf("Track Status: %s\n\n", strings.Join(report.TrackStatus, ",")))

	writer.WriteString(fmt.Sprintf("Total Errors: %d\n", report.TotalErrors))
	if !brief {
		for _, err := range report.Errors {
			writer.WriteString(fmt.Sprintf("\tError: %s\n", err))
		}
	}

	writer.WriteString(fmt.Sprintf("Total Tickets: %d\n", report.TotalTickets))
	writer.WriteString(fmt.Sprintf("Total AssigneeChanges: %d\n", report.TotalAssigneeChanges))
	writer.WriteString(fmt.Sprintf("Total StatusChanges: %d\n", report.TotalStatusChanges))
	writer.WriteString(fmt.Sprintf("Total Comments: %d\n", report.TotalComments))
	writer.WriteString(fmt.Sprintf("Total WorkSpaces: %d\n", report.TotalWorkSpaces))
	writer.WriteString(fmt.Sprintf("Total Repositories: %d\n", report.TotalRepos))
	writer.WriteString(fmt.Sprintf("Total Pull Requests: %d\n", report.TotalPRs))
	writer.WriteString(fmt.Sprintf("Total Commits: %d\n\n", report.TotalCommits))

	sort.Slice(report.Users, func(i, j int) bool {
		return report.Users[i].UserKey < report.Users[j].UserKey
	})
	for _, user := range report.Users {
		if userFilter != "" && !checkUser(user) {
			continue
		}

		writer.WriteString(fmt.Sprintln(user.UserKey))
		writer.WriteString(fmt.Sprintf("\tEmail Address: %s\n", user.EmailAddress))
		writer.WriteString(fmt.Sprintf("\tDisplay Name: %s\n", user.DisplayName))
		writer.WriteString(fmt.Sprintf("\tNickname: %s\n", user.NickName))
		writer.WriteString(fmt.Sprintf("\tAccount ID: %s\n", user.AccountID))
		writer.WriteString(fmt.Sprintf("\tOther Identifiers: %s\n", user.OtherIDs))

		if !brief {
			outputTextIssues("Assigned", user.AssignedIssues, writer)
			outputTextIssues("Commented", user.CommentedIssues, writer)
			outputTextIssues("Assigned Inactive", user.AssignedInactiveIssues, writer)
			outputTextIssues("Other", user.OtherIssues, writer)
		}

		for _, repo := range sortedKeys(user.PullRequests) {
			prs := user.PullRequests[repo]
			writer.WriteString(fmt.Sprintf("\tPull Requests for: %s (%d)\n", repo, len(prs)))
			if !brief {
				sort.Strings(prs)
				for _, pr := range prs {
					writer.WriteString(fmt.Sprintf("\t\t%s\n", pr))
				}
			}
		}

		for _, repo := range sortedKeys(user.Commits) {
			commits := user.Commits[repo]
			writer.WriteString(fmt.Sprintf("\tCommits for: %s (%d)\n", repo, commits))
		}

		writer.WriteString(fmt.Sprintf("\tTotal number of Tickets: %d\n", user.TotalTickets))
		for _, role := range sortedKeys(user.TicketsByStatus) {
			writer.WriteString(fmt.Sprintf("\tNumber of Tickets in %s: %d\n", role, user.TicketsByStatus[role]))
		}
		writer.WriteString(fmt.Sprintf("\tTotal number of PRs: %d\n", user.TotalPullRequests))
		writer.WriteString(fmt.Sprintf("\tTotal number of Commits: %d\n", user.TotalCommits))

		writer.WriteString(fmt.Sprintln())
		writer.WriteString(fmt.Sprintln())
		writer.Flush()
	}
	writer.Flush()
}

func outputTextIssues(title string, issues map[string][]string, writer *bufio.Writer) {
	writer.WriteString(fmt.Sprintf("\t%s Issues\n", title))
	for _, issue := range sortedKeys(issues) {
		sort.Strings(issues[issue])
		writer.WriteString(fmt.Sprintf("\t\t%s (%s)\n", issue, strings.Join(issues[issue], ",")))
	}
}

func sortedKeys(m interface{}) []string {
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		fmt.Printf("input is not a map, it is a %s", v.Kind())
		return nil
	}

	keys := make([]string, v.Len())
	for i, key := range v.MapKeys() {
		if key.Kind() != reflect.String {
			fmt.Printf("map key is not a string, it is a %s", key.Kind())
			return nil
		}
		keys[i] = key.String()
	}

	sort.Strings(keys)
	return keys
}
