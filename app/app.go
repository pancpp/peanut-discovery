package app

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/pnet"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/pancpp/peanut-disc/conf"
	"github.com/pancpp/peanut-disc/peerstore"
)

var (
	gHost      host.Host
	gPeerStore *peerstore.PeerStore
)

func Init(ctx context.Context) error {
	// peer store
	if err := initPeerStore(ctx); err != nil {
		return err
	}

	// init discovery service
	if err := initDiscovery(ctx); err != nil {
		return err
	}

	return nil
}

func initPeerStore(ctx context.Context) error {
	pstore := peerstore.NewPeerStore(
		ctx,
		time.Duration(conf.GetInt("p2p.peer_ttl"))*time.Second)
	gPeerStore = pstore

	return nil
}

func initDiscovery(ctx context.Context) error {
	// p2p opts
	var opts []libp2p.Option

	// private key
	privateKeyPath := conf.GetString("p2p.private_key_path")
	privateKeyB64, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Printf("reading private key err: %v, path: %s", err, privateKeyPath)
		return err
	}
	privateKeyBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(privateKeyB64)))
	if err != nil {
		log.Printf("base64 unmarshal err: %v, string: %s", err, string(privateKeyB64))
		return err
	}
	privateKey, err := crypto.UnmarshalPrivateKey(privateKeyBytes)
	if err != nil {
		log.Printf("invalid private key, err: %v, string: %s", err, string(privateKeyBytes))
		return err
	}
	opts = append(opts, libp2p.Identity(privateKey))

	// listen addresses
	opts = append(opts, libp2p.Transport((quic.NewTransport)))
	listenAddrs := conf.GetStringSlice("p2p.listen_multiaddrs")
	if len(listenAddrs) > 0 {
		opts = append(opts, libp2p.ListenAddrStrings(listenAddrs...))
	}

	opts = append(opts,
		libp2p.ForceReachabilityPublic(),
		libp2p.DisableRelay(),
	)

	// NAT service
	opts = append(opts, libp2p.EnableNATService())

	// pnet psk
	pskPath := conf.GetString("p2p.pnet_psk_path")
	if pskPath != "" {
		pskFile, err := os.Open(pskPath)
		if err != nil {
			return err
		}
		defer pskFile.Close()

		psk, err := pnet.DecodeV1PSK(pskFile)
		if err != nil {
			return err
		}

		opts = append(opts, libp2p.PrivateNetwork(psk))
		log.Println("private network is enabled")
	}

	// create libp2p host
	h, err := libp2p.New(opts...)
	if err != nil {
		return err
	}

	// save variables to global
	gHost = h

	// register stream handlers
	gHost.SetStreamHandler(HEARTBEAT_TOPIC, handleHeartbeat)
	gHost.SetStreamHandler(DISCOVERY_TOPIC, handleDiscovery)

	log.Println("PeerID:", h.ID())
	log.Println("Listen Addrs:", h.Addrs())

	return nil
}
