# Dependency and Module Release Design

## Context

The repository is a Go workspace containing the root `v2` module, four
published writer modules, and several command, example, and issue-reproduction
modules. Published versions lag behind `main`, and some workspace modules still
require older internal releases. The workspace masks those stale requirements
during local development.

## Goals

- Update external dependencies to the highest versions compatible with each
  module's existing Go version.
- Update every internal `github.com/yeqown/go-qrcode` requirement to the new
  release versions.
- Preserve the current Go directives: Go 1.19 for most modules and Go 1.20 for
  `writer/file`.
- Verify every module and the combined workspace before tagging.
- Create and atomically push the new module tags to `origin`.

## Non-Goals

- Do not raise the Go or CI baseline.
- Do not tag `cmd/qrcode`, `cmd/wasm`, `example`, or `.issues`.
- Do not create GitHub Releases.
- Do not refactor production code outside changes required by dependency
  compatibility.

## External Dependency Policy

Direct dependencies will be upgraded only when their module `go` directive is
compatible with the consuming module:

| Dependency | Selected version | Reason |
| --- | --- | --- |
| `github.com/stretchr/testify` | `v1.11.1` | Latest release; supports Go 1.17 |
| `github.com/urfave/cli/v2` | `v2.27.7` | Latest release; supports Go 1.18 |
| `github.com/mattn/go-runewidth` | `v0.0.17` | Highest release supporting Go 1.19 |
| `golang.org/x/image` | `v0.24.0` | Already at the ceiling; `v0.25.0` requires Go 1.23 |

Transitive dependencies will be recalculated by `go mod tidy` using the
module's existing Go compatibility level. No dependency update may implicitly
raise a module's `go` directive.

## Internal Module Versions

The release versions are:

| Module | New tag |
| --- | --- |
| `github.com/yeqown/go-qrcode/v2` | `v2.3.0` |
| `github.com/yeqown/go-qrcode/writer/standard` | `writer/standard/v1.4.0` |
| `github.com/yeqown/go-qrcode/writer/terminal` | `writer/terminal/v1.1.3` |
| `github.com/yeqown/go-qrcode/writer/file` | `writer/file/v1.0.1` |
| `github.com/yeqown/go-qrcode/writer/compressed` | `writer/compressed/v1.0.2` |

The root and standard modules receive minor versions because `main` contains
new exported functionality since their last tags. The remaining writer modules
receive patch versions for dependency metadata updates.

All workspace consumers will use these versions:

- Writer modules will require root `v2.3.0`.
- `cmd/qrcode` will require root `v2.3.0`, standard `v1.4.0`, and terminal
  `v1.1.3`.
- `cmd/wasm` will require root `v2.3.0` and standard `v1.4.0`.
- `example` will require root `v2.3.0`, standard `v1.4.0`, terminal `v1.1.3`,
  and file `v1.0.1`.
- `.issues` will require root `v2.3.0` and compressed `v1.0.2`.

## Verification

Verification will run in two layers:

1. Run module tests with the module's baseline toolchain and local workspace
   resolution so the not-yet-published internal versions resolve to this
   checkout.
2. Run the full workspace build, tests, race-enabled tests where currently used
   by CI, `go vet`, and module consistency checks.

The final checks must confirm that Go directives are unchanged, no dependency
requires a newer Go version, all internal requirements match the release table,
and the worktree contains only intended changes.

## Commit and Tag Publication

Work is performed on `chore/update-dependencies`. The implementation commit is
the release commit, and all five lightweight tags point to it, matching the
repository's recent tag convention.

Before publication, fetch remote tag state and fail if any target tag exists
locally or remotely with a different object. After every verification passes,
push all five tags to `origin` in one atomic Git push. An atomic push prevents a
partial release if any tag is rejected. The branch itself is not pushed and no
GitHub Release is created as part of this work.

## Failure Handling

- Do not create tags when tests or consistency checks fail.
- Do not push tags when the worktree is dirty after the release commit.
- Do not overwrite or force-update an existing tag.
- If the atomic push is rejected, leave remote tag state unchanged and report
  the rejection before taking further action.
