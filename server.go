package main

import (
	"bitbucket.org/chrj/smtpd"
	"crypto/tls"
	"fmt"
	"github.com/jbrodriguez/mlog"
	"path"
)

func prepareServer() *smtpd.Server {
	mlog.Info("Loading TLS certificate")
	cert, err := tls.LoadX509KeyPair(path.Join(dataDir, "cert.pem"), path.Join(dataDir, "key.pem"))
	if err != nil {
		mlog.Fatalf("Cert load failed: %v", err)
	}
	server = &smtpd.Server{
		ConnectionChecker: handleConnection,
		Handler:           handle,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}
	return server
}

func handleConnection(peer smtpd.Peer) error {
	mlog.Info("Connection accepted: remote_host=%s", peer.Addr)
	return nil
}

func handle(peer smtpd.Peer, env smtpd.Envelope) error {
	fmt.Printf("Sender: %s\nRecipients: %s\nContent:\n%s\n-----\n", env.Sender, env.Recipients[0], env.Data)
	mlog.Info("%d", getAddressId(env.Sender))
	for _, recp := range env.Recipients {
		mlog.Info("%d", getAddressId(recp))
	}
	return nil
}
