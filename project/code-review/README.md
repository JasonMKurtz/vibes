# Code Review Project

This project provides a simple code review tool with a Golang backend and a React frontend.

## Backend

The backend API is implemented in Go using Gin and GORM with an SQLite database. The main data model is a `PR` that keeps track of the pull request author, a list of reviewers, and which actor needs to take the next action. The next actor must match either the PR author or one of the reviewers.

Run the backend with:

```bash
cd backend
go run ./...
```

Run tests with:

```bash
go test ./...
```

> **Note**: Dependencies are referenced in `go.mod` but are not vendored. You may need network access to download them.

Frontend is not yet implemented.
