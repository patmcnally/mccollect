# mccollect

A terminal application for managing your Marvel Champions: The Card Game collection. Browse and toggle pack ownership in an interactive TUI, or script everything via CLI flags for use with agents and automation.

## Features

- **Full card database** — imports all cards from [marvelsdb-json-data](https://github.com/zzorba/marvelsdb-json-data) into a local SQLite database
- **Incremental updates** — git pulls the data repo and re-imports only changed files
- **Collection tracking** — toggle pack ownership with space/enter in the TUI, or set it directly via CLI
- **Multiple collections** — name your collection (e.g. `Pat`, `Friend`) to track multiple owners
- **Agent-friendly** — every TUI action has a CLI equivalent; pass `--json` for structured output

## Installation

```bash
# From the project root
go build -o mccollect .

# Or install to $GOPATH/bin
go install .
```

Requires Go 1.26+. Uses [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) (no CGo required).

## Quick Start

```bash
# 1. Clone card data
git clone https://github.com/zzorba/marvelsdb-json-data ~/src/marvelsdb-json-data

# 2. Build the database
mccollect import --data ~/src/marvelsdb-json-data

# 3. Import your collection from marvelcdb.com (see below)
mccollect collection import-html --name Pat collection.html

# 4. Launch the TUI
mccollect --name Pat
```

## Commands

### `mccollect` (TUI)

Launches the interactive terminal UI. Use `--name` to select which collection to edit.

```bash
mccollect
mccollect --name Pat
mccollect --name Pat --db /path/to/other/cards.db
```

**TUI keys:**

| Key | Action |
|---|---|
| `j` / `k` or `↑` / `↓` | Navigate |
| `space` / `enter` | Toggle pack owned/not-owned |
| `tab` | Switch between Packs and Stats views |
| `q` / `ctrl+c` | Quit |

---

### `mccollect import`

Full rebuild of the database from a local marvelsdb-json-data clone. Wipes and recreates all card data (collection ownership is preserved).

```bash
mccollect import --data ~/src/marvelsdb-json-data
mccollect import --data ~/src/marvelsdb-json-data --db /path/to/cards.db
mccollect import --data ~/src/marvelsdb-json-data --json
```

---

### `mccollect update`

Incremental update — git pulls the data repo and re-imports only changed pack files.

```bash
mccollect update --data ~/src/marvelsdb-json-data
mccollect update --data ~/src/marvelsdb-json-data --dry-run   # preview only
mccollect update --data ~/src/marvelsdb-json-data --json
```

---

### `mccollect collection import-html`

Import pack ownership from a saved marvelcdb.com collection page.

1. Log in to [marvelcdb.com](https://marvelcdb.com) and visit `/collection/packs`
2. Save the page as HTML: **File → Save Page As → Webpage, HTML Only**
3. Run:

```bash
mccollect collection import-html --name Pat collection.html
mccollect collection import-html --name Pat collection.html --json
```

---

### `mccollect collection list`

List all packs with owned (`*`) or not-owned status, grouped by type.

```bash
mccollect collection list --name Pat
mccollect collection list --name Pat --json
```

---

### `mccollect collection set`

Set the ownership status of a specific pack by its pack code.

```bash
mccollect collection set --name Pat --pack core --owned
mccollect collection set --name Pat --pack cyclops --owned
mccollect collection set --name Pat --pack green_goblin --not-owned
mccollect collection set --name Pat --pack core --owned --json
```

Pack codes match the identifiers in marvelsdb-json-data (e.g. `core`, `cyclops`, `mut_gen`).

---

## Global Flags

| Flag | Default | Description |
|---|---|---|
| `--db` | `~/.config/mccollect/cards.db` | Path to the SQLite database |
| `--json` | `false` | Output JSON instead of human-readable text |

## Agent Usage

All commands support `--json` for machine-readable output. Example workflow:

```bash
# Check which packs are owned (returns JSON array)
mccollect collection list --name Pat --json

# Set ownership from a script
mccollect collection set --name Pat --pack cyclops --owned --json

# Get import stats
mccollect import --data ~/src/marvelsdb-json-data --json
# → {"packs":60,"sets":369,"cards":4200,"commit":"8dbe3ab..."}
```

## Database

The SQLite database at `~/.config/mccollect/cards.db` (default) contains:

| Table | Contents |
|---|---|
| `cards` | All ~4,200 cards with every field |
| `packs` | All packs with code, name, type, release date |
| `sets` | All card sets |
| `collections` | Named collections (one per player) |
| `collection_packs` | Pack ownership per collection |
| `_meta` | Last import commit, data root path |

The database is the single source of truth. See `../skills/marvel-champions/references/db_schema.md` for the full schema and common query patterns.

## Project Layout

```
mccollect/
  main.go
  cmd/                         # Cobra CLI commands
    root.go                    # --db, --json flags; launches TUI by default
    import.go                  # mccollect import
    update.go                  # mccollect update
    tui.go                     # TUI launcher
    collection.go              # collection parent command
    collection_list.go
    collection_set.go
    collection_import_html.go
  db/                          # SQLite layer (all SQL)
    db.go, schema.go, meta.go
    packs.go, sets.go, cards.go, collection.go
  model/                       # Shared structs
    model.go
  importer/                    # JSON parsing, git ops, HTML scraping
    importer.go, packs.go, sets.go, cards.go, git.go, html.go
  tui/                         # Bubbletea views
    app.go, packs.go, stats.go, style.go
```

## Dependencies

- [cobra](https://github.com/spf13/cobra) — CLI framework
- [bubbletea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [lipgloss](https://github.com/charmbracelet/lipgloss) — terminal styling
- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) — CGo-free SQLite driver
