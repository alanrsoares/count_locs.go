# 🐹 Count LOCs 🐹

A blazing-fast™ command-line tool to recursively count lines of code in a directory, supporting custom glob patterns for file matching. Built with **Go**, leveraging concurrency for speed and efficiency.

## Features

- **Recursive directory traversal**
- **Glob pattern matching** for file extensions and paths
- **Concurrent processing** for fast line counting
- **Breakdown of LOC** by glob pattern (when multiple patterns are used)
- Lightweight and easy to build and install

## Installation

### Prerequisites

- [Go](https://golang.org/dl/) installed

### Steps

1. **Clone the repository:**

   ```bash
   git clone git@github.com:yourusername/count_locs.git
   cd count_locs
   ```

2. **Download dependencies:**

   ```bash
   go mod tidy
   ```

3. **Build the binary:**

   ```bash
   go build -o count_locs
   ```

4. **(Optional) Install Globally:**

   Move the binary to a directory that's in your PATH. For example:

   ```bash
   sudo mv count_locs /usr/local/bin/
   ```

5. **Verify Installation:**

   ```bash
   count_locs --help
   ```

## Usage

### Basic Command

```bash
count_locs <directory> <glob-patterns>...
```

### Examples

#### Count all Go files in the current directory:

```bash
count_locs . "**/*.go"
```

#### Count all Go, Python, and Markdown files in a project:

```bash
count_locs ./ "**/*.go" "**/*.py" "**/*.md"
```

#### Example Output

If multiple glob patterns are provided, you'll see a breakdown:

```plaintext
Breakdown of Lines of Code by Glob:
  **/*.go: 1200
  **/*.py: 800
  **/*.md: 150

Total:   2150 lines of code

Elapsed time: 150ms
```

If only one glob pattern is used, only the total is shown:

```plaintext
Total:   1500 lines of code

Elapsed time: 100ms
```

## Development

### Prerequisites

Ensure you have Go installed. You can download it from [golang.org](https://golang.org/dl/).

### Build

Build the project with:

```bash
go build -o count_locs
```

### Run Directly

You can also run the program without building a binary:

```bash
go run main.go <directory> <glob-patterns>...
```

### Run Tests

If tests are provided, run:

```bash
go test ./...
```

## Makefile Support

A `Makefile` is included for convenience. Here are some example targets:

- **Build the project:**
  ```bash
  make build
  ```
- **Run tests:**
  ```bash
  make test
  ```
- **Install globally:**
  ```bash
  make install
  ```
- **Clean build artifacts:**
  ```bash
  make clean
  ```
- **Run the binary:**
  ```bash
  make run
  ```

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a feature branch.
3. Commit your changes.
4. Push to your branch.
5. Open a pull request.

## License

This project is licensed under the [Unlicense](LICENSE).
