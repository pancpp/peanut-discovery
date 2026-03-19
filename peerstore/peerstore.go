package peerstore

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

type PeerInfo struct {
	MultiAddrs []ma.Multiaddr
	lastSeen   time.Time
}

type PeerStore struct {
	mtx   sync.RWMutex
	peers map[peer.ID]PeerInfo
	ttl   time.Duration
}

func NewPeerStore(ctx context.Context, ttl time.Duration) *PeerStore {
	ps := &PeerStore{
		peers: make(map[peer.ID]PeerInfo),
		ttl:   ttl,
	}

	go ps.cleanupLoop(ctx)

	return ps
}

func (ps *PeerStore) Update(id peer.ID, multiAddrs []ma.Multiaddr) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	ps.peers[id] = PeerInfo{
		MultiAddrs: multiAddrs,
		lastSeen:   time.Now(),
	}
}

func (ps *PeerStore) Get(id peer.ID) (PeerInfo, bool) {
	ps.mtx.RLock()
	defer ps.mtx.RUnlock()

	info, ok := ps.peers[id]
	if !ok {
		return PeerInfo{}, false
	}

	if time.Since(info.lastSeen) > ps.ttl {
		return PeerInfo{}, false
	}

	return info, true
}

func (ps *PeerStore) Remove(id peer.ID) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	delete(ps.peers, id)
}

func (ps *PeerStore) GetAll() map[peer.ID]PeerInfo {
	ps.mtx.RLock()
	defer ps.mtx.RUnlock()

	now := time.Now()
	result := make(map[peer.ID]PeerInfo, len(ps.peers))
	for id, info := range ps.peers {
		if now.Sub(info.lastSeen) <= ps.ttl {
			result[id] = info
		}
	}

	return result
}

func (ps *PeerStore) GetLastSeen(id peer.ID) (time.Time, bool) {
	ps.mtx.RLock()
	defer ps.mtx.RUnlock()

	info, ok := ps.peers[id]
	if !ok {
		return time.Time{}, false
	}

	return info.lastSeen, true
}

func (ps *PeerStore) cleanupLoop(ctx context.Context) {
	ticker := time.NewTicker(ps.ttl / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			ps.removeExpired()
		}
	}
}

func (ps *PeerStore) removeExpired() {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()

	for id, info := range ps.peers {
		if time.Since(info.lastSeen) > ps.ttl {
			log.Printf("peer expired: %s", id)
			delete(ps.peers, id)
		}
	}
}
