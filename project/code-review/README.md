# Code Review Project

This project provides a simple code review tool with a Golang backend and a React frontend.

## Backend

The backend API is implemented in Go using Gin and GORM with an SQLite database. The main data model is a `PR` that keeps track of the pull request author, reviewer, and which actor needs to take the next action.

Run the backend with:

```bash
cd backend
go run ./...
```

> **Note**: Dependencies are referenced in `go.mod` but are not vendored. You may need network access to download them.

Frontend is not yet implemented.
