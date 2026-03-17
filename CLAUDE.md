# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

peanut-discovery is a P2P discovery server written in Go (1.25.0), module path `github.com/pancpp/peanut-discovery`. It uses libp2p-style networking (QUIC transport, multiaddrs) and is designed to run as a Linux systemd service.

## Build & Run

```bash
# Build
go build

# Build with version info (linker flags inject into conf package vars)
go build -ldflags "-X github.com/pancpp/peanut-discovery/conf.gVersion=1.0.0 -X github.com/pancpp/peanut-discovery/conf.gBuildTime=$(date -u +%Y%m%d%H%M%S) -X github.com/pancpp/peanut-discovery/conf.gGitHash=$(git rev-parse --short HEAD)"

# Run
./peanut-discovery [--config /path/to/discovery.yaml] [--version]
```

Default config path: `/etc/peanut/discovery.yaml`. The binary auto-creates the config directory and file if missing.

## Architecture

Initialization flows sequentially in `main.go`: **conf → logger → app**, then blocks on OS signals (SIGTERM/SIGINT/SIGHUP) for graceful shutdown.

- **conf/** — Configuration via Viper + pflag. Defaults are set in `conf.init()`. Other packages read config through `conf.GetString()`, `conf.GetBool()`, etc. wrappers.
- **logger/** — Log rotation via lumberjack. Reads settings from Viper (`log.*` keys).
- **app/** — Application logic entry point. Receives a `context.Context` from main for shutdown propagation.

## Key Configuration Keys

| Key | Default | Purpose |
|-----|---------|---------|
| `log.path` | `/var/log/peanut/discovery.log` | Log file location |
| `log.enable_console_log` | `false` | Also log to stdout |
| `p2p.private_key_path` | `/etc/peanut/discovery-private-key.b64` | Node identity key |
| `p2p.pnet_psk_path` | `""` | Optional private network PSK |
| `p2p.listen_multiaddrs` | `["/ip4/0.0.0.0/udp/19880/quic-v1"]` | Listen addresses |
