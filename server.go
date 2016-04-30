package main

import (
	"bitbucket.org/chrj/smtpd"
	"crypto/tls"
	"fmt"
	"github.com/jbrodriguez/mlog"
	"github.com/nochso/smtpd/models"
	"path"
)

func prepareServer() *smtpd.Server {
	mlog.Info("Loading TLS certificate")
	cert, err := tls.LoadX509KeyPair(path.Join(dataDir, "cert.pem"), path.Join(dataDir, "key.pem"))
	if err != nil {
		mlog.Fatalf("Cert load failed: %v", err)
	}
	server = &smtpd.Server{
		Handler: handle,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}
	return server
}

func handle(peer smtpd.Peer, env smtpd.Envelope) error {
	fmt.Printf("Sender: %s\nRecipients: %s\nContent:\n%s\n-----\n", env.Sender, env.Recipients[0], env.Data)
	mlog.Info("%d", getAddressId(env.Sender))
	for _, recp := range env.Recipients {
		mlog.Info("%d", getAddressId(recp))
	}
	return nil
}

func getAddressId(address string) int {
	addr, err := models.AddressByAddress(db, address)
	if err != nil {
		addr = &models.Address{Address: address}
		addr.Insert(db)
		return addr.ID
	}
	return addr.ID
}
