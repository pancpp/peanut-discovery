package app

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"strings"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/pnet"
	quic "github.com/libp2p/go-libp2p/p2p/transport/quic"
	"github.com/pancpp/peanut-discovery/conf"
)

func newHost(_ context.Context) (host.Host, error) {
	// p2p opts
	var opts []libp2p.Option

	// private key
	privateKeyPath := conf.GetString("p2p.private_key_path")
	privateKeyB64, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Printf("reading private key err: %v, path: %s", err, privateKeyPath)
		return nil, err
	}
	privateKeyBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(privateKeyB64)))
	if err != nil {
		log.Printf("base64 unmarshal err: %v, string: %s", err, string(privateKeyB64))
		return nil, err
	}
	privateKey, err := crypto.UnmarshalPrivateKey(privateKeyBytes)
	if err != nil {
		log.Printf("invalid private key, err: %v, string: %s", err, string(privateKeyBytes))
		return nil, err
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
			return nil, err
		}
		defer pskFile.Close()

		psk, err := pnet.DecodeV1PSK(pskFile)
		if err != nil {
			return nil, err
		}

		opts = append(opts, libp2p.PrivateNetwork(psk))
		log.Println("private network is enabled")
	}

	// create libp2p host
	h, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}

	return h, nil
}
