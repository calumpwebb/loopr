# Loopr Testing Plan - proxlog PRD

## Test Objectives

Validate that loopr can:
1. Initialize a new project with templates
2. Import tasks from a PRD
3. Run planning loops that refine tasks
4. Run build loops that implement functionality
5. Create git commits after each iteration
6. Handle interactive prompts correctly
7. Produce working Go code

## Test Environment

- **Project**: loopr-test-bed (clean git repo)
- **PRD**: PRD.md (proxlog HTTP proxy logger)
- **Binary**: ./loopr (freshly built)
- **Method**: tmux for interactive testing

## Pre-Test Setup

1. Build latest loopr binary
2. Clear all Docker sandboxes (fresh auth state)
3. Kill any existing tmux sessions
4. Ensure loopr-test-bed is clean git repo

## Test Cases

### TC1: Initialize Project
**Command**: `loopr init`
**Expected**:
- Creates `.loopr/` directory
- Extracts templates: tasks.md, context.md
- Extracts CLAUDE.md to project root
- Prompts before overwriting if exists

**Validation**:
- [ ] `.loopr/tasks.md` exists and has template content
- [ ] `.loopr/context.md` exists
- [ ] `CLAUDE.md` exists in project root

### TC2: Import Tasks from PRD
**Command**: `loopr import PRD.md`
**Expected**:
- Runs Claude to extract tasks from PRD
- Updates `.loopr/tasks.md` with checkboxes
- Tasks have priorities (high/medium/low)
- Creates git commit with import results

**Validation**:
- [ ] `.loopr/tasks.md` has tasks related to proxlog
- [ ] Tasks have proper format: `- [ ] Task (priority: X)`
- [ ] Git commit created with message about import
- [ ] Tasks cover: CLI setup, proxy server, logging, viewer

### TC3: Planning Loop (2-3 iterations)
**Command**: `loopr plan` (with 2-3 iterations)
**Expected**:
- Prompts for authentication (if needed)
- Prompts for iterations count
- Each iteration:
  - Analyzes existing tasks
  - Refines/adds/removes tasks
  - Updates priorities
  - Creates git commit
- Does NOT implement code (planning only)

**Validation**:
- [ ] 2-3 git commits created (one per iteration)
- [ ] Tasks in `.loopr/tasks.md` are more detailed/refined
- [ ] No Go code files created yet
- [ ] Commit messages describe planning changes

### TC4: Build Loop (3-5 iterations)
**Command**: `loopr build` (with 3-5 iterations)
**Expected**:
- Each iteration:
  - Picks highest priority unchecked task
  - Implements the task
  - Checks off task in tasks.md
  - Creates git commit
- Creates Go files (main.go, cmd/, internal/)
- Code should compile

**Validation**:
- [ ] 3-5 git commits created
- [ ] Tasks checked off in `.loopr/tasks.md`
- [ ] Go files created with proper structure
- [ ] `go build` succeeds (at least by end)
- [ ] Commit messages describe what was built

### TC5: Verify Functionality
**Manual checks after build iterations**
**Expected**:
- proxlog binary can be built
- Has expected commands (start, view)
- Code quality is reasonable

**Validation**:
- [ ] `go build` produces binary
- [ ] Binary has help text
- [ ] Code has proper package structure
- [ ] No obvious syntax errors

## Test Execution Strategy

### Phase 1: Setup (Manual)
```bash
cd /Users/calum/Development/loopr
go build -o loopr main.go
docker sandbox ls | awk 'NR>1 {print $1}' | xargs -I {} docker sandbox rm {}
tmux kill-session -t loopr-test 2>/dev/null
```

### Phase 2: Run Tests (Automated via subagent)
- Launch tmux session in loopr-test-bed
- Execute each test case sequentially
- Capture output and validate
- Monitor git commits

### Phase 3: Validation (Automated)
- Check file existence
- Parse tasks.md for expected format
- Verify git commit count
- Attempt to build Go binary

### Phase 4: Cleanup & Report
- Kill tmux session
- Generate test report
- Push results to git (if successful)

## Success Criteria

✅ All test cases pass
✅ At least 5-8 git commits created total
✅ proxlog has basic structure (main.go, commands)
✅ No crashes or hangs during execution
✅ Tasks are properly formatted and tracked

## Failure Scenarios to Handle

- Authentication fails → Need Anthropic API key
- Sandbox not found → Docker not running
- Git push fails → Branch doesn't exist
- Prompt timeout → Interactive response needed
- Claude produces invalid Go code → Check error handling

## Timeline Estimate

- Setup: 1 min
- TC1 (init): 30 sec
- TC2 (import): 2-3 min
- TC3 (plan 3x): 6-9 min
- TC4 (build 5x): 15-25 min
- Validation: 2 min

**Total**: ~30-40 minutes end-to-end
