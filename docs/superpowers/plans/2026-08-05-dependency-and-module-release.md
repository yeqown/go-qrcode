# Dependency and Module Release Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Upgrade external and internal Go module dependencies without raising the existing Go baselines, then publish five synchronized module tags to `origin`.

**Architecture:** Treat each `go.mod` as an independently versioned module while using temporary Go workspaces to test unpublished internal versions against this checkout. Update compatible external dependencies first, align internal version requirements second, and publish lightweight tags only after clean baseline-specific and workspace-wide verification.

**Tech Stack:** Go modules and workspaces, Go 1.19.13, Go 1.20.14, Git, GitHub remote tags

---

### Task 1: Establish the Release Baseline

**Files:**
- Verify: `go.mod`
- Verify: `go.work`
- Verify: `.github/workflows/go.yml`
- Verify: all nested `go.mod` files

- [ ] **Step 1: Confirm the branch and clean worktree**

Run from the repository root:

```bash
git branch --show-current
git status --short
```

Expected: the branch is `chore/update-dependencies` and `git status --short`
prints nothing.

- [ ] **Step 2: Confirm every target tag is unused locally**

```bash
git tag --list \
  v2.3.0 \
  writer/standard/v1.4.0 \
  writer/terminal/v1.1.3 \
  writer/file/v1.0.1 \
  writer/compressed/v1.0.2
```

Expected: no output.

- [ ] **Step 3: Confirm every target tag is unused on `origin`**

```bash
git ls-remote --tags origin \
  refs/tags/v2.3.0 \
  refs/tags/writer/standard/v1.4.0 \
  refs/tags/writer/terminal/v1.1.3 \
  refs/tags/writer/file/v1.0.1 \
  refs/tags/writer/compressed/v1.0.2
```

Expected: no output. Stop if any ref is returned; never overwrite a release
tag.

- [ ] **Step 4: Make the baseline toolchains available**

```bash
GOTOOLCHAIN=go1.19.13 go version
GOTOOLCHAIN=go1.20.14 go version
```

Expected: `go version go1.19.13` and `go version go1.20.14`.

- [ ] **Step 5: Prepare ignored test fixture directories**

```bash
mkdir -p writer/standard/testdata writer/terminal/testdata
```

Expected: both directories exist and `git status --short` remains empty.

- [ ] **Step 6: Run the local-workspace baseline tests**

```bash
for module_dir in . .issues cmd/qrcode example writer/compressed writer/file writer/standard writer/terminal; do
  (
    cd "$module_dir"
    go test ./...
  )
done
```

```bash
installed_go_root=$(go env GOROOT)
(
  cd cmd/wasm
  GOOS=js GOARCH=wasm go test \
    -exec="$installed_go_root/misc/wasm/go_js_wasm_exec" \
    ./...
)
```

Expected: every command exits 0 using local workspace resolution. The workspace
baseline is required because `.issues` and `example` already consume APIs from
the current checkout that do not exist in their stale published requirements,
and nested modules intentionally do not track individual `go.sum` files.

### Task 2: Upgrade Compatible External Dependencies

**Files:**
- Modify: `go.mod`
- Modify: `go.sum`
- Modify: `writer/standard/go.mod`
- Modify: `writer/terminal/go.mod`
- Modify: `cmd/qrcode/go.mod`
- Possibly modify: `go.work.sum`

- [ ] **Step 1: Upgrade root and standard test dependencies**

```bash
GOTOOLCHAIN=go1.19.13 GOWORK=off go get github.com/stretchr/testify@v1.11.1
```

```bash
(
  cd writer/standard
  GOTOOLCHAIN=go1.19.13 GOWORK=off go get github.com/stretchr/testify@v1.11.1
)
```

Expected: root and standard now require `github.com/stretchr/testify v1.11.1`.

- [ ] **Step 2: Upgrade the terminal width dependency to the Go 1.19 ceiling**

```bash
(
  cd writer/terminal
  GOTOOLCHAIN=go1.19.13 GOWORK=off go get github.com/mattn/go-runewidth@v0.0.17
)
```

Expected: terminal requires `github.com/mattn/go-runewidth v0.0.17`, not a
newer release requiring Go 1.20 or Go 1.23.

- [ ] **Step 3: Upgrade the CLI dependency**

```bash
(
  cd cmd/qrcode
  GOTOOLCHAIN=go1.19.13 GOWORK=off go get github.com/urfave/cli/v2@v2.27.7
)
```

Expected: `cmd/qrcode/go.mod` requires `github.com/urfave/cli/v2 v2.27.7` and
its compatible transitive dependency versions.

- [ ] **Step 4: Tidy each externally changed module at Go 1.19 compatibility**

```bash
GOTOOLCHAIN=go1.19.13 GOWORK=off go mod tidy -go=1.19 -compat=1.19
```

