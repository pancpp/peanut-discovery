package app

import (
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	ma "github.com/multiformats/go-multiaddr"
)

const (
	HEARTBEAT_TOPIC       = "/peanut/heartbeat/1.0"
	heartbeatReadDeadline = 10 * time.Second
	heartbeatMaxReadBytes = 8192
)

type HeartbeatMessage struct {
	Addrs []string `json:"addrs"`
}

func handleHeartbeat(stream network.Stream) {
	defer stream.Close()

	remotePeer := stream.Conn().RemotePeer()
	remoteObservedAddr := stream.Conn().RemoteMultiaddr()

	if err := stream.SetReadDeadline(time.Now().Add(heartbeatReadDeadline)); err != nil {
		log.Printf("heartbeat: set read deadline err: %v, peer: %s", err, remotePeer)
		return
	}

	// read JSON message from the peer
	data, err := io.ReadAll(io.LimitReader(stream, heartbeatMaxReadBytes))
	if err != nil {
		log.Printf("heartbeat: read err: %v, peer: %s", err, remotePeer)
		return
	}

	var msg HeartbeatMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("heartbeat: json decode err: %v, peer: %s", err, remotePeer)
		return
	}

	// parse reported multiaddrs
	var reportedAddrs []ma.Multiaddr
	for _, s := range msg.Addrs {
		addr, err := ma.NewMultiaddr(s)
		if err != nil {
			log.Printf("heartbeat: invalid multiaddr from peer %s: %q", remotePeer, s)
			continue
		}
		reportedAddrs = append(reportedAddrs, addr)
	}

	// combine reported addresses with the observed relay address
	allAddrs := make([]ma.Multiaddr, 0, len(reportedAddrs)+1)
	allAddrs = append(allAddrs, remoteObservedAddr)
	allAddrs = append(allAddrs, reportedAddrs...)

	gPeerStore.Update(remotePeer, allAddrs)
	log.Printf("heartbeat: updated peer %s, addrs: %v", remotePeer, allAddrs)
}
