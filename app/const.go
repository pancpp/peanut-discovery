package app

import "time"

const (
	P2P_DIAL_TIMEOUT   = 10 * time.Second
	P2P_READ_TIMEOUT   = 10 * time.Second
	P2P_WRITE_TIMEOUT  = 10 * time.Second
	P2P_MAX_READ_BYTES = 8192
)

const (
	PROTOCOL_HEARTBEAT = "/peanut/heartbeat/1.0"
	PROTOCOL_DISCOVERY = "/peanut/discovery/1.0"
)
