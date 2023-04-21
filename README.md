# Got
[![Go](https://github.com/rameshvarun/got/actions/workflows/go.yml/badge.svg)](https://github.com/rameshvarun/got/actions/workflows/go.yml)

A VCS written in Go, made for practice.

## Example Usage
```bash
got init # In an empty folder

# Create an initial commit
echo "Test File" > test.txt
got commit -m "Initial commit." -a "Nobody"

# Modify the file
echo "Test File Modified" > test.txt
got commit -m "Second commit." -a "Nobody"

# See repository history
got log
```

## Concepts
HEADS - A list of all heads in the repository

CURRENT - The current revision

## Commands
### got init
### got log
Shows the commit history of the current repository, sorted by time. This command works by starting at the heads, then progressively traversing parent pointers.

### got commit -m "message" -a "author"
### got status
### got merge