```bash
for module_dir in writer/standard writer/terminal; do
  (
    cd "$module_dir"
    GOTOOLCHAIN=go1.19.13 GOWORK=off go mod tidy -go=1.19 -compat=1.19
  )
done
```

Expected: all commands exit 0 and retain `go 1.19`. Defer `cmd/qrcode` tidy
until Task 3 because its stale published writer/core combination cannot be
loaded consistently outside the workspace.

- [ ] **Step 5: Verify the selected direct versions and unchanged image ceiling**

```bash
rg -n \
  'github.com/stretchr/testify v1.11.1|github.com/urfave/cli/v2 v2.27.7|github.com/mattn/go-runewidth v0.0.17|golang.org/x/image v0.24.0' \
  go.mod writer/standard/go.mod writer/terminal/go.mod cmd/qrcode/go.mod
```

Expected: both testify requirements are `v1.11.1`, CLI is `v2.27.7`,
runewidth is `v0.0.17`, and x/image remains `v0.24.0`.

- [ ] **Step 6: Re-run tests for externally changed modules**

```bash
GOTOOLCHAIN=go1.19.13 GOWORK=off go test ./...
```

```bash
for module_dir in writer/standard writer/terminal; do
  (
    cd "$module_dir"
    GOTOOLCHAIN=go1.19.13 GOWORK=off go test ./...
  )
done
```

```bash
(
  cd cmd/qrcode
  GOTOOLCHAIN=go1.19.13 GOWORK=off go test -mod=mod ./...
)
```

Expected: every test command exits 0.

- [ ] **Step 7: Review and commit the external dependency update**

```bash
git diff --check
git diff -- go.mod go.sum writer/standard/go.mod writer/terminal/go.mod cmd/qrcode/go.mod go.work.sum
git status --short
```

Expected: only dependency metadata has changed; ignored nested `go.sum` files
must not be staged.

```bash
git add go.mod go.sum writer/standard/go.mod writer/terminal/go.mod cmd/qrcode/go.mod go.work.sum
git commit -m "chore(deps): update compatible external dependencies"
```

Expected: one dependency commit is created. If `go.work.sum` did not change,
omit it from `git add`.

### Task 3: Align Internal Module Requirements

**Files:**
- Modify: `.issues/go.mod`
- Modify: `cmd/qrcode/go.mod`
- Modify: `cmd/wasm/go.mod`
- Modify: `example/go.mod`
- Modify: `writer/compressed/go.mod`
- Modify: `writer/file/go.mod`
- Modify: `writer/standard/go.mod`
- Modify: `writer/terminal/go.mod`

- [ ] **Step 1: Point every writer module at root `v2.3.0`**

```bash
for module_dir in writer/compressed writer/file writer/standard writer/terminal; do
  (
    cd "$module_dir"
    go mod edit -require=github.com/yeqown/go-qrcode/v2@v2.3.0
  )
done
```

Expected: all four writer `go.mod` files require root `v2.3.0`.

- [ ] **Step 2: Align command module requirements**

```bash
(
  cd cmd/qrcode
  go mod edit \
    -require=github.com/yeqown/go-qrcode/v2@v2.3.0 \
    -require=github.com/yeqown/go-qrcode/writer/standard@v1.4.0 \
    -require=github.com/yeqown/go-qrcode/writer/terminal@v1.1.3
)
```

```bash
(
  cd cmd/wasm
  go mod edit \
    -require=github.com/yeqown/go-qrcode/v2@v2.3.0 \
    -require=github.com/yeqown/go-qrcode/writer/standard@v1.4.0
)
```

Expected: command modules contain exactly the approved internal versions.

- [ ] **Step 3: Align example and issue-reproduction requirements**

```bash
(
  cd example
  go mod edit \
    -require=github.com/yeqown/go-qrcode/v2@v2.3.0 \
    -require=github.com/yeqown/go-qrcode/writer/file@v1.0.1 \
    -require=github.com/yeqown/go-qrcode/writer/standard@v1.4.0 \
    -require=github.com/yeqown/go-qrcode/writer/terminal@v1.1.3
)
```

```bash
(
  cd .issues
  go mod edit \
    -require=github.com/yeqown/go-qrcode/v2@v2.3.0 \
    -require=github.com/yeqown/go-qrcode/writer/compressed@v1.0.2
)
```

Expected: example and issue modules contain exactly the approved versions;
`example` retains `github.com/yeqown/go-qrcode v1.5.10` for v1 examples.

- [ ] **Step 4: Verify there are no stale internal requirements**

```bash
rg -n 'github.com/yeqown/go-qrcode' --hidden -g 'go.mod'
```

Expected internal release requirements:

