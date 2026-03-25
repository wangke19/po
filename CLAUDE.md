# Claude Code Instructions

This repo contains `po`, a Polarion ALM CLI tool (modeled after `gh`).

## Working with Polarion

When the user asks about Polarion work items, test cases, test runs, or any Polarion data, use the `po` binary. Always prefer `po` over raw `curl`.

**Build the binary if `/tmp/po` is missing:**
```
go build -o /tmp/po ./cmd/po
```

**Always prefix with:**
```
source ~/.bashrc && POLARION_VERIFY_SSL=false /tmp/po <command>
```
The server uses a self-signed cert. `POLARION_URL`, `POLARION_PROJECT`, and `POLARION_TOKEN` are in `~/.bashrc`.

## Common commands

```bash
po workitem view OCP-84783
po case view OCP-84783
po testcase steps OCP-84783
po testrun list --limit 10
po testrun records 20240426-1546 --result failed
po comment list OCP-84783
po search "type:testcase" --limit 10
```

IDs accept both `OCP-84783` and `OSE/OCP-84783`.

## Project structure

- `pkg/polarion/` — REST API client
- `pkg/cmd/` — Cobra commands, one subdir per group
- `pkg/cmdutil/` — Factory/dependency injection
- `internal/config/` — Config file + env var handling
- `pkg/jsonfields/` — `--json` flag field filtering
