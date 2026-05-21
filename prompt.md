# Go Teacher Prompt Templates

Use these prompts when invoking the Go Teacher agent. Prefer Vietnamese if you write in Vietnamese.

---

## Lesson / Explanation
- "Explain `<topic>` in Go with a short example and a 2-exercise drill." 
- Example: "Explain `channels` and `select` in Go with a minimal example and two exercises."

## Code Review
- "Review this function and suggest idiomatic Go improvements (show code)."
- Provide the function code after the prompt. Ask which Go version to target if unclear.

## Small Example Request
- "Give a small, runnable example showing `context.WithTimeout` usage and explain each line."
- "Show a minimal HTTP handler that decodes JSON input and validates fields (no full app)."

## Debugging / Fix Suggestion
- "This test/example fails with `<error>` — suggest likely causes and a minimal fix or diagnostic steps." 

## Exercises / Practice
- "Give 3 exercises to practice goroutines and channels, ordered by difficulty, with brief success criteria."

## Scaffolding Requests (small)
- "Provide a minimal file scaffold for `<purpose>` (e.g., `auth middleware`, `service layer`) with comments explaining each part."

## Teaching Preferences (include in prompt when relevant)
- Preferred Go version: `<e.g., 1.20>`
- Preference: `concise` or `step-by-step`
- Apply edits directly? `yes` / `no` (agent will ask for permission before writing files)

---

Tips for using the agent:
- Ask for minimal, runnable snippets rather than full applications.
- Include the relevant file or function in the prompt when requesting reviews.
- If you want the agent to apply edits, say so explicitly and grant permission.
