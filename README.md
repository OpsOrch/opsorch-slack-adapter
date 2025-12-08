# OpsOrch Slack Adapter

[![Version](https://img.shields.io/github/v/release/opsorch/opsorch-slack-adapter)](https://github.com/opsorch/opsorch-slack-adapter/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/opsorch/opsorch-slack-adapter)](https://github.com/opsorch/opsorch-slack-adapter/blob/main/go.mod)
[![License](https://img.shields.io/github/license/opsorch/opsorch-slack-adapter)](https://github.com/opsorch/opsorch-slack-adapter/blob/main/LICENSE)
[![CI](https://github.com/opsorch/opsorch-slack-adapter/workflows/CI/badge.svg)](https://github.com/opsorch/opsorch-slack-adapter/actions)

This adapter integrates OpsOrch with Slack, enabling rich message delivery to Slack channels using [Block Kit](https://api.slack.com/block-kit).

## Capabilities

This adapter provides one capability:

1. **Messaging Provider**: Send rich messages to Slack channels

## Features

- **Rich Messaging**: Supports Headers, Sections (including field sets), and Dividers via OpsOrch's generic Block model
- **Markdown Support**: Automatically converts standard Markdown links (`[text](url)`) to Slack's format (`<url|text>`)
- **Block Kit Integration**: Maps OpsOrch blocks to the subset of Slack Block Kit elements used by the adapter
- **Channel Targeting**: Send messages to specific channels by ID (public, private, or IM)
- **Plugin Architecture**: Runs as a standalone binary communicating via JSON-RPC

### Version Compatibility

- **Adapter Version**: 0.1.0
- **Requires OpsOrch Core**: >=0.1.0
- **Slack API**: Web API
- **Go Version**: 1.21+

## Configuration

The messaging adapter requires the following configuration:

| Field | Type | Required | Description | Default |
|-------|------|----------|-------------|---------|
| `token` | string | Yes | Slack Bot User OAuth Token (starts with `xoxb-`) | - |

### Authentication Setup

#### 1. Create a Slack App

1. Go to [api.slack.com/apps](https://api.slack.com/apps)
2. Click **Create New App**
3. Choose **From scratch**
4. Give your app a name (e.g., `OpsOrch Bot`)
5. Select your workspace
6. Click **Create App**

#### 2. Add Bot Token Scopes

1. In your app settings, go to **OAuth & Permissions**
2. Scroll down to **Scopes** → **Bot Token Scopes**
3. Click **Add an OAuth Scope**
4. Add the following scope:
   - `chat:write` - Send messages as the bot

#### 3. Install App to Workspace

1. In your app settings, go to **OAuth & Permissions**
2. Click **Install to Workspace**
3. Review the permissions and click **Allow**
4. Copy the **Bot User OAuth Token** (starts with `xoxb-`)

⚠️ **Important**: Store the token securely in your secret manager (Vault, AWS SSM, Kubernetes secrets, etc.) and inject it into OpsOrch config at runtime.

#### 4. Add Bot to Channel

Before the bot can send messages to a channel, you must invite it:

1. In Slack, go to the channel where you want to send messages
2. Type `/invite @YourAppName` (replace with your app name)
3. The bot will join the channel and can now send messages

#### 5. Find Channel ID

Messages must be addressed by channel ID. To find the channel ID:

1. Right-click on the channel name in Slack
2. Select **Copy link**
3. The channel ID is in the URL: `https://yourworkspace.slack.com/archives/C1234567890`
4. The ID is `C1234567890`

### Example Configuration

**JSON format:**
```json
{
  "token": "xoxb-your-bot-token-here"
}
```

**Environment variables:**
```bash
export OPSORCH_MESSAGING_PLUGIN=/path/to/bin/messagingplugin
export OPSORCH_MESSAGING_CONFIG='{"token":"xoxb-..."}'
```

## Field Mapping

### Message Block Mapping

OpsOrch's generic Block model maps to Slack's Block Kit:

| OpsOrch Block Type | Slack Block Type | Transformation | Notes |
|-------------------|------------------|----------------|-------|
| `header` | `header` | Direct mapping | Text limited to 150 characters |
| `section` | `section` | Direct mapping | Supports text and fields |
| `divider` | `divider` | Direct mapping | Visual separator |
| `section` (with fields map) | `section` with fields | Converted to section block | Key-value pairs displayed in columns |

### Markdown Conversion

The adapter automatically converts standard Markdown links to Slack's format:

| Standard Markdown | Slack Format | Example |
|------------------|--------------|---------|
| `[text](url)` | `<url\|text>` | `[Dashboard](https://example.com)` → `<https://example.com\|Dashboard>` |

### Message Result Fields

The provider returns the standard OpsOrch `schema.MessageResult` fields:

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Slack message timestamp returned by the API |
| `channel` | string | Channel ID where message was sent |
| `sentAt` | string (RFC3339) | Timestamp recorded by the adapter when the send completed |
| `metadata` | object | Optional additional fields (unused by the current implementation) |

## Usage

### In-Process Mode

Import the adapter for side effects to register it with OpsOrch Core:

```go
import _ "github.com/opsorch/opsorch-slack-adapter/messaging"
```

Configure via environment variables:

```bash
export OPSORCH_MESSAGING_PROVIDER=slack
export OPSORCH_MESSAGING_CONFIG='{"token":"xoxb-..."}'
```

### Plugin Mode

Build the plugin binary:

```bash
make plugin
```

Configure OpsOrch Core to use the plugin:

```bash
export OPSORCH_MESSAGING_PLUGIN=/path/to/bin/messagingplugin
export OPSORCH_MESSAGING_CONFIG='{"token":"xoxb-..."}'
```

### Docker Deployment

Download pre-built plugin binaries from [GitHub Releases](https://github.com/opsorch/opsorch-slack-adapter/releases):

```dockerfile
FROM ghcr.io/opsorch/opsorch-core:latest
WORKDIR /opt/opsorch

# Download plugin binary
ADD https://github.com/opsorch/opsorch-slack-adapter/releases/download/v0.1.0/messagingplugin-linux-amd64 ./plugins/messagingplugin
RUN chmod +x ./plugins/messagingplugin

# Configure plugin
ENV OPSORCH_MESSAGING_PLUGIN=/opt/opsorch/plugins/messagingplugin
```

## Development

### Prerequisites

- Go 1.21 or later
- Slack workspace with admin access
- Slack Bot User OAuth Token

### Building

```bash
# Download dependencies
go mod download

# Run unit tests
make test

# Build all packages
make build

# Build plugin binary
make plugin

# Run integration tests (requires SLACK_TOKEN and SLACK_CHANNEL)
make integ
```

### Testing

**Unit Tests:**
```bash
make test
```

**Integration Tests:**

Integration tests send real messages to a Slack channel.

**Prerequisites:**
- A Slack workspace with a bot installed
- A Slack Bot User OAuth Token (`xoxb-...`)
- A channel ID where the bot has been invited

**Setup:**
```bash
# Set required environment variables
export SLACK_TOKEN="xoxb-your-bot-token"
export SLACK_CHANNEL="C1234567890"  # Your channel ID

# Run integration tests
make integ

# Or run specific messaging tests
make integ-message
```

**What the tests do:**
- Send a test message with headers, sections, and fields to the specified channel
- Verify the message was sent successfully
- Capture message metadata (timestamp, channel ID)

**Expected behavior:**
- A test message will appear in your Slack channel
- The message will contain formatted blocks (header, sections, divider)
- Tests verify the Slack API returns success

**Note:** Integration tests send real messages to your Slack channel. Use a test channel to avoid cluttering production channels.

### Project Structure

```
opsorch-slack-adapter/
├── messaging/                  # Messaging provider implementation
│   ├── slack_provider.go      # Core provider logic
│   └── slack_provider_test.go # Unit tests
├── cmd/
│   └── messagingplugin/       # Plugin entrypoint
│       └── main.go
├── integ/                      # Integration tests
│   └── messaging.go
├── Makefile
└── README.md
```

**Key Components:**

- **messaging/slack_provider.go**: Implements messaging.Provider interface, handles Slack Block Kit conversion and API calls
- **cmd/messagingplugin**: JSON-RPC plugin wrapper for messaging provider
- **integ/messaging.go**: End-to-end integration tests against live Slack workspace

## CI/CD & Pre-Built Binaries

The repository includes GitHub Actions workflows:

- **CI** (`ci.yml`): Runs tests and linting on every push/PR to main
- **Release** (`release.yml`): Manual workflow that:
  - Runs tests and linting
  - Creates version tags (patch/minor/major)
  - Builds multi-arch binaries (linux-amd64, linux-arm64, darwin-amd64, darwin-arm64)
  - Publishes binaries as GitHub release assets

### Downloading Pre-Built Binaries

Pre-built plugin binaries are available from [GitHub Releases](https://github.com/opsorch/opsorch-slack-adapter/releases).

**Supported platforms:**
- Linux (amd64, arm64)
- macOS (amd64, arm64)

## Plugin RPC Contract

OpsOrch Core communicates with the plugin over stdin/stdout using JSON-RPC.

### Message Format

**Request:**
```json
{
  "method": "messaging.send",
  "config": { /* decrypted configuration */ },
  "payload": { /* method-specific request body */ }
}
```

**Response:**
```json
{
  "result": { /* method-specific result */ },
  "error": "optional error message"
}
```

### Configuration Injection

The `config` field contains the decrypted configuration map from `OPSORCH_MESSAGING_CONFIG`. The plugin receives this on every request, so it never stores secrets on disk.

### Supported Methods

#### messaging.send

Send a message to a Slack channel.

**Request:**
```json
{
  "method": "messaging.send",
  "config": {"token": "xoxb-..."},
  "payload": {
    "channel": "C1234567890",
    "blocks": [
      {
        "type": "header",
        "text": "Incident Alert"
      },
      {
        "type": "section",
        "text": "Database connection timeout detected"
      },
      {
        "type": "section",
        "fields": {
          "Severity": "Critical",
          "Service": "api-backend",
          "Environment": "production"
        }
      },
      {
        "type": "divider"
      },
      {
        "type": "section",
        "text": "[View Dashboard](https://dashboard.example.com)"
      }
    ]
  }
}
```

**Response:**
```json
{
  "result": {
    "id": "1234567890.123456",
    "channel": "C1234567890",
    "sentAt": "2024-01-01T12:00:00Z"
  }
}
```

## Message Examples

### Simple Text Message

```json
{
  "channel": "C1234567890",
  "blocks": [
    {
      "type": "section",
      "text": "Deployment completed successfully"
    }
  ]
}
```

### Rich Incident Alert

```json
{
  "channel": "C1234567890",
  "blocks": [
    {
      "type": "header",
      "text": "🚨 Critical Incident"
    },
    {
      "type": "section",
      "text": "High error rate detected in payment service"
    },
    {
      "type": "section",
      "fields": {
        "Severity": "Critical",
        "Service": "payment-api",
        "Error Rate": "15%",
        "Started": "2024-01-01 10:00 UTC"
      }
    },
    {
      "type": "divider"
    },
    {
      "type": "section",
      "text": "[View Incident](https://incidents.example.com/123) | [View Logs](https://logs.example.com)"
    }
  ]
}
```

## Security Considerations

1. **Never log the bot token**: Avoid logging the config or token in application logs
2. **Rotate tokens regularly**: Rotate the bot token at the cadence required by your organization's security policy
3. **Use environment variables**: Store the `OPSORCH_MESSAGING_CONFIG` in a secure environment variable or secrets management system
4. **Restrict file permissions**: If storing config in files, ensure proper file permissions (e.g., 0600)
5. **Limit bot permissions**: Only grant the `chat:write` scope; avoid unnecessary permissions
6. **Use private channels**: For sensitive alerts, use private channels and carefully manage membership

## License

Apache 2.0

See LICENSE file for details.
