<p align="center">
<img src="./.github/logo.png" width="150px" align="center"/>
<div align="center">DepsHub MCP Server for Effortless Dependency Updates</div>
</br>
</p>

- 🔧 Find and fix breaking changes in seconds
- ✨ Works with VS Code, Cursor, Windsurf, Zed and more
- 🧠 Identify the best upgrade path for your dependencies
- 😌 Provide your editor with all the necesarry context from the release notes and changelogs
- 🌍 Supports 2M+ packages out of the box
- 🏎️ Fast respones
  
**Discord: [https://discord.gg/NuEXZwNDtN](https://discord.gg/NuEXZwNDtN)**

## Installation
You can use DepsHub MCP in two modes:
- Connect to a remotely running instance (recommended)
- Run locally, using pre-build docker container `ghcr.io/depshubhq/mcp`

### VS Code
```json
{
  "mcpServers": {
    "depshub": {
      "url": "https://mcp.depshub.com/mcp"
    }
  }
}
```

### Cursor
```json
{
  "mcpServers": {
    "depshub": {
      "url": "https://mcp.depshub.com/mcp"
    }
  }
}
```

### Windsurf
Windsurt doesn't support streamable HTTP, so you have to Docker locally.

```json
{
  "mcpServers": {
    "depshub": {
      "command": "docker",
      "args": ["run", "-i", "--rm", "--init", "ghcr.io/depshubhq/mcp"]
    }
  }
}
```

### Zed
```json
{
  "context_servers": {
    "depshub": {
      "command": {
        "path": "docker",
        "args": ["run", "-i", "--rm", "--init", "ghcr.io/depshubhq/mcp"]
      },
      "settings": {}
    }
  }
}
```

### Other editors
Any editor that supports MCP protocol should be able to work with DepsHub. 
You can either use our [official](ghcr.io/depshubhq/mcp) Docker container or just point to the remote MCP URL: `https://mcp.depshub.com/mcp`.

## Supported ecosystems
- JavaScript/TypeScript - npm, yarn, pnpm
- Go
- Ruby gems
- Rust - cargo
