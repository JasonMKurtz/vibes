# Code Review Project

This project provides a simple code review tool with a Golang backend and a React frontend.

## Backend

The backend API is implemented in Go using Gin and GORM with an SQLite database. The main data models are:

* `PR` – tracks the pull request author, a list of reviewers and which actor needs to take the next action.
* `Review` – represents a review of a PR and can hold an overall state such as approval or request for changes.
* `Comment` – individual comments attached to a PR or review.

The next actor on a PR must match either the PR author or one of the reviewers.

Run the backend with:

```bash
cd backend
go run ./...
```

The server can be configured with environment variables:

* `PORT` - HTTP listen address, defaults to `:8080`.
* `DB_PATH` - SQLite database file path, defaults to `code_review.db`.

Run tests with:

```bash
go test ./...
```

> **Note**: Dependencies are referenced in `go.mod` but are not vendored. You may need network access to download them.

Frontend is not yet implemented.
