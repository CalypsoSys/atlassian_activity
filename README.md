# atlassian_activity
Program gets and correlates activity in Atlassian Jira and Bitbucket for a period of time 


# Tool Usage Documentation
## Overview
This tool is designed to gather and report on work items from Atlassian products, specifically Jira and Bitbucket, within a specified date range. Users can specify various options either through command-line flags or a configuration file to customize the tool's behavior according to their needs. The program `atlassian_activity` creates JSON data files. ``atlassian_reporter` reads JSON data created by `atlassian_activity` and outputs text files.

# Command-Line Options for `atlassian_activity`
The tool supports the following command-line options:

* `-ad` or `--atlassian-domain`: Sets the Atlassian domain to query.
* `-ju` or `--jira-username`: The username for Jira authentication.
* `-jt` or `--jira-token`: The authentication token for Jira.
* `-bu` or `--bitbucket-username`: The username for Bitbucket authentication.
* `-bt` or `--bitbucket-token`: The authentication token for Bitbucket.
* `-fd` or `--from-date`: The start date for the report (inclusive).
* `-td` or `--to-date`: The end date for the report (inclusive).
* `-ts` or `--track-status`: Comma-separated Jira status codes to track.
* `-o` or `--output`: The location for the output file(s).

## Configuration File
Alternatively, you can use a configuration file to specify the options. This is useful for avoiding the need to enter sensitive information directly into the command line or to reuse configurations. Here is an example of the configuration file structure:
* `-cf` or `--config-file`: Specifies the path to the configuration input file.

## ini
```
[Credentials]
jira-username = your_email@example.com
jira-token = your_jira_api_token
bitbucket-username = your_email@example.com
bitbucket-token = your_bitbucket_api_token

[Settings]
atlassian-domain = your_domain.atlassian.net
from-date = 2023-04-01
to-date = 2023-04-07
track-status = Closed,Done
output = .
```

* `[Credentials]`: This section includes authentication information for Jira and Bitbucket.
    * `jira-username` and `jira-token`: Credentials for Jira.
    * `bitbucket-username` and `bitbucket-token`: Credentials for Bitbucket.
* `[Settings]`: This section defines operational parameters for the tool.
    * `atlassian-domain`: Your Atlassian domain.
    * `from-date` and `to-date`: Define the report's time range.
    * `track-status`: Status codes changes in Jira issues to include in the report.
    * `output`: The directory where output files will be saved.

## Usage Example
To run the tool with command-line options:

```
./atlassian_activity  -fd "2024-01-07" -td "2024-01-13" -ju "your_email@example.com" -jt "your_jira_api_token" -bu "your_email@example.com" -bt "your_bitbucket_api_token" -ad "your_domain.atlassian.net" -ts "Closed, Done" -o "."
```

Or using the configuration file only, ensure the file is properly set up as described above and then run:

```
atlassian_activity --config-file="config.ini"
```
This documentation should guide you in configuring and using the tool effectively for your needs. If you encounter any issues or have further questions, please refer to the tool's help command or consult the detailed user guide.

