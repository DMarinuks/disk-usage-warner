# Disk Usage Warner
It will send a warning email when disk usage is above your desired percentage.

## Usage
The email will be sent only when percentage is specified.
```bash
go run ./cmd/duw check -h
go run ./cmd/duw check --verbose
go run ./cmd/duw check --path=/ --percentage 7
```
### Using env file
Create `.env` file in the root folder, disk usage warner will load `.env` file automatically.

| Variable      | Description | Example |
| ----------- | ----------- | ----------- |
| DUW_LOG_LEVEL | Set log level | warn |
| DUW_PATHS | Paths to check, separated by comma | /foo,/bar |
| DUW_PERCENTAGE | At what usage percentage an email will be sent | 80 |
| DUW_ADMINS | List of admin emails that should get the warning email. Separated by comma. | admin1@bar.com,admin2@foo.com |

The DUW_MAIL_ variables are self explanatory.