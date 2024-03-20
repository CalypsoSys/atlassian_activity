package main

import (
	"atlassian_activity/internal/results"
	"bufio"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
)

type userTotals struct {
	totalTickets      int
	totalAssigned     int
	totalCommented    int
	totalPullRequests int
	totalCommits      int
}

var (
	periods = []string{}
	users   = map[string]map[string]userTotals{}
)

func addPeriod(period string) {
	periods = append(periods, period)
}

func addUser(period string, user *results.User) {
	if _, exists := users[user.UserKey]; !exists {
		users[user.UserKey] = map[string]userTotals{}
	}

	users[user.UserKey][period] = userTotals{
		totalTickets:      user.TotalTickets,
		totalAssigned:     len(user.AssignedIssues),
		totalCommented:    len(user.CommentedIssues),
		totalPullRequests: user.TotalPullRequests,
		totalCommits:      user.TotalCommits,
	}
}

func writeByUsers() {
	file, _ := os.Create(path.Join(outputFolder, "atlassian_totals_by_user.csv"))
	defer file.Close()

	writer := bufio.NewWriter(file)

	sort.Strings(periods)
	writer.WriteString(fmt.Sprintf("\t%s\n", strings.Join(periods, "\t\t\t\t\t")))
	writer.WriteString(fmt.Sprintf("user%s\n", strings.Repeat("\tTickets\tAssigned\tCommented\tPRs\tCommits", len(periods))))

	for _, user := range sortedKeys(users) {
		writer.WriteString(user)
		for _, period := range periods {
			totals := userTotals{}
			if _, exists := users[user][period]; exists {
				totals = users[user][period]
			}
			writer.WriteString(fmt.Sprintf("\t%d\t%d\t%d\t%d\t%d", totals.totalTickets, totals.totalAssigned, totals.totalCommented, totals.totalPullRequests, totals.totalCommits))
		}
		writer.WriteString("\n")
		writer.Flush()
	}
	writer.Flush()
}