## Sample JSON output `atlassian_work_report_20240303_to_20240309.json`
```
{
    "title": "Atlassian Work Report 2024-03-03 to 2024-03-09",
    "domain": "yourcompany.atlassian.net",
    "from_date": "2024-03-03T00:00:00Z",
    "to_date": "2024-03-10T00:00:00Z",
    "track_status": [
        "Pending Code Review",
        "Fixed",
        "Closed",
        "Unresolved",
        "Done"
    ],
    "errors": [],
    "total_errors": 0,
    "total_tickets": 670,
    "total_assignee_changes": 435,
    "total_status_changes": 915,
    "total_comments": 767,
    "total_workSpaces": 19,
    "total_repositoies": 80,
    "total_pull_requests": 87,
    "total_commits": 544,
    "users": [
        {
            "user_key": "developer@yourcompany.com",
            "email_address": "developer@yourcompany.com",
            "display_name": "First Last",
            "nick_name": "Name",
            "account_id": "Unique Account Identifier",
            "other_identifiers": "First Last,Name",
            "assigned_issues": {
                "PROJA-11553 - Fix aggregation of account balances": [
                    "Fixed",
                    "Pending Code Review",
                    "Assigned"
                ],
                "PROJA-19984 - Correct aggregation in C++ module to store balances.": [
                    "Assigned Inactive",
                    "Assigned"
                ],
                "PROJA-18699 - Allow enhanced balance calculation": [
                    "Assigned",
                    "Fixed",
                    "Pending Code Review"
                ],
                "PROJA-18687 - Document enhanced balance calculation": [
                    "Commented",
                    "Assigned",
                    "Fixed",
                    "Pending Code Review"
                ],
                "PROJA-18711 - Using the authoritative ID or username for user instances": [
                    "Assigned",
                    "Fixed",
                    "Pending Code Review",
                    "Commented"
                ],
                "PROJA-18712 - Using the authoritative ID or username for customer instances": [
                    "Assigned"
                ],
                "PROJA-18713 -Using the authoritative ID or username for system instances": [
                    "Assigned"
                ],
                "PROJA-18714 - Using the authoritative ID or username for logging instances": [
                    "Assigned"
                ],
                "PROJA-18715 - Using the authoritative ID or username for virtual instances": [
                    "Assigned"
                ],
                "PROJA-18737 - Fallback code for user login recovery": [
                    "Assigned Inactive",
                    "Assigned"
                ],
                "PROJA-18781 - clustering all instances in windows throws a error after initial startup": [
                    "Assigned"
                ],
                "PROJA-18800 - Update system tool for user security": [
                    "Assigned",
                    "Done"
                ],
                "PROJA-18863 - Diagnostics code for Windows creates no diagnostics files": [
                    "Assigned"
                ],
                "PROJB-9955 - system failing to capture login failures": [
                    "Assigned",
                    "Commented"
                ]
            },
            "commented_issues": {
                "PROJA-14086 - Detailed reports fail after row 1000": [
                    "Commented"
                ],
                "PROJA-13715 - Security upgrade for ssh library in production": [
                    "Commented"
                ],
                "PROJA-13764 - Bar charts do not render correctly": [
                    "Commented"
                ],
                "PROJA-18672 - Can we provide a list of combinations that would result in all 6 login statuses?": [
                    "Assigned Inactive",
                    "Done",
                    "Commented"
                ],
                "PROJA-18692 - Have internal mapping numbers to controls": [
                    "Fixed",
                    "Pending Code Review",
                    "Commented"
                ],
                "PROJA-18736 - Balance calculations rounding enhancements": [
                    "Commented"
                ],
                "PROJC-1054 - Bill Matching Scenarios - utilities": [
                    "Commented"
                ],
                "PROJC-1056 - Need 2 logins with same name": [
                    "Commented"
                ],
                "PROJC-1057 - 2 logins with different names": [
                    "Commented"
                ]
            },
            "assigned_inactive_issues": {
                "PROJA-55779 - Using the user id to increase the balance fails in API": [
                    "Assigned Inactive"
                ],
                "PROJA-12237 - User linking multiple account to 1 bill": [
                    "Assigned Inactive"
                ],
                "PROJA-11423 - Creating DB records for bills that are not actual bills": [
                    "Assigned Inactive"
                ],
                "PROJA-13123 - Login name uniqueness - as we redo authentication": [
                    "Assigned Inactive"
                ],
                "PROJA-13928 - Vulnerability - check 3rd party libraries": [
                    "Assigned Inactive"
                ],
                "PROJA-18157 - New shell for root access": [
                    "Assigned Inactive"
                ],
                "PROJA-18169 - create a unique id of an user for log coorelation": [
                    "Assigned Inactive"
                ],
                "PROJA-18542 - system Scan - How many attempts to authenticate failure": [
                    "Assigned Inactive"
                ]
            },
            "other_issues": {
                "PROJA-11057 - UI automation failed": [
                    "Fixed"
                ],
                "PROJA-28699 - Slow rendering takes approximately 120 seconds.": [
                    "Fixed"
                ],
                "PROJA-11161 - Create API doc": [
                    "Fixed"
                ],
                "PROJA-11553 - Aggregation of balances": [
                    "Fixed"
                ],
                "PROJA-11748 - Convert template page to jquery": [
                    "Fixed"
                ],
                "PROJA-11773 - Convert menu selection to jquery": [
                    "Fixed"
                ]
            },
            "pull_requests": {
                "workspace/repo1": [
                    "PROJA-18800 - disable task service"
                ],
                "workspace/repo2": [
                    "PROJA-18863 - fix diagnosis code for solaris",
                    "PROJA-18672 add comment for user state code",
                    "PROJA-18692 - fix navigation on menu with"
                ]
            },
            "commits": {
                "workspace/repo3": 2,
                "workspace/repo1": 2,
                "workspace/repo2": 5
            },
            "total_tickets": 11,
            "tickets_by_status": {
                "Assigned": 14,
                "Assigned Inactive": 11,
                "Commented": 12,
                "Done": 2,
                "Fixed": 5,
                "Pending Code Review": 5
            },
            "total_pull_requests": 4,
            "total_commits": 9
        }
    ]
}
```

