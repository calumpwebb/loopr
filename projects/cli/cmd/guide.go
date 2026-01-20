package cmd

import (
	"fmt"
)

// Guide outputs a minimal "Hello World" guide for getting started with loopr
func Guide() {
	guide := `# Loopr Hello World

## Quick Start

1. **Initialize loopr** in your project:
   ` + "```bash" + `
   loopr init
   ` + "```" + `

2. **Add tasks** to .loopr/tasks.md:
   ` + "```markdown" + `
   - [ ] Add user authentication (priority: high)
   - [ ] Create dashboard UI (priority: medium)
   ` + "```" + `

3. **Run a planning loop** (analyzes code, refines tasks):
   ` + "```bash" + `
   loopr plan
   ` + "```" + `

4. **Run a build loop** (implements highest priority tasks):
   ` + "```bash" + `
   loopr build
   ` + "```" + `

Each iteration automatically commits to git. That's it!

**Need more?** Run ` + "`loopr <command> --help`" + ` for details on any command.
`
	fmt.Println(guide)
}
