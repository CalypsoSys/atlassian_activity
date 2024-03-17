package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"time"
)

func callAPI(url string, userName string, token string) []byte {
	var sleepMinutes time.Duration = 10
	maxRetr := 10
	retry := true
	var results []byte
	for i := 2; i <= maxRetr && retry; i++ {
		retry, results = callAPIRaw(url, userName, token)
		if retry {
			log.Printf("%d of %d Retrying %s", i, maxRetr, url)
			log.Printf("Sleeping for %d minutes", sleepMinutes)
			time.Sleep(time.Minute * sleepMinutes)
		} else {
			if i > 2 {
				log.Printf("Retry after for %d*%d minutes successful", i-1, sleepMinutes)
			}
			return results
		}
	}

	logError("retry to %s has been exhausted", url)

	return nil
}

func callAPIRaw(url string, userName string, token string) (bool, []byte) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logError("Error creating request: %v", err)
		return false, nil
	}
	auth := base64.StdEncoding.EncodeToString([]byte(userName + ":" + token))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		logError("Error on request: %v", err)
		return false, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		bodyMessage := ""
		if err == nil {
			bodyMessage = string(body)
		} else {
			logError("Getting repsonse body %v", err)
		}
		if resp.StatusCode == http.StatusTooManyRequests || bodyMessage == "Rate limit for this resource has been exceeded" {
			return true, nil
		}
		logError("Request failed with status code: %d", resp.StatusCode)
		logError("Response body: %s", string(body))
		return false, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logError("Error reading response body: %v", err)
		return false, nil
	}

	return false, body

}

func getUserActivity(userId int) *activity {
	if userId == -1 {
		return nil
	}

	if _, exists := activities[userId]; !exists {
		activities[userId] = &activity{
			tickets:      map[string]map[string]bool{},
			commits:      map[string]int{},
			pullRequests: map[string]map[string]bool{},
		}
	}

	return activities[userId]
}

func addTicketActivity(userId int, key string, summary string, assigned string) {
	if activity := getUserActivity(userId); activity != nil {
		issue := fmt.Sprintf("%s - %s", key, summary)
		if _, exists := activity.tickets[issue]; !exists {
			activity.tickets[issue] = map[string]bool{}
		}

		if _, exists := activity.tickets[issue][assigned]; !exists {
			activity.tickets[issue][assigned] = true
		}
	}
}

func addPullRequest(userId int, repo string, title string) {
	if activity := getUserActivity(userId); activity != nil {
		if _, exists := activity.pullRequests[repo]; !exists {
			activity.pullRequests[repo] = map[string]bool{}
		}
		exists := false
		orgTitle := title
		for i := 2; !exists; i++ {
			if _, exists = activity.pullRequests[repo][title]; !exists {
				activity.pullRequests[repo][title] = true
				exists = true
			}
			title = fmt.Sprintf("%s (%d)", orgTitle, i)
		}
	}
}

func addCommit(userId int, repo string) {
	if activity := getUserActivity(userId); activity != nil {
		if _, exists := activity.commits[repo]; !exists {
			activity.commits[repo] = 1
		} else {
			activity.commits[repo]++
		}
	}
}

func wasInRangeWeek(datetimeStr string) int {
	// Parse the string into a time.Time object
	// Note: The layout string is a reference time in Go's specific format, used to interpret the given datetime string
	layout := "2006-01-02T15:04:05.999-0700"
	datetime, err := time.Parse(layout, datetimeStr)
	if err != nil {
		datetime, err = time.Parse(time.RFC3339, datetimeStr)
		if err != nil {
			logError("Error parsing datetime: %v", err)
			return -1
		}
	}

	if datetime.After(toDate) {
		return 1
	} else if datetime.Before(fromDate) {
		return -1
	}

	return 0
}

func sortMapStringMapStringBool(maap map[string]map[string]bool) []string {
	// Extract keys to a slice
	keys := make([]string, 0, len(maap))
	for k := range maap {
		keys = append(keys, k)
	}
	// Sort keys alphabetically
	sort.Strings(keys)

	return keys
}

func sortMapStringBool(maap map[string]bool) []string {
	// Extract keys to a slice
	keys := make([]string, 0, len(maap))
	for k := range maap {
		keys = append(keys, k)
	}
	// Sort keys alphabetically
	sort.Strings(keys)

	return keys
}

func logError(format string, a ...any) {
	err := fmt.Sprintf(format, a...)
	errors = append(errors, err)
	log.Print(err)
}
