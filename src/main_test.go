package main

import (
	"bufio"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestCountLines creates a temporary file with known content and
// ensures that countLines returns the expected count of non-empty lines.
func TestCountLines(t *testing.T) {
	// Prepare a temporary file with a mix of non-empty and whitespace lines.
	content := "Line one\n\n   \nLine two\nLine three\n"
	tmpFile, err := os.CreateTemp("", "testfile*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Expected: "Line one", "Line two", and "Line three" (3 lines).
	count := countLines(tmpFile.Name())
	if count != 3 {
		t.Errorf("expected 3 non-empty lines, got %d", count)
	}
}

// TestCountLocs creates a temporary directory structure with files that match
// (and do not match) a glob pattern and verifies that countLocs returns the correct counts.
func TestCountLocs(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// File 1: matching pattern **/*.txt in root directory.
	file1Path := filepath.Join(tempDir, "file1.txt")
	if err := os.WriteFile(file1Path, []byte("Hello\nWorld\n\n"), 0644); err != nil {
		t.Fatalf("failed to write file1: %v", err)
	}
	// Expected count: 2 (ignoring the empty line)

	// File 2: matching pattern **/*.txt inside a subdirectory.
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}
	file2Path := filepath.Join(subDir, "file2.txt")
	if err := os.WriteFile(file2Path, []byte("One\nTwo\nThree\n"), 0644); err != nil {
		t.Fatalf("failed to write file2: %v", err)
	}
	// Expected count: 3

	// File 3: a file that does not match the glob pattern.
	file3Path := filepath.Join(tempDir, "file3.log")
	if err := os.WriteFile(file3Path, []byte("Log line\n"), 0644); err != nil {
		t.Fatalf("failed to write file3: %v", err)
	}

	patterns := []string{"**/*.txt"}
	results := countLocs(tempDir, patterns)
	total := results["**/*.txt"]
	expectedTotal := 2 + 3
	if total != expectedTotal {
		t.Errorf("expected %d lines total, got %d", expectedTotal, total)
	}
}

// TestProcessInputMultiplePatterns creates a temporary directory with two files
// matching different glob patterns and then calls processInput.
// It captures standard output and verifies that the printed breakdown and totals are present.
func TestProcessInputMultiplePatterns(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "testprocess")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a Rust file matching **/*.rs.
	file1Path := filepath.Join(tempDir, "file1.rs")
	if err := os.WriteFile(file1Path, []byte("fn main() {}\n\n// comment\n"), 0644); err != nil {
		t.Fatalf("failed to write file1.rs: %v", err)
	}
	// Expected count: 2 non-empty lines.

	// Create a TypeScript file matching **/*.ts.
	file2Path := filepath.Join(tempDir, "file2.ts")
	if err := os.WriteFile(file2Path, []byte("console.log('Hello');\n"), 0644); err != nil {
		t.Fatalf("failed to write file2.ts: %v", err)
	}
	// Expected count: 1 non-empty line.

	patterns := []string{"**/*.rs", "**/*.ts"}

	// Redirect standard output.
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	// Call processInput. (Note that processInput prints elapsed time so we allow a little slack.)
	start := time.Now()
	processInput(tempDir, patterns)
	elapsed := time.Since(start)
	_ = elapsed // Not asserting on elapsed time; just illustrating that it is printed.

	// Restore stdout and capture output.
	w.Close()
	os.Stdout = oldStdout
	var outputBuilder strings.Builder
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		outputBuilder.WriteString(scanner.Text() + "\n")
	}
	output := outputBuilder.String()

	// Check that breakdown and total are printed.
	if !strings.Contains(output, "Breakdown of Lines of Code by Glob:") {
		t.Errorf("expected breakdown header in output, got: %s", output)
	}
	if !strings.Contains(output, "**/*.rs:") || !strings.Contains(output, "**/*.ts:") {
		t.Errorf("expected glob patterns to be listed in output, got: %s", output)
	}
	if !strings.Contains(output, "Total:") {
		t.Errorf("expected total line count in output, got: %s", output)
	}
}

// TestPrintHelp captures the output of printHelp and verifies it contains expected help text.
func TestPrintHelp(t *testing.T) {
	// Redirect stdout.
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	printHelp()

	// Restore stdout and capture output.
	w.Close()
	os.Stdout = oldStdout
	var outputBuilder strings.Builder
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		outputBuilder.WriteString(scanner.Text() + "\n")
	}
	output := outputBuilder.String()

	if !strings.Contains(output, "Usage: count_locs") {
		t.Errorf("help message does not contain expected text, got: %s", output)
	}
}

// TestPrintVersion captures the output of printVersion and verifies it contains the version string.
func TestPrintVersion(t *testing.T) {
	// Redirect stdout.
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	printVersion()

	// Restore stdout and capture output.
	w.Close()
	os.Stdout = oldStdout
	var outputBuilder strings.Builder
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		outputBuilder.WriteString(scanner.Text() + "\n")
	}
	output := outputBuilder.String()

	if !strings.Contains(output, version) {
		t.Errorf("version output does not contain expected version string %q, got: %s", version, output)
	}
}

// TestCountLocsConcurrency verifies countLocs works even when many files are processed concurrently.
func TestCountLocsConcurrency(t *testing.T) {
	// Create a temporary directory.
	tempDir, err := os.MkdirTemp("", "testconcurrency")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create several files in a loop.
	numFiles := runtime.NumCPU() * 2
	totalExpected := 0
	for i := 0; i < numFiles; i++ {
		filePath := filepath.Join(tempDir, "file_"+strconv.Itoa(i)+".txt")
		// Each file will have a fixed 2 non-empty lines.
		content := "Line 1\nLine 2\n\n"
		totalExpected += 2
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write %s: %v", filePath, err)
		}
	}

	patterns := []string{"**/*.txt"}
	results := countLocs(tempDir, patterns)
	if results["**/*.txt"] != totalExpected {
		t.Errorf("expected %d lines total, got %d", totalExpected, results["**/*.txt"])
	}
}
