package main

import "fmt"
import (
	"bitbucket.org/chrj/smtpd"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var db *sql.DB
var server *smtpd.Server
var host = "noch.so,loggle.tv"
var smtpPort = 25

var Version = ""
var BuildDate = ""

func main() {
	printVersion()
	db = openDatabase()
	defer db.Close()
	prepareDatabase(db)
	prepareCert()
	server := prepareServer()
	log.Println("Starting smtpd server")
	server.ListenAndServe(fmt.Sprintf(":%d", smtpPort))
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
