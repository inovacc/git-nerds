# TODO: Map of Shell Commands to Go Implementations

This file tracks the mapping of original `git-quick-stats` shell commands to their Go equivalents for the Go port of the project.

## General Approach
- For each shell command or function in the Bash script, identify the Go file/module where it should be implemented.
- Track progress and notes for each item.

---

## Command Mapping

| Shell Command / Feature                | Go File(s)           | Status      | Notes |
|----------------------------------------|----------------------|-------------|-------|
| Detailed git stats (main menu)         | stats.go             | TODO        |       |
| Commits per author                     | autor.go             | TODO        |       |
| Changelogs                             | changelogs.go        | TODO        |       |
| Changelogs by author                   | changelogs.go        | TODO        |       |
| Daily stats                            | day.go               | TODO        |       |
| Commits per hour                       | hour.go              | TODO        |       |
| Commits per month                      | day.go               | TODO        |       |
| Commits per year                       | day.go               | TODO        |       |
| Commits per weekday                    | day.go               | TODO        |       |
| Commits by author by weekday           | autor.go, day.go     | TODO        |       |
| Commits by hour by author              | autor.go, hour.go    | TODO        |       |
| Commits by timezone                    | day.go               | TODO        |       |
| Commits by author by timezone          | autor.go, day.go     | TODO        |       |
| Branch tree                            | branches.go          | TODO        |       |
| Branches by date                       | branches.go          | TODO        |       |
| Contributors                           | autor.go             | TODO        |       |
| New contributors                       | autor.go             | TODO        |       |
| Suggest reviewers                      | autor.go             | TODO        |       |
| Export CSV                             | stats.go             | TODO        |       |
| Export JSON                            | stats.go             | TODO        |       |
| Heatmap/calendar                       | day.go, hour.go      | TODO        |       |
| Interactive menu                       | stats.go             | TODO        |       |
| Environment variable support           | stats.go             | TODO        |       |
| Backend abstraction (exec/go-git)      | backend.go, exec_backend.go, git.go | IN PROGRESS | Fallback logic added |
| Unix tool replacements (sort, grep...) | unixcompat/          | IN PROGRESS |       |

---

## How to Use
- Mark each feature as `IN PROGRESS` or `DONE` as you implement it.
- Add notes for tricky shell logic or edge cases.
- Reference the original Bash function or line if needed.

---

*This file is a living document. Update as you port features from Bash to Go.*

