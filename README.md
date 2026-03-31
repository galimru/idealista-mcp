<div align="center">

**MCP server for the Idealista property search platform**

*Ask your AI assistant to search for properties, filter by price or size, and find location IDs — right from the chat*

<p>
  <img src="https://img.shields.io/github/v/release/galimru/idealista-mcp" alt="Latest release">
  <img src="https://img.shields.io/badge/MCP-compatible-blueviolet" alt="MCP">
  <img src="https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go&logoColor=white" alt="Go">
  <img src="https://img.shields.io/badge/license-MIT-blue" alt="License">
</p>

</div>

---

Connect your AI assistant to Idealista. Search for properties for sale or rent, filter by price, size, and rooms, and resolve location names to IDs — all through natural conversation.

## Quick Start

**1. Install**

Download a binary from the [releases page](https://github.com/galimru/idealista-mcp/releases).

Or build from source:

```bash
git clone https://github.com/galimru/idealista-mcp.git
cd idealista-mcp
make install
```

**2. Connect to Claude Desktop**

Add to `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "idealista": {
      "command": "/path/to/idealista-mcp",
      "env": {
        "IDEALISTA_CLIENT_KEY": "your-client-key",
        "IDEALISTA_CLIENT_SECRET": "your-client-secret",
        "IDEALISTA_SIGNING_SECRET": "your-signing-secret"
      }
    }
  }
}
```

## Tools

| Tool | What it does |
|------|--------------|
| `search_locations` | Search for location IDs by name prefix (needed as input for `search_ads`) |
| `search_ads` | Search property listings with filters for price, size, rooms, and more |

### Example workflow

1. Use `search_locations` with `prefix=Valencia` to get the location ID for Valencia
2. Use `search_ads` with the returned `location_id` to browse listings

## Configuration

**Environment variables**

| Variable | Required | Description |
|----------|----------|-------------|
| `IDEALISTA_CLIENT_KEY` | Yes | OAuth client key (client ID) |
| `IDEALISTA_CLIENT_SECRET` | Yes | OAuth client secret |
| `IDEALISTA_SIGNING_SECRET` | Yes | Raw HMAC-SHA256 signing secret for request authentication |
| `IDEALISTA_DEBUG` | No | Set to any non-empty value to log HTTP traffic to stderr |

A persistent session file is stored at `~/.config/idealista-mcp/session.json`. It holds a stable device identifier (generated once on first run) and the cached OAuth token so restarts reuse a valid token without re-fetching.

## Notes

- This project is not affiliated with Idealista or its parent company.

## Contributing

Bug fixes and clear improvements are welcome. Open an issue first for anything non-trivial.

## License

MIT
