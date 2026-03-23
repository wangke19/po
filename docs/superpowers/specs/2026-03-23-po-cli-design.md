# po CLI Design Spec

**Date:** 2026-03-23
**Status:** Approved

## Overview

`po` is a command-line tool for Siemens Polarion ALM, modeled after `gh` (GitHub CLI). It enables AI Agent invocation through structured JSON output, inherited authentication, and a rich command set. Any Polarion instance is supported — the hostname is configured at login time.

## Goals

- AI Agent-friendly: structured `--json field1,field2` output on all read commands
- `po api` passthrough for arbitrary REST calls with automatic auth
- Secure credential storage (system keyring + env var fallback)
- Scriptable: non-interactive flags, clear exit codes, stdin support
- Rich subcommands covering test case, test run, and generic work item workflows

## Non-Goals

- SOAP API support (Bearer token only for now)
- Web browser OAuth flow
- Plugin/extension system
- Delete operations for work items, test cases, or test runs (Polarion soft-deletion semantics are complex; use `po api` for now)

---

## Project Structure

```
po/
├── cmd/po/main.go
├── internal/
│   ├── build/              # version/date injection via ldflags
│   ├── config/             # YAML config + keyring abstraction
│   └── pocmd/              # Main() entry, root wiring
├── pkg/
│   ├── cmd/
│   │   ├── root/           # NewCmdRoot, registers all subcommands
│   │   ├── auth/           # login, logout, status, token
│   │   ├── api/            # po api <endpoint> passthrough
│   │   ├── case/           # list, view, create, edit
│   │   ├── testrun/        # list, view, create, result
│   │   └── workitem/       # list, view, create, edit (generic)
│   ├── cmdutil/            # Factory, error helpers, auth check
│   ├── iostreams/          # color, pager, TTY detection
│   └── polarion/           # typed REST API client
├── go.mod                  # module: github.com/wangke19/po
└── Makefile
```

---

## Configuration

**File:** `~/.config/po/config.yml`

```yaml
hosts:
  <hostname>:
    default_project: <project-id>
    verify_ssl: true
# token stored in system keyring under key: po:<hostname>
```

**Env var overrides** (higher priority than config file):

| Variable | Purpose |
|---|---|
| `POLARION_URL` | Polarion host URL |
| `POLARION_TOKEN` | Bearer token |
| `POLARION_PROJECT` | Default project ID |
| `POLARION_VERIFY_SSL` | `true`/`false` |

---

## Core Interfaces

### Factory

Injected into every command constructor:

```go
type Factory struct {
    AppVersion     string
    IOStreams       *iostreams.IOStreams
    Config         func() (config.Config, error)
    HttpClient     func() (*http.Client, error)
    PolarionClient func() (*polarion.Client, error)
}
```

### Polarion REST Client (`pkg/polarion/client.go`)

Thin, typed wrapper over `{host}/polarion/rest/v1`. The `project` field is baked into the base URL path for all project-scoped endpoints (e.g., `GET /projects/{project}/workitems/{id}`).

```go
type Client struct {
    baseURL    string // e.g. https://polarion.example.com/polarion/rest/v1
    token      string
    project    string // default project; embedded in URL paths
    httpClient *http.Client
}

// Key domain types (JSON tags drive FilterFields):
type WorkItem struct {
    ID          string `json:"id"`
    Title       string `json:"title"`
    Type        string `json:"type"`
    Status      string `json:"status"`
    Author      string `json:"author"`
    Description string `json:"description"`
    URL         string `json:"url"`
}

type WorkItemInput struct {
    Title       string `json:"title,omitempty"`
    Type        string `json:"type,omitempty"`
    Status      string `json:"status,omitempty"`
    Description string `json:"description,omitempty"`
}

type TestRun struct {
    ID       string `json:"id"`
    Title    string `json:"title"`
    Status   string `json:"status"`
    Template string `json:"template"`
    URL      string `json:"url"`
}

type TestRunInput struct {
    Title    string `json:"title,omitempty"`
    Template string `json:"template,omitempty"`
}

type TestResult struct {
    Result  string `json:"result"` // passed|failed|blocked
    Comment string `json:"comment,omitempty"`
}

// All methods use project-scoped URL paths, e.g.:
//   GET  /projects/{project}/workitems?query=...
//   GET  /projects/{project}/workitems/{id}
//   POST /projects/{project}/workitems
//   PATCH /projects/{project}/workitems/{id}
func (c *Client) ListWorkItems(ctx context.Context, query string, limit int) ([]WorkItem, error)
func (c *Client) GetWorkItem(ctx context.Context, id string) (*WorkItem, error)
func (c *Client) CreateWorkItem(ctx context.Context, in WorkItemInput) (*WorkItem, error)
func (c *Client) UpdateWorkItem(ctx context.Context, id string, in WorkItemInput) (*WorkItem, error)
func (c *Client) ListTestRuns(ctx context.Context, query string, limit int) ([]TestRun, error)
func (c *Client) GetTestRun(ctx context.Context, id string) (*TestRun, error)
func (c *Client) CreateTestRun(ctx context.Context, in TestRunInput) (*TestRun, error)
func (c *Client) UpdateTestRunResult(ctx context.Context, runID, caseID string, result TestResult) error
```

