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
	flag.Parse()

	userFilter = strings.ToLower(userFilter)
}

func main() {
	file, _ := os.Create(path.Join(outputFolder, "atlassian_totals.csv"))
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString("period\ttickets\tassigned changes\tstatus changes\tcomments\tpull requests\tcommits\tworkspaces\trepositories\n")

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

				writer.WriteString(fmt.Sprintf("%s\t%d\t%d\t%d\t%d\t%d\t%d\t%d\t%d\n",
					period,
					report.TotalTickets,
					report.TotalAssigneeChanges,
					report.TotalStatusChanges,
					report.TotalComments,
					report.TotalPRs,
					report.TotalCommits,
					report.TotalWorkSpaces,
					report.TotalRepos))

				outputTextFile(strings.ReplaceAll(name, ".json", ".txt"), &report)

				for _, user := range report.Users {
					if userFilter == "" || checkUser(user) {
						addUser(period, user)
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

	writer.Flush()

	writeByUsers()
}

func checkUser(user *results.User) bool {
	return strings.Contains(strings.ToLower(user.UserKey), userFilter) ||
		strings.Contains(strings.ToLower(user.AccountID), userFilter) ||
		strings.Contains(strings.ToLower(user.DisplayName), userFilter) ||
		strings.Contains(strings.ToLower(user.EmailAddress), userFilter) ||
		strings.Contains(strings.ToLower(user.NickName), userFilter) ||
		strings.Contains(strings.ToLower(user.OtherIDs), userFilter)
}
