package app

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	ma "github.com/multiformats/go-multiaddr"
)

type HeartbeatMessage struct {
	IPAddr     string   `json:"ip_addr"`
	MultiAddrs []string `json:"multi_addrs"`
}

func handleHeartbeat(stream network.Stream) {
	defer stream.Close()

	remotePeer := stream.Conn().RemotePeer()

	if err := stream.SetReadDeadline(time.Now().Add(P2P_READ_TIMEOUT)); err != nil {
		log.Printf("[heartbeat] set read deadline err: %v, peer: %s", err, remotePeer)
		return
	}

	// read JSON message from the peer
	data, err := io.ReadAll(io.LimitReader(stream, P2P_MAX_READ_BYTES))
	if err != nil {
		log.Printf("[heartbeat] read err: %v, peer: %s", err, remotePeer)
		return
	}

	var msg HeartbeatMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("[heartbeat] json decode err: %v, peer: %s", err, remotePeer)
		return
	}

	// parse IP address
	ipAddr := net.ParseIP(msg.IPAddr)
	if ipAddr == nil {
		log.Println("[heartbeat] invalid IP addr:", msg.IPAddr)
		return
	}

	// parse multiaddrs
	var multiAddrs []ma.Multiaddr
	for _, s := range msg.MultiAddrs {
		addr, err := ma.NewMultiaddr(s)
		if err != nil {
			log.Printf("[heartbeat] invalid multiaddr from peer %s: %q", remotePeer, s)
			continue
		}
		multiAddrs = append(multiAddrs, addr)
	}
	if len(multiAddrs) == 0 {
		log.Println("[heartbeat] invalid multi-addresses:", msg.MultiAddrs)
		return
	}

	gPeerStore.Update(remotePeer, ipAddr, multiAddrs)

	log.Printf("[heartbeat] updated peer %s, IP: %v, multiAddrs: %v",
		remotePeer, ipAddr, multiAddrs)
}
