package main

import (
	"bitbucket.org/chrj/smtpd"
	"crypto/tls"
	"fmt"
	"github.com/nochso/smtpd/models"
	"log"
)

func prepareServer() *smtpd.Server {
	log.Println("Loading TLS certificate")
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("Cert load failed: %v", err)
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
	log.Println(getAddressId(env.Sender))
	for _, recp := range env.Recipients {
		log.Println(getAddressId(recp))
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
