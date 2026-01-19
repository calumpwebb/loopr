package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/calumpwebb/loopr/internal/ui"
)

// Archive moves completed tasks from tasks.md to .loopr/completed/YYYY-MM-DD.md
func Archive() {
	// Check if .loopr directory exists
	if _, err := os.Stat(".loopr"); os.IsNotExist(err) {
		fmt.Println(ui.ErrorStyle.Render("✗ .loopr directory not found"))
		fmt.Println("\nRun 'loopr init' first to set up your project.")
		os.Exit(1)
	}

	// Check if tasks.md exists
	tasksFile := ".loopr/tasks.md"
	if _, err := os.Stat(tasksFile); os.IsNotExist(err) {
		fmt.Println(ui.ErrorStyle.Render("✗ .loopr/tasks.md not found"))
		fmt.Println("\nRun 'loopr init' first to set up your project.")
		os.Exit(1)
	}

	// Read tasks.md
	file, err := os.Open(tasksFile)
	if err != nil {
		fmt.Println(ui.ErrorStyle.Render(fmt.Sprintf("✗ Failed to read tasks.md: %v", err)))
		os.Exit(1)
	}
	defer file.Close()

	var remainingLines []string
	var completedLines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Check if line is a completed task: starts with "- [x]" or "- [X]"
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "- [x]") || strings.HasPrefix(trimmedLine, "- [X]") {
			completedLines = append(completedLines, line)
		} else {
			remainingLines = append(remainingLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(ui.ErrorStyle.Render(fmt.Sprintf("✗ Failed to read tasks.md: %v", err)))
		os.Exit(1)
	}

	// Check if there are any completed tasks
	if len(completedLines) == 0 {
		fmt.Println(ui.WarningStyle.Render("No completed tasks to archive."))
		os.Exit(0)
	}

	// Create archive file path with today's date
	today := time.Now().Format("2006-01-02")
	archiveDir := ".loopr/completed"
	archiveFile := fmt.Sprintf("%s/%s.md", archiveDir, today)

	// Create completed directory if it doesn't exist
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		fmt.Println(ui.ErrorStyle.Render(fmt.Sprintf("✗ Failed to create completed directory: %v", err)))
		os.Exit(1)
	}

	// Check if archive file already exists
	var existingContent []string
	if _, err := os.Stat(archiveFile); err == nil {
		// File exists, read existing content
		existingFile, err := os.Open(archiveFile)
		if err != nil {
			fmt.Println(ui.ErrorStyle.Render(fmt.Sprintf("✗ Failed to read existing archive: %v", err)))
			os.Exit(1)
		}
		defer existingFile.Close()

		scanner := bufio.NewScanner(existingFile)
		for scanner.Scan() {
			existingContent = append(existingContent, scanner.Text())
		}
	}

	// Write archive file
	archiveOutput, err := os.Create(archiveFile)
	if err != nil {
		fmt.Println(ui.ErrorStyle.Render(fmt.Sprintf("✗ Failed to create archive file: %v", err)))
		os.Exit(1)
	}
	defer archiveOutput.Close()

	writer := bufio.NewWriter(archiveOutput)

	// If existing content, append to it
	if len(existingContent) > 0 {
		for _, line := range existingContent {
			fmt.Fprintln(writer, line)
		}
		fmt.Fprintln(writer, "") // Add blank line separator
	} else {
		// Write header for new file
		fmt.Fprintf(writer, "# Completed Tasks - %s\n\n", today)
	}

	// Write completed tasks
	for _, line := range completedLines {
		fmt.Fprintln(writer, line)
	}

	if err := writer.Flush(); err != nil {
		fmt.Println(ui.ErrorStyle.Render(fmt.Sprintf("✗ Failed to write archive file: %v", err)))
		os.Exit(1)
	}

	// Write remaining tasks back to tasks.md
	tasksOutput, err := os.Create(tasksFile)
	if err != nil {
		fmt.Println(ui.ErrorStyle.Render(fmt.Sprintf("✗ Failed to write tasks.md: %v", err)))
		os.Exit(1)
	}
	defer tasksOutput.Close()

	tasksWriter := bufio.NewWriter(tasksOutput)
	for _, line := range remainingLines {
		fmt.Fprintln(tasksWriter, line)
	}

	if err := tasksWriter.Flush(); err != nil {
		fmt.Println(ui.ErrorStyle.Render(fmt.Sprintf("✗ Failed to write tasks.md: %v", err)))
		os.Exit(1)
	}

	// Show success message
	fmt.Println()
	fmt.Println(ui.SuccessStyle.Render(fmt.Sprintf("✓ Archived %d completed task(s)", len(completedLines))))
	fmt.Println()
	fmt.Printf("Archive location: %s\n", archiveFile)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  - Review .loopr/tasks.md to see remaining tasks")
	fmt.Println("  - Run 'loopr plan' to refine tasks")
	fmt.Println("  - Run 'loopr build' to implement tasks")
}
