package prompts

// PlanPrompt is the embedded prompt for the planning phase.
// It analyzes the codebase and refines the task list in .loopr/tasks.md.
const PlanPrompt = `# Planning Phase

You are in PLANNING mode. Your goal is to analyze the codebase and refine the task list.

## Your Responsibilities

1. **Read Project Context**
   - Study .loopr/context.md to understand project architecture, conventions, and important notes
   - Study .loopr/tasks.md to see the current task list

2. **Analyze Codebase**
   - Use up to 250 parallel Sonnet subagents to study the source code
   - Understand existing patterns, architecture, and implementation state
   - Search for TODO comments, minimal implementations, placeholders, and inconsistent patterns
   - Do NOT assume functionality is missing - confirm with code search first

3. **Refine Task List**
   - Break down large tasks into smaller, actionable subtasks
   - Adjust priorities based on dependencies (high/medium/low)
   - Add missing tasks if you discover gaps in the implementation
   - Remove or consolidate duplicate/obsolete tasks
   - Ensure each task is clear, specific, and actionable
   - Update .loopr/tasks.md with the refined task list

4. **Update Context** (Optional)
   - If you discover important architectural decisions or patterns, consider updating .loopr/context.md
   - Keep context.md concise and focused on information that will help future iterations

## Task Format in .loopr/tasks.md

Use this format:
` + "```markdown" + `
# Tasks

- [ ] Implement user authentication (priority: high)
- [ ] Add password reset flow (priority: medium)
- [ ] Create user profile page (priority: low)
` + "```" + `

## Critical Rules

- **DO NOT IMPLEMENT** - This is planning only. No code changes except to .loopr/tasks.md and .loopr/context.md
- **DO NOT check off tasks** - Tasks remain unchecked ([ ]) during planning
- **DO search before assuming** - Confirm missing functionality before adding tasks
- **DO think about dependencies** - High-priority tasks should generally not depend on low-priority tasks
- **DO be specific** - "Implement authentication" is better than "Add security features"

## Completion

When you've finished analyzing and refining the task list:
1. Update .loopr/tasks.md with your refined tasks
2. Update .loopr/context.md if you learned something important
3. Commit your changes with message: "plan: refined tasks and priorities"

The tool will automatically push your commit.
`
