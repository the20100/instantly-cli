# instantly-cli

A CLI tool for the [Instantly.ai](https://instantly.ai) API v2.

Outputs JSON when piped (for agent use) and human-readable tables in a terminal.

---

## Install

```bash
git clone https://github.com/the20100/instantly-cli
cd instantly-cli
go build -o instantly .
mv instantly /usr/local/bin/
```

## Authentication

Get your API key from: https://app.instantly.ai/app/settings/integrations

```bash
# Save key to config file
instantly auth set-key <your-api-key>

# Or set env var
export INSTANTLY_API_KEY=<your-api-key>
```

Resolution order: env var → config file.

---

## Global Flags

| Flag | Description |
|------|-------------|
| `--json` | Force JSON output |
| `--pretty` | Force pretty-printed JSON (implies --json) |

Output is auto-detected: JSON when piped, tables in terminal.

---

## Commands

### auth

```bash
instantly auth set-key <api-key>   # Save API key
instantly auth status              # Show auth status
instantly auth logout              # Remove saved key
```

### campaign

```bash
instantly campaign list                          # List all campaigns
instantly campaign list --status active          # Filter by status
instantly campaign list --search "outreach"
instantly campaign get <id>
instantly campaign create "My Campaign" --accounts email@domain.com --daily-limit 50
instantly campaign update <id> --name "New Name" --daily-limit 100
instantly campaign delete <id>
instantly campaign activate <id>
instantly campaign pause <id>
instantly campaign analytics <id>
instantly campaign analytics-overview
instantly campaign duplicate <id>
```

### account

```bash
instantly account list
instantly account get <id>
instantly account create user@domain.com --first-name John --last-name Doe \
  --smtp-host smtp.domain.com --smtp-port 587 --smtp-user user@domain.com --smtp-pass secret \
  --imap-host imap.domain.com --imap-port 993 --imap-user user@domain.com --imap-pass secret
instantly account update <id> --daily-limit 100 --warmup-enabled
instantly account delete <id>
instantly account pause <email>
instantly account resume <email>
instantly account warmup enable <email>
instantly account warmup disable <email>
instantly account warmup analytics --email user@domain.com
```

### lead

```bash
instantly lead list
instantly lead list --campaign-id <id>
instantly lead list --status interested
instantly lead get <id>
instantly lead create john@acme.com --first-name John --last-name Doe --company Acme
instantly lead update <id> --company "New Corp"
instantly lead delete <id>
instantly lead update-interest <id> --status interested
# Statuses: interested, not_interested, meeting_booked, meeting_completed, closed
```

### leadlist

```bash
instantly leadlist list
instantly leadlist get <id>
instantly leadlist create "My Lead List"
instantly leadlist update <id> --name "New Name"
instantly leadlist delete <id>
```

### email

```bash
instantly email list
instantly email list --campaign-id <id>
instantly email list --type reply --is-unread
instantly email get <id>
instantly email reply <id> --body "Thanks for reaching out!"
instantly email forward <id> --to colleague@domain.com
instantly email mark-read <thread-id>
```

### webhook

```bash
instantly webhook list
instantly webhook get <id>
instantly webhook create "My Hook" --url https://myapp.com/hook --events reply_received,email_sent
instantly webhook update <id> --url https://newurl.com
instantly webhook delete <id>
instantly webhook test <id>
instantly webhook resume <id>
instantly webhook event-types     # List all available event types
```

### customtag

```bash
instantly customtag list
instantly customtag get <id>
instantly customtag create "Hot Lead" --color "#FF0000"
instantly customtag update <id> --name "Warm Lead" --color "#FFA500"
instantly customtag delete <id>
instantly customtag toggle --tag-id <id> --resource-id <id> --resource-type campaign
```

### blocklist

```bash
instantly blocklist list
instantly blocklist get <id>
instantly blocklist create --value spam@example.com --type email
instantly blocklist create --value spammydomain.com --type domain
instantly blocklist update <id> --value newvalue@example.com
instantly blocklist delete <id>
```

### apikey

```bash
instantly apikey list
instantly apikey create "My Integration Key"
instantly apikey delete <id>
```

### workspace

```bash
instantly workspace get
instantly workspace update --name "New Workspace Name"
instantly workspace member list
instantly workspace member get <id>
instantly workspace member create --email user@domain.com --role member
instantly workspace member update <id> --role admin
instantly workspace member delete <id>
```

### analytics

```bash
instantly analytics campaign                          # Overview of all campaigns
instantly analytics campaign --campaign-id <id>       # Single campaign
instantly analytics campaign --start-date 2024-01-01 --end-date 2024-12-31
instantly analytics warmup
instantly analytics warmup --email user@domain.com
```

### subsequence

```bash
instantly subsequence list --campaign-id <id>
instantly subsequence get <id>
instantly subsequence create --campaign-id <id> --name "Follow-up"
instantly subsequence update <id> --name "New Name"
instantly subsequence delete <id>
instantly subsequence pause <id>
instantly subsequence resume <id>
instantly subsequence duplicate <id>
```

### verify

```bash
instantly verify create user@domain.com
instantly verify check <id>
```

### job

```bash
instantly job list
instantly job get <id>
```

### update

```bash
instantly update    # Pull latest from GitHub, rebuild, replace binary
```

### info

```bash
instantly info    # Show binary path, config path, auth status
```

---

## Agent Usage

All commands output JSON when piped:

```bash
# Get all active campaign IDs
instantly campaign list --status active | jq '.[].id'

# Get all leads in a campaign
instantly lead list --campaign-id <id> | jq '.[] | {email, status: .lt_interest_status}'

# Get campaign analytics as JSON
instantly campaign analytics <id> --json

# List all webhook event types
instantly webhook event-types --json | jq '.[].name'
```

---

## API Reference

Full API documentation: https://developer.instantly.ai/api/v2/
