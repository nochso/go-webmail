package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/mail"
	"path"
	"strings"
	"time"

	"bitbucket.org/porkbonk/smtpd"
	"github.com/nochso/go-webmail/model"
	"github.com/nochso/mlog"
)

func prepareServer() *smtpd.Server {
	mlog.Trace("Loading TLS certificate")
	cert, err := tls.LoadX509KeyPair(path.Join(dataDir, "cert.pem"), path.Join(dataDir, "key.pem"))
	if err != nil {
		mlog.Fatalf("Cert load failed: %v", err)
	}
	server = &smtpd.Server{
		ConnectionChecker: handleConnection,
		RecipientChecker:  handleRcpt,
		Handler:           handle,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
		ProtocolLogger: mlog.Logger.Trace,
	}
	return server
}

func handleConnection(peer smtpd.Peer) error {
	mlog.Trace("Connection accepted: remote_host=%s", peer.Addr)
	return nil
}

// handleRcpt makes sure to reject unknown recipients. Error 550 when rejected.
func handleRcpt(peer smtpd.Peer, addr string) error {
	if addressBelongsToHost(&mail.Address{Address: addr}) {
		return nil
	}
	mlog.Warning("Rejecting mail addressed to '%s' with error 550", addr)
	return smtpd.Error{Code: 550, Message: "No such user here"}
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
	mlog.Trace(
		"Incoming mail: remote_host=%s protocol=%s helo_name=%s%s",
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
	recipients, err := mail.ParseAddressList(header.Get("To"))
	if err != nil {
		mlog.Warning("Unable to parse 'To' header '%s': %s", header.Get("to"), err)
	}
	allowedRecipients := filterAddressesByAllowedHosts(recipients)
	if len(allowedRecipients) == 0 {
		mlog.Warning("Ignoring mail: none of the recipient domains are allowed: %v", recipients)
		return nil
	}
	mlog.Trace("Saving %d mail(s) for recipient(s): %v", len(allowedRecipients), allowedRecipients)
	for _, recipient := range allowedRecipients {
		mailRow := model.Mail{
			Content:    string(env.Data),
			TsReceived: time.Now().Unix(),
		}
		err = mailRow.Save(db)
		toRow := model.AddressTo{MailID: mailRow.ID, AddressID: getAddressId(recipient)}
		toRow.Insert(db)
		fromRow := model.AddressFrom{MailID: mailRow.ID, AddressID: getAddressId(sender)}
		fromRow.Save(db)
		if err != nil {
			mlog.Error(errors.New("Unable to insert mail in database: " + err.Error()))
		}
	}
	return nil
}

// filterAddressesByAllowedHosts returns all mail addresses are allowed according to smtp.accept.domains
func filterAddressesByAllowedHosts(addresses []*mail.Address) []*mail.Address {
	allowed := make([]*mail.Address, 0)
	for _, address := range addresses {
		if addressBelongsToHost(address) {
			allowed = append(allowed, address)
		}
	}
	return allowed
}

// addressBelongsToHost returns true if the address belongs to one of the allowed hosts
func addressBelongsToHost(addr *mail.Address) bool {
	for _, domain := range cfg.Domains {
		if strings.HasPrefix(domain, ".") {
			if strings.HasSuffix(addr.Address, domain) {
				return true
			}
		} else if strings.HasSuffix(addr.Address, "@"+domain) {
			return true
		} else if domain == "*" {
			return true
		}
	}
	return false
}
