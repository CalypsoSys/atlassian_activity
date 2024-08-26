package main

import (
	"atlassian_activity/internal/results"
	"fmt"
	"time"
)

func gatherResults() *results.Report {
	output := results.Report{}
	output.Title = fmt.Sprintf("Atlassian Work Report %s to %s", fromDate.Format("2006-01-02"), toDate.Add(-24*time.Hour).Format("2006-01-02"))
	output.Domain = atlasianDomain
	output.FromDate = fromDate
	output.ToDate = toDate

	output.TrackStatus = trackStatusIn

	output.Errors = errors
	output.TotalErrors = len(errors)

	output.TotalTickets = totalTickets
	output.TotalAssigneeChanges = totalAssigneeChanges
	output.TotalStatusChanges = totalStatusChanges
	output.TotalComments = totalComments
	output.TotalWorkSpaces = totalWorkSpaces
	output.TotalRepos = totalRepos
	output.TotalPRs = totalPRs
	output.TotalCommits = totalCommits

	output.Users = []*results.User{}
	for _, user := range allUsers {
		outputUser := results.User{}
		outputUser.UserKey = user.getUserIdentifier()
		outputUser.EmailAddress = user.emailAddress
		outputUser.DisplayName = user.displayName
		outputUser.NickName = user.nickName
		outputUser.AccountID = user.accountID
		outputUser.OtherIDs = user.otherIdentifer()

		userActivities := activities[user.key]
		if userActivities == nil {
			logError("No activities found for user %s", user.getUserIdentifier())
		}

		outputUser.AssignedIssues = map[string][]string{}
		outputUser.CommentedIssues = map[string][]string{}
		outputUser.AssignedInactiveIssues = map[string][]string{}
		outputUser.OtherIssues = map[string][]string{}
		outputUser.TicketsByStatus = map[string]int{}
		if userActivities != nil {
			outputUser.TotalTickets = len(userActivities.tickets)
		}

		issueProcessed := map[string]bool{}
		for _, role := range assignmentStatuses {
			var issues map[string][]string
			switch role {
			case assigned:
				issues = outputUser.AssignedIssues
			case commented:
				issues = outputUser.CommentedIssues
			case inactive:
				issues = outputUser.AssignedInactiveIssues
			default:
				issues = outputUser.OtherIssues
			}

			if userActivities != nil {
				for issue, assigned := range userActivities.tickets {
					if _, exists := issueProcessed[issue]; !exists {
						if _, exists := assigned[role]; exists || role == other {
							statuses := []string{}
							for activity := range assigned {
								outputUser.TicketsByStatus[activity]++
								statuses = append(statuses, activity)
							}
							issues[issue] = statuses
							issueProcessed[issue] = true
						}
					}
				}
			}
		}

		outputUser.PullRequests = map[string][]string{}
		if userActivities != nil {
			for repo, prs := range userActivities.pullRequests {
				outputUser.TotalPullRequests += len(prs)
				outputUser.PullRequests[repo] = []string{}
				for pr := range prs {
					outputUser.PullRequests[repo] = append(outputUser.PullRequests[repo], pr)
				}
			}
		}

		outputUser.Commits = map[string]int{}
		if userActivities != nil {
			for repo, commits := range userActivities.commits {
				outputUser.Commits[repo] = commits
				outputUser.TotalCommits += commits
			}
		}

		output.Users = append(output.Users, &outputUser)
	}

	return &output
}
