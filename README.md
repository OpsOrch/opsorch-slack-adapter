# OpsOrch Slack Adapter

This is an **OpsOrch** adapter for **Slack**. It implements the `messaging.Provider` interface to send rich messages to Slack channels using [Block Kit](https://api.slack.com/block-kit).

## Features

- **Rich Messaging**: Supports Headers, Sections, Fields, and Dividers via OpsOrch's generic Block model.
- **Markdown Support**: Automatically converts standard Markdown links (`[text](url)`) to Slack's format (`<url|text>`).
- **Plugin Architecture**: Runs as a standalone binary communicating via JSON-RPC.

## Quick Start

1.  **Create a Slack App**:
    - Go to [api.slack.com/apps](https://api.slack.com/apps).
    - Create a new app and select your workspace.
2.  **Add Scopes**:
    - Go to **OAuth & Permissions**.
    - Add the following **Bot Token Scopes**:
        - `chat:write`
3.  **Install App**:
    - Click **Install to Workspace**.
    - Copy the **Bot User OAuth Token** (`xoxb-...`).
4.  **Add Bot to Channel**:
    - In Slack, go to the channel you want to use.
    - Type `/invite @YourApp`.

## Configuration

The adapter requires a configuration map with the following fields:

```json
{
  "token": "xoxb-your-slack-token"
}
```

## Repository Layout

- `messaging/slack_provider.go`: Implements the `messaging.Provider` interface using the Slack API.
- `cmd/messagingplugin/main.go`: JSON-RPC plugin entrypoint.
- `integ/messaging.go`: Integration tests.
- `Makefile`: Build and test targets.

## Building

```bash
make test           # Run unit tests
make build          # Build all packages
make plugin         # Build ./bin/messagingplugin
make integ-message  # Run integration tests (requires SLACK_TOKEN and SLACK_CHANNEL env vars)
```

## CI/CD

The repository includes GitHub Actions workflows:

- **CI** (`ci.yml`): Runs tests and linting on every push/PR to main
- **Release** (`release.yml`): Manual workflow that:
  - Runs tests and linting
  - Creates version tags (patch/minor/major)
  - Builds multi-arch binaries (linux-amd64, linux-arm64, darwin-amd64, darwin-arm64)
  - Publishes binaries as GitHub release assets

### Pre-Built Binaries

Download pre-built plugin binaries from [GitHub Releases](https://github.com/opsorch/opsorch-slack-adapter/releases):

```dockerfile
# Use in custom Docker images
FROM ghcr.io/opsorch/opsorch-core:latest
WORKDIR /opt/opsorch

ADD https://github.com/opsorch/opsorch-slack-adapter/releases/download/v0.1.0/messagingplugin-linux-amd64 ./plugins/messagingplugin
RUN chmod +x ./plugins/messagingplugin

ENV OPSORCH_MESSAGING_PLUGIN=/opt/opsorch/plugins/messagingplugin
```

## Integration Tests

To run the integration tests, you need to provide your token and a channel ID that your app is invited:

```bash
export SLACK_TOKEN="xoxb-..."
export SLACK_CHANNEL="C12345..."
make integ
```
