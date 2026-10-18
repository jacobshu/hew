# Bubble Tea v2 Upgrade Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Move hew from bubbletea v0.27.0 to bubbletea v2.1.0 (latest, 2026-10-08), with matching bubbles and lipgloss.

**Architecture:** v2 is a new module path (`charm.land/...`), not a bump. Bubbletea, bubbles and lipgloss must move together or the build breaks. Move one package at a time, smallest first, keep `go build ./...` green at each commit.

**Tech Stack:** `charm.land/bubbletea/v2 v2.1.0`, `charm.land/bubbles/v2 v2.2.1`, `charm.land/lipgloss/v2 v2.0.6` (versions checked with `go list -m` on 2026-10-18). Needs Go 1.26 (bubbletea v2.1.0 `go` directive). Local toolchain is go1.26.5, fine.

**Spec:** None. Request: "upgrade this golang cli to use the most recent version of bubbletea". Source of truth for API changes: the upgrade guide in the bubbletea v2 repo (UPGRADE_GUIDE_V2.md). Signatures below are from memory of that guide. Confirm each against the module source in `$(go env GOMODCACHE)` before editing.

## Global Constraints

- `go.mod` `go` directive becomes `1.26` or higher (currently `1.22.3`).
- Old `github.com/charmbracelet/{bubbletea,bubbles,lipgloss}` imports must be gone at the end. `go.mod` must not list them.
- `github.com/charmbracelet/*` and `charm.land/*` lipgloss types must never mix in one package. Mixed types fail to compile.
- Behavior stays the same. No new features, no tidy-ups beyond what the API forces.
- No test files exist in the repo. Verification is `go build ./...`, `go vet ./...`, and a manual run of each TUI.

## Review Focus

- Ctrl+C / Esc still quit `hew cht` (key match was `msg.Type`, now `msg.String()`).
- Tab, Shift+Tab, Ctrl+N, Ctrl+P still cycle inputs in `cht`.
- Alt screen and mouse-wheel scrolling still work in `cht` (options moved from `NewProgram` to `View`).
- Terminal resize before first render in `cht` does not panic (viewport built on first `WindowSizeMsg`).
- `hew load` spinner still ticks and `q` still quits.
- Theme colors in `forestfox` still render (lipgloss `Color` is now a function returning `color.Color`).

## File Structure

Files that change. None created, none deleted.

- `go.mod`, `go.sum`: new module paths, Go version.
- `pkg/forestfox/forestfox.go`: `Lipgloss` field type.
- `pkg/load/load.go`, `pkg/load/cmd.go`: spinner model, `View()`, key msg.
- `pkg/cht/cht.go`, `pkg/cht/cmd.go`: textinput, viewport, `View()`, key msgs, program options.
- `colors.go`: comment only, import path in a comment. Update for tidiness.
- `pkg/du`, `pkg/hash`, `pkg/shared`: no charm imports. Do not touch.

---

### Task 1: Baseline and branch

**Files:** none

- [ ] **Step 1: Branch.** `git switch -c chore/bubbletea-v2`
- [ ] **Step 2: Baseline build.** Run `go build ./... && go vet ./...`. Expected: no errors. Record any existing vet warnings so they are not blamed on the upgrade.
- [ ] **Step 3: Baseline run.** Run `go run . cht` and `go run . load` (check `cmd/` for the real subcommand names). Note how each looks and behaves. This is the "before" to compare with.

### Task 2: Swap module paths and Go version

**Files:**
- Modify: `go.mod`, `go.sum`
- Modify: all five importing files (imports only, no logic)

**Interfaces:**
- Produces: imports `tea "charm.land/bubbletea/v2"`, `charm.land/bubbles/v2/{textinput,viewport,spinner}`, `lipgloss "charm.land/lipgloss/v2"` in every file that imported the old paths.

- [ ] **Step 1: Rewrite imports.** In `colors.go` (comment), `pkg/cht/{cht,cmd}.go`, `pkg/forestfox/forestfox.go`, `pkg/load/{load,cmd}.go`, replace `github.com/charmbracelet/bubbletea` with `charm.land/bubbletea/v2`, `.../bubbles/` with `charm.land/bubbles/v2/`, `.../lipgloss` with `charm.land/lipgloss/v2`.
- [ ] **Step 2: Fetch.** Run `go get charm.land/bubbletea/v2@v2.1.0 charm.land/bubbles/v2@v2.2.1 charm.land/lipgloss/v2@v2.0.6 && go mod tidy`. Expected: `go.mod` shows `go 1.26` or higher and no `github.com/charmbracelet/{bubbletea,bubbles,lipgloss}` lines.
- [ ] **Step 3: See the damage.** Run `go build ./... 2>&1 | head -50`. Expected: FAIL in `forestfox`, `load`, `cht`. This list is the work for Tasks 3–5. Do not commit a red build; carry on to Task 3 in the same working tree, commit after Task 5 if the three cannot be green separately. (They can: `forestfox` first, then `load`, then `cht` each build alone.)

### Task 3: forestfox color type

**Files:**
- Modify: `pkg/forestfox/forestfox.go:17` and the `lipgloss.Color(...)` assignments near line 85

**Interfaces:**
- Produces: field `Lipgloss` of type `color.Color` (`image/color`), set from `lipgloss.Color(hex)`. `load` and `cht` pass it to `Foreground(...)` and `BorderForeground(...)`, which take `color.Color` in v2.

