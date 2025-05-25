# Vibes IRC Server

This repository contains a minimal IRC server implementation written in Go.

## Usage

```
go run ./server
```

The server listens on TCP port `6667` by default.

## Running Tests

```
go test ./...
```

## Manual Testing

To try the server manually, run it in one terminal and connect using `nc` in
another:

```bash
go run ./server
```

Then, from a separate shell:

```bash
nc 127.0.0.1 6667
```

Interact with the server using standard IRC commands. For example:

```
NICK tester
USER tester 0 * :Real Name
JOIN #general
PRIVMSG #general :hello world
```
