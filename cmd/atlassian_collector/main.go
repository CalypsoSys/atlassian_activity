package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/hashicorp/go-cleanhttp"
)

type activity struct {
	tickets      map[string]map[string]bool
	pullRequests map[string]map[string]bool
	commits      map[string]int
}

const (
	assigned  = "Assigned"
	commented = "Commented"
	inactive  = "Assigned Inactive"
	other     = "Other"
)

var (
	atlasianDomain    string
	jiraUserName      string
	jiraToken         string
	bitbucketUserName string
	bitbucketToken    string
	trackStatusIn     []string
	trackStatus       = map[string]bool{}
	fromDate          time.Time
	toDate            time.Time
	outputFolder      string

	client     *http.Client
	activities = map[int]*activity{}

	assignmentStatuses = []string{assigned, commented, inactive, other}

	totalTickets         int
	totalAssigneeChanges int
	totalStatusChanges   int
	totalComments        int
	totalWorkSpaces      int
	totalRepos           int
	totalPRs             int
	totalCommits         int

	errors = []string{}
)

func init() {
	var configFile string
	flag.StringVar(&configFile, "cf", "", "Config Input file")
	flag.StringVar(&configFile, "config-file", "", "Config Input file")

	flag.StringVar(&atlasianDomain, "ad", "", "Atlasian Domain")
	flag.StringVar(&atlasianDomain, "atlassian-domain", "", "Atlasian Domain")

	flag.StringVar(&jiraUserName, "ju", "", "Jira User Name")
	flag.StringVar(&jiraUserName, "jira-username", "", "Jira User Name")
	flag.StringVar(&jiraToken, "jt", "", "Jira Authentication Token")
	flag.StringVar(&jiraToken, "jira-token", "", "Jira Authentication Token")
	flag.StringVar(&bitbucketUserName, "bu", "", "Bitbucket User Name")
	flag.StringVar(&bitbucketUserName, "bitbucket-username", "", "Bitbucket User Name")
	flag.StringVar(&bitbucketToken, "bt", "", "Bitbucket Authentication Token")
	flag.StringVar(&bitbucketToken, "bitbucket-token", "", "Bitbucket Authentication Token")

	var fromDateStr, toDateStr, trackStatusStr string
	flag.StringVar(&fromDateStr, "fd", "", "From Date")
	flag.StringVar(&fromDateStr, "from-date", "", "From Date")
	flag.StringVar(&toDateStr, "td", "", "To Date")
	flag.StringVar(&toDateStr, "to-date", "", "To Date")
	flag.StringVar(&trackStatusStr, "ts", "", "Jira Status Codes to track")
	flag.StringVar(&trackStatusStr, "track-status", "", "Jira Status Codes to track")
	flag.StringVar(&outputFolder, "o", "", "Output file(s) location")
	flag.StringVar(&outputFolder, "output", "", "Output file(s) location")
	flag.Parse()

	useParams := false
	if configFile == "" {
		for _, check := range []string{atlasianDomain, jiraUserName, jiraToken, bitbucketUserName, bitbucketToken, fromDateStr, toDateStr, trackStatusStr, outputFolder} {
			if check != "" {
				useParams = true
			}
		}
	}

	if !useParams {
		if configFile == "" {
			configFile = "config.ini"
		}
		_, err := os.Stat(configFile)
		if os.IsNotExist(err) {
			logError("No config file found")
			os.Exit(1)
		}
		loadFromConfigFile(configFile)
	} else {
		if fromDateStr != "" {
			fromDate = getDate(fromDateStr, false)
		}
		if toDateStr != "" {
			toDate = getDate(toDateStr, true)
		}

		trackStatusIn = csvToMap(trackStatusStr, trackStatus)
	}

	validateParameters()
}

func main() {
	client = cleanhttp.DefaultClient()
	client.Timeout = 30 * time.Second // Set the timeout to 30 seconds
	if jiraUserName != "" && jiraToken != "" {
		loadIssues()
	}
	if bitbucketUserName != "" && bitbucketToken != "" {
		bitbucketActivity()
	}

	results := gatherResults()
	name := path.Join(outputFolder, strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(results.Title, "-", ""), " ", "_")))

	data, err := json.Marshal(results)
	if err == nil {
		os.WriteFile(name+".json", data, 0666)
	} else {
		logError("Count not marshall json output %v", err)
	}
}
