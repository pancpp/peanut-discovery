package app

import (
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	ma "github.com/multiformats/go-multiaddr"
)

type AnnounceMessage struct {
	MultiAddrs []string `json:"multi_addrs"`
}

func handleAnnouncement(stream network.Stream) {
	defer stream.Close()

	remotePeer := stream.Conn().RemotePeer()

	if err := stream.SetReadDeadline(time.Now().Add(P2P_READ_TIMEOUT)); err != nil {
		log.Printf("[announce] set read deadline err: %v, peer: %s", err, remotePeer)
		return
	}

	// read JSON message from the peer
	data, err := io.ReadAll(io.LimitReader(stream, P2P_MAX_READ_BYTES))
	if err != nil {
		log.Printf("[announce] read err: %v, peer: %s", err, remotePeer)
		return
	}

	var msg AnnounceMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("[announce] json decode err: %v, peer: %s", err, remotePeer)
		return
	}

	// parse multiaddrs
	var multiAddrs []ma.Multiaddr
	for _, s := range msg.MultiAddrs {
		addr, err := ma.NewMultiaddr(s)
		if err != nil {
			log.Printf("[announce] invalid multiaddr from peer %s: %q", remotePeer, s)
			continue
		}
		multiAddrs = append(multiAddrs, addr)
	}
	if len(multiAddrs) == 0 {
		log.Println("[announce] invalid multi-addresses:", msg.MultiAddrs)
		return
	}

	gPeerStore.Update(remotePeer, multiAddrs)

	log.Printf("[announce] updated peer %s, multiAddrs: %v", remotePeer, multiAddrs)
}
