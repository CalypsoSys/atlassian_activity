package main

import (
	"os"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

func loadFromConfigFile(configFile string) {
	cfg, err := ini.Load(configFile)
	if err != nil {
		logError("Fail to read %s file: %v", configFile, err)
		os.Exit(1)
	}

	jiraUserName = cfg.Section("Credentials").Key("jira-username").String()
	jiraToken = cfg.Section("Credentials").Key("jira-token").String()
	bitbucketUserName = cfg.Section("Credentials").Key("bitbucket-username").String()
	bitbucketToken = cfg.Section("Credentials").Key("bitbucket-token").String()
	atlasianDomain = cfg.Section("Settings").Key("atlassian-domain").String()

	fromDate = getDate(cfg.Section("Settings").Key("from-date").String(), false)
	toDate = getDate(cfg.Section("Settings").Key("to-date").String(), true)

	trackStatusIn = csvToMap(cfg.Section("Settings").Key("track-status").String(), trackStatus)
	outputFolder = cfg.Section("Settings").Key("output").String()
}

func getDate(strDate string, toDate bool) time.Time {
	date, err := time.Parse("2006-01-02", strDate)
	if err != nil {
		var typ string
		if toDate {
			typ = "to"
		} else {
			typ = "from"
		}
		logError("Failed to parse %s date: %v", typ, err)
		os.Exit(1)
	}
	if toDate {
		date = date.Add(24 * time.Hour)
	}

	return date
}

func csvToMap(csv string, mapp map[string]bool) []string {
	original := []string{}
	for _, token := range strings.Split(csv, ",") {
		token = strings.TrimSpace(token)
		original = append(original, token)
		token := strings.ToLower(token)
		mapp[token] = true
	}

	return original
}

func validateParameters() {
	if atlasianDomain == "" {
		logError("You must supply a Atlassian domain")
		os.Exit(1)
	}

	if (jiraUserName != "" && jiraToken == "") || (jiraUserName == "" && jiraToken != "") {
		logError("You must supply a Jira username and token")
		os.Exit(1)
	}

	if (bitbucketUserName != "" && bitbucketToken == "") || (bitbucketUserName == "" && bitbucketToken != "") {
		logError("You must supply a Bitbucket username and token")
		os.Exit(1)
	}

	if jiraUserName == "" && jiraToken == "" && bitbucketUserName == "" && bitbucketToken == "" {
		logError("No credentials supplied for Jira or Bitbucket")
		os.Exit(1)
	}

	if fromDate.After(toDate) {
		logError("From date is after To date")
		os.Exit(1)
	}
}
