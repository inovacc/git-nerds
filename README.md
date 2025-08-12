# git-quick-stats (Go Port)

Simple and efficient way to access various statistics in a git repository.

---

## Overview
This project is a Go port of the original [`git-quick-stats`](https://github.com/git-quick-stats/git-quick-stats) Bash script. It provides a fast, scriptable, and interactive way to extract nerdy statistics from any local git repository.

## Features
- Contribution stats per author (commits, insertions, deletions, lines changed, files)
- Changelogs (overall and by author)
- Daily, monthly, yearly, weekday, and hourly commit stats
- Calendar and heatmap visualizations
- Branch tree and branch-by-date
- Contributors and new contributors
- Reviewer suggestions
- Export to CSV and JSON
- Interactive and non-interactive CLI

## Usage

### Non-interactive mode
```sh
git-quick-stats [OPTIONS]
```

### Interactive mode
```sh
git-quick-stats
```

## Options (from original Bash/man)

### Generate
- `-T`, `--detailed-git-stats` — detailed list of git stats
- `-R`, `--git-stats-by-branch` — stats by branch
- `-c`, `--changelogs` — changelogs
- `-L`, `--changelogs-by-author` — changelogs by author
- `-S`, `--my-daily-stats` — your current daily stats
- `-V`, `--csv-output-by-branch` — daily stats by branch in CSV
- `-j`, `--json-output` — save git log as JSON

### List
- `-b`, `--branch-tree` — ASCII graph of branch history
- `-D`, `--branches-by-date` — branches by date
- `-C`, `--contributors` — list contributors
- `-n`, `--new-contributors` — new contributors since date
- `-a`, `--commits-per-author` — commits per author
- `-d`, `--commits-per-day` — commits per day
- `-m`, `--commits-by-month` — commits per month
- `-Y`, `--commits-by-year` — commits per year
- `-w`, `--commits-by-weekday` — commits per weekday
- `-W`, `--commits-by-author-by-weekday` — commits per weekday by author
- `-o`, `--commits-by-hour` — commits per hour
- `-A`, `--commits-by-author-by-hour` — commits per hour by author
- `-z`, `--commits-by-timezone` — commits per timezone
- `-Z`, `--commits-by-author-by-timezone` — commits per timezone by author

### Calendar
- `-k`, `--commits-calendar-by-author` — calendar heatmap by author
- `-H`, `--commits-heatmap` — heatmap of commits per day/hour (last 30 days)

### Suggest
- `-r`, `--suggest-reviewers` — suggest code reviewers
- `-h`, `-?`, `--help` — show help

## Environment Variables
- `_GIT_SINCE`, `_GIT_UNTIL` — limit git time log (e.g. `export _GIT_SINCE="2017-01-20"`)
- `_GIT_LIMIT` — limit output log (e.g. `export _GIT_LIMIT=20`)
- `_GIT_LOG_OPTIONS` — extra git log options
- `_GIT_PATHSPEC` — exclude files/dirs (e.g. `export _GIT_PATHSPEC=':!pattern'`)
- `_GIT_MERGE_VIEW` — show/hide merge commits (`enable`, `exclusive`)
- `_GIT_SORT_BY` — sort stats by field/order (e.g. `export _GIT_SORT_BY="commits-desc"`)
- `_MENU_THEME` — color theme (`default`, `legacy`, `none`)
- `_GIT_BRANCH` — branch to analyze
- `_GIT_IGNORE_AUTHORS` — regex to filter authors
- `_GIT_DAYS` — number of days for heatmap

## Example
```sh
export _GIT_SINCE="2023-01-01"
export _GIT_LIMIT=50
git-quick-stats --detailed-git-stats
```

## Project Structure
- `cmd/` — CLI entrypoint (planned)
- `pkg/` — Go logic (stats, backend, unixcompat, etc.)
- `testdata/` — test repositories and fixtures
- `git-quick-stats` — original Bash script (reference)
- `git-quick-stats.1` — original man page (reference)

## Migration Notes
- This Go module aims to replicate all features of the Bash version.
- The Bash script and man page are kept for reference and as documentation.

## See Also
- [git-quick-stats (original Bash)](https://github.com/git-quick-stats/git-quick-stats)
- `man ./git-quick-stats.1`

---

*This README is generated from the original Bash script and man page for clarity and completeness.*

