# peanut-discovery

A peer discovery server for the [peanut](https://github.com/pancpp/peanut) P2P virtual network. It maintains a registry of peer addresses so that nodes can find each other.

Different from Kad DHT service, this discovery server is a simple example implementation for small size private networks.

In practice, a single discovery server is sufficient to host a private network. Nevertheless, multiple discovery servers can also be deployed for better robustness. 

Built with [go-libp2p](https://github.com/libp2p/go-libp2p) using QUIC transport.

## How It Works

Peanut nodes interact with the discovery server through two protocols:

- **Heartbeat** (`/peanut/heartbeat/1.0`) — Peers periodically send their IP address and multiaddrs to register themselves. The server stores this information in a TTL-based peer store.
- **Discovery** (`/peanut/discovery/1.0`) — Peers send a list of peer IDs they want to look up. The server responds with known IP addresses and multiaddrs for each requested peer.

Peer entries expire automatically after a configurable TTL (default 300s). A background cleanup loop removes stale entries.

## Build

```bash
go build
```

With version info:
```bash
./jenkins.sh
```
or
```bash
go build -ldflags "\
  -X github.com/pancpp/peanut-discovery/conf.gVersion=1.0.0 \
  -X github.com/pancpp/peanut-discovery/conf.gBuildTime=$(date -u +%Y%m%d%H%M%S) \
  -X github.com/pancpp/peanut-discovery/conf.gGitHash=$(git rev-parse --short HEAD)"
```

## Usage

```bash
./peanut-discovery [--config /path/to/discovery.yaml] [--version]
```

The default config path is `/srv/peanut-discovery/discovery.yaml`. The binary auto-creates the config directory and file if they don't exist.

## Configuration

All settings can be specified in the YAML config file.

| Key | Default | Description |
|-----|---------|-------------|
| `log.path` | `/srv/peanut-discovery/discovery.log` | Log file location |
| `log.enable_console_log` | `false` | Also log to stdout |
| `log.max_size` | `500` | Max log file size in MB |
| `log.max_backups` | `3` | Number of old log files to keep |
| `log.local_time` | `true` | Use local time in log filenames |
| `log.compress` | `true` | Compress rotated log files |
| `p2p.private_key_path` | `/srv/peanut-discovery/private-key.b64` | Base64-encoded libp2p private key |
| `p2p.pnet_psk_path` | `""` | Private network pre-shared key (optional) |
| `p2p.listen_multiaddrs` | `["/ip4/0.0.0.0/udp/19880/quic-v1"]` | Listen addresses |
| `p2p.peer_ttl` | `300` | Peer entry TTL in seconds |

## Project Structure

```
├── main.go           # Entry point, signal handling
├── conf/             # Configuration (Viper + pflag)
├── logger/           # Log rotation (lumberjack)
├── app/
│   ├── app.go        # Initialization, protocol registration
│   ├── host.go       # libp2p host setup (QUIC, identity, pnet)
│   ├── heartbeat.go  # Heartbeat protocol handler
│   ├── discovery.go  # Discovery protocol handler
│   └── const.go      # Protocol IDs and timeouts
└── peerstore/        # Thread-safe peer store with TTL expiration
```

## Protocol Messages

### Heartbeat (peer → server)

```json
{
  "multi_addrs": ["/ip4/203.0.113.1/udp/19880/quic-v1"]
}
```

### Discovery Request (peer → server)

```json
{
  "peer_ids": ["12D3KooW...", "12D3KooW..."]
}
```

### Discovery Response (server → peer)

```json
{
  "peer_info": {
    "12D3KooW...": {
      "multi_addrs": ["/ip4/203.0.113.1/udp/19880/quic-v1"]
    }
  }
}
```

## License

[MIT](LICENSE)
