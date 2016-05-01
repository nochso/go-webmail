package main

import (
	"bitbucket.org/porkbonk/smtpd"
	"crypto/tls"
	"fmt"
	"github.com/nochso/mlog"
	"github.com/nochso/smtpd/models"
	"net/mail"
	"path"
	"strings"
	"time"
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
		ProtocolLogger: mlog.Logger.Trace,
	}
	return server
}

func handleConnection(peer smtpd.Peer) error {
	mlog.Info("Connection accepted: remote_host=%s", peer.Addr)
	return nil
}

func handle(peer smtpd.Peer, env smtpd.Envelope) error {
	tlsInfo := ""
	if peer.TLS != nil {
		tlsVersions := map[uint16]string{
			tls.VersionSSL30: "SSL3.0",
			tls.VersionTLS10: "TLS1.0",
			tls.VersionTLS11: "TLS1.1",
			tls.VersionTLS12: "TLS1.2",
		}
		tlsInfo = fmt.Sprintf(" tls_version=%s tls_cypher=0x%x", tlsVersions[peer.TLS.Version], peer.TLS.CipherSuite)
	}
	mlog.Info(
		"Accepting mail: remote_host=%s protocol=%s helo_name=%s%s",
		peer.Addr,
		peer.Protocol,
		peer.HeloName,
		tlsInfo,
	)
	r := strings.NewReader(string(env.Data))
	m, err := mail.ReadMessage(r)
	if err != nil {
		mlog.Fatal(err)
	}
	header := m.Header
	sender, err := mail.ParseAddress(header.Get("From"))
	if err != nil {
		sender = &mail.Address{Address: header.Get("From")}
	}
	mailRow := models.Mail{
		SenderID:   getAddressId(sender),
		Content:    string(env.Data),
		TsReceived: time.Now().Unix(),
		Subject:    header.Get("Subject"),
	}
	mailRow.Save(db)
	recipients, err := mail.ParseAddressList(header.Get("To"))
	if err != nil {
		mlog.Warning("Unable to parse 'To' header '%s': %s", header.Get("to"), err)
	}
	for _, recipient := range recipients {
		recipientRow := &models.MailRecipient{
			MailID:      mailRow.ID,
			RecipientID: getAddressId(recipient),
		}
		err := recipientRow.Save(db)
		if err != nil {
			mlog.Warning("Error saving recipient %s", err)
		}
	}
	return nil
}