# Command-Line Options for `atlassian_reporter`
The tool supports the following command-line options for configuring its operation:

* `-i` or `--input`: Specifies the location of input file(s).
* `-o` or `--output`: Specifies the location for output file(s).
* `-uf` or `--user-filter`: Applies a filter to match specific user fields.
* `-b` or `--brief`: Generates a brief version of the output.
* `-a` or `--all`: Generates all available types of output reports.
* `-s` or `--summary`: Generates a summary report.
* `-us` or `--user-summary`: Generates a summary report focused on user information.
* `-ud` or `--user-detail`: Generates a detailed report containing user information.

## Behavior and Dependencies
### Output Types
The `--summary`, `--user-summary`, and `--user-detail` flags control the types of reports the program generates. If none of these flags are specified, the program will prompt the user to select at least one output type.

The `--all` flag is a convenience option that enables all available report types, overriding individual report type flags.

### User Filtering
The `--user-filter` option allows targeting specific users within the input files. This filter applies to fields associated with user information. The filtering is case-insensitive.

### Input and Output
Input and output folders must be specified for the program to know where to find the input files and where to save the generated reports. If these are not provided, the program may not function as expected.

## Usage Example
To run the tool with command-line options:

```
./atlassian_reporter -i "./output" -o "./output" -a
./bin/atlassian_reporter -i "./output/last_year" -o "./output/last_year" -a -b -uf "@example.com"
```

## Sample output `atlassian_totals.csv`
```
period	tickets	assigned changes	status changes	comments	pull requests	commits	workspaces	repositories
2024-03-03 to 2024-03-10	670	435	915	767	87	544	19	80
```

## Sample output `atlassian_totals_by_user.csv`
```
	2024-03-03 to 2024-03-10
user	Tickets	Assigned	Commented	PRs	Commented
developer@yourcompany.com	11	14	9	4	9
```

