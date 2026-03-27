package app

import "time"

const (
	P2P_DIAL_TIMEOUT   = 10 * time.Second
	P2P_READ_TIMEOUT   = 10 * time.Second
	P2P_WRITE_TIMEOUT  = 10 * time.Second
	P2P_MAX_READ_BYTES = 8192
)

const (
	PROTOCOL_ANNOUNCE = "/peanut/announce/1.0"
	PROTOCOL_DISCOVER = "/peanut/discover/1.0"
)
