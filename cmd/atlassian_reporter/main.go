package main

import (
	"atlassian_activity/internal/results"
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	inputFolder  string
	outputFolder string
	userFilter   string
	brief        bool
	all          bool
	summary      bool
	userSummary  bool
	userDetails  bool
)

func init() {
	flag.StringVar(&inputFolder, "i", "", "Input file(s) location")
	flag.StringVar(&inputFolder, "input", "", "Input file(s) location")
	flag.StringVar(&outputFolder, "o", "", "Output file(s) location")
	flag.StringVar(&outputFolder, "output", "", "Output file(s) location")
	flag.StringVar(&userFilter, "uf", "", "User filter (match in user fields)")
	flag.StringVar(&userFilter, "user-filter", "", "User filter (match in user fields)")
	flag.BoolVar(&brief, "b", false, "Brief output")
	flag.BoolVar(&brief, "brief", false, "Brief output")
	flag.BoolVar(&all, "a", false, "All output report types")
	flag.BoolVar(&all, "all", false, "All output report types")
	flag.BoolVar(&summary, "s", false, "Summary report")
	flag.BoolVar(&summary, "summary", false, "Summary Report")
	flag.BoolVar(&userSummary, "us", false, "User summary report")
	flag.BoolVar(&userSummary, "user-summary", false, "User summary Report")
	flag.BoolVar(&userDetails, "ud", false, "User detail report")
	flag.BoolVar(&userDetails, "user-detail", false, "User detail report")
	flag.Parse()

	if all {
		summary = true
		userSummary = true
		userDetails = true 
	}

	if !summary && !userSummary && !userDetails {
		fmt.Println("No output type chosen, please select at least 1 output type")
		os.Exit(1)
	}

	userFilter = strings.ToLower(userFilter)
}

func main() {
	var summaryWriter *bufio.Writer
	if summary {
		file, _ := os.Create(path.Join(outputFolder, "atlassian_totals.csv"))
		defer file.Close()

		summaryWriter = bufio.NewWriter(file)
		summaryWriter.WriteString("period\ttickets\tassigned changes\tstatus changes\tcomments\tpull requests\tcommits\tworkspaces\trepositories\n")
	}

	// Walk the directory
	err := filepath.Walk(inputFolder, func(name string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error accessing path:", name, err)
			return err
		}

		if strings.HasPrefix(info.Name(), "atlassian_work_report_") && strings.HasSuffix(info.Name(), ".json") {
			data, err := os.ReadFile(name)
			if err == nil {
				report := results.Report{}
				json.Unmarshal(data, &report)

				period := fmt.Sprintf("%s to %s", report.FromDate.Format("2006-01-02"), report.ToDate.Format("2006-01-02"))
				addPeriod(period)

				if summaryWriter != nil {
					summaryWriter.WriteString(fmt.Sprintf("%s\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n",
						period,
						report.TotalTickets,
						report.TotalAssigneeChanges,
						report.TotalStatusChanges,
						report.TotalComments,
						report.TotalPRs,
						report.TotalCommits,
						report.TotalWorkSpaces,
						report.TotalRepos))
				}

				if userDetails {
					outputTextFile(strings.ReplaceAll(name, ".json", ".txt"), &report)
				}

				if userSummary {
					for _, user := range report.Users {
						if userFilter == "" || checkUser(user) {
							addUser(period, user)
						}
					}
				}
			} else {
				fmt.Printf("error reading %s %v", info.Name(), err)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error walking through directory:", err)
	}

	summaryWriter.Flush()

	if userSummary {
		writeByUsers()
	}
}

func checkUser(user *results.User) bool {
	return strings.Contains(strings.ToLower(user.UserKey), userFilter) ||
		strings.Contains(strings.ToLower(user.AccountID), userFilter) ||
		strings.Contains(strings.ToLower(user.DisplayName), userFilter) ||
		strings.Contains(strings.ToLower(user.EmailAddress), userFilter) ||
		strings.Contains(strings.ToLower(user.NickName), userFilter) ||
		strings.Contains(strings.ToLower(user.OtherIDs), userFilter)
}