## Sample output `atlassian_work_report_20240303_to_20240309.txt`
```
Atlassian Work Report 2024-03-03 to 2024-03-09
Domain: yourcompany.atlassian.net
Track Status: Pending Code Review,Fixed,Closed,Unresolved,Done

Total Errors: 0
Total Tickets: 670
Total AssigneeChanges: 435
Total StatusChanges: 915
Total Comments: 767
Total WorkSpaces: 19
Total Repositories: 80
Total Pull Requests: 87
Total Commits: 544

developer@yourcompany.com
	Email Address: developer@yourcompany.com
	Display Name: First Last
	Nickname: Name
	Account ID: Unique Account Identifier
	Other Identifiers: First Last,Name
	Assigned Issues
		PROJA-11553 - Fix aggregation of account balances (Assigned,Fixed,Pending Code Review)
		PROJA-18687 - Document enhanced balance calculation (Assigned,Commented,Fixed,Pending Code Review)
		PROJA-18699 - Allow enhanced balance calculation (Assigned,Fixed,Pending Code Review)
		PROJA-18711 - Using the authoritative ID or username for user instances (Assigned,Commented,Fixed,Pending Code Review)
		PROJA-18712 - Using the authoritative ID or username for customer instances (Assigned)
		PROJA-18713 -Using the authoritative ID or username for system instances (Assigned)
		PROJA-18714 - Using the authoritative ID or username for logging instances (Assigned)
		PROJA-18715 - Using the authoritative ID or username for virtual instances (Assigned)
		PROJA-18737 - Fallback code for user login recovery (Assigned,Assigned Inactive)
		PROJA-18781 - clustering all instances in windows throws a error after initial startup (Assigned)
		PROJA-18800 - Update system tool for user security (Assigned,Done)
		PROJA-18863 - Diagnostics code for Windows creates no diagnostics files (Assigned)
		PROJA-19984 - Correct aggregation in C++ module to store balances. (Assigned,Assigned Inactive)
		PROJB-9955 - system failing to capture login failures (Assigned,Commented)
	Commented Issues
		PROJA-13715 - Security upgrade for ssh library in production (Commented)
		PROJA-13764 - Bar charts do not render correctly (Commented)
		PROJA-14086 - Detailed reports fail after row 1000 (Commented)
		PROJA-18672 - Can we provide a list of combinations that would result in all 6 login statuses? (Assigned Inactive,Commented,Done)
		PROJA-18692 - Have internal mapping numbers to controls (Commented,Fixed,Pending Code Review)
		PROJA-18736 - Balance calculations rounding enhancements (Commented)
		PROJC-1054 - Bill Matching Scenarios - utilities (Commented)
		PROJC-1056 - Need 2 logins with same name (Commented)
		PROJC-1057 - 2 logins with different names (Commented)
	Assigned Inactive Issues
		PROJA-11423 - Creating DB records for bills that are not actual bills (Assigned Inactive)
		PROJA-12237 - User linking multiple account to 1 bill (Assigned Inactive)
		PROJA-13123 - Login name uniqueness - as we redo authentication (Assigned Inactive)
		PROJA-13928 - Vulnerability - check 3rd party libraries (Assigned Inactive)
		PROJA-18157 - New shell for root access (Assigned Inactive)
		PROJA-18169 - create a unique id of an user for log coorelation (Assigned Inactive)
		PROJA-18542 - system Scan - How many attempts to authenticate failure (Assigned Inactive)
		PROJA-55779 - Using the user id to increase the balance fails in API (Assigned Inactive)
	Other Issues
		PROJA-11057 - UI automation failed (Fixed)
		PROJA-11161 - Create API doc (Fixed)
		PROJA-11553 - Aggregation of balances (Fixed)
		PROJA-11748 - Convert template page to jquery (Fixed)
		PROJA-11773 - Convert menu selection to jquery (Fixed)
		PROJA-28699 - Slow rendering takes approximately 120 seconds. (Fixed)
	Pull Requests for: workspace/repo1 (1)
		PROJA-18800 - disable task service
	Pull Requests for: workspace/repo2 (3)
		PROJA-18672 add comment for user state code
		PROJA-18692 - fix navigation on menu with
		PROJA-18863 - fix diagnosis code for solaris
	Commits for: workspace/repo1 (2)
	Commits for: workspace/repo2 (5)
	Commits for: workspace/repo3 (2)
	Total number of Tickets: 11
	Number of Tickets in Assigned: 14
	Number of Tickets in Assigned Inactive: 11
	Number of Tickets in Commented: 12
	Number of Tickets in Done: 2
	Number of Tickets in Fixed: 5
	Number of Tickets in Pending Code Review: 5
	Total number of PRs: 4
	Total number of Commits: 9

```
