# Got
A VCS written in Go, made for practice.

## Example Usage
```bash
got init # In an empty folder

# Create an initial commit
echo "Test File" > test.txt
got commit

# Modify the file
echo "Test File Modified" > test.txt
got commit

# Checkout the initial revision
got checkout ...
```

## Concepts
HEADS - A list of all heads in the repository

CURRENT - The current revision

## Commands
### got init
### got log
### got commit -m "message" -a "author"
### got status
### got merge
