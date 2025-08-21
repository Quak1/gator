# gator

A simple terminal blog agreggator. Made following Boot.dev's [Build a Blog Aggregator in Go](https://www.boot.dev/courses/build-blog-aggregator-golang) course.

## Features

- **CLI interface**: easy to use command line interface
- **RSS Feed Management**: Add, list, and follow RSS feeds
- **User Management**: Create users with different following lists
- **Posts storage**: Automatically parse and store posts from feeds

## Prerequisites

- Go 1.25+
- PostgreSQL instance
- psql
- [goose](https://github.com/pressly/goose)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/Quak1/gator
cd gator
```

2. Install dependencies:

```bash
go mod download
```

3. Set up PostgreSQL database

```bash
psql -U postgres -c "CREATE DATABASE gator;"
```

4. Copy `.gatorconfig` to your home directory, and set your connection string

5. Run database migrations

```bash
cd sql/schema
goose postgres <connection_string> up
```

6. Build or Install the application

```bash
go build
go install
```

## Usage

1. Run the executable

```bash
./gator
```

- Alternatively if installed run from anywhere

```bash
gator <command> [...args]
```

### Available commands

| command     | params                   | Action                                                                          |
| ----------- | ------------------------ | ------------------------------------------------------------------------------- |
| `login`     | `username`               | Login into a given `username`                                                   |
| `register`  | `username`               | Register a new user                                                             |
| `reset`     |                          | Deletes all registered users and related items                                  |
| `users`     |                          | List all registered users                                                       |
| `agg`       | `<interval>` eg. `1h 5m` | Start the agreggator. Pulls new posts every `interval`                          |
| `addfeed`   | `<feed_name> <feed_url>` | Add new feed to the database                                                    |
| `feeds`     |                          | List all current registered feeds                                               |
| `follow`    | `<feed_url>`             | Follow a feed                                                                   |
| `unfollow`  | `<feed_url>`             | Unfollow a feed                                                                 |
| `following` |                          | List following feeds                                                            |
| `browse`    | `[post_count]`           | List posts from following feeds. Shows up to `post_count` posts. Defaults to 2. |
