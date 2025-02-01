package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bmatcuk/doublestar/v4"
)

const version = "0.1.0"

const helpMessage = `Usage: count_locs <directory> <glob-patterns>...

Options:
  -h, --help       Show this help message
  -v, --version    Show version information

Examples:
  count_locs ./src "**/*.rs" "**/*.ts"
  count_locs ./ "**/*.css"
`

func main() {
	args := os.Args

	// When no arguments are provided, show error
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: count_locs <directory> <glob-patterns>...")
		os.Exit(1)
	}

	// Handle help and version flags.
	// Note: In Go, os.Args[0] is the program name.
	if len(args) == 2 {
		switch args[1] {
		case "-h", "--help":
			printHelp()
			return
		case "-v", "--version":
			printVersion()
			return
		}
		// If only one argument is given (and not help/version),
		// show usage error.
		fmt.Fprintln(os.Stderr, "Usage: count_locs <directory> <glob-patterns>...")
		os.Exit(1)
	}

	// Otherwise the first argument is the directory and the remaining arguments are glob patterns.
	dir := args[1]
	patterns := args[2:]

	processInput(dir, patterns)
}

// processInput performs the work: it canonicalizes the directory,
// counts lines for each glob pattern concurrently, and prints the result.
func processInput(dir string, patterns []string) {
	startTime := time.Now()

	absDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to resolve directory: %v\n", err)
		os.Exit(1)
	}

	results := countLocs(absDir, patterns)

	// Sum the total lines across all patterns.
	totalLines := 0
	for _, count := range results {
		totalLines += count
	}

	// If more than one pattern, show the breakdown.
	if len(patterns) > 1 {
		fmt.Println("Breakdown of Lines of Code by Glob:")
		fmt.Println()

		for pattern, count := range results {
			fmt.Printf("-  %s: %d\n", pattern, count)
		}

		fmt.Println()
	}

	elapsed := time.Since(startTime)
	fmt.Printf("Total:\t%d lines of code\n\nElapsed time: %v\n", totalLines, elapsed)
}

// countLocs processes each glob pattern under the given root directory.
// It uses the doublestar.Glob function to support recursive globbing.
// For each matching file, it counts the non-empty (and non-whitespace) lines concurrently.
func countLocs(root string, patterns []string) map[string]int {
	results := make(map[string]int)

	for _, pattern := range patterns {
		// Create a filesystem rooted at 'root'
		fsys := os.DirFS(root)

		// Use doublestar.Glob with the fs.FS argument and the pattern.
		matches, err := doublestar.Glob(fsys, pattern)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing pattern %s: %v\n", pattern, err)
			continue
		}

		var total int64
		var wg sync.WaitGroup
		sem := make(chan struct{}, runtime.NumCPU())

		// Process each matching file concurrently.
		for _, file := range matches {
			// Since the file paths are relative to the filesystem root, build the absolute path.
			fullPath := filepath.Join(root, file)
			info, err := os.Stat(fullPath)
			if err != nil || info.IsDir() {
				continue
			}

			wg.Add(1)
			go func(filePath string) {
				defer wg.Done()
				sem <- struct{}{}
				count := countLines(filePath)
				atomic.AddInt64(&total, int64(count))
				<-sem
			}(fullPath)
		}
		wg.Wait()
		results[pattern] = int(total)
	}

	return results
}

// countLines opens the file at filePath and counts the number of non-empty, non-whitespace lines.
func countLines(filePath string) int {
	file, err := os.Open(filePath)
	if err != nil {
		// If the file cannot be opened, count it as 0 lines.
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			count++
		}
	}
	return count
}

// printHelp displays the help message.
func printHelp() {
	fmt.Print(helpMessage)
}

// printVersion displays the version information.
func printVersion() {
	fmt.Printf("count_locs version %s\n", version)
}
