# po — Polarion CLI

`po` is a command-line tool for [Polarion ALM](https://polarion.plm.automation.siemens.com/), modeled after the [GitHub CLI (`gh`)](https://cli.github.com/). It lets you work with work items, test cases, test runs, and more directly from your terminal.

## Installation

```bash
go install github.com/wangke19/po@latest
```

Or build from source:

```bash
git clone https://github.com/wangke19/po
cd po
go build -o /usr/local/bin/po ./cmd/po
```

## Authentication

### Option 1: Environment variables (recommended for CI/scripts)

```bash
export POLARION_URL=https://polarion.example.com
export POLARION_PROJECT=MYPROJECT
export POLARION_TOKEN=<your-bearer-token>
export POLARION_VERIFY_SSL=false  # if using self-signed cert
```

### Option 2: Interactive login (stores token in system keyring)

```bash
po auth login
```

Check status:

```bash
po auth status
```

## Usage

### Work Items

```bash
po workitem list --limit 20
po workitem list --type testcase --status approved
po workitem view OCP-84783
po workitem create --type testcase --title "My test" --status draft
po workitem edit OCP-84783 --status approved
po workitem delete OCP-84783 --confirm
```

### Test Cases

```bash
po case list --limit 20
po case view OCP-84783
po testcase steps OCP-84783
po testcase step-add OCP-84783 --action "Do X" --expected-result "See Y"
po testcase step-edit OCP-84783 1 --action "Updated action"
po testcase step-remove OCP-84783 1
```

### Test Runs

```bash
po testrun list --limit 20
po testrun view TR-123
po testrun create --title "Sprint 42 Run" --template MyTemplate
po testrun records TR-123
po testrun records TR-123 --result failed
po testrun records TR-123 --not-run
po testrun result TR-123 OCP-84783 --result passed --comment "All good"
po run start TR-123
po run pause TR-123
po run finish TR-123
```

### Comments

```bash
po comment list OCP-84783
po comment add OCP-84783 --body "Verified on 4.14"
echo "Verified" | po comment add OCP-84783 --body -
```

### Links

```bash
po link list OCP-84783
po link add OCP-84783 OCP-99999 --role relates_to
po link remove OCP-84783 OCP-99999 --role relates_to
```

### Attachments

```bash
po attachment list OCP-84783
po attachment upload OCP-84783 ./report.txt
po attachment download OCP-84783 <attachment-id> -o report.txt
```

### Search

```bash
po search "type:testcase status:approved" --limit 50
po search "type:testcase" --author jdoe
```

### Other

```bash
po project list
po project view MYPROJECT
po clone workitem OCP-84783
po clone testrun TR-123
po open OCP-84783        # open in browser
po whoami
po api /projects/{project}/workitems   # raw API call
```

## JSON output

Most commands support `--json` for machine-readable output:

```bash
po workitem list --json id,title,status
po testrun records TR-123 --json caseId,result
po workitem view OCP-84783 --json        # all fields
```

## Environment variables

| Variable | Description |
|----------|-------------|
| `POLARION_URL` | Polarion server URL (e.g. `https://polarion.example.com`) |
| `POLARION_PROJECT` | Default project ID (e.g. `MYPROJECT`) |
| `POLARION_TOKEN` | Bearer token for authentication |
| `POLARION_VERIFY_SSL` | Set to `false` to skip TLS certificate verification |

Environment variables take precedence over the config file.

## Shell completion

```bash
po completion bash > /etc/bash_completion.d/po
po completion zsh > ~/.zsh/completions/_po
```