### JSON Field Selector

Shared utility used by all read commands. It works on the JSON-marshalled form of a value: marshal to JSON first, then filter top-level keys. This avoids reflection on unexported struct fields.

```go
// FilterFields marshals v to JSON, then returns a new JSON object containing
// only the keys named in fields. If fields is nil/empty, all keys are returned.
func FilterFields(v any, fields []string) (json.RawMessage, error)
```

Usage: `po case list --json id,title,status,url`

On write commands (`create`, `edit`, `result`), `--json` with no field argument prints the full created/updated resource as JSON (all fields). The field selector is optional on write commands.

---

## Command Surface

### `po auth`

```
po auth login   --hostname <url> --project <id> [--with-token]
po auth logout  [--hostname <url>]
po auth status
po auth token
```

- `login`: prompts for token interactively or reads from stdin (`--with-token`), validates via `GET /projects/{project}`, stores token in system keyring
- `status`: prints host, project, token validity
- `token`: prints raw token for scripting

### `po api`

```
po api <endpoint-path>
       [--method GET|POST|PATCH|DELETE]
       [-f key=value ...]
       [-H "Header: value" ...]
       [--paginate]
       [--input <file>]
```

- Substitutes `{project}` with the configured default project (reads from env var `POLARION_PROJECT` first, then `~/.config/po/config.yml`)
- Adds `Authorization: Bearer <token>` automatically (same resolution order: env var → keyring)
- Pretty-prints JSON on TTY, raw on pipe
- `--paginate`: GET only; follows Polarion REST pagination by reading `links.next` from each response envelope and concatenating `data` arrays until no `next` link is present. No-op for non-paginated endpoints. Not valid with POST/PATCH/DELETE.

### `po case` (workitem type=testcase)

```
po case list   [--query <lucene>] [--limit N] [--json id,title,status,url]
po case view   <id>               [--json id,title,status,url,description]
po case create -t <title> [-d <desc>] [--status draft|approved] [--json [field,...]]
po case edit   <id> [-t <title>] [-d <desc>] [--status <s>] [--json [field,...]]
```

On `create`/`edit`, `--json` with no field list prints all fields of the created/updated resource.

### `po testrun`

```
po testrun list    [--query <lucene>] [--limit N] [--json id,title,status,url]
po testrun view    <id>               [--json id,title,status,url,template]
po testrun create  -t <title> [--template <tpl>] [--json [field,...]]
po testrun result  <run-id> <case-id> --result passed|failed|blocked [--comment <s>]
```

`po testrun result` prints nothing on success (exit 0) and an error message on stderr on failure. Add `--json` to print the updated test record as JSON.

### `po workitem` (generic)

```
po workitem list   --type <type> [--query <lucene>] [--limit N] [--json id,title,type,status,url]
po workitem view   <id>          [--json id,title,type,status,url,description]
po workitem create --type <type> -t <title> [-d <desc>] [--json [field,...]]
po workitem edit   <id>          [-t <title>] [-d <desc>] [--json [field,...]]
```

---

## Authentication Flow

```
po auth login --hostname https://polarion.example.com --project MY_PROJECT
  1. Prompt: "Token: " (hidden input) or read from stdin with --with-token
  2. GET /polarion/rest/v1/projects/MY_PROJECT  → validate token
  3. Store token in system keyring (key: po:polarion.example.com)
  4. Write hostname + default_project to ~/.config/po/config.yml
  5. Print: "Logged in to polarion.example.com as project MY_PROJECT"
```

**Hostname normalization:** the value passed to `--hostname` is stripped of scheme and trailing slashes before storage. `https://polarion.example.com/` → stored key `polarion.example.com`, keyring key `po:polarion.example.com`. All lookups normalize the same way.

Config resolution order (highest to lowest priority):
1. `POLARION_*` env vars
2. System keyring (token) + `~/.config/po/config.yml` (hostname, project, ssl)

---

## Exit Codes

| Code | Meaning |
|---|---|
| 0 | Success |
| 1 | General error |
| 2 | User cancelled |
| 3 | Reserved |
| 4 | Auth error |

---

## Delivery Plan (one PR per milestone)

| PR | Scope |
|---|---|
| 1 | Project scaffold, `go.mod`, Makefile, `iostreams`, `cmdutil/Factory`, `internal/config`, `po auth` (login/logout/status/token). Auth login uses `net/http` directly for token validation — the typed client is not yet available. |
| 2 | `pkg/polarion` typed client, `po api` passthrough + JSON field selector utility |
| 3 | `po case` (list/view/create/edit) |
| 4 | `po testrun` (list/view/create/result) |
| 5 | `po workitem` (list/view/create/edit) |