```text
github.com/yeqown/go-qrcode/v2 v2.3.0
github.com/yeqown/go-qrcode/writer/standard v1.4.0
github.com/yeqown/go-qrcode/writer/terminal v1.1.3
github.com/yeqown/go-qrcode/writer/file v1.0.1
github.com/yeqown/go-qrcode/writer/compressed v1.0.2
```

The only retained old-major requirement is
`github.com/yeqown/go-qrcode v1.5.10` in `example/go.mod`.

- [ ] **Step 5: Tidy the CLI with local modules and synchronize the workspace**

```bash
(
  cd cmd/qrcode
  go mod tidy -go=1.19 -compat=1.19
)

go work sync
```

Expected: local workspace modules satisfy the unpublished internal version
requirements, CLI metadata is tidy, and no permanent `replace` directives are
added.

- [ ] **Step 6: Review and commit the internal dependency alignment**

```bash
git diff --check
git diff -- .issues/go.mod cmd/qrcode/go.mod cmd/wasm/go.mod example/go.mod writer/compressed/go.mod writer/file/go.mod writer/standard/go.mod writer/terminal/go.mod go.work.sum
```

Expected: only approved version and checksum changes.

```bash
git add .issues/go.mod cmd/qrcode/go.mod cmd/wasm/go.mod example/go.mod writer/compressed/go.mod writer/file/go.mod writer/standard/go.mod writer/terminal/go.mod go.work.sum
git commit -m "chore(deps): align internal module releases"
```

Expected: the internal dependency commit is created. If `go.work.sum` is
unchanged, omit it from `git add`.

### Task 4: Verify Go Baselines and Workspace Integration

**Files:**
- Verify: every workspace module
- Possibly modify: `go.work.sum`

- [ ] **Step 1: Create a temporary Go 1.19 workspace without `writer/file`**

```bash
go119_workspace_dir=$(mktemp -d /private/tmp/go-qrcode-go119.XXXXXX)
(
  cd "$go119_workspace_dir"
  GOTOOLCHAIN=go1.19.13 go work init \
    /Users/apple/projects/opensource/go-qrcode \
    /Users/apple/projects/opensource/go-qrcode/.issues \
    /Users/apple/projects/opensource/go-qrcode/cmd/qrcode \
    /Users/apple/projects/opensource/go-qrcode/cmd/wasm \
    /Users/apple/projects/opensource/go-qrcode/writer/compressed \
    /Users/apple/projects/opensource/go-qrcode/writer/standard \
    /Users/apple/projects/opensource/go-qrcode/writer/terminal
)
```

Expected: `$go119_workspace_dir/go.work` exists with `go 1.19` and local uses.

- [ ] **Step 2: Test and vet every Go 1.19 release consumer**

```bash
for module_dir in . .issues cmd/qrcode writer/compressed writer/standard writer/terminal; do
  (
    cd "$module_dir"
    GOTOOLCHAIN=go1.19.13 GOWORK="$go119_workspace_dir/go.work" go test ./...
    GOTOOLCHAIN=go1.19.13 GOWORK="$go119_workspace_dir/go.work" go vet ./...
  )
done
```

```bash
go119_root=$(GOTOOLCHAIN=go1.19.13 go env GOROOT)
(
  cd cmd/wasm
  GOOS=js GOARCH=wasm GOTOOLCHAIN=go1.19.13 GOWORK="$go119_workspace_dir/go.work" go test \
    -exec="$go119_root/misc/wasm/go_js_wasm_exec" \
    ./...
  GOOS=js GOARCH=wasm GOTOOLCHAIN=go1.19.13 GOWORK="$go119_workspace_dir/go.work" go vet ./...
)
```

Expected: all tests and vet commands exit 0 under Go 1.19.13, with the WASM
module compiled and tested for its declared `js/wasm` build target.

- [ ] **Step 3: Run CI-equivalent race tests under Go 1.19**

```bash
for module_dir in . writer/standard writer/terminal; do
  (
    cd "$module_dir"
    GOTOOLCHAIN=go1.19.13 GOWORK="$go119_workspace_dir/go.work" go test -race ./...
  )
done
```

Expected: all three race test commands exit 0.

- [ ] **Step 4: Create a temporary Go 1.20 workspace for file and example**

```bash
go120_workspace_dir=$(mktemp -d /private/tmp/go-qrcode-go120.XXXXXX)
(
  cd "$go120_workspace_dir"
  GOTOOLCHAIN=go1.20.14 go work init \
    /Users/apple/projects/opensource/go-qrcode \
    /Users/apple/projects/opensource/go-qrcode/example \
    /Users/apple/projects/opensource/go-qrcode/writer/file \
    /Users/apple/projects/opensource/go-qrcode/writer/standard \
    /Users/apple/projects/opensource/go-qrcode/writer/terminal
)
```

