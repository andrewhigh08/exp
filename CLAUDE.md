# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A personal Go learning/experiment repository (`github.com/andrewhigh08/exp`, Go 1.27.2).
It is a catalog of small, self-contained programs grouped by topic: algorithms,
concurrency patterns, design patterns, language features, runnable examples, and benchmarks.

`README.md` is the canonical, exhaustive index of every example (in Russian) — consult
it to locate a topic rather than re-scanning the tree.

## Repository shape (non-obvious)

- **Almost every leaf directory is its own `package main` with a single `main.go`.** These
  are not part of a shared library; they are independent programs meant to be run one at a
  time. There is no application entry point or service to "start".

- **There are two separate `go.mod` files:**
  - root: `github.com/andrewhigh08/exp` (Go 1.27.2, requires `golang.org/x/sync`)
  - `code_generation/`: module `repogen` (Go 1.27.2, requires `golang.org/x/tools` and
    `gorm.io/gorm`; declares `tool repogen/cmd/repogen`)

  Consequence: `go build ./...` / `go vet ./...` from the repo root **only covers the root
  module**. The nested `code_generation/` module is excluded — you must `cd` into it to build,
  run, vet, or manage its dependencies. Do not add its imports to the root `go.mod`.

- `aaa/` is a gitignored scratch directory; ignore it.

## Common commands

Run any single example directly (most common workflow):

```bash
go run algorithms/zeros_to_the_right/main.go
go run concurrency/producer_consumer/main.go
go run design_patterns/decorator/main.go
```

Benchmarks (`benchmarks/`, `package bench` — the only tests in the repo):

```bash
cd benchmarks && go test -bench . -benchmem
```

Build / vet the root module (does NOT touch nested modules — see above):

```bash
go build ./...
go vet ./...
gofmt -l .        # no linter config (golangci-lint/Makefile) exists; gofmt is the standard
```

## Code generation (`code_generation/`, module `repogen`)

`repogen` is a CRUD-repository generator. It parses a Go file via `go/ast`, finds structs
tagged with the comment `//repogen:entity`, and emits a GORM v2 (`gorm.io/gorm`) backed
`*_gen.go` repository (Get/Create/Update/Delete) from a single whole-file `text/template`,
formatted with `golang.org/x/tools/imports`. Primary key and column names are resolved with
GORM's own tag parser and default naming strategy (`gorm.io/gorm/schema`).

It is driven entirely by `go generate`: it reads the `GOFILE` env var that `go generate`
sets and **fails if run standalone**. The generator is a module tool (`tool repogen/cmd/repogen`
in `code_generation/go.mod`), and `code_generation/gen.go` runs it with
`//go:generate go tool repogen`, so no `go install` / `PATH` setup is needed. The generated
`gen_gen.go` is committed; regenerate it after changing the template or `gen.go`.

```bash
cd code_generation && go generate ./... && go build ./... && go vet ./...
```
