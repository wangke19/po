# po — Polarion CLI

This repo contains `po`, a CLI tool for Polarion ALM (modeled after GitHub CLI `gh`).

## Using Polarion

When the user asks about Polarion work items, test cases, test runs, or any Polarion data, use the `po` binary to interact with the live Polarion instance. Always prefer `po` over raw `curl` calls.

**Build the binary first if `/tmp/po` is missing:**
```
go build -o /tmp/po ./cmd/po
```

**Required env var for SSL:**
```
POLARION_VERIFY_SSL=false
```
Always prefix `po` commands with `POLARION_VERIFY_SSL=false` since the Polarion server uses a self-signed cert.

The `POLARION_URL`, `POLARION_PROJECT`, and `POLARION_TOKEN` are set in `~/.bashrc` — source it before running commands.

## Common commands

```bash
# Work items / test cases
po workitem list --limit 10
po workitem view OCP-84783
po case view OCP-84783
po testcase steps OCP-84783

# Test runs
po testrun list --limit 10
po testrun view 20240426-1546
po testrun records 20240426-1546 --result failed
po testrun records 20240426-1546 --not-run

# Comments & links
po comment list OCP-84783
po link list OCP-84783

# Search
po search "type:testcase status:approved" --limit 10
```

IDs can be passed with or without project prefix (`OCP-84783` or `OSE/OCP-84783` both work).

## Project structure

- `pkg/polarion/` — Polarion REST API client (workitems, testruns, teststeps, comments, links, attachments)
- `pkg/cmd/` — Cobra command implementations, one subdir per command group
- `pkg/cmdutil/` — Factory pattern for dependency injection
- `internal/config/` — Config file + env var handling
- `pkg/jsonfields/` — JSON field filtering for `--json` flag

## Development workflow

- One PR per feature, merged before starting the next
- Build: `go build ./...`
- Test: `go test ./...`
- Integration test: build `/tmp/po` and run against real Polarion (see env vars above)
