package app

import (
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

type DiscoveryRequestMsg struct {
	PeerIDs []string `json:"peer_ids"`
}

type DiscoveryPeerMsg struct {
	Multiaddrs []string `json:"multi_addrs"`
}

type DiscoveryResponseMsg struct {
	PeerInfo map[string]DiscoveryPeerMsg `json:"peer_info"`
}

func handleDiscovery(stream network.Stream) {
	defer stream.Close()

	remotePeer := stream.Conn().RemotePeer()

	if err := stream.SetReadDeadline(time.Now().Add(P2P_READ_TIMEOUT)); err != nil {
		log.Printf("[discovery] set read deadline err: %v, peer: %s", err, remotePeer)
		return
	}

	// read request
	data, err := io.ReadAll(io.LimitReader(stream, P2P_MAX_READ_BYTES))
	if err != nil {
		log.Printf("[discovery] read err: %v, peer: %s", err, remotePeer)
		return
	}

	var req DiscoveryRequestMsg
	if err := json.Unmarshal(data, &req); err != nil {
		log.Printf("[discovery] json decode err: %v, peer: %s", err, remotePeer)
		return
	}

	// look up each requested peer
	resp := DiscoveryResponseMsg{
		PeerInfo: make(map[string]DiscoveryPeerMsg),
	}

	for _, pidStr := range req.PeerIDs {
		targetID, err := peer.Decode(pidStr)
		if err != nil {
			log.Printf("[discovery] invalid peer_id %q from peer %s", pidStr, remotePeer)
			continue
		}

		info, ok := gPeerStore.Get(targetID)
		if !ok {
			continue
		}

		addrs := make([]string, len(info.MultiAddrs))
		for i, addr := range info.MultiAddrs {
			addrs[i] = addr.String()
		}

		resp.PeerInfo[pidStr] = DiscoveryPeerMsg{
			Multiaddrs: addrs,
		}
	}

	// write response
	respData, err := json.Marshal(resp)
	if err != nil {
		log.Printf("[discovery] json encode err: %v, peer: %s", err, remotePeer)
		return
	}

	if err := stream.SetWriteDeadline(time.Now().Add(P2P_WRITE_TIMEOUT)); err != nil {
		log.Printf("[discovery] set write deadline err: %v, peer: %s", err, remotePeer)
		return
	}

	if _, err := stream.Write(respData); err != nil {
		log.Printf("[discovery] write err: %v, peer: %s", err, remotePeer)
		return
	}

	log.Printf("[discovery] peer %s query response: %v", remotePeer, string(respData))
}
