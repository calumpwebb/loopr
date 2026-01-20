package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/calumpwebb/loopr/internal/ui"
)

// Status shows the current task status
func Status() {
	// Check if .loopr directory exists
	if _, err := os.Stat(".loopr"); os.IsNotExist(err) {
		fmt.Println(ui.ErrorStyle.Render("✗ ERROR: .loopr directory not found"))
		fmt.Println("\nRun 'loopr init' first to set up your project.")
		os.Exit(1)
	}

	// Check if tasks.md exists
	tasksFile := ".loopr/tasks.md"
	if _, err := os.Stat(tasksFile); os.IsNotExist(err) {
		fmt.Println(ui.ErrorStyle.Render("✗ ERROR: .loopr/tasks.md not found"))
		fmt.Println("\nRun 'loopr init' first to set up your project.")
		os.Exit(1)
	}

	// Read tasks.md
	file, err := os.Open(tasksFile)
	if err != nil {
		fmt.Println(ui.ErrorStyle.Render(fmt.Sprintf("✗ ERROR: Failed to read tasks.md: %v", err)))
		os.Exit(1)
	}
	defer file.Close()

	var uncheckedCount int
	var checkedCount int
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Check if line is a task
		if strings.HasPrefix(line, "- [x]") || strings.HasPrefix(line, "- [X]") {
			checkedCount++
		} else if strings.HasPrefix(line, "- [ ]") {
			uncheckedCount++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(ui.ErrorStyle.Render(fmt.Sprintf("✗ ERROR: Failed to read tasks.md: %v", err)))
		os.Exit(1)
	}

	// Get last commit message
	lastCommitCmd := exec.Command("git", "log", "-1", "--pretty=format:%s")
	lastCommitOutput, err := lastCommitCmd.Output()
	lastCommit := "No commits yet"
	if err == nil && len(lastCommitOutput) > 0 {
		lastCommit = string(lastCommitOutput)
	}

	// Get git status
	gitStatusCmd := exec.Command("git", "status", "--porcelain")
	gitStatusOutput, err := gitStatusCmd.Output()
	gitStatus := "clean working tree"
	if err == nil && len(gitStatusOutput) > 0 {
		lines := strings.Split(strings.TrimSpace(string(gitStatusOutput)), "\n")
		gitStatus = fmt.Sprintf("%d file(s) modified", len(lines))
	}

	// Print status
	fmt.Println()
	fmt.Println(ui.SuccessStyle.Render("=== Loopr Status ==="))
	fmt.Println()

	// Task counts
	totalTasks := uncheckedCount + checkedCount
	if totalTasks == 0 {
		fmt.Println(ui.WarningStyle.Render("No tasks found"))
	} else {
		fmt.Printf("Tasks:     %d total (%d unchecked, %d completed)\n", totalTasks, uncheckedCount, checkedCount)

		// Progress bar
		if totalTasks > 0 {
			percentage := (float64(checkedCount) / float64(totalTasks)) * 100
			barLength := 20
			filled := int((float64(checkedCount) / float64(totalTasks)) * float64(barLength))
			bar := strings.Repeat("█", filled) + strings.Repeat("░", barLength-filled)
			fmt.Printf("Progress:  [%s] %.1f%%\n", bar, percentage)
		}
	}

	fmt.Println()

	// Git info
	fmt.Printf("Last commit: %s\n", lastCommit)
	fmt.Printf("Git status:  %s\n", gitStatus)

	fmt.Println()

	// Next steps
	if uncheckedCount > 0 {
		fmt.Println("Next steps:")
		fmt.Println("  - Run 'loopr plan' to refine tasks")
		fmt.Println("  - Run 'loopr build' to implement tasks")
	} else if checkedCount > 0 {
		fmt.Println("All tasks complete! 🎉")
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  - Run 'loopr archive' to clean up completed tasks")
		fmt.Println("  - Add more tasks to .loopr/tasks.md")
		fmt.Println("  - Import new tasks with 'loopr import <file>'")
	}
}
