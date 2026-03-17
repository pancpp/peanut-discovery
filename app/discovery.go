package app

import (
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	DISCOVERY_TOPIC       = "/peanut/discovery/1.0"
	discoveryReadDeadline = 10 * time.Second
	discoveryMaxReadBytes = 4096
)

type DiscoveryRequestMsg struct {
	PeerID string `json:"peer_id"`
}

type DiscoveryResponseMsg struct {
	Addrs []string `json:"addrs"`
}

func handleDiscovery(stream network.Stream) {
	defer stream.Close()

	remotePeer := stream.Conn().RemotePeer()

	if err := stream.SetReadDeadline(time.Now().Add(discoveryReadDeadline)); err != nil {
		log.Printf("discovery: set read deadline err: %v, peer: %s", err, remotePeer)
		return
	}

	// read request
	data, err := io.ReadAll(io.LimitReader(stream, discoveryMaxReadBytes))
	if err != nil {
		log.Printf("discovery: read err: %v, peer: %s", err, remotePeer)
		return
	}

	var req DiscoveryRequestMsg
	if err := json.Unmarshal(data, &req); err != nil {
		log.Printf("discovery: json decode err: %v, peer: %s", err, remotePeer)
		return
	}

	// look up the requested peer
	targetID, err := peer.Decode(req.PeerID)
	if err != nil {
		log.Printf("discovery: invalid peer_id %q from peer %s", req.PeerID, remotePeer)
		return
	}

	resp := DiscoveryResponseMsg{}
	if info, ok := gPeerStore.Get(targetID); ok {
		resp.Addrs = make([]string, len(info.Addrs))
		for i, addr := range info.Addrs {
			resp.Addrs[i] = addr.String()
		}
	}

	// write response
	respData, err := json.Marshal(resp)
	if err != nil {
		log.Printf("discovery: json encode err: %v, peer: %s", err, remotePeer)
		return
	}

	if _, err := stream.Write(respData); err != nil {
		log.Printf("discovery: write err: %v, peer: %s", err, remotePeer)
		return
	}

	log.Printf("discovery: peer %s queried for %s, returned %d addrs", remotePeer, targetID, len(resp.Addrs))
}
