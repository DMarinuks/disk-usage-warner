# Disk Usage Warner
It will send an email, using `mailer/warning.html`, template when disk usage is above your desired percentage.

## Usage
The email will be sent only when percentage is specified.
```bash
go run ./cmd/duw check -h
go run ./cmd/duw check --verbose
go run ./cmd/duw check --path=/ --percentage 7
```
### Using env file
Create `.env` file in the root folder, disk usage warner will load `.env` file automatically.