package prompts

import "fmt"

// GetImportPrompt returns the prompt for importing tasks from a PRD/spec file.
// It takes the file path as a parameter to include in the prompt.
func GetImportPrompt(sourceFile string) string {
	return fmt.Sprintf(`# Import Tasks from Specification

You are in IMPORT mode. Your goal is to extract actionable tasks from a specification document.

## Your Responsibilities

1. **Read the Source Document**
   - Read the file: %s
   - This may be a PRD (Product Requirements Document), feature spec, or technical design doc
   - Understand the scope and requirements

2. **Extract Tasks**
   - Identify all actionable implementation tasks
   - Break down features into specific, concrete tasks
   - Each task should be implementable in a reasonable iteration
   - Assign priorities based on:
     - Critical/foundational features → high
     - Standard features → medium
     - Nice-to-have features → low
     - Dependencies (if task B depends on task A, task A should be high/medium)

3. **Format Tasks**
   - Use checkbox format: `+"`- [ ] Task description (priority: level)`"+`
   - Be specific and actionable
   - Example: "Implement user authentication with JWT" not "Add security"
   - Example: "Create database schema for users table" not "Database stuff"

4. **Append to Task List**
   - Read existing .loopr/tasks.md
   - Append new tasks to the end (do NOT overwrite existing tasks)
   - If a task already exists in tasks.md, skip it (no duplicates)
   - Maintain consistent formatting

5. **Commit Changes**
   - Commit with message: `+"`import: added tasks from %s`"+`

## Task Format

`+"```markdown"+`
# Tasks

## Existing tasks (preserve these)
- [ ] Existing task 1 (priority: high)
- [x] Completed task (priority: medium)

## Newly imported tasks (add these)
- [ ] New task from PRD (priority: high)
- [ ] Another new task (priority: medium)
`+"```"+`

## Priority Guidelines

**High Priority:**
- Foundational/infrastructure tasks that other tasks depend on
- Core features that are critical to the product
- Security or data integrity tasks
- Blocking issues

**Medium Priority:**
- Standard features
- Enhancements to existing functionality
- Performance improvements
- Non-critical bug fixes

**Low Priority:**
- Nice-to-have features
- UI polish
- Documentation
- Refactoring (unless blocking other work)

## Critical Rules

- **DO preserve existing tasks** - Never delete or modify tasks already in tasks.md
- **DO avoid duplicates** - Check if a task already exists before adding
- **DO be specific** - Each task should be clear and actionable
- **DO consider dependencies** - If task B needs task A, make task A higher priority
- **DO NOT use emojis** - Keep all output, commit messages, and task descriptions emoji-free
- **DO NOT implement anything** - This is import only, no code changes
- **DO NOT check off tasks** - All imported tasks should be unchecked [ ]

## Example

If the PRD says:
`+"```"+`
Feature: User Management
- Users can register with email/password
- Users can reset their password
- Admin dashboard to view all users
`+"```"+`

Extract tasks like:
`+"```markdown"+`
- [ ] Create database schema for users table (priority: high)
- [ ] Implement user registration endpoint (priority: high)
- [ ] Implement password hashing with bcrypt (priority: high)
- [ ] Create password reset flow with email tokens (priority: medium)
- [ ] Build admin dashboard UI (priority: medium)
- [ ] Add admin API endpoint to list users (priority: medium)
`+"```"+`

## Completion

After extracting and appending tasks:
1. Verify all tasks are in correct format with priorities
2. Verify no existing tasks were modified or deleted
3. Commit changes with: `+"`import: added tasks from %s`"+`
4. The tool will automatically push your commit
`, sourceFile, sourceFile, sourceFile)
}
