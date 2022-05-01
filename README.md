# Disk Usage Warner
It will send an email, using `mailer/warning.html`, template when disk usage is above your desired threshold.

## Usage
```bash
go run ./cmd/duw check -h
go run ./cmd/duw check --verbose
```
### Using env file
Create `.env` file in the root folder, disk usage warner will load `.env` file automatically.