./bin/atlassian_collector  -fd "2024-01-07" -td "2024-01-13" -ju "your_email@example.com" -jt "your_jira_api_token" -bu "your_email@example.com" -bt "your_bitbucket_api_token" -ad "your_domain.atlassian.net" -ts "Closed, Done" -o "."
./bin/atlassian_collector  -fd "2024-01-14" -td "2024-01-20" -ju "your_email@example.com" -jt "your_jira_api_token" -bu "your_email@example.com" -bt "your_bitbucket_api_token" -ad "your_domain.atlassian.net" -ts "Closed, Done" -o "."

./bin/atlassian_reporter -i "." -o "."