- [ ] **Step 1: Change the field type** from `lipgloss.Color` to `color.Color`; add `"image/color"` import.
- [ ] **Step 2: Build the package.** `go build ./pkg/forestfox`. Expected: PASS. The `lipgloss.Color(x)` assignments should compile unchanged since it is now a function returning `color.Color`.
- [ ] **Step 3: Commit.** `git add pkg/forestfox go.mod go.sum && git commit -m "chore: move forestfox to lipgloss v2"` (only if the package builds alone; otherwise fold into Task 5's commit).

### Task 4: Port `hew load`

**Files:**
- Modify: `pkg/load/load.go:86-92` (key msg), `:115` (`View`), `:69-71` (spinner)
- Modify: `pkg/load/cmd.go:17`

**Interfaces:**
- Consumes: `forestfox.GetTheme().X.Lipgloss` as `color.Color` (Task 3).
- Produces: `func (m loadModel) View() tea.View`.

- [ ] **Step 1: Key message.** Change `case tea.KeyMsg:` to `case tea.KeyPressMsg:`. The `msg.String()` checks for `"q"`, `"esc"`, `"ctrl+c"` stay.
- [ ] **Step 2: View.** Change signature to `View() tea.View`. Return `tea.NewView(s)` where it returned `s`.
- [ ] **Step 3: Spinner.** Check `spinner.New()`, `spinner.Points`, `m.spinner.Tick` and `spinner.TickMsg` against the v2 source. Fix whatever the compiler flags. Expect `s.Style` to still take a lipgloss style.
- [ ] **Step 4: NewProgram.** `tea.NewProgram(NewLoadModel())` in `cmd.go` should stand. Fix if the compiler objects.
- [ ] **Step 5: Verify.** `go build ./pkg/load && go vet ./pkg/load`, then run `hew load` and confirm: spinner animates, links get made, `q` quits.
- [ ] **Step 6: Commit.** `git commit -am "chore: port load to bubbletea v2"`

### Task 5: Port `hew cht`

**Files:**
- Modify: `pkg/cht/cht.go:75-150` (`Update`), `:153` (`View`), `:116-135` (viewport setup)
- Modify: `pkg/cht/cmd.go:17-20`

**Interfaces:**
- Consumes: `forestfox` colors (Task 3).
- Produces: `func (m chtModel) View() tea.View` with `AltScreen = true` and `MouseMode = tea.MouseModeCellMotion` set on the returned view.

- [ ] **Step 1: Keys.** Replace `case tea.KeyMsg: switch msg.Type {...}` with `case tea.KeyPressMsg: switch msg.String() {...}`. Map: `"enter"`, `"ctrl+c"`/`"esc"`, `"shift+tab"`/`"ctrl+p"`, `"tab"`/`"ctrl+n"`.
- [ ] **Step 2: Program options.** In `cmd.go` delete `tea.WithAltScreen()` and `tea.WithMouseCellMotion()`. Set both on the view in `View()` (see Interfaces). Keep the comments, move them with the lines.
- [ ] **Step 3: Viewport.** `viewport.New(40, 20)` becomes `viewport.New(viewport.WithWidth(40), viewport.WithHeight(20))`. Width and height are getters now, so `m.viewport.Width` becomes `m.viewport.Width()` in `View()` (the debug `log.Printf` near line 159). Delete `HighPerformanceRendering`, `useHighPerformanceRenderer`, and the `viewport.Sync` branch: v2 dropped them. Check `YPosition` still exists; if not, drop it and its `headerHeight` uses.
- [ ] **Step 4: Style frame sizes.** `m.viewport.Style.GetHorizontalFrameSize()` and `GetWidth()` should compile. Fix if the viewport field was renamed.
- [ ] **Step 5: textinput.** Check `textinput.New()`, `.Focus()` (v2 may return a `tea.Cmd`; collect it into `cmds` if so), `.Blur()`, `.Value()`, `.View()`, and any style fields at `cht.go:48-60`.
- [ ] **Step 6: View.** Change to `View() tea.View`; wrap the final string with `tea.NewView(...)`; set `v.AltScreen` and `v.MouseMode` on it.
- [ ] **Step 7: Verify.** `go build ./... && go vet ./...`. Then run `hew cht` and walk the Review Focus list: type in both inputs, Tab/Shift+Tab/Ctrl+N/Ctrl+P, Enter to fetch, scroll with wheel and keys, resize the window, quit with Esc and with Ctrl+C. Check `logfile` for panics.
- [ ] **Step 8: Commit.** `git commit -am "chore: port cht to bubbletea v2"`

### Task 6: Cleanup and final check

**Files:**
- Modify: `go.mod`, `README.md` (only if it names the old library or Go version)

- [ ] **Step 1:** `go mod tidy && go build ./... && go vet ./...`. Expected: clean.
- [ ] **Step 2:** `grep -rn "github.com/charmbracelet" --include='*.go' --include=go.mod .` Expected: no output.
- [ ] **Step 3:** Rebuild the binary the way you normally install it and run both TUIs once more.
- [ ] **Step 4: Commit** `git commit -am "chore: tidy after bubbletea v2 upgrade"`

## Risks

- **Go 1.26 floor.** Anyone building hew needs Go 1.26+. Fine on this machine. Check CI and other boxes.
- **libsql.** `go-libsql` is cgo. The Go bump should not touch it, but confirm `go build ./pkg/du` still links.
- **API guesses.** The signatures above come from memory. The compiler is the judge. Where a step says "check", read the module source before editing.
- **No tests.** Manual runs are the only safety net. Do not skip Task 5 Step 7.