Expected: `$go120_workspace_dir/go.work` exists and resolves every dependency
used by `writer/file` and `example` locally.

- [ ] **Step 5: Test and vet file and example under Go 1.20**

```bash
for module_dir in writer/file example; do
  (
    cd "$module_dir"
    GOTOOLCHAIN=go1.20.14 GOWORK="$go120_workspace_dir/go.work" go test ./...
    GOTOOLCHAIN=go1.20.14 GOWORK="$go120_workspace_dir/go.work" go vet ./...
  )
done
```

Expected: all commands exit 0 under Go 1.20.14.

- [ ] **Step 6: Run combined workspace integration with the installed toolchain**

```bash
for module_dir in . .issues cmd/qrcode example writer/compressed writer/file writer/standard writer/terminal; do
  (
    cd "$module_dir"
    go test ./...
  )
done
```

```bash
installed_go_root=$(go env GOROOT)
(
  cd cmd/wasm
  GOOS=js GOARCH=wasm go test \
    -exec="$installed_go_root/misc/wasm/go_js_wasm_exec" \
    ./...
)
```

```bash
go build -race ./...
```

Expected: every module test and the root race build exit 0.

- [ ] **Step 7: Verify directives, module pins, and formatting**

```bash
rg -n '^go ' --hidden -g 'go.mod' -g 'go.work'
rg -n 'github.com/yeqown/go-qrcode' --hidden -g 'go.mod'
git diff --check
git status --short
```

Expected: all existing Go directives are unchanged, all internal versions match
Task 3, and only an intentional `go.work.sum` refresh may remain uncommitted.

- [ ] **Step 8: Commit a checksum refresh only when required**

If `go.work.sum` changed after verification:

```bash
git add go.work.sum
git commit -m "chore(deps): refresh workspace checksums"
```

Expected: `git status --short` prints nothing after the optional commit. Do not
create an empty commit when `go.work.sum` is unchanged.

### Task 5: Create and Atomically Publish Module Tags

**Files:**
- Verify: Git commit and tag metadata only

- [ ] **Step 1: Run the final release gate**

```bash
git status --short
git diff --check origin/main...HEAD
git log --oneline --decorate origin/main..HEAD
```

Expected: the worktree is clean, the diff check exits 0, and the log contains
the approved design, plan, external dependency, and internal dependency commits
plus an optional workspace checksum commit.

- [ ] **Step 2: Recheck remote tag absence immediately before tagging**

```bash
git ls-remote --tags origin \
  refs/tags/v2.3.0 \
  refs/tags/writer/standard/v1.4.0 \
  refs/tags/writer/terminal/v1.1.3 \
  refs/tags/writer/file/v1.0.1 \
  refs/tags/writer/compressed/v1.0.2
```

Expected: no output. Stop if any target ref exists.

- [ ] **Step 3: Create lightweight local tags on the verified release commit**

```bash
git tag v2.3.0
git tag writer/standard/v1.4.0
git tag writer/terminal/v1.1.3
git tag writer/file/v1.0.1
git tag writer/compressed/v1.0.2
```

Expected: all commands exit 0 without moving or replacing another tag.

- [ ] **Step 4: Confirm every tag points to `HEAD` and is lightweight**

```bash
git rev-parse HEAD \
  v2.3.0 \
  writer/standard/v1.4.0 \
  writer/terminal/v1.1.3 \
  writer/file/v1.0.1 \
  writer/compressed/v1.0.2
```

Expected: the same commit SHA is printed six times.

```bash
git cat-file -t v2.3.0
git cat-file -t writer/standard/v1.4.0
git cat-file -t writer/terminal/v1.1.3
git cat-file -t writer/file/v1.0.1
git cat-file -t writer/compressed/v1.0.2
```

Expected: each command prints `commit`.

- [ ] **Step 5: Atomically push all tags to `origin`**

```bash
git push --atomic origin \
  refs/tags/v2.3.0 \
  refs/tags/writer/standard/v1.4.0 \
  refs/tags/writer/terminal/v1.1.3 \
  refs/tags/writer/file/v1.0.1 \
  refs/tags/writer/compressed/v1.0.2
```

Expected: Git reports all five new tags. If the push fails, the atomic operation
must leave all five remote refs absent; do not retry with `--force`.

- [ ] **Step 6: Verify remote tag SHAs**

```bash
git ls-remote --tags origin \
  refs/tags/v2.3.0 \
  refs/tags/writer/standard/v1.4.0 \
  refs/tags/writer/terminal/v1.1.3 \
  refs/tags/writer/file/v1.0.1 \
  refs/tags/writer/compressed/v1.0.2
```

Expected: five refs are returned and every SHA equals `git rev-parse HEAD`.
No branch push or GitHub Release creation is part of this task.
