# hello

A simple command-line tool built with Go and `urfave/cli`, part of the [command-line-go](https://github.com/mnishiguchi/command-line-go) project.

## Features

- Prints a motivational quote in Japanese.
- Demonstrates basic usage of `urfave/cli`.
- Follows a clean project structure for CLI tools.

## Installation

First, clone the main repository and navigate into this subproject:

```sh
git clone https://github.com/mnishiguchi/command-line-go.git
cd command-line-go/hello
```

Then, install dependencies:

```sh
go mod tidy
```

## Usage

### **Build the program**

```sh
go build -o build/hello ./cmd
```

### **Run the CLI**

```sh
./build/hello
```

#### **Expected Output**

```
元氣が一番、元氣があれば何でもできる！
```

## Testing

Run tests using:

```sh
go test ./tests
```

## Project Structure

```
command-line-go/
├── hello/               # This subproject
│   ├── .gitignore       # Ignore build files and cache
│   ├── README.md        # Project documentation
│   ├── go.mod           # Go module file
│   ├── go.sum           # Go module dependencies
│   ├── cmd/
│   │   └── main.go      # CLI entry point
│   ├── build/           # Compiled binary (ignored in .gitignore)
│   └── tests/           # CLI integration tests
└── [other subprojects]  # Other CLI tools in the same repo
```

## License

This is just an example project and is not meant to be used as a package.
