package main

import "fmt"
import (
	"bitbucket.org/chrj/smtpd"
	"crypto/tls"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nochso/smtpd/models"
	"log"
)

var db *sql.DB
var server *smtpd.Server
var host = "noch.so,loggle.tv"

var Version = ""
var BuildDate = ""

func main() {
	printVersion()
	log.Println("Opening SQLite database")
	dbPath := "./mail.sqlite"
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Unable to open or create SQLite database file '%s': %s", dbPath, err)
	}
	defer db.Close()
	prepareDatabase(db)
	prepareCert()
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
	log.Println("Starting smtpd server")
	server.ListenAndServe(":25")
}

func printVersion() {
	fmt.Print("noch.so smtpd")
	if Version != "" {
		fmt.Printf(" %s", Version)
	}
	if BuildDate != "" {
		fmt.Printf(" built %s", BuildDate)
	}
	fmt.Print("\n\n")
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

func handle(peer smtpd.Peer, env smtpd.Envelope) error {
	fmt.Printf("Sender: %s\nRecipients: %s\nContent:\n%s\n-----\n", env.Sender, env.Recipients[0], env.Data)
	log.Println(getAddressId(env.Sender))
	for _, recp := range env.Recipients {
		log.Println(getAddressId(recp))
	}
	return nil
}
