package results

import "time"

type User struct {
	UserKey                string              `json:"user_key"`
	EmailAddress           string              `json:"email_address"`
	DisplayName            string              `json:"display_name"`
	NickName               string              `json:"nick_name"`
	AccountID              string              `json:"account_id"`
	OtherIDs               string              `json:"other_identifiers"`
	AssignedIssues         map[string][]string `json:"assigned_issues"`
	CommentedIssues        map[string][]string `json:"commented_issues"`
	AssignedInactiveIssues map[string][]string `json:"assigned_inactive_issues"`
	OtherIssues            map[string][]string `json:"other_issues"`
	PullRequests           map[string][]string `json:"pull_requests"`
	Commits                map[string]int      `json:"commits"`
	TotalTickets           int                 `json:"total_tickets"`
	TicketsByStatus        map[string]int      `json:"tickets_by_status"`
	TotalPullRequests      int                 `json:"total_pull_requests"`
	TotalCommits           int                 `json:"total_commits"`
}

type Report struct {
	Title                string    `json:"title"`
	Domain               string    `json:"domain"`
	FromDate             time.Time `json:"from_date"`
	ToDate               time.Time `json:"to_date"`
	TrackStatus          []string  `json:"track_status"`
	Errors               []string  `json:"errors"`
	TotalErrors          int       `json:"total_errors"`
	TotalTickets         int       `json:"total_tickets"`
	TotalAssigneeChanges int       `json:"total_assignee_changes"`
	TotalStatusChanges   int       `json:"total_status_changes"`
	TotalComments        int       `json:"total_comments"`
	TotalWorkSpaces      int       `json:"total_workSpaces"`
	TotalRepos           int       `json:"total_repositoies"`
	TotalPRs             int       `json:"total_pull_requests"`
	TotalCommits         int       `json:"total_commits"`

	Users []*User `json:"users"`
}
