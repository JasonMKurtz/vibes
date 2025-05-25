# Vibes IRC Server

This repository contains a minimal IRC server implementation written in Go.

## Usage

Run the server from the repository root:

```
go run ./project/irc/server
```

The server listens on TCP port `6667` by default.

## Running Tests

Change into the `project/irc` directory and run:

```
go test ./...
```

## Manual Testing

You can exercise the server manually using `telnet` or `nc`:

1. Start the server:

   ```
   go run ./project/irc/server
   ```

2. In another terminal connect to the server:

   ```
   nc localhost 6667
   ```

3. Issue some basic IRC commands:

   ```
   NICK alice
   USER alice 0 * :Alice
   JOIN #chat
   PRIVMSG #chat :hello there!
   QUIT
   ```

The server will write connection and channel activity to `server.log` while
errors continue to appear on stderr.

## Client CLI

This repository includes a small command line client built on top of the
SDK. Run it from the project root:

```
go run ./project/irc/cli
```

Use the `-server` flag to specify the IRC server address if it differs from
the default `localhost:6667`. The CLI accepts simple commands:

- `login <name>` – set the nickname and username
- `join <channel>` – join a channel
- `msg <target> <message>` – send a message
- `read` – read one line from the server
- `quit` – exit the program

Each command is typed interactively after starting the client.
