# echog

A simple Go implementation of the `echo` command.

## Features

- Prints text to standard output.
- Supports `-n` flag to omit the newline.
- Mimics Unix `echo` behavior.

## Installation

Clone the repository and navigate to this project:

```sh
git clone https://github.com/mnishiguchi/command-line-go.git
cd command-line-go/echog
```

Then, install dependencies:

```sh
go mod tidy
```

## Usage

### Run directly (without building)

```sh
go run cmd/main.go 元氣ですかー
```

#### Expected Output

```
元氣ですかー
```

### Omit newline

```sh
go run cmd/main.go -n 元氣ですかー
```

```
元氣ですかー%  # No newline
```

### Build and run

```sh
go build -o build/echog ./cmd
./build/echog 元氣ですかー
```

## Testing

Run CLI tests:

```sh
go test ./tests
```

## Project Structure

```
echog/
├── .gitignore       # Ignore build files
├── README.md        # Documentation
├── go.mod           # Go module file
├── go.sum           # Dependencies
├── cmd/
│   └── main.go      # CLI entry point
├── build/           # Compiled binary (ignored)
└── tests/           # CLI tests
    └── cli_test.go  # CLI integration tests
```
