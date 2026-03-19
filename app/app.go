package app

import (
	"context"
	"log"
	"time"

	"github.com/pancpp/peanut-discovery/conf"
	"github.com/pancpp/peanut-discovery/peerstore"
)

var (
	gPeerStore *peerstore.PeerStore
)

func Init(ctx context.Context) error {
	// create peer store
	pstoreTTL := time.Duration(conf.GetInt("p2p.peer_ttl")) * time.Second
	gPeerStore = peerstore.NewPeerStore(ctx, pstoreTTL)

	// create p2p host
	p2pHost, err := newHost(ctx)
	if err != nil {
		return err
	}
	log.Println("[app] PeerID:", p2pHost.ID())
	log.Println("[app] listen addr:", p2pHost.Addrs())

	// init heartbeat service
	p2pHost.SetStreamHandler(PROTOCOL_HEARTBEAT, handleHeartbeat)

	// init discovery service
	p2pHost.SetStreamHandler(PROTOCOL_DISCOVERY, handleDiscovery)

	return nil
}
