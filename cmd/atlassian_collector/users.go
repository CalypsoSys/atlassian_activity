package main

import (
	"regexp"
	"sort"
	"strings"
)

type user struct {
	key              int
	accountID        string
	displayName      string
	emailAddress     string
	nickName         string
	otherIdentifiers map[string]bool
}

var (
	allUsers = []*user{}
)

func sortUsers() {
	// Sort the slice
	sort.Slice(allUsers, func(i, j int) bool {
		return allUsers[i].getUserIdentifier() < allUsers[j].getUserIdentifier()
	})
}

func (user *user) getUserIdentifier() string {
	if user.emailAddress != "" {
		return user.emailAddress
	} else if user.displayName != "" {
		return user.displayName
	} else if user.nickName != "" {
		return user.nickName
	} else if user.accountID != "" {
		return user.accountID
	}

	return user.otherIdentifer()
}

func (user *user) otherIdentifer() string {
	other := []string{}
	for k := range user.otherIdentifiers {
		other = append(other, k)
	}

	sort.Strings(other)
	return strings.Join(other, ",")
}

func addFindUser(accountID string, displayName string, emailAddressRaw string, nickName string) int {
	if accountID == "" && displayName == "" && emailAddressRaw == "" && nickName == "" {
		return -1
	}

	re := regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b`)
	matches := re.FindAllString(emailAddressRaw, -1)

	emailAddress := ""
	otherIdentifers := map[string]bool{}
	for _, match := range matches {
		if emailAddress == "" {
			emailAddress = match
		} else if emailAddress != match {
			otherIdentifers[match] = true
			logError("bad email match %s vs %s", emailAddress, match)
		}

		emailAddressRaw = strings.ReplaceAll(emailAddressRaw, match, "")
	}

	emailAddressRaw = strings.TrimSpace(strings.ReplaceAll(emailAddressRaw, "<>", ""))
	if emailAddressRaw != "" {
		otherIdentifers[emailAddressRaw] = true
	}

	for _, user := range allUsers {
		if compareNotEmpty(accountID, user.accountID) || compareNotEmpty(displayName, user.displayName) ||
			compareNotEmpty(emailAddress, user.emailAddress) || compareNotEmpty(nickName, user.nickName) {
			user.enhanceUser(accountID, displayName, emailAddress, nickName, otherIdentifers)
			return user.key
		}
	}

	user := user{
		key:              len(allUsers),
		accountID:        accountID,
		displayName:      displayName,
		emailAddress:     emailAddress,
		nickName:         nickName,
		otherIdentifiers: otherIdentifers,
	}

	allUsers = append(allUsers, &user)

	return user.key
}

func (user *user) enhanceUser(accountID string, displayName string, emailAddress string, nickName string, otherIdentifers map[string]bool) {
	user.accountID = user.checkSetUser(user.accountID, accountID)
	user.displayName = user.checkSetUser(user.displayName, displayName)
	user.emailAddress = user.checkSetUser(user.emailAddress, emailAddress)
	user.nickName = user.checkSetUser(user.nickName, nickName)

	for k := range otherIdentifers {
		if _, exists := user.otherIdentifiers[k]; !exists {
			user.otherIdentifiers[k] = true
		}
	}
}

func (user *user) checkSetUser(set string, check string) string {
	if set == "" {
		return check
	}
	if check != "" && set != check {
		user.otherIdentifiers[check] = true
		logError("Mismatch %s vs %s", set, check)
	}

	return set
}

func compareNotEmpty(string1 string, string2 string) bool {
	if string1 == "" || string2 == "" {
		return false
	}

	return strings.EqualFold(string1, string2)
}
